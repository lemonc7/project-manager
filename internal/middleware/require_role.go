package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/pkg/util"
)

func RequireRoleAtLeast(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roles := util.GetRoles(c)

			if !util.HasAtLeastRole(roles, role) {
				return echo.NewHTTPError(http.StatusForbidden, model.ErrNoPermission)
			}

			return next(c)
		}
	}
}
