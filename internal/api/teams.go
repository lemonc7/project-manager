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

type TeamService interface {
	CreateTeam(ctx context.Context, req model.CreateTeamRequest) (repo.CreateTeamRow, error)
	UpdateTeam(ctx context.Context, id int32, req model.UpdateTeamRequest) (repo.UpdateTeamRow, error)
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
func (t *TeamHandler) CreateTeam(c echo.Context) error {
	// 绑定请求
	var req model.CreateTeamRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	// 调用业务函数
	team, err := t.teamService.CreateTeam(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusCreated, team)
}

// 更新小组
func (t *TeamHandler) UpdateTeam(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	// 绑定请求
	var req model.UpdateTeamRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// 调用业务函数
	team, err := t.teamService.UpdateTeam(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, team)
}

// 删除小组
func (t *TeamHandler) DeleteTeam(c echo.Context) error {
	// 获取id
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	// 删除数据
	if err := t.teamService.DeleteTeam(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功删除小组",
	})
}

// 获取小组
func (t *TeamHandler) ListTeams(c echo.Context) error {
	// 获取数据
	teams, err := t.teamService.ListTeams(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 返回响应
	return c.JSON(http.StatusOK, teams)
}
