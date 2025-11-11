package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
)

type ProjectService interface {
	CreateProject(ctx context.Context, req model.CreateProjectRequest) error
	UpdateProject(ctx context.Context, req model.UpdateProjectRequest) error
	DeleteProject(ctx context.Context, req model.DeleteProjectRequest) error
	ListProjectsFiltered(ctx context.Context, req model.ListProjectsRequest) ([]repo.ListProjectsFilteredRow, error)
}

var _ ProjectService = (*service.ProjectService)(nil)

type ProjectHandler struct {
	projectService ProjectService
}

func (p *ProjectHandler) CreateProject(c echo.Context) error {
	// 绑定请求
	var req model.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 业务函数
	if err := p.projectService.CreateProject(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "成功创建项目",
	})
}

func (p *ProjectHandler) UpdateProject(c echo.Context) error {
	// 绑定请求
	var req model.UpdateProjectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 业务函数
	if err := p.projectService.UpdateProject(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功更新项目",
	})
}

func (p *ProjectHandler) DeleteProject(c echo.Context) error {
	// 绑定请求
	var req model.DeleteProjectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 删除数据
	if err := p.projectService.DeleteProject(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功删除项目",
	})
}

func (p *ProjectHandler) ListProjectsFiltered(c echo.Context) error {
	// 绑定请求
	var req model.ListProjectsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 获取数据
	projects, err := p.projectService.ListProjectsFiltered(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"data": projects,
	})
}
