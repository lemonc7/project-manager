package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/pkg/jwt"
)

func JWTAuther(jwt *jwt.JWT) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("缺失token"))
			}

			ls := strings.Split(authHeader, " ")
			if len(ls) != 2 || ls[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("无效token"))
			}

			tokenString := ls[1]
			claims, err := jwt.ParseToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("解码token: %w", err))
			}

			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("roles", claims.Roles)
			return next(c)
		}
	}
}
