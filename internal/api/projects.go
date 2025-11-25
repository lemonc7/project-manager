package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/util"
)

type ProjectService interface {
	CreateProject(ctx context.Context, req model.CreateProjectRequest) (model.Project, error)
	UpdateProject(ctx context.Context, id int32, req model.UpdateProjectRequest) (model.Project, error)
	DeleteProject(ctx context.Context, id int32) error
	ListProjectsFiltered(ctx context.Context, req model.ListProjectsRequest) ([]model.Project, error)
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
	project, err := p.projectService.CreateProject(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusCreated, project)
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

	// 业务函数
	project, err := p.projectService.UpdateProject(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, project)
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

func (p *ProjectHandler) ListProjectsFiltered(c echo.Context) error {
	// 绑定请求
	var req model.ListProjectsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// 获取数据
	projects, err := p.projectService.ListProjectsFiltered(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, projects)
}
