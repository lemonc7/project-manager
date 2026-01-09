package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/util"
	"github.com/lemonc7/zest"
)

type TodoService interface {
	ListTodosByReporter(ctx context.Context, reporterID int32, req model.ListTodosRequest) ([]repo.ListTodosByReporterRow, error)
	ListTodosByAssignee(ctx context.Context, assigneeID int32, req model.ListTodosRequest) ([]repo.ListTodosByAssigneeRow, error)
	CreateTodo(ctx context.Context, req model.CreateTodoRequest) error
	UpdateTodo(ctx context.Context, id int32, req model.UpdateTodoRequest) error
	DeleteTodo(ctx context.Context, id int32) error
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

func (t *TodoHandler) ListTodosByReporter(c *zest.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	var req model.ListTodosRequest
	if err := c.Bind(model.ListTodoRequestSchema, &req); err != nil {
		return err
	}

	todos, err := t.todoService.ListTodosByReporter(c.Context(), id, req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取todos: %w", err))
	}

	return c.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) ListTodosByAssignee(c *zest.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	var req model.ListTodosRequest
	if err := c.Bind(model.ListTodoRequestSchema, &req); err != nil {
		return zest.NewHTTPError(http.StatusBadRequest, fmt.Errorf("绑定参数: %w", err))
	}

	todos, err := t.todoService.ListTodosByAssignee(c.Context(), id, req)
	if err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("获取todos: %w", err))
	}

	return c.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) CreateTodo(c *zest.Context) error {
	var req model.CreateTodoRequest
	if err := c.Bind(model.CreateTodoRequestSchema, &req); err != nil {
		return err
	}

	if err := t.todoService.CreateTodo(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("创建todo: %w", err))
	}

	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功创建todo",
	})
}

func (t *TodoHandler) UpdateTodo(c *zest.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	var req model.UpdateTodoRequest
	if err := c.Bind(model.UpdateTodoRequestSchema, &req); err != nil {
		return err
	}

	if err := t.todoService.UpdateTodo(c.Context(), id, req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("更新todo: %w", err))
	}

	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功更新todo",
	})
}

func (t *TodoHandler) DeleteTodo(c *zest.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}

	if err := t.todoService.DeleteTodo(c.Context(), id); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("删除todo: %w", err))
	}

	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功删除todo",
	})
}
