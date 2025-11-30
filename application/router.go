package application

import "github.com/lemonc7/project-manager/internal/middleware"

func (a *Application) initRouter() {
	// 检查服务器
	a.e.GET("/ping", a.serverHandler.Ping)

	// 用户认证
	jwtMW := middleware.JWTAuther(a.jwt)
	auth := a.e.Group("/api/auth")
	authJWT := auth.Group("", jwtMW)

	auth.POST("/login", a.authHandler.Login)
	auth.POST("/init_admin", a.authHandler.InitAdmin)
	authJWT.POST("/password", a.authHandler.ResetPassword)
	authJWT.PUT("/password", a.authHandler.ChangePassword)

	// 业务
	api := a.e.Group("/api", jwtMW)

	employeeRouter := api.Group("/employee")
	teamRouter := api.Group("/team")
	projectRouter := api.Group("/project")
	processRouter := api.Group("/process")

	// 员工
	employeeRouter.GET("", a.employeeHandler.ListEmployees)
	employeeRouter.GET("/me", a.employeeHandler.GetEmployee)
	employeeRouter.POST("", a.employeeHandler.CreateEmployee)
	employeeRouter.PUT("/:id", a.employeeHandler.UpdateEmployee)
	employeeRouter.DELETE("/:id", a.employeeHandler.DeleteEmployee)

	// 团队
	teamRouter.POST("", a.teamHandler.CreateTeam)
	teamRouter.GET("", a.teamHandler.ListTeams)
	teamRouter.PUT("/:id", a.teamHandler.UpdateTeam)
	teamRouter.DELETE("/:id", a.teamHandler.DeleteTeam)

	// 项目
	projectRouter.GET("/:id", a.processHandler.ListProcesses)
	projectRouter.GET("/status/:filter", a.projectHandler.ListProjects)
	projectRouter.POST("", a.projectHandler.CreateProject)
	projectRouter.PUT("/:id", a.projectHandler.UpdateProject)
	projectRouter.DELETE("/:id", a.projectHandler.DeleteProject)

	// 流程
	processRouter.POST("", a.processHandler.InsertProcessLast)
	processRouter.POST("/:id", a.processHandler.InsertProcessAt)
	processRouter.PUT("/:id", a.processHandler.UpdateProcess)
	processRouter.DELETE("/:id", a.processHandler.DeleteProcess)
	processRouter.GET("/task/:id", a.processHandler.ListTasks)
}
