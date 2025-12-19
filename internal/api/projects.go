package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

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
	GetProjectByID(ctx context.Context, id int32) (repo.GetProjectByIDRow, error)
	ResetStatus(ctx context.Context, id int32) error
	CopyProject(ctx context.Context, id int32, req model.CopyProjectRequest) error
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

func (p *ProjectHandler) CreateProject(c echo.Context) error {
	// 绑定请求
	var req model.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("绑定参数: %w", err))
	}

	// 验证数据
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Errorf("校验参数: %w", err))
	}

	// 业务函数
	err := p.projectService.CreateProject(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("创建项目: %w", err))
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("更新项目: %w", err))
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("删除项目: %w", err))
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取项目: %w", err))
	}
	return c.JSON(http.StatusOK, projects)

}

func (p *ProjectHandler) StartProject(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	if err := p.projectService.StartProject(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("开始项目: %w", err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功开始项目",
	})
}

func (p *ProjectHandler) GetProjectByID(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	project, err := p.projectService.GetProjectByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取项目: %w", err))
	}

	return c.JSON(http.StatusOK, project)
}

func (p *ProjectHandler) ResetStatus(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	if err := p.projectService.ResetStatus(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("重置项目和流程状态: %w", err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功重置项目状态",
	})
}

func (p *ProjectHandler) CopyProject(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	var req model.CopyProjectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("绑定参数: %w", err))
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Errorf("校验参数: %w", err))
	}
	if err := p.projectService.CopyProject(c.Request().Context(), id, req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("复制项目与相关流程: %w", err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功复制项目与相关流程",
	})
}

func (p *ProjectHandler) IsProjectExists(c echo.Context) error {
	var req model.IsProjectExistsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("绑定参数: %w", err))
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Errorf("校验参数: %w", err))
	}

	exists, err := p.projectService.IsProjectExists(c.Request().Context(), req.NewName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("判断项目是否存在: %w", err))
	}
	return c.JSON(http.StatusOK, exists)
}

func (p *ProjectHandler) ExportProject(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("project_%d_%d", id, time.Now().Unix())
	// 对文件名进行编码，防止中文或特殊字符导致浏览器解析错误
	encodedFilename := url.QueryEscape(filename)

	// 设置响应头
	c.Response().Header().Set(echo.HeaderContentType, "text/csv; charset=utf-8")
	// 告诉浏览器这是一个附件，需要下载
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", encodedFilename))

	err = p.projectService.ExportProject(c.Request().Context(), id, c.Response())
	if err != nil {
		c.Logger().Errorf("导出CSV失败: %w", err)
		return nil
	}

	return nil
}
