package api

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/util"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, req model.CreateEmployeeRequest) error
	UpdateEmployee(ctx context.Context, id int32, req model.UpdateEmployeeRequest) error
	DeleteEmployee(ctx context.Context, id int32) error
	ListEmployees(ctx context.Context) ([]repo.ListEmployeesRow, error)
	GetEmployee(ctx context.Context, username string) (string, error)
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
	// 获取jwt中的角色
	roles, ok := c.Get("roles").([]string)
	if !ok || !slices.Contains(roles, "admin") {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("需要管理员权限"))
	}
	// 绑定请求
	var req model.CreateEmployeeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("绑定参数: %w", err))
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Errorf("校验参数: %w", err))
	}

	// 没传roles设置为空数组
	if req.Roles == nil {
		req.Roles = []string{}
	}

	// 业务函数
	err := e.employeeService.CreateEmployee(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("创建employee: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "成功创建员工",
	})
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
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("绑定参数: %w", err))
	}

	// 业务函数
	err = e.employeeService.UpdateEmployee(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("更新employee: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功更新员工",
	})
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("删除employee: %w", err))
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取employees: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusOK, employees)
}

// 获取用户
func (e *EmployeeHandler) GetEmployeeName(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("未授权"))
	}
	employeeName, err := e.employeeService.GetEmployee(c.Request().Context(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取employee姓名: %w", err))
	}

	return c.JSON(http.StatusOK, employeeName)
}
