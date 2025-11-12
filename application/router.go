package application

func (a *Application) initRouter() {
	employeeRouter := a.e.Group("/api/employee")
	teamRouter := a.e.Group("/api/team")
	projectRouter := a.e.Group("/api/project")
	processRouter := a.e.Group("/api/process")

	// 员工
	employeeRouter.POST("", a.employeeHandler.CreateEmployee)
	employeeRouter.GET("", a.employeeHandler.ListEmployees)
	employeeRouter.PUT("", a.employeeHandler.UpdateEmployee)
	employeeRouter.DELETE("", a.employeeHandler.DeleteEmployee)

	// 团队
	teamRouter.POST("", a.teamHandler.CreateTeam)
	teamRouter.GET("", a.teamHandler.ListTeams)
	teamRouter.PUT("", a.teamHandler.UpdateTeam)
	teamRouter.DELETE("", a.teamHandler.DeleteTeam)

	// 项目
	projectRouter.POST("", a.projectHandler.CreateProject)
	projectRouter.GET("", a.projectHandler.ListProjectsFiltered)
	projectRouter.PUT("", a.projectHandler.UpdateProject)
	projectRouter.DELETE("", a.projectHandler.DeleteProject)

	// 流程
	processRouter.POST("", a.processHandler.CreateProcess)
	processRouter.GET("", a.processHandler.ListProcessFiltered)
	processRouter.PUT("", a.processHandler.UpdateProcess)
	processRouter.DELETE("", a.processHandler.DeleteProcess)
}
