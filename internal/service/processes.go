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
	querier repo.Querier
}

func NewProcessService(db *pgxpool.Pool) *ProcessService {
	return &ProcessService{
		db:      db,
		querier: repo.New(db),
	}
}

func (p *ProcessService) InsertProcessLast(ctx context.Context, req model.InsertProcessRequest) (model.Process, error) {
	process, err := p.querier.InsertProcessLast(ctx, repo.InsertProcessLastParams{
		Name:        req.Name,
		ProjectID:   req.ProjectID,
		EmployeeID:  req.EmployeeID,
		Description: req.Description,
		Remark:      req.Remark,
	})

	if err != nil {
		return model.Process{}, err
	}

	return model.Process{
		ID:          process.ID,
		ProjectID:   process.ProjectID,
		Name:        process.Name,
		Description: process.Description,
		Remark:      process.Remark,
		Status:      process.Status,
		EmployeeID:  process.EmployeeID,
		Sequence:    process.Sequence,
		StartedAt:   process.StartedAt,
		FinishedAt:  process.FinishedAt,
	}, nil
}

func (p *ProcessService) InsertProcessAt(ctx context.Context, sequence int32, req model.InsertProcessRequest) error {
	// 验证sequence
	maxSeq, err := p.querier.GetMaxSequence(ctx, req.ProjectID)
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
	qtx := repo.New(p.db).WithTx(tx)

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

func (p *ProcessService) UpdateProcess(ctx context.Context, id int32, req model.UpdateProcessRequest) (model.Process, error) {
	process, err := p.querier.UpdateProcess(ctx, repo.UpdateProcessParams{
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
		return model.Process{}, err
	}

	return model.Process{
		ID:          process.ID,
		ProjectID:   process.ProjectID,
		Name:        process.Name,
		Description: process.Description,
		Remark:      process.Remark,
		Status:      process.Status,
		EmployeeID:  process.EmployeeID,
		Sequence:    process.Sequence,
		StartedAt:   process.StartedAt,
		FinishedAt:  process.FinishedAt,
	}, nil
}

func (p *ProcessService) DeleteProcess(ctx context.Context, id int32) error {
	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := repo.New(p.db).WithTx(tx)

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

func (p *ProcessService) ListProcesses(ctx context.Context, projectID int32) ([]model.Process, error) {
	data, err := p.querier.ListProcesses(ctx, projectID)

	if err != nil {
		return nil, err
	}
	processes := make([]model.Process, len(data))

	for i, p := range data {
		processes[i] = model.Process{
			ID:          p.ID,
			ProjectID:   p.ProjectID,
			Name:        p.Name,
			Description: p.Description,
			Remark:      p.Remark,
			Status:      p.Status,
			EmployeeID:  p.EmployeeID,
			Sequence:    p.Sequence,
			StartedAt:   p.StartedAt,
			FinishedAt:  p.FinishedAt,
		}
	}
	return processes, nil

}

func (p *ProcessService) ListTasksTodo(ctx context.Context, employeeID int32) ([]model.Task, error) {
	data, err := p.querier.ListTasksTodo(ctx, &employeeID)
	if err != nil {
		return nil, err
	}
	tasks := make([]model.Task, len(data))

	for i, t := range data {
		tasks[i] = model.Task{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			ProjectID:   t.ProjectID,
			ProjectName: t.ProjectName,
		}
	}

	return tasks, nil
}

func (p *ProcessService) ListTasksReady(ctx context.Context, employeeID int32) ([]model.Task, error) {
	data, err := p.querier.ListTasksReady(ctx, &employeeID)
	if err != nil {
		return nil, err
	}

	tasks := make([]model.Task, len(data))
	for i, t := range data {
		tasks[i] = model.Task{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			ProjectID:   t.ProjectID,
			ProjectName: t.ProjectName,
		}
	}

	return tasks, nil
}

func (p *ProcessService) ListTasksDoing(ctx context.Context, employeeID int32) ([]model.Task, error) {
	ts, err := p.querier.ListTasksDoing(ctx, &employeeID)
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
