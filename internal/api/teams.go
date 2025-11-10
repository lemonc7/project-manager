package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/model"
)

type TeamService interface {
	CreateTeam(ctx context.Context, req model.CreateTeamRequest) error
}

type TeamHandler struct {
	teamService TeamService
}

// 创建小组
func (t *TeamHandler) CreateTeam(c echo.Context) error {
	// 提取数据
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
		"message": "成功创建小组-" + req.Name,
	})
}

// 更新小组

// 删除小组

// 获取小组
