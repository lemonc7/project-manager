package application

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lemonc7/project-manager/config"
	mw "github.com/lemonc7/project-manager/internal/middleware"
)

func NewEcho(cfg config.ServerConfig, validator echo.Validator) *echo.Echo {
	e := echo.New()

	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.ReadTimeout = cfg.ReadTimeout

	// 全局错误处理
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// 已响应，不做处理
		if c.Response().Committed {
			return
		}

		code, errMsg := parseError(err)
		c.JSON(code, echo.Map{
			"error": errMsg,
		})
	}

	e.Validator = validator

	e.Use(mw.Logger)
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return e
}

func parseError(err error) (code int, msg string) {
	// 已经响应，直接返回code 200
	if err == nil {
		return http.StatusOK, ""
	}

	// 处理*echo.HTTPError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		switch msgType := he.Message.(type) {
		case error:
			msg = msgType.Error()
		case string:
			msg = msgType
		default:
			msg = fmt.Sprintf("%v", msgType)
		}
		return code, msg
	}

	// 处理标准错误error
	return http.StatusInternalServerError, err.Error()
}
