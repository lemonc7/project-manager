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

func (t *TodoService) ListTodosForCreator(ctx context.Context, id int32) ([]repo.ListTodosForCreatorRow, error) {
	return t.querier.ListTodosForCreator(ctx, id)
}

func (t *TodoService) ListTodosByAssignee(ctx context.Context, id int32) ([]repo.ListTodosForAssigneeRow, error) {
	return t.querier.ListTodosForAssignee(ctx, id)
}

func (t *TodoService) CreateTodo(ctx context.Context, req model.CreateTodoRequest) error {
	return t.querier.CreateTodo(ctx, repo.CreateTodoParams{
		Name:        req.Name,
		Description: req.Description,
		Remark:      req.Remark,
		Priority:    req.Priority,
		StatusID:    int32(model.TodoStatusPending),
		CreatorID:   req.CreatorID,
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
		StatusID:    req.StatusID,
		CreatorID:   req.CreatorID,
		AssigneeID:  req.AssigneeID,
		StartedAt:   req.StartedAt,
		FinishedAt:  req.FinishedAt,
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
