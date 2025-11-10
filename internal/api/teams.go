package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
)

type TeamService interface {
	CreateTeam(ctx context.Context, req model.CreateTeamRequest) error
	UpdateTeam(ctx context.Context, req model.UpdateTeamRequest) error
	DeleteTeam(ctx context.Context, req model.DeleteTeamRequest) error
	ListTeams(ctx context.Context) ([]repo.Team, error)
}

var _ TeamService = (*service.TeamService)(nil)

type TeamHandler struct {
	teamService TeamService
}

// 创建小组
func (t *TeamHandler) CreateTeam(c echo.Context) error {
	// 绑定请求
	var req model.CreateTeamRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 调用业务函数
	if err := t.teamService.CreateTeam(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "成功创建小组",
	})
}

// 更新小组
func (t *TeamHandler) UpdateTeam(c echo.Context) error {
	// 绑定请求
	var req model.UpdateTeamRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 调用业务函数
	if err := t.teamService.UpdateTeam(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功更新小组",
	})
}

// 删除小组
func (t *TeamHandler) DeleteTeam(c echo.Context) error {
	// 绑定请求
	var req model.DeleteTeamRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 验证数据
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 删除数据
	if err := t.teamService.DeleteTeam(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"teams": teams,
	})
}
