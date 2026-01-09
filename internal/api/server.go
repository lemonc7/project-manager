package api

import (
	"context"
	"net/http"

	"github.com/lemonc7/zest"
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

func (s *ServerHandler) Ping(c *zest.Context) error {
	if err := s.serverService.Ping(c.Context()); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, zest.Map{
		"message": "pong",
	})
}
