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

func (p *ProcessService) InsertProcessAt(ctx context.Context, req model.InsertProcessRequest) error {
	// 验证sequence
	maxSeq, err := p.querier.GetMaxSequence(ctx, req.ProjectID)
	if err != nil {
		return fmt.Errorf("获取最大序号: %w", err)
	}

	sequence := req.Sequence
	if sequence < 1 || sequence > maxSeq+1 {
		return fmt.Errorf("序号无效: %d", sequence)
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
		return fmt.Errorf("更改后续流程的序号: %w", err)
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

func (p *ProcessService) UpdateProcess(ctx context.Context, req model.UpdateProcessRequest) error {
	rows, err := p.querier.UpdateProcess(ctx, repo.UpdateProcessParams{
		Name:        req.Name,
		Description: req.Description,
		Remark:      req.Remark,
		Status:      req.Status,
		EmployeeID:  req.EmployeeID,
		SpendTime:   req.SpendTime,
		ID:          req.ID,
	})

	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("没有找到被更新的流程")
	}

	return nil
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
		return fmt.Errorf("删除流程: %w", err)
	}

	// 重新排序
	if err := qtx.ReorderProcesses(ctx, projectID); err != nil {
		return fmt.Errorf("重新排序流程: %w", err)
	}

	return tx.Commit(ctx)
}

func (p *ProcessService) DeleteProcesses(ctx context.Context, req model.DeleteProcessesRequest) error {
	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := repo.New(p.db).WithTx(tx)

	// 先删除processes
	rows, err := qtx.DeleteProcesses(ctx, repo.DeleteProcessesParams{
		ProjectID: req.ProjectID,
		Column2:   req.Sequences,
	})

	if err != nil {
		return fmt.Errorf("删除流程: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("没有找到被删除的流程")
	}

	// 重新排序
	if err := qtx.ReorderProcesses(ctx, req.ProjectID); err != nil {
		return fmt.Errorf("重新排序流程: %w", err)
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

func (p *ProcessService) FinishProcess(ctx context.Context, req model.FinishProcessRequest) error {
	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := repo.New(p.db).WithTx(tx)

	// 现将process设置为"完成"状态
	rows, err := qtx.FinishProcess(ctx, repo.FinishProcessParams{
		ID:        req.ID,
		SpendTime: &req.SpendTime,
		Remark:    req.Remark,
	})
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("流程状态未更新: 流程的初始状态必须为'进行'")
	}

	// 获取后两步process
	processes, err := qtx.GetNextTwoProcesses(ctx, repo.GetNextTwoProcessesParams{
		ProjectID: req.ProjectID,
		Sequence:  req.Sequence,
	})

	if err != nil {
		return fmt.Errorf("获取后两步流程: %w", err)
	}

	// 设置下一步process的状态为"进行"
	if len(processes) > 0 {
		rows, err := qtx.UpdateProcessStatus(ctx, repo.UpdateProcessStatusParams{
			ID:            processes[0].ID,
			TargetStatus:  "进行",
			CurrentStatus: "准备",
		})

		if err != nil {
			return fmt.Errorf("更新下一步流程状态: %w", err)
		}

		if rows == 0 {
			return fmt.Errorf("流程状态未更新: 下一步流程的初始状态必须为'准备'")
		}
	}

	// 设置下下步流程状态为"准备"
	if len(processes) > 1 {
		rows, err := qtx.UpdateProcessStatus(ctx, repo.UpdateProcessStatusParams{
			ID:            processes[1].ID,
			TargetStatus:  "准备",
			CurrentStatus: "等待",
		})

		if err != nil {
			return fmt.Errorf("更新下下一步流程状态: %w", err)
		}

		if rows == 0 {
			return fmt.Errorf("流程状态未更新: 下下一步流程的初始状态必须为'等待'")
		}
	}

	return tx.Commit(ctx)
}

func (p *ProcessService) ReworkProcess(ctx context.Context, req model.ReworkProcessRequest) error {
	// 查询当前process
	process, err := p.querier.GetProcessByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("获取目前的流程: %w", err)
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
		return fmt.Errorf("更改后续流程的序号: %w", err)
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
