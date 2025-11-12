package application

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lemonc7/project-manager/config"
	mw "github.com/lemonc7/project-manager/internal/middleware"
)

func NewEcho(cfg config.ServerConfig, validator echo.Validator) *echo.Echo {
	e := echo.New()

	// e.HideBanner = true
	// e.HidePort = true
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.ReadTimeout = cfg.ReadTimeout

	e.Validator = validator

	e.Use(mw.Logger)
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return e
}
