package api

import (
	"context"
	"net/http"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/zest"
)

type ProcessService interface {
	InsertProcessLast(ctx context.Context, req model.InsertProcessRequest) error
	InsertProcessAt(ctx context.Context, req model.InsertProcessRequest) error
	UpdateProcess(ctx context.Context, req model.UpdateProcessRequest) error
	DeleteProcess(ctx context.Context, req model.DeleteProcessRequest) error
	DeleteProcesses(ctx context.Context, req model.DeleteProcessesRequest) error
	ListProcesses(ctx context.Context, req model.ListProcessesRequest) ([]repo.ListProcessesRow, error)
	ListTasks(ctx context.Context, req model.ListTasksRequest) ([]repo.ListTasksRow, error)
	FinishProcess(ctx context.Context, req model.FinishProcessRequest) error
	ReworkProcess(ctx context.Context, req model.ReworkProcessRequest) error
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

func (p *ProcessHandler) InsertProcessLast(c *zest.Context) error {
	// 绑定请求
	var req model.InsertProcessRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// 业务函数
	err := p.processService.InsertProcessLast(c.Context(), req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "创建流程失败").Wrap(err)
	}

	// 返回响应
	return c.JSON(http.StatusCreated, zest.Map{
		"message": "成功创建流程",
	})
}

func (p *ProcessHandler) InsertProcessAt(c *zest.Context) error {
	// 绑定请求
	var req model.InsertProcessRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// 业务函数
	if err := p.processService.InsertProcessAt(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "插入流程失败").Wrap(err)
	}

	return c.JSON(http.StatusCreated, zest.Map{
		"message": "成功创建流程",
	})
}

func (p *ProcessHandler) UpdateProcess(c *zest.Context) error {
	// 绑定请求
	var req model.UpdateProcessRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// 业务函数
	if err := p.processService.UpdateProcess(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "更新流程失败").Wrap(err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功更新流程",
	})
}

func (p *ProcessHandler) DeleteProcess(c *zest.Context) error {
	var req model.DeleteProcessRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// 业务函数
	if err := p.processService.DeleteProcess(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "删除流程失败").Wrap(err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功删除流程",
	})
}

func (p *ProcessHandler) DeleteProcesses(c *zest.Context) error {
	var req model.DeleteProcessesRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := p.processService.DeleteProcesses(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "删除流程失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功删除流程",
	})
}

func (p *ProcessHandler) ListProcesses(c *zest.Context) error {
	var req model.ListProcessesRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// 获取数据
	processes, err := p.processService.ListProcesses(c.Context(), req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "获取流程信息失败").Wrap(err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, processes)
}

func (p *ProcessHandler) ListTasks(c *zest.Context) error {
	var req model.ListTasksRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	tasks, err := p.processService.ListTasks(c.Context(), req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "获取todo任务失败").Wrap(err)
	}
	return c.JSON(http.StatusOK, tasks)
}

func (p *ProcessHandler) FinishProcess(c *zest.Context) error {
	var req model.FinishProcessRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := p.processService.FinishProcess(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "完成流程失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功完成流程",
	})
}

func (p *ProcessHandler) ReworkProcess(c *zest.Context) error {
	// 绑定请求
	var req model.ReworkProcessRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// 业务函数
	if err := p.processService.ReworkProcess(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "返工流程失败").Wrap(err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功创建返工流程",
	})
}
