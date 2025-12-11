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

type ProjectService interface {
	CreateProject(ctx context.Context, req model.CreateProjectRequest) error
	UpdateProject(ctx context.Context, id int32, req model.UpdateProjectRequest) error
	DeleteProject(ctx context.Context, id int32) error
	ListProjectsFiltered(ctx context.Context, req model.ListProjectsRequest) ([]repo.ListProjectsFilteredRow, error)
	StartProject(ctx context.Context, id int32) error
	FinishProject(ctx context.Context, id int32) error
	GetProjectByID(ctx context.Context, id int32) (repo.GetProjectByIDRow, error)
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
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("绑定参数: %w", err))
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Errorf("校验参数: %w", err))
	}

	// 业务函数
	err := p.projectService.CreateProject(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("创建project: %w", err))
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
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("绑定参数: %w", err))
	}

	// 业务函数
	err = p.projectService.UpdateProject(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("更新project: %w", err))
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("删除project: %w", err))
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
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("绑定参数: %w", err))
	}

	projects, err := p.projectService.ListProjectsFiltered(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取project: %w", err))
	}
	return c.JSON(http.StatusOK, projects)

}

func (p *ProjectHandler) StartProject(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	if err := p.projectService.StartProject(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("开始project: %w", err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功开始项目",
	})
}

func (p *ProjectHandler) FinishProject(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	if err := p.projectService.FinishProject(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("完成project: %w", err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功完成项目",
	})
}

func (p *ProjectHandler) GetProjectByID(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	
	project, err := p.projectService.GetProjectByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取project: %w", err))
	}

	return c.JSON(http.StatusOK, project)
}
