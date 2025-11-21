package service

import (
	"context"
	"fmt"

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

func (p *ProcessService) InsertProcessLast(ctx context.Context, req model.InsertProcessRequest) (repo.Process, error) {
	process, err := p.queries.InsertProcessLast(ctx, repo.InsertProcessLastParams{
		Name:        req.Name,
		ProjectID:   req.ProjectID,
		EmployeeID:  req.EmployeeID,
		Description: req.Description,
		Remark:      req.Remark,
	})

	if err != nil {
		return repo.Process{}, err
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

func (p *ProcessService) UpdateProcess(ctx context.Context, id int32, req model.UpdateProcessRequest) (repo.Process, error) {
	process, err := p.queries.UpdateProcess(ctx, repo.UpdateProcessParams{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		EmployeeID:  req.EmployeeID,
		Remark:      req.Remark,
		ID:          id,
	})

	if err != nil {
		return repo.Process{}, err
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

func (p *ProcessService) ListProcessesFiltered(ctx context.Context, req model.ListProcessesRequest) ([]repo.ListProcessesFilteredRow, error) {
	processes, err := p.queries.ListProcessesFiltered(ctx, repo.ListProcessesFilteredParams{
		ProjectID:  req.ProjectID,
		Status:     req.Status,
		EmployeeID: req.EmployeeID,
	})

	if err != nil {
		return nil, err
	}

	return processes, nil
}
