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

type TeamService interface {
	CreateTeam(ctx context.Context, req model.CreateTeamRequest) error
	UpdateTeam(ctx context.Context, id int32, req model.UpdateTeamRequest) error
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
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("绑定参数: %w", err))
	}

	// 验证数据
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Errorf("校验参数: %w", err))
	}

	// 调用业务函数
	err := t.teamService.CreateTeam(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("创建team: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "成功创建小组",
	})
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
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("绑定参数: %w", err))
	}

	// 调用业务函数
	err = t.teamService.UpdateTeam(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("更新team: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功更新小组",
	})
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("删除team: %w", err))
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取team: %w", err))
	}

	// 返回响应
	return c.JSON(http.StatusOK, teams)
}
