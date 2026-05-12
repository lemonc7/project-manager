package middleware

import (
	"net/http"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/pkg/util"
	"github.com/lemonc7/zest"
)

func RequireRoleAtLeast(role model.Role) zest.MiddlewareFunc {
	return func(next zest.HandlerFunc) zest.HandlerFunc {
		return func(c *zest.Context) error {
			roles := util.GetRoles(c)

			if !util.HasAtLeastRole(roles, role) {
				return zest.NewHTTPError(http.StatusForbidden, "没有权限").Wrap(model.ErrNoPermission)
			}

			return next(c)
		}
	}
}
