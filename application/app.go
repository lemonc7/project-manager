package application

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/config"
	"github.com/lemonc7/project-manager/database"
	"github.com/lemonc7/project-manager/internal/api"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/logger"
	"go.uber.org/zap"
)

type Application struct {
	e   *echo.Echo
	db  *sql.DB
	cfg *config.Config

	employeeHandler *api.EmployeeHandler
	teamHandler     *api.TeamHandler
	projectHandler  *api.ProjectHandler
	processHandler  *api.ProcessHandler
}

type CustomValidator struct {
	validator *validator.Validate
}

func (c *CustomValidator) Validate(i any) error {
	return c.validator.Struct(i)
}

func InitApp(filePath string) (*Application, error) {
	// 读取配置
	cfg, err := config.NewConfig(filePath)
	if err != nil {
		return nil, err
	}

	// 初始化logger
	logger.InitLogger(cfg.Logger)

	// 连接数据库
	db, err := database.NewDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	// 创建echo
	validator := &CustomValidator{validator: validator.New()}

	e := NewEcho(cfg.Server, validator)

	// handler
	employeeHandler := api.NewEmployeeHandler(
		service.NewEmployeeService(db),
	)
	teamHandler := api.NewTeamHandler(
		service.NewTeamService(db),
	)
	projectHandler := api.NewProjectHandler(
		service.NewProjectService(db),
	)
	processHandler := api.NewProcessHandler(
		service.NewProcessService(db),
	)

	app := &Application{
		e:               e,
		db:              db,
		cfg:             cfg,
		employeeHandler: employeeHandler,
		teamHandler:     teamHandler,
		projectHandler:  projectHandler,
		processHandler:  processHandler,
	}
	// 注册路由
	app.initRouter()
	return app, nil
}

func (a *Application) Start() {
	go func() {
		if err := a.e.Start(":" + strconv.Itoa(a.cfg.Server.Port)); err != nil {
			log.Fatal(err.Error())
		}
	}()

	a.shutdown()
}

func (a *Application) shutdown() {
	defer func() {
		if err := a.db.Close(); err != nil {
			logger.Error("close database", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()
	if err := a.e.Shutdown(ctx); err != nil {
		logger.Error("shutdown server", zap.Error(err))
	}else {
		logger.Info("HTTP server shutdown successfully")
	}

	logger.Info("shutdown successfully")
}
