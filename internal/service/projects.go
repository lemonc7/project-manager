package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

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
		Name:        req.Name,
		Description: req.Description,
		ManagerID:   req.ManagerID,
		Remark:      req.Remark,
		Priority:    req.Priority,
		Duration:    req.Duration,
	})
}

func (p *ProjectService) UpdateProject(ctx context.Context, id int32, req model.UpdateProjectRequest) error {
	rows, err := p.querier.UpdateProject(ctx, repo.UpdateProjectParams{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Remark:      req.Remark,
		Status:      req.Status,
		ManagerID:   req.ManagerID,
		Priority:    req.Priority,
		Duration:    req.Duration,
	})
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("没有找到被更新的项目")
	}

	return nil
}

func (p *ProjectService) DeleteProject(ctx context.Context, id int32) error {
	return p.querier.DeleteProject(ctx, id)
}

func (p *ProjectService) ListProjectsFiltered(ctx context.Context, req model.ListProjectsRequest) ([]repo.ListProjectsFilteredRow, error) {
	// 获取数据
	projects, err := p.querier.ListProjectsFiltered(ctx, repo.ListProjectsFilteredParams{
		Status:    req.Status,
		ManagerID: req.ManagerID,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
	})
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *ProjectService) StartProject(ctx context.Context, id int32) error {
	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := repo.New(p.db).WithTx(tx)

	// 获取项目前两个process
	processes, err := p.querier.GetNextTwoProcesses(ctx, repo.GetNextTwoProcessesParams{
		ProjectID: id,
		Sequence:  0,
	})
	if err != nil {
		return fmt.Errorf("获取process: %w", err)
	}

	// 项目没有process,开始project没有意义
	if len(processes) == 0 {
		return fmt.Errorf("项目没有process")
	}

	// 先将项目设置为"进行"状态
	rows, err := p.querier.UpdateProjectStatus(ctx, repo.UpdateProjectStatusParams{
		ID:            id,
		TargetStatus:  "进行",
		CurrentStatus: "排单",
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
		TargetStatus:  "进行",
		CurrentStatus: "等待",
	})

	if err != nil {
		return fmt.Errorf("更新第一步process: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("第一步process未被更新: 第一步process初始状态应为'等待'")
	}

	// 如果有第二步
	if len(processes) > 1 {
		rows, err := qtx.UpdateProcessStatus(ctx, repo.UpdateProcessStatusParams{
			ID:            processes[1].ID,
			TargetStatus:  "准备",
			CurrentStatus: "等待",
		})

		if err != nil {
			return fmt.Errorf("更新第二步process: %w", err)
		}

		if rows == 0 {
			return fmt.Errorf("第二步process未被更新: 第二步process初始状态应为'等待'")
		}
	}

	return tx.Commit(ctx)
}

func (p *ProjectService) GetProjectByID(ctx context.Context, id int32) (repo.GetProjectByIDRow, error) {
	return p.querier.GetProjectByID(ctx, id)
}

func (p *ProjectService) ResetStatus(ctx context.Context, id int32) error {
	// 开启事务
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := repo.New(p.db).WithTx(tx)

	// 重置project的状态
	rows, err := qtx.ResetProject(ctx, id)

	if err != nil {
		return fmt.Errorf("重置项目状态: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("未找到对应项目")
	}

	// 重置process的状态
	rows, err = qtx.ResetProcessByProjectID(ctx, id)
	if err != nil {
		return fmt.Errorf("重置流程状态: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("该项目中未找到对应流程或无流程")
	}

	return tx.Commit(ctx)
}

func (p *ProjectService) CopyProject(ctx context.Context, id int32, req model.CopyProjectRequest) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := repo.New(p.db).WithTx(tx)

	// 复制项目并重置
	newProjectID, err := qtx.CopyProject(ctx, repo.CopyProjectParams{
		ID:        id,
		Name:      req.NewName,
		ManagerID: req.ManagerID,
	})

	if err != nil {
		return fmt.Errorf("复制项目: %w", err)
	}

	// 复制process并重置
	rows, err := qtx.CopyProcessesByProjectID(ctx, repo.CopyProcessesByProjectIDParams{
		TargetProject:   newProjectID,
		OriginalProject: id,
	})

	if err != nil {
		return fmt.Errorf("复制流程: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("未找到对应流程或无流程")
	}

	return tx.Commit(ctx)
}

func (p *ProjectService) IsProjectExists(ctx context.Context, name string) (bool, error) {
	exists, err := p.querier.IsProjectExists(ctx, name)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p *ProjectService) ExportProject(ctx context.Context, id int32, w io.Writer) error {
	// excel中文支持
	if _, err := w.Write([]byte("\xEF\xBB\xBF")); err != nil {
		return fmt.Errorf("设置Excel支持中文: %w", err)
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	project, err := p.querier.GetProjectByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取项目: %w", err)
	}

	// 项目名
	if err := writer.Write([]string{project.Name, project.Description}); err != nil {
		return fmt.Errorf("写入项目名: %w", err)
	}

	processes, err := p.querier.ListProcesses(ctx, id)
	if err != nil {
		return fmt.Errorf("获取流程: %w", err)
	}

	// 表头
	if err := writer.Write([]string{
		"序号",
		"名称",
		"描述",
		"机时",
		"备注",
	}); err != nil {
		return fmt.Errorf("写入表头: %w", err)
	}

	for _, p := range processes {
		spendTime := ""
		if p.SpendTime != nil {
			spendTime = fmt.Sprintf("%.1f", *p.SpendTime)
		}
		if err := writer.Write([]string{
			strconv.Itoa(int(p.Sequence)),
			p.Name,
			p.Description,
			spendTime,
			p.Remark,
		}); err != nil {
			return fmt.Errorf("写入流程: %w", err)
		}
	}

	return writer.Error()
}
