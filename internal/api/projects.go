package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/util"
	"github.com/lemonc7/zest"
)

type ProjectService interface {
	CreateProject(ctx context.Context, req model.CreateProjectRequest) error
	UpdateProject(ctx context.Context, req model.UpdateProjectRequest) error
	DeleteProject(ctx context.Context, id int32) error
	ListProjectsFiltered(ctx context.Context, req model.ListProjectsRequest) ([]repo.ListProjectsFilteredRow, error)
	StartProject(ctx context.Context, id int32) error
	GetProjectByID(ctx context.Context, id int32) (repo.GetProjectByIDRow, error)
	ResetStatus(ctx context.Context, id int32) error
	CopyProject(ctx context.Context,req model.CopyProjectRequest) error
	IsProjectExists(ctx context.Context, name string) (bool, error)
	ExportProject(ctx context.Context, id int32, w io.Writer) error
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
	// 获取id
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	// 删除数据
	if err := p.projectService.DeleteProject(c.Context(), id); err != nil {
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

	projects, err := p.projectService.ListProjectsFiltered(c.Context(), req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "获取项目失败").Wrap(err)
	}
	return c.JSON(http.StatusOK, projects)

}

func (p *ProjectHandler) StartProject(c *zest.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	if err := p.projectService.StartProject(c.Context(), id); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "开始项目失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功开始项目",
	})
}

func (p *ProjectHandler) GetProjectByID(c *zest.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	project, err := p.projectService.GetProjectByID(c.Context(), id)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "获取项目失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, project)
}

func (p *ProjectHandler) ResetStatus(c *zest.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	if err := p.projectService.ResetStatus(c.Context(), id); err != nil {
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

	exists, err := p.projectService.IsProjectExists(c.Context(), req.NewName)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "判断项目存在失败").Wrap(err)
	}
	return c.JSON(http.StatusOK, exists)
}

func (p *ProjectHandler) ExportProject(c *zest.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("project_%d_%d", id, time.Now().Unix())

	// 设置响应头
	c.SetHeader(zest.HeaderContentType, "text/csv; charset=utf-8")
	// 告诉浏览器这是一个附件，需要下载
	c.SetHeader(zest.HeaderContentDisposition, `attachment; filename*=UTF-8''`+url.PathEscape(filename))
	err = p.projectService.ExportProject(c.Context(), id, c.Response())
	if err != nil {
		log.Printf("[%v] 导出CSV失败: %v", c.Get("requestID"), err)
	}

	return nil
}
