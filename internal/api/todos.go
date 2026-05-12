package api

import (
	"context"
	"net/http"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/zest"
)

type TodoService interface {
	ListTodosForCreator(ctx context.Context, req model.ListTodosRequest) ([]repo.ListTodosForCreatorRow, error)
	ListTodosForAssignee(ctx context.Context, req model.ListTodosRequest) ([]repo.ListTodosForAssigneeRow, error)
	CreateTodo(ctx context.Context, req model.CreateTodoRequest) error
	UpdateTodo(ctx context.Context, req model.UpdateTodoRequest) error
	DeleteTodo(ctx context.Context, req model.DeleteTodoRequest) error
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

func (t *TodoHandler) ListTodosForCreator(c *zest.Context) error {
	var req model.ListTodosRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	todos, err := t.todoService.ListTodosForCreator(c.Context(), req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "获取todo失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) ListTodosForAssignee(c *zest.Context) error {
	var req model.ListTodosRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	todos, err := t.todoService.ListTodosForAssignee(c.Context(), req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "获取分配todo失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) CreateTodo(c *zest.Context) error {
	var req model.CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := t.todoService.CreateTodo(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "创建todo失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功创建todo",
	})
}

func (t *TodoHandler) UpdateTodo(c *zest.Context) error {
	var req model.UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := t.todoService.UpdateTodo(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "更新todo失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功更新todo",
	})
}

func (t *TodoHandler) DeleteTodo(c *zest.Context) error {
	var req model.DeleteTodoRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := t.todoService.DeleteTodo(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "更新todo失败").Wrap(err)
	}

	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功删除todo",
	})
}
