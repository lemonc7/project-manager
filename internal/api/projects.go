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

type ProjectService interface {
	CreateProject(ctx context.Context, req model.CreateProjectRequest) error
	UpdateProject(ctx context.Context, id int32, req model.UpdateProjectRequest) error
	DeleteProject(ctx context.Context, id int32) error
	ListProjectsFiltered(ctx context.Context, req model.ListProjectsRequest) ([]model.Project, error)
	ListProjectsTodo(ctx context.Context, req model.ListProjectsRequest) ([]model.Project, error)
	ListProjectsDoing(ctx context.Context, req model.ListProjectsRequest) ([]model.Project, error)
	ListProjectsDone(ctx context.Context, req model.ListProjectsRequest) ([]model.Project, error)
}

var _ ProjectService = (*service.ProjectService)(nil)

type ProjectHandler struct {
	projectService ProjectService
}

func NewProjectHandler(service ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: service,
	}
}

func (p *ProjectHandler) CreateProject(c echo.Context) error {
	// 绑定请求
	var req model.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	// 业务函数
	err := p.projectService.CreateProject(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "成功创建项目",
	})
}

func (p *ProjectHandler) UpdateProject(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	// 绑定请求
	var req model.UpdateProjectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	// 业务函数
	err = p.projectService.UpdateProject(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功更新项目",
	})
}

func (p *ProjectHandler) DeleteProject(c echo.Context) error {
	// 获取id
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	// 删除数据
	if err := p.projectService.DeleteProject(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功删除项目",
	})
}

func (p *ProjectHandler) ListProjects(c echo.Context) error {
	// 绑定请求
	var req model.ListProjectsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// 获取状态
	filter := c.Param("filter")

	switch filter {
	case "todo":
		projects, err := p.projectService.ListProjectsTodo(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, projects)
	case "doing":
		projects, err := p.projectService.ListProjectsDoing(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, projects)
	case "done":
		projects, err := p.projectService.ListProjectsDone(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, projects)
	case "filter":
		projects, err := p.projectService.ListProjectsFiltered(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, projects)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("filter参数错误: %s", filter))
	}
}
