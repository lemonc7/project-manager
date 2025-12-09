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
		return fmt.Errorf("获取最大sequence: %w", err)
	}

	if sequence < 1 || sequence > maxSeq+1 {
		return fmt.Errorf("sequence参数无效: %d", sequence)
	}

	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
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
		return fmt.Errorf("更改后续process的sequence: %w", err)
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
	err := p.querier.UpdateProcess(ctx, repo.UpdateProcessParams{
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

	return err
}

func (p *ProcessService) DeleteProcess(ctx context.Context, id int32) error {
	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := repo.New(p.db).WithTx(tx)

	// 先删除process
	projectID, err := qtx.DeleteProcess(ctx, id)
	if err != nil {
		return fmt.Errorf("删除process: %w", err)
	}

	// 重新排序
	err = qtx.ReorderProcesses(ctx, projectID)
	if err != nil {
		return fmt.Errorf("重新排序process: %w", err)
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
			return fmt.Errorf("当前状态必须为'进行'")
		}
		return err
	}

	return nil
}

func (p *ProcessService) ReworkProcess(ctx context.Context, id int32, req model.ReworkProcessRequest) error {
	// 查询当前process
	process, err := p.querier.GetProcessByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取目前的process: %w", err)
	}

	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := repo.New(p.db).WithTx(tx)

	// 先将后续的sequence都加1
	err = qtx.ShiftProcessesUp(ctx, repo.ShiftProcessesUpParams{
		ProjectID: process.ProjectID,
		Sequence:  process.Sequence,
	})

	if err != nil {
		return fmt.Errorf("更改后续process的sequence: %w", err)
	}

	// 添加返工的process
	err = qtx.InsertReworkProcess(ctx, repo.InsertReworkProcessParams{
		ProjectID:   process.ProjectID,
		Name:        process.Name + "(返工)",
		Description: process.Description,
		Remark:      req.Remark,
		Status:      "返工",
		EmployeeID:  process.EmployeeID,
		Sequence:    process.Sequence,
		SpendTime:   &req.SpendTime,
		StartedAt:   process.StartedAt,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
