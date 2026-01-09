package application

import (
	mw "github.com/lemonc7/project-manager/internal/middleware"
	"github.com/lemonc7/zest/middleware"
)

func (a *Application) initRouter() {
	// 检查服务器
	a.zest.GET("/ping", a.serverHandler.Ping)

	// 用户登录
	a.zest.POST("/login", a.authHandler.Login)
	a.zest.POST("/init_admin", a.authHandler.InitAdmin)

	// 业务
	api := a.zest.Group("/api", mw.RequireRoleAtLeast("employee"))

	api.PUT("/password", a.authHandler.ChangePassword)
	api.POST("/password/{id}", a.authHandler.ResetPassword, mw.RequireRoleAtLeast("admin"))

	employeeRouter := api.Group("/employee")
	teamRouter := api.Group("/team")
	projectRouter := api.Group("/project")
	processRouter := api.Group("/process")
	todoRouter := api.Group("/todo")

	// 员工
	employeeRouter.GET("", a.employeeHandler.ListEmployees)
	employeeRouter.GET("/me", a.employeeHandler.GetEmployeeName)
	employeeRouter.POST("", a.employeeHandler.CreateEmployee, mw.RequireRoleAtLeast("admin"))
	employeeRouter.PATCH("/{id}", a.employeeHandler.UpdateEmployee, mw.RequireRoleAtLeast("admin"))
	employeeRouter.DELETE("/{id}", a.employeeHandler.DeleteEmployee, mw.RequireRoleAtLeast("admin"))

	// 团队
	teamRouter.GET("", a.teamHandler.ListTeams)
	teamRouter.POST("", a.teamHandler.CreateTeam, mw.RequireRoleAtLeast("manager"))
	teamRouter.PATCH("/{id}", a.teamHandler.UpdateTeam, mw.RequireRoleAtLeast("manager"))
	teamRouter.DELETE("/{id}", a.teamHandler.DeleteTeam, mw.RequireRoleAtLeast("manager"))

	// 项目
	projectRouter.GET("/{id}", a.projectHandler.GetProjectByID)
	projectRouter.POST("/list", a.projectHandler.ListProjects)
	projectRouter.GET("/exist", a.projectHandler.IsProjectExists)
	projectRouter.GET("/export/{id}", a.projectHandler.ExportProject)
	projectRouter.POST("", a.projectHandler.CreateProject, mw.RequireRoleAtLeast("manager"))
	projectRouter.POST("/start/{id}", a.projectHandler.StartProject, mw.RequireRoleAtLeast("manager"))
	projectRouter.POST("/reset/{id}", a.projectHandler.ResetStatus, mw.RequireRoleAtLeast("manager"))
	projectRouter.POST("/copy/{id}", a.projectHandler.CopyProject, mw.RequireRoleAtLeast("manager"))
	projectRouter.PATCH("/{id}", a.projectHandler.UpdateProject, mw.RequireRoleAtLeast("manager"))
	projectRouter.DELETE("/{id}", a.projectHandler.DeleteProject, mw.RequireRoleAtLeast("manager"))

	// 流程
	processRouter.GET("/{id}", a.processHandler.ListProcesses)
	processRouter.GET("/task/{id}", a.processHandler.ListTasks)
	processRouter.POST("", a.processHandler.InsertProcessLast)
	processRouter.POST("/{id}", a.processHandler.InsertProcessAt)
	processRouter.POST("/rework/{id}", a.processHandler.ReworkProcess)
	processRouter.POST("/finish/{id}", a.processHandler.FinishProcess)
	processRouter.PATCH("/{id}", a.processHandler.UpdateProcess)
	processRouter.DELETE("/{id}", a.processHandler.DeleteProcess)
	processRouter.DELETE("/series/{id}", a.processHandler.DeleteProcesses)

	// todo
	todoRouter.GET("/todos/reporter/{id}", a.todoHandler.ListTodosByReporter)
	todoRouter.GET("/todos/assignee/{id}", a.todoHandler.ListTodosByAssignee)
	todoRouter.POST("", a.todoHandler.CreateTodo)
	todoRouter.PATCH("/{id}", a.todoHandler.UpdateTodo)
	todoRouter.DELETE("/{id}", a.todoHandler.DeleteTodo)

	// // 挂载静态文件
	a.zest.Use(middleware.Static(middleware.StaticConfig{
		Root:   "./dist",
		HTML5:  true,
		Browse: false,
	}))
}
