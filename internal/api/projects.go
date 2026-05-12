package api

import (
	"context"
	"net/http"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/zest"
)

type ProjectService interface {
	CreateProject(ctx context.Context, req model.CreateProjectRequest) error
	UpdateProject(ctx context.Context, req model.UpdateProjectRequest) error
	DeleteProject(ctx context.Context, req model.DeleteProjectRequest) error
	ListProjects(ctx context.Context, req model.ListProjectsRequest) ([]repo.ListProjectsRow, error)
	StartProject(ctx context.Context, req model.StartProjectRequest) error
	GetProjectByID(ctx context.Context, req model.GetProjectRequest) (repo.GetProjectByIDRow, error)
	ResetStatus(ctx context.Context, req model.ResetStatusRequest) error
	CopyProject(ctx context.Context, req model.CopyProjectRequest) error
	IsProjectExists(ctx context.Context, req model.IsProjectExistsRequest) (bool, error)
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

func (p *ProjectHandler) CreateProject(c *zest.Context) error {
	// 绑定请求
	var req model.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// 业务函数
	err := p.projectService.CreateProject(c.Context(), req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "创建项目失败").Wrap(err)
	}

	// 返回响应
	return c.JSON(http.StatusCreated, zest.Map{
		"message": "成功创建项目",
	})
}

func (p *ProjectHandler) UpdateProject(c *zest.Context) error {
	// 绑定请求
	var req model.UpdateProjectRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// 业务函数
	if err := p.projectService.UpdateProject(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "更新项目失败").Wrap(err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功更新项目",
	})
}

func (p *ProjectHandler) DeleteProject(c *zest.Context) error {
	var req model.DeleteProjectRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	// 删除数据
	if err := p.projectService.DeleteProject(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "删除项目失败").Wrap(err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功删除项目",
	})
}

func (p *ProjectHandler) ListProjects(c *zest.Context) error {
	// 绑定请求
	var req model.ListProjectsRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	projects, err := p.projectService.ListProjects(c.Context(), req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "获取项目失败").Wrap(err)
	}
	return c.JSON(http.StatusOK, projects)
}

func (p *ProjectHandler) StartProject(c *zest.Context) error {
	var req model.StartProjectRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := p.projectService.StartProject(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "开始项目失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功开始项目",
	})
}

func (p *ProjectHandler) GetProjectByID(c *zest.Context) error {
	var req model.GetProjectRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	project, err := p.projectService.GetProjectByID(c.Context(), req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "获取项目失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, project)
}

func (p *ProjectHandler) ResetStatus(c *zest.Context) error {
	var req model.ResetStatusRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := p.projectService.ResetStatus(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "重置状态失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功重置项目状态",
	})
}

func (p *ProjectHandler) CopyProject(c *zest.Context) error {
	var req model.CopyProjectRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := p.projectService.CopyProject(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "复制项目失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功复制项目与相关流程",
	})
}

func (p *ProjectHandler) IsProjectExists(c *zest.Context) error {
	var req model.IsProjectExistsRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	exists, err := p.projectService.IsProjectExists(c.Context(), req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "判断项目存在失败").Wrap(err)
	}
	return c.JSON(http.StatusOK, exists)
}
