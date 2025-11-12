package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
)

type ProcessService interface {
	CreateProcess(ctx context.Context, req model.CreateProcessRequest) error
	UpdateProcess(ctx context.Context, req model.UpdateProcessRequest) error
	DeleteProcess(ctx context.Context, req model.DeleteProcessRequest) error
	ListProcessesFiltered(ctx context.Context, req model.ListProcessesRequest) ([]repo.ListProcessesFilteredRow, error)
}

var _ ProcessService = (*service.ProcessService)(nil)

type ProcessHandler struct {
	processService ProcessService
}

func (p *ProcessHandler) CreateProcess(c echo.Context) error {
	// 绑定请求
	var req model.CreateProcessRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 业务函数
	if err := p.processService.CreateProcess(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "成功创建流程",
	})
}

func (p *ProcessHandler) UpdateProcess(c echo.Context) error {
	// 绑定请求
	var req model.UpdateProcessRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 业务函数
	if err := p.processService.UpdateProcess(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功更新流程",
	})
}

func (p *ProcessHandler) DeleteProcess(c echo.Context) error {
	// 绑定请求
	var req model.DeleteProcessRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 业务函数
	if err := p.processService.DeleteProcess(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功删除流程",
	})
}

func (p *ProcessHandler) ListProcessFiltered(c echo.Context) error {
	// 绑定请求
	var req model.ListProcessesRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 获取数据
	processes, err := p.processService.ListProcessesFiltered(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"data": processes,
	})
}
