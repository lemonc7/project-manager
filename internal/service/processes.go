package service

import (
	"context"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
)

type ProcessService struct {
	db      *pgxpool.Pool
	queries *repo.Queries
}

func NewProcessService(db *pgxpool.Pool) *ProcessService {
	return &ProcessService{
		db:      db,
		queries: repo.New(db),
	}
}

func (p *ProcessService) InsertProcessLast(ctx context.Context, req model.InsertProcessRequest) (repo.InsertProcessLastRow, error) {
	process, err := p.queries.InsertProcessLast(ctx, repo.InsertProcessLastParams{
		Name:        req.Name,
		ProjectID:   req.ProjectID,
		EmployeeID:  req.EmployeeID,
		Description: req.Description,
		Remark:      req.Remark,
	})

	if err != nil {
		return repo.InsertProcessLastRow{}, err
	}

	return process, nil
}

func (p *ProcessService) InsertProcessAt(ctx context.Context, sequence int32, req model.InsertProcessRequest) error {
	// 验证sequence
	maxSeq, err := p.queries.GetMaxSequence(ctx, req.ProjectID)
	if err != nil {
		return fmt.Errorf("fail to get max sequence: %w", err)
	}

	if sequence < 1 || sequence > maxSeq+1 {
		return fmt.Errorf("invalid sequence: %d", sequence)
	}

	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	// 保证出错时回滚
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := p.queries.WithTx(tx)

	// 先将后续的process sequence加1
	err = qtx.ShiftProcessesUp(ctx, repo.ShiftProcessesUpParams{
		ProjectID: req.ProjectID,
		Sequence:  sequence,
	})
	if err != nil {
		return err
	}

	// 插入新的process
	err = qtx.InsertProcessAt(ctx, repo.InsertProcessAtParams{
		ProjectID:   req.ProjectID,
		Sequence:    sequence,
		Name:        req.Name,
		EmployeeID:  req.EmployeeID,
		Description: req.Description,
		Remark:      req.Remark,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *ProcessService) UpdateProcess(ctx context.Context, id int32, req model.UpdateProcessRequest) (repo.UpdateProcessRow, error) {
	process, err := p.queries.UpdateProcess(ctx, repo.UpdateProcessParams{
		Name:        req.Name,
		Description: req.Description,
		Remark:      req.Remark,
		Status:      req.Status,
		EmployeeID:  req.EmployeeID,
		StartedAt:   req.StartedAt,
		FinishedAt:  req.FinishedAt,
		ID:          id,
	})

	if err != nil {
		return repo.UpdateProcessRow{}, err
	}

	return process, nil
}

func (p *ProcessService) DeleteProcess(ctx context.Context, id int32) error {
	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := p.queries.WithTx(tx)

	// 先删除process
	projectID, err := qtx.DeleteProcess(ctx, id)
	if err != nil {
		return err
	}

	// 重新排序
	err = qtx.ReorderProcesses(ctx, projectID)
	if err != nil {
		return err
	}

	// 提交事务
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (p *ProcessService) ListProcesses(ctx context.Context, projectID int32) ([]repo.ListProcessesRow, error) {
	processes, err := p.queries.ListProcesses(ctx, projectID)

	if err != nil {
		return nil, err
	}

	return processes, nil
}

func (p *ProcessService) ListTasksTodo(ctx context.Context, employeeID int32) ([]repo.ListTasksTodoRow, error) {
	tasks, err := p.queries.ListTasksTodo(ctx, &employeeID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (p *ProcessService) ListTasksReady(ctx context.Context, employeeID int32) ([]repo.ListTasksReadyRow, error) {
	tasks, err := p.queries.ListTasksReady(ctx, &employeeID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (p *ProcessService) ListTasksDoing(ctx context.Context, employeeID int32) ([]model.Task, error) {
	ts, err := p.queries.ListTasksDoing(ctx, &employeeID)
	if err != nil {
		return nil, err
	}

	tasks := make([]model.Task, len(ts))
	// 计算剩余天数
	today := time.Now().Truncate(24 * time.Hour)
	for i, t := range ts {
		tasks[i] = model.Task{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			StartedAt:   t.StartedAt,
			ProjectID:   t.ProjectID,
			ProjectName: t.ProjectName,
		}

		if t.Duration != nil && *t.Duration > 0 && t.StartedAt != nil {
			endDate := t.StartedAt.AddDate(0, 0, int(*t.Duration))
			remainingTime := endDate.Sub(today)

			days := int32(math.Ceil(remainingTime.Hours() / 24))
			tasks[i].RemainingDays = &days
		}
	}

	// 排序
	slices.SortFunc(tasks, func(a, b model.Task) int {
		switch {
		case a.RemainingDays != nil && b.RemainingDays == nil:
			return -1
		case a.RemainingDays == nil && b.RemainingDays != nil:
			return 1
		case a.RemainingDays != nil && b.RemainingDays != nil:
			if *a.RemainingDays < *b.RemainingDays {
				return -1
			} else if *a.RemainingDays > *b.RemainingDays {
				return 1
			}
		}
		return 0
	})
	return tasks, nil
}
