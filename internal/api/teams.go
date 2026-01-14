package api

import (
	"context"
	"net/http"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/util"
	"github.com/lemonc7/zest"
)

type TeamService interface {
	CreateTeam(ctx context.Context, req model.CreateTeamRequest) error
	UpdateTeam(ctx context.Context, req model.UpdateTeamRequest) error
	DeleteTeam(ctx context.Context, id int32) error
	ListTeams(ctx context.Context) ([]repo.ListTeamsRow, error)
}

var _ TeamService = (*service.TeamService)(nil)

type TeamHandler struct {
	teamService TeamService
}

func NewTeamHandler(service TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: service,
	}
}

// 创建小组
func (t *TeamHandler) CreateTeam(c *zest.Context) error {
	// 绑定请求
	var req model.CreateTeamRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// 调用业务函数
	err := t.teamService.CreateTeam(c.Context(), req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "创建小组失败").Wrap(err)
	}

	// 返回响应
	return c.JSON(http.StatusCreated, zest.Map{
		"message": "成功创建小组",
	})
}

// 更新小组
func (t *TeamHandler) UpdateTeam(c *zest.Context) error {
	// 绑定请求
	var req model.UpdateTeamRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// 调用业务函数
	if err := t.teamService.UpdateTeam(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "更新小组失败").Wrap(err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功更新小组",
	})
}

// 删除小组
func (t *TeamHandler) DeleteTeam(c *zest.Context) error {
	// 获取id
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	// 删除数据
	if err := t.teamService.DeleteTeam(c.Context(), id); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "删除小组失败").Wrap(err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功删除小组",
	})
}

// 获取小组
func (t *TeamHandler) ListTeams(c *zest.Context) error {
	// 获取数据
	teams, err := t.teamService.ListTeams(c.Context())
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "获取小组失败").Wrap(err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, teams)
}
