package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/pkg/util"
)

type ProcessService struct {
	db      *sql.DB
	queries *repo.Queries
}

func NewProcessService(db *sql.DB) *ProcessService {
	return &ProcessService{
		db:      db,
		queries: repo.New(db),
	}
}

func (p *ProcessService) InsertProcessLast(ctx context.Context, req model.InsertProcessRequest) error {
	err := p.queries.InsertProcessLast(ctx, repo.InsertProcessLastParams{
		Name:        req.Name,
		ProjectID:   req.ProjectID,
		EmployeeID:  req.EmployeeID,
		Description: req.Description,
		Remark:      req.Remark,
	})

	if err != nil {
		return err
	}

	return nil
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
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	// 保证出错时回滚
	defer tx.Rollback()

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

	return tx.Commit()
}

func (p *ProcessService) UpdateProcess(ctx context.Context, id int32, req model.UpdateProcessRequest) error {
	var status repo.NullProcessStatus
	if req.Status != nil {
		status.Valid = true
		status.ProcessStatus = repo.ProcessStatus(*req.Status)
	}
	err := p.queries.UpdateProcess(ctx, repo.UpdateProcessParams{
		Name:        util.ToNullString(req.Name),
		Description: util.ToNullString(req.Description),
		Status:      status,
		EmployeeID:  util.ToNullInt32(req.EmployeeID),
		Remark:      util.ToNullString(req.Remark),
		ID:          id,
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *ProcessService) DeleteProcess(ctx context.Context, id int32) error {
	// 开启事务
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

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

	return tx.Commit()
}

func (p *ProcessService) ListProcesses(ctx context.Context, id int32) ([]repo.ListProcessesFilteredRow, error) {
	processes, err := p.queries.ListProcessesFiltered(ctx, repo.ListProcessesFilteredParams{
		ProjectID: util.ToNullInt32(&id),
	})
	if err != nil {
		return nil, err
	}
	return processes, nil
}

func (p *ProcessService) ListProcessesFiltered(ctx context.Context, req model.ListProcessesRequest) ([]repo.ListProcessesFilteredRow, error) {
	var status repo.NullProcessStatus
	if req.Status != nil {
		status.Valid = true
		status.ProcessStatus = repo.ProcessStatus(*req.Status)
	}

	processes, err := p.queries.ListProcessesFiltered(ctx, repo.ListProcessesFilteredParams{
		ProjectID:  util.ToNullInt32(req.ProjectID),
		Status:     status,
		EmployeeID: util.ToNullInt32(req.EmployeeID),
	})

	if err != nil {
		return nil, err
	}

	return processes, nil
}
