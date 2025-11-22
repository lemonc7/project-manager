package application

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/config"
	"github.com/lemonc7/project-manager/database"
	"github.com/lemonc7/project-manager/internal/api"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/hasher"
	"github.com/lemonc7/project-manager/pkg/jwt"
	"github.com/lemonc7/project-manager/pkg/logger"
	"go.uber.org/zap"
)

type Application struct {
	e      *echo.Echo
	dbPool *pgxpool.Pool
	cfg    *config.Config

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
	dbPool, err := database.NewDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	// 创建echo
	validator := &CustomValidator{validator: validator.New()}

	e := NewEcho(cfg.Server, validator)

	// handler
	passwordHash := hasher.NewPasswordHash()
	jwt := jwt.NewJWT(cfg.JWT)
	employeeHandler := api.NewEmployeeHandler(
		service.NewEmployeeService(dbPool, passwordHash, jwt),
	)
	teamHandler := api.NewTeamHandler(
		service.NewTeamService(dbPool),
	)
	projectHandler := api.NewProjectHandler(
		service.NewProjectService(dbPool),
	)
	processHandler := api.NewProcessHandler(
		service.NewProcessService(dbPool),
	)

	app := &Application{
		e:               e,
		dbPool:          dbPool,
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
		a.dbPool.Close()
		logger.Info("database connection pool closed")
	}()

	// 监听系统信号
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
	} else {
		logger.Info("HTTP server shutdown successfully")
	}

	logger.Info("shutdown successfully")
}
