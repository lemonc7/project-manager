package api

import (
	"context"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
)

type TodoService interface {
	ListTodosFiltered(ctx context.Context, req model.ListTodosFilteredRequest) ([]repo.ListTodosFilteredRow, error)
}

var _ TodoService = (*service.TodoService)(nil)

type TodoHandler struct {
	todoService TodoService
}

func NewTodoHandler(service TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: service,
	}
}
