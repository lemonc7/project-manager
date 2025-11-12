package application

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/api"
)

type Application struct {
	e  *echo.Echo
	db *sql.DB

	employeeHandler *api.EmployeeHandler
	teamHandler     *api.TeamHandler
	projectHandler  *api.ProjectHandler
	processHandler  *api.ProcessHandler
}

func InitApp(filePath string) {

}
