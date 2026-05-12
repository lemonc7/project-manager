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
		StatusID:    int32(model.ProcessStatusWaiting),
	})
}

func (p *ProcessService) InsertProcessAt(ctx context.Context, req model.InsertProcessRequest) error {
	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	// 保证出错时回滚
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := repo.New(p.db).WithTx(tx)

	// 验证sequence
	maxSeq, err := qtx.GetMaxOrder(ctx, req.ProjectID)
	if err != nil {
		return fmt.Errorf("获取最大序号: %w", err)
	}

	order := req.StepOrder
	if order < 1 || order > maxSeq+1 {
		return fmt.Errorf("序号无效: %d", order)
	}

	// 先将后续的process sequence加1
	err = qtx.ShiftProcessesUp(ctx, repo.ShiftProcessesUpParams{
		ProjectID: req.ProjectID,
		StepOrder: order,
	})
	if err != nil {
		return fmt.Errorf("更改后续流程的序号: %w", err)
	}

	// 插入新的process
	err = qtx.InsertProcessAt(ctx, repo.InsertProcessAtParams{
		ProjectID:   req.ProjectID,
		StepOrder:   order,
		Name:        req.Name,
		EmployeeID:  req.EmployeeID,
		Description: req.Description,
		Remark:      req.Remark,
		StatusID:    int32(model.ProcessStatusWaiting),
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
		StatusID:    req.StatusID,
		EmployeeID:  req.EmployeeID,
		SpendHours:  req.SpendHours,
		ID:          req.ID,
	})

	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("未找到被更新的流程")
	}

	return nil
}

func (p *ProcessService) DeleteProcess(ctx context.Context, req model.DeleteProcessRequest) error {
	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := repo.New(p.db).WithTx(tx)

	// 先删除process
	projectID, err := qtx.DeleteProcess(ctx, req.ID)
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
		ProjectID:  req.ProjectID,
		StepOrders: req.StepOrders,
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

func (p *ProcessService) ListProcesses(ctx context.Context, req model.ListProcessesRequest) ([]repo.ListProcessesRow, error) {
	return p.querier.ListProcesses(ctx, req.ID)
}

func (p *ProcessService) ListTasks(ctx context.Context, req model.ListTasksRequest) ([]repo.ListTasksRow, error) {
	return p.querier.ListTasks(ctx, repo.ListTasksParams{
		EmployeeID: req.EmployeeID,
		Category:   req.Status,
	})
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

	rows, err := qtx.FinishProcess(ctx, repo.FinishProcessParams{
		ID:         req.ID,
		SpendHours: &req.SpendHours,
		Remark:     req.Remark,
	})
	if err != nil {
		return fmt.Errorf("更新流程信息: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("流程状态未更新: 流程的初始状态必须为'进行'")
	}

	// 获取后两步process
	processes, err := qtx.GetNextTwoProcesses(ctx, repo.GetNextTwoProcessesParams{
		ProjectID: req.ProjectID,
		StepOrder: req.StepOrder,
	})

	if err != nil {
		return fmt.Errorf("获取后两步流程: %w", err)
	}

	// 设置下一步process的状态为"进行"
	if len(processes) > 0 {
		rows, err := qtx.UpdateProcessStatus(ctx, repo.UpdateProcessStatusParams{
			ID:            processes[0].ID,
			TargetStatus:  int32(model.ProcessStatusDoing),
			CurrentStatus: int32(model.ProcessStatusReady),
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
			TargetStatus:  int32(model.ProcessStatusReady),
			CurrentStatus: int32(model.ProcessStatusWaiting),
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
		StepOrder: process.StepOrder,
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
		StatusID:    int32(model.ProcessStatusRework),
		EmployeeID:  process.EmployeeID,
		StepOrder:   process.StepOrder,
		SpendHours:  &req.SpendHours,
		StartedAt:   process.StartedAt,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
