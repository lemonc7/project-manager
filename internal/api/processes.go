package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/util"
)

type ProcessService interface {
	InsertProcessLast(ctx context.Context, req model.InsertProcessRequest) (model.Process, error)
	InsertProcessAt(ctx context.Context, sequence int32, req model.InsertProcessRequest) error
	UpdateProcess(ctx context.Context, id int32, req model.UpdateProcessRequest) (model.Process, error)
	DeleteProcess(ctx context.Context, id int32) error
	ListProcesses(ctx context.Context, projectID int32) ([]model.Process, error)
	ListTasksTodo(ctx context.Context, employeeID int32) ([]model.Task, error)
	ListTasksReady(ctx context.Context, employeeID int32) ([]model.Task, error)
	ListTasksDoing(ctx context.Context, employeeID int32) ([]model.Task, error)
}

var _ ProcessService = (*service.ProcessService)(nil)

type ProcessHandler struct {
	processService ProcessService
}

func NewProcessHandler(service ProcessService) *ProcessHandler {
	return &ProcessHandler{
		processService: service,
	}
}

func (p *ProcessHandler) InsertProcessLast(c echo.Context) error {
	// 绑定请求
	var req model.InsertProcessRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	// 业务函数
	process, err := p.processService.InsertProcessLast(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusCreated, process)
}

func (p *ProcessHandler) InsertProcessAt(c echo.Context) error {
	sequence, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	// 绑定请求
	var req model.InsertProcessRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	// 业务函数
	if err := p.processService.InsertProcessAt(c.Request().Context(), sequence, req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "成功创建流程",
	})
}

func (p *ProcessHandler) UpdateProcess(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	// 绑定请求
	var req model.UpdateProcessRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// 业务函数
	process, err := p.processService.UpdateProcess(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, process)
}

func (p *ProcessHandler) DeleteProcess(c echo.Context) error {
	// 获取id
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	// 业务函数
	if err := p.processService.DeleteProcess(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功删除流程",
	})
}

func (p *ProcessHandler) ListProcesses(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	// 获取数据
	processes, err := p.processService.ListProcesses(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, processes)
}

func (p *ProcessHandler) ListTasks(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	statusStr := c.QueryParam("status")
	switch statusStr {
	case "todo":
		tasks, err := p.processService.ListTasksTodo(c.Request().Context(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, tasks)
	case "ready":
		tasks, err := p.processService.ListTasksReady(c.Request().Context(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, tasks)
	case "doing":
		tasks, err := p.processService.ListTasksDoing(c.Request().Context(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, tasks)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("status参数错误: %s", statusStr))
	}
}
