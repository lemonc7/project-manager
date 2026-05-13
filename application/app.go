package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/config"
	"github.com/lemonc7/project-manager/database"
	"github.com/lemonc7/project-manager/internal/api"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/hasher"
	"github.com/lemonc7/project-manager/pkg/jwt"
	"github.com/lemonc7/zest"
	"github.com/lemonc7/zest/middleware"
)

type Application struct {
	zest   *zest.Zest
	server *http.Server
	dbPool *pgxpool.Pool
	cfg    *config.Config
	jwt    *jwt.JWT

	serverHandler   *api.ServerHandler
	authHandler     *api.AuthHandler
	employeeHandler *api.EmployeeHandler
	teamHandler     *api.TeamHandler
	projectHandler  *api.ProjectHandler
	processHandler  *api.ProcessHandler
	todoHandler     *api.TodoHandler
}

func InitApp(filePath string) (*Application, error) {
	// 读取配置
	cfg, err := config.NewConfig(filePath)
	if err != nil {
		return nil, err
	}

	// 连接数据库
	dbPool, err := database.NewDB(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	// 迁移schema
	m, err := migrate.New(
		"file://database/migrate",
		cfg.Database.DSN(),
	)
	if err != nil {
		return nil, fmt.Errorf("schema migration connect: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("schema migration up: %w", err)
	}

	jwt := jwt.NewJWT(cfg.JWT)

	// 注册Zest Handler
	z := zest.New()
	z.Use(middleware.RequestID())
	z.Use(middleware.Logger(middleware.LoggerConfig{
		TZ: cfg.Logger.TZ,
	}))
	z.Use(middleware.Recovery())
	z.Use(middleware.CORS())
	z.Use(middleware.JWT(jwt, func(c *zest.Context) bool {
		return c.Path == "/login" || c.Path == "/init_admin" || c.Path == "/ping"
	}))

	// 注册Http Server
	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.Server.Port),
		Handler:      z,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// handler
	passwordHash := hasher.NewPasswordHash()

	authHandler := api.NewAuthHandler(
		service.NewAuthService(dbPool, passwordHash, jwt),
	)
	employeeHandler := api.NewEmployeeHandler(
		service.NewEmployeeService(dbPool, passwordHash),
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
	todoHandler := api.NewTodoHandler(
		service.NewTodoService(dbPool),
	)

	serverHandler := api.NewServerService(
		service.NewServerService(dbPool),
	)

	app := &Application{
		zest:   z,
		dbPool: dbPool,
		cfg:    cfg,
		jwt:    jwt,
		server: srv,

		serverHandler:   serverHandler,
		authHandler:     authHandler,
		employeeHandler: employeeHandler,
		teamHandler:     teamHandler,
		projectHandler:  projectHandler,
		processHandler:  processHandler,
		todoHandler:     todoHandler,
	}
	// 注册路由
	app.initRouter()
	return app, nil
}

func (a *Application) Start() {
	// 最后关闭数据库连接
	defer func() {
		a.dbPool.Close()
		log.Println("成功关闭数据库连接")
	}()

	// 启动HTTP服务
	log.Printf("HTTP服务已启动，监听端口 :%d", a.cfg.Server.Port)
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("启动HTTP服务失败", err.Error())
			}
		}
	}()

	// 阻塞监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("收到退出信号，开始关闭服务...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	// 关闭HTTP服务
	if err := a.server.Shutdown(ctx); err != nil {
		log.Println("关闭HTTP服务失败: ", err)
	} else {
		log.Println("成功关闭HTTP服务")
	}
}
