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

type ProcessService interface {
	InsertProcessLast(ctx context.Context, req model.InsertProcessRequest) error
	InsertProcessAt(ctx context.Context, id int32, req model.InsertProcessRequest) error
	UpdateProcess(ctx context.Context, id int32, req model.UpdateProcessRequest) error
	DeleteProcess(ctx context.Context, id int32) error
	ListProcessesFiltered(ctx context.Context, req model.ListProcessesRequest) ([]repo.ListProcessesFilteredRow, error)
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// 业务函数
	if err := p.processService.InsertProcessLast(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	// 绑定请求
	var req model.InsertProcessRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// 业务函数
	if err := p.processService.InsertProcessAt(c.Request().Context(), sequence, req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 业务函数
	if err := p.processService.UpdateProcess(c.Request().Context(), id, req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 获取数据
	processes, err := p.processService.ListProcessesFiltered(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusOK, processes)
}
