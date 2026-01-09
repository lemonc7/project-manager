package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/util"
	"github.com/lemonc7/zest"
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
func (e *EmployeeHandler) CreateEmployee(c *zest.Context) error {
	// 绑定请求
	var req model.CreateEmployeeRequest
	if err := c.Bind(model.CreateEmployeeSchema, &req); err != nil {
		return err
	}

	// 业务函数
	err := e.employeeService.CreateEmployee(c.Context(), req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("创建employee: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusCreated, zest.Map{
		"message": "成功创建员工",
	})
}

// 更新员工
func (e *EmployeeHandler) UpdateEmployee(c *zest.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	// 绑定请求
	var req model.UpdateEmployeeRequest
	if err := c.Bind(model.UpdateEmployeeSchema, &req); err != nil {
		return err
	}

	// 业务函数
	err = e.employeeService.UpdateEmployee(c.Context(), id, req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("更新employee: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功更新员工",
	})
}

// 删除员工
func (e *EmployeeHandler) DeleteEmployee(c *zest.Context) error {
	// 获取员工id
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	// 业务函数
	if err := e.employeeService.DeleteEmployee(c.Context(), id); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("删除employee: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功删除员工",
	})
}

// 获取员工
func (e *EmployeeHandler) ListEmployees(c *zest.Context) error {
	// 获取数据
	employees, err := e.employeeService.ListEmployees(c.Context())
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取employees: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusOK, employees)
}

// 获取用户
func (e *EmployeeHandler) GetEmployeeName(c *zest.Context) error {
	username, ok := c.Get("username").(string)
	if !ok {
		return zest.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("未授权"))
	}
	employeeName, err := e.employeeService.GetEmployee(c.Context(), username)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取employee姓名: %w", err))
	}

	return c.JSON(http.StatusOK, employeeName)
}
