package api

import (
	"context"
	"net/http"

	"github.com/lemonc7/zest"
)

type ServerService interface {
	Health(ctx context.Context) error
}

type ServerHandler struct {
	serverService ServerService
}

func NewServerService(service ServerService) *ServerHandler {
	return &ServerHandler{
		serverService: service,
	}
}

func (s *ServerHandler) Health(c *zest.Context) error {
	if err := s.serverService.Health(c.Context()); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError).Wrap(err)
	}
	return c.String(http.StatusOK, "Welcome to the Project-Manager!")
}
