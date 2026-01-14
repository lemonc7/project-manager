package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
)

type TodoService struct {
	querier repo.Querier
}

func NewTodoService(db *pgxpool.Pool) *TodoService {
	return &TodoService{
		querier: repo.New(db),
	}
}

func (t *TodoService) ListTodosByReporter(ctx context.Context, req model.ListTodosRequest) ([]repo.ListTodosByReporterRow, error) {
	todos, err := t.querier.ListTodosByReporter(ctx, repo.ListTodosByReporterParams{
		Status:     req.Status,
		ReporterID: req.ID,
	})

	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (t *TodoService) ListTodosByAssignee(ctx context.Context, req model.ListTodosRequest) ([]repo.ListTodosByAssigneeRow, error) {
	todos, err := t.querier.ListTodosByAssignee(ctx, repo.ListTodosByAssigneeParams{
		Status:     req.Status,
		AssigneeID: req.ID,
	})

	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (t *TodoService) CreateTodo(ctx context.Context, req model.CreateTodoRequest) error {
	return t.querier.CreateTodo(ctx, repo.CreateTodoParams{
		Name:        req.Name,
		Description: req.Description,
		Remark:      req.Remark,
		Priority:    req.Priority,
		ReporterID:  req.ReporterID,
		AssigneeID:  req.AssigneeID,
	})
}

func (t *TodoService) UpdateTodo(ctx context.Context, req model.UpdateTodoRequest) error {
	rows, err := t.querier.UpdateTodo(ctx, repo.UpdateTodoParams{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Remark:      req.Remark,
		Priority:    req.Priority,
		Status:      req.Status,
		ReporterID:  req.ReporterID,
		AssigneeID:  req.AssigneeID,
	})

	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("没有找到被更新的todo")
	}

	return nil
}

func (t *TodoService) DeleteTodo(ctx context.Context, id int32) error {
	rows, err := t.querier.DeleteTodo(ctx, id)
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("没有找到被删除的todo")
	}

	return nil
}
