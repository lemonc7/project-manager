package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (p *ProcessService) InsertProcessLast(ctx context.Context, req model.InsertProcessRequest) error {
	return p.querier.InsertProcessLast(ctx, repo.InsertProcessLastParams{
		Name:        req.Name,
		ProjectID:   req.ProjectID,
		EmployeeID:  req.EmployeeID,
		Description: req.Description,
		Remark:      req.Remark,
	})
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

func (p *ProcessService) UpdateProcess(ctx context.Context, id int32, req model.UpdateProcessRequest) error {
	return p.querier.UpdateProcess(ctx, repo.UpdateProcessParams{
		Name:        req.Name,
		Description: req.Description,
		Remark:      req.Remark,
		Status:      req.Status,
		EmployeeID:  req.EmployeeID,
		SpendTime:   req.SpendTime,
		StartedAt:   req.StartedAt,
		FinishedAt:  req.FinishedAt,
		ID:          id,
	})
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

	return tx.Commit(ctx)
}

func (p *ProcessService) ListProcesses(ctx context.Context, projectID int32) ([]repo.ListProcessesRow, error) {
	processes, err := p.querier.ListProcesses(ctx, projectID)

	if err != nil {
		return nil, err
	}

	return processes, nil
}

func (p *ProcessService) ListTasksTodo(ctx context.Context, employeeID int32) ([]repo.ListTasksTodoRow, error) {
	tasks, err := p.querier.ListTasksTodo(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (p *ProcessService) ListTasksReady(ctx context.Context, employeeID int32) ([]repo.ListTasksReadyRow, error) {
	tasks, err := p.querier.ListTasksReady(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (p *ProcessService) ListTasksDoing(ctx context.Context, employeeID int32) ([]repo.ListTasksDoingRow, error) {
	tasks, err := p.querier.ListTasksDoing(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (p *ProcessService) FinishProcess(ctx context.Context, id int32, req model.FinishProcessRequest) error {
	_, err := p.querier.FinishProcess(ctx, repo.FinishProcessParams{
		ID:        id,
		SpendTime: &req.SpendTime,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("当前任务状态必须为'进行'")
		}
		return fmt.Errorf("数据库错误: %w", err)
	}

	return nil
}
