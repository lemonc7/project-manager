package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ServerService interface {
	Ping(ctx context.Context) error
}

type ServerHandler struct {
	serverService ServerService
}

func NewServerService(service ServerService) *ServerHandler {
	return &ServerHandler{
		serverService: service,
	}
}

func (s *ServerHandler) Ping(c echo.Context) error {
	ctx := c.Request().Context()

	if err := s.serverService.Ping(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "pong",
	})
}
