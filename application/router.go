package application

func (a *Application) initRouter() {
	employeeRouter := a.e.Group("/api/employee")
	teamRouter := a.e.Group("/api/team")
	projectRouter := a.e.Group("/api/project")
	processRouter := a.e.Group("/api/process")

	// 员工
	employeeRouter.POST("", a.employeeHandler.CreateEmployee)
	employeeRouter.GET("", a.employeeHandler.ListEmployees)
	employeeRouter.PUT("/:id", a.employeeHandler.UpdateEmployee)
	employeeRouter.DELETE("/:id", a.employeeHandler.DeleteEmployee)

	// 团队
	teamRouter.POST("", a.teamHandler.CreateTeam)
	teamRouter.GET("", a.teamHandler.ListTeams)
	teamRouter.PUT("/:id", a.teamHandler.UpdateTeam)
	teamRouter.DELETE("/:id", a.teamHandler.DeleteTeam)

	// 项目
	projectRouter.GET("/:id", a.projectHandler.ListProcesses)
	projectRouter.GET("/filter", a.projectHandler.ListProjectsFiltered)
	projectRouter.POST("", a.projectHandler.CreateProject)
	projectRouter.PUT("/:id", a.projectHandler.UpdateProject)
	projectRouter.DELETE("/:id", a.projectHandler.DeleteProject)

	// 流程
	processRouter.GET("/filter", a.processHandler.ListProcessesFiltered)
	processRouter.POST("", a.processHandler.InsertProcessLast)
	processRouter.POST("/:id", a.processHandler.InsertProcessAt)
	processRouter.PUT("/:id", a.processHandler.UpdateProcess)
	processRouter.DELETE("/:id", a.processHandler.DeleteProcess)
	processRouter.GET("/task/:id", a.processHandler.ListTasks)
}
