package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
)

type ProjectService struct {
	db      *pgxpool.Pool
	querier repo.Querier
}

func NewProjectService(db *pgxpool.Pool) *ProjectService {
	return &ProjectService{
		db:      db,
		querier: repo.New(db),
	}
}

func (p *ProjectService) CreateProject(ctx context.Context, req model.CreateProjectRequest) error {
	return p.querier.CreateProject(ctx, repo.CreateProjectParams{
		Name:         req.Name,
		Description:  req.Description,
		Remark:       req.Remark,
		StatusID:     int32(model.ProjectStatusPending),
		ManagerID:    req.ManagerID,
		Priority:     req.Priority,
		DurationDays: req.DurationDays,
	})
}

func (p *ProjectService) UpdateProject(ctx context.Context, req model.UpdateProjectRequest) error {
	rows, err := p.querier.UpdateProject(ctx, repo.UpdateProjectParams{
		ID:           req.ID,
		Name:         req.Name,
		Description:  req.Description,
		Remark:       req.Remark,
		StatusID:     req.StatusID,
		ManagerID:    req.ManagerID,
		Priority:     req.Priority,
		StartedAt:    req.StartedAt,
		DurationDays: req.DurationDays,
		FinishedAt:   req.FinishedAt,
	})
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("没有找到被更新的项目")
	}

	return nil
}

func (p *ProjectService) DeleteProject(ctx context.Context, req model.DeleteProjectRequest) error {
	return p.querier.DeleteProject(ctx, req.ID)
}

func (p *ProjectService) ListProjects(ctx context.Context, req model.ListProjectsRequest) ([]repo.ListProjectsRow, error) {
	return p.querier.ListProjects(ctx, repo.ListProjectsParams{
		StatusIds: req.StatusIDs,
		ManagerID: req.ManagerID,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
	})
}

func (p *ProjectService) StartProject(ctx context.Context, req model.StartProjectRequest) error {
	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := repo.New(p.db).WithTx(tx)

	// 获取项目前两个process
	processes, err := qtx.GetNextTwoProcesses(ctx, repo.GetNextTwoProcessesParams{
		ProjectID: req.ID,
		StepOrder: 0,
	})
	if err != nil {
		return fmt.Errorf("获取相关流程: %w", err)
	}

	// 项目没有process,开始project没有意义
	if len(processes) == 0 {
		return fmt.Errorf("项目没有相关流程")
	}

	// 先将项目设置为"进行"状态
	rows, err := qtx.UpdateProjectStatus(ctx, repo.UpdateProjectStatusParams{
		ID:            req.ID,
		TargetStatus:  int32(model.TodoStatusDoing),
		CurrentStatus: int32(model.TodoStatusPending),
	})
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("项目未被更新: 项目初始状态应为'排单'")
	}

	// 更新第一步process状态
	rows, err = qtx.UpdateProcessStatus(ctx, repo.UpdateProcessStatusParams{
		ID:            processes[0].ID,
		TargetStatus:  int32(model.ProcessStatusDoing),
		CurrentStatus: int32(model.ProcessStatusWaiting),
	})

	if err != nil {
		return fmt.Errorf("更新第一步流程状态: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("第一步流程状态未被更新: 第一步流程初始状态应为'等待'")
	}

	// 如果有第二步
	if len(processes) > 1 {
		rows, err := qtx.UpdateProcessStatus(ctx, repo.UpdateProcessStatusParams{
			ID:            processes[1].ID,
			TargetStatus:  int32(model.ProcessStatusReady),
			CurrentStatus: int32(model.ProcessStatusWaiting),
		})

		if err != nil {
			return fmt.Errorf("更新第二步流程状态: %w", err)
		}

		if rows == 0 {
			return fmt.Errorf("第二步流程状态未被更新: 第二步流程初始状态应为'等待'")
		}
	}

	return tx.Commit(ctx)
}

func (p *ProjectService) GetProjectByID(ctx context.Context, req model.GetProjectRequest) (repo.GetProjectByIDRow, error) {
	return p.querier.GetProjectByID(ctx, req.ID)
}

func (p *ProjectService) ResetStatus(ctx context.Context, req model.ResetStatusRequest) error {
	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := repo.New(p.db).WithTx(tx)

	// 重置project的状态
	rows, err := qtx.ResetProject(ctx, req.ID)

	if err != nil {
		return fmt.Errorf("重置项目状态: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("未找到对应项目")
	}

	// 重置process的状态
	rows, err = qtx.ResetProcessByProjectID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("重置流程状态: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("该项目中未找到对应流程或无流程")
	}

	return tx.Commit(ctx)
}

func (p *ProjectService) CopyProject(ctx context.Context, req model.CopyProjectRequest) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := repo.New(p.db).WithTx(tx)

	// 复制项目并重置
	newProjectID, err := qtx.CopyProject(ctx, repo.CopyProjectParams{
		ID:        req.ID,
		Name:      req.NewName,
		ManagerID: req.ManagerID,
	})

	if err != nil {
		return fmt.Errorf("复制项目: %w", err)
	}

	// 复制process并重置
	rows, err := qtx.CopyProcessesByProjectID(ctx, repo.CopyProcessesByProjectIDParams{
		TargetProject:   newProjectID,
		OriginalProject: req.ID,
	})

	if err != nil {
		return fmt.Errorf("复制流程: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("未找到对应流程或无流程")
	}

	return tx.Commit(ctx)
}

func (p *ProjectService) IsProjectExists(ctx context.Context, req model.IsProjectExistsRequest) (bool, error) {
	exists, err := p.querier.IsProjectExists(ctx, req.NewName)
	if err != nil {
		return false, err
	}
	return exists, nil
}
