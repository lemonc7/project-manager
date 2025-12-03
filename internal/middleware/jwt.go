package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/pkg/jwt"
)

func JWTAuther(jwt *jwt.JWT) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			ls := strings.Split(authHeader, " ")

			if len(ls) != 2 || ls[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header")
			}

			tokenString := ls[1]
			claims, err := jwt.ParseToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("roles", claims.Roles)
			if !slices.Contains(claims.Roles, "admin") &&
				!slices.Contains(claims.Roles, "manager") &&
				!slices.Contains(claims.Roles, "employee") {
				return echo.NewHTTPError(http.StatusForbidden, model.ErrNoPermission)
			}
			return next(c)
		}
	}
}
