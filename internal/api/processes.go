package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/util"
)

type ProcessService interface {
	InsertProcessLast(ctx context.Context, req model.InsertProcessRequest) error
	InsertProcessAt(ctx context.Context, sequence int32, req model.InsertProcessRequest) error
	UpdateProcess(ctx context.Context, id int32, req model.UpdateProcessRequest) error
	DeleteProcess(ctx context.Context, id int32) error
	DeleteProcesses(ctx context.Context, projectID int32, req model.DeleteProcessesRequest) error
	ListProcesses(ctx context.Context, projectID int32) ([]repo.ListProcessesRow, error)
	ListTasksTodo(ctx context.Context, employeeID int32) ([]repo.ListTasksTodoRow, error)
	ListTasksReady(ctx context.Context, employeeID int32) ([]repo.ListTasksReadyRow, error)
	ListTasksDoing(ctx context.Context, employeeID int32) ([]repo.ListTasksDoingRow, error)
	FinishProcess(ctx context.Context, id int32, req model.FinishProcessRequest) error
	ReworkProcess(ctx context.Context, id int32, req model.ReworkProcessRequest) error
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
	if err := util.BindRequest(
		c.Request().Body,
		model.InsertProcessRequestSchema,
		&req,
	); err != nil {
		return err
	}

	// 业务函数
	err := p.processService.InsertProcessLast(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("插入process: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "成功创建流程",
	})
}

func (p *ProcessHandler) InsertProcessAt(c echo.Context) error {
	sequence, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	if sequence < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("sequence参数错误"))
	}

	// 绑定请求
	var req model.InsertProcessRequest
	if err := util.BindRequest(
		c.Request().Body,
		model.InsertProcessRequestSchema,
		&req,
	); err != nil {
		return err
	}

	// 业务函数
	if err := p.processService.InsertProcessAt(c.Request().Context(), sequence, req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("插入process: %w", err))
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
	if err := util.BindRequest(
		c.Request().Body,
		model.UpdateEmployeeSchema,
		&req,
	); err != nil {
		return err
	}

	// 业务函数
	err = p.processService.UpdateProcess(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("更新process: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功更新流程",
	})
}

func (p *ProcessHandler) DeleteProcess(c echo.Context) error {
	// 获取id
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	// 业务函数
	if err := p.processService.DeleteProcess(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("删除process: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功删除流程",
	})
}

func (p *ProcessHandler) DeleteProcesses(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	var req model.DeleteProcessesRequest
	if err := util.BindRequest(
		c.Request().Body,
		model.DeleteProcessesRequestSchema,
		&req,
	); err != nil {
		return err
	}

	if err := p.processService.DeleteProcesses(c.Request().Context(), id, req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("删除流程: %w", err))
	}

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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取process: %w", err))
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
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取task todo: %w", err))
		}
		return c.JSON(http.StatusOK, tasks)
	case "ready":
		tasks, err := p.processService.ListTasksReady(c.Request().Context(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取task ready: %w", err))
		}
		return c.JSON(http.StatusOK, tasks)
	case "doing":
		tasks, err := p.processService.ListTasksDoing(c.Request().Context(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取task doing: %w", err))
		}
		return c.JSON(http.StatusOK, tasks)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("status参数错误: %s", statusStr))
	}
}

func (p *ProcessHandler) FinishProcess(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	var req model.FinishProcessRequest
	if err := util.BindRequest(
		c.Request().Body,
		model.FinishProcessRequestSchema,
		&req,
	); err != nil {
		return err
	}

	if err := p.processService.FinishProcess(c.Request().Context(), id, req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("完成process: %w", err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功完成流程",
	})
}

func (p *ProcessHandler) ReworkProcess(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	// 绑定请求
	var req model.ReworkProcessRequest
	if err := util.BindRequest(
		c.Request().Body,
		model.ReworkProcessRequestSchema,
		&req,
	); err != nil {
		return err
	}

	// 业务函数
	if err := p.processService.ReworkProcess(c.Request().Context(), id, req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("返工process: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功创建返工流程",
	})
}
