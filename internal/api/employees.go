package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/util"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, req model.CreateEmployeeRequest) (repo.Employee, error)
	UpdateEmployee(ctx context.Context, id int32, req model.UpdateEmployeeRequest) (repo.Employee, error)
	DeleteEmployee(ctx context.Context, id int32) error
	ListEmployees(ctx context.Context) ([]repo.ListEmployeesRow, error)
}

var _ EmployeeService = (*service.EmployeeService)(nil)

type EmployeeHandler struct {
	employeeService EmployeeService
}

func NewEmployeeHandler(service EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: service,
	}
}

// 创建员工
func (e *EmployeeHandler) CreateEmployee(c echo.Context) error {
	// 绑定请求
	var req model.CreateEmployeeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// 业务函数
	employee, err := e.employeeService.CreateEmployee(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusCreated, employee)
}

// 更新员工
func (e *EmployeeHandler) UpdateEmployee(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	// 绑定请求
	var req model.UpdateEmployeeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 业务函数
	employee, err := e.employeeService.UpdateEmployee(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusOK, employee)
}

// 删除员工
func (e *EmployeeHandler) DeleteEmployee(c echo.Context) error {
	// 获取员工id
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	// 业务函数
	if err := e.employeeService.DeleteEmployee(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功删除员工",
	})
}

// 获取员工
func (e *EmployeeHandler) ListEmployees(c echo.Context) error {
	// 获取数据
	employees, err := e.employeeService.ListEmployees(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusOK, employees)
}
