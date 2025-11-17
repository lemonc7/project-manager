package service

import (
	"context"
	"database/sql"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/pkg/util"
)

type ProjectService struct {
	querier repo.Querier
}

func NewProjectService(db *sql.DB) *ProjectService {
	return &ProjectService{
		querier: repo.New(db),
	}
}

func (p *ProjectService) CreateProject(ctx context.Context, req model.CreateProjectRequest) error {
	// 插入数据
	err := p.querier.CreateProject(ctx, repo.CreateProjectParams{
		Name:        req.Name,
		Description: req.Description,
		ManagerID:   req.ManagerID,
		Remark:      req.Remark,
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) UpdateProject(ctx context.Context, id int32, req model.UpdateProjectRequest) error {
	var status repo.NullProjectStatus
	if req.Status != nil {
		status.Valid = true
		status.ProjectStatus = repo.ProjectStatus(*req.Status)
	}
	// 更新数据
	err := p.querier.UpdateProject(ctx, repo.UpdateProjectParams{
		ID:          id,
		Name:        util.ToNullString(req.Name),
		Description: util.ToNullString(req.Description),
		Status:      status,
		ManagerID:   util.ToNullInt32(req.ManagerID),
		StartTime:   util.ToNullTime(req.StartTime),
		Duration:    util.ToNullInt32(req.Duration),
		Remark:      util.ToNullString(req.Remark),
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) DeleteProject(ctx context.Context, id int32) error {
	// 删除数据
	if err := p.querier.DeleteProject(ctx, id); err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) ListProjectsFiltered(ctx context.Context, req model.ListProjectsRequest) ([]repo.ListProjectsFilteredRow, error) {
	var status repo.NullProjectStatus
	if req.Status != nil {
		status.Valid = true
		status.ProjectStatus = repo.ProjectStatus(*req.Status)
	}
	// 获取数据
	projects, err := p.querier.ListProjectsFiltered(ctx, repo.ListProjectsFilteredParams{
		Status:    status,
		ManagerID: util.ToNullInt32(req.ManagerID),
		StartTime: util.ToNullTime(req.StartTime),
		EndTime:   util.ToNullTime(req.EndTime),
	})
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *ProjectService) ListProcesses(ctx context.Context, id int32) ([]repo.ListProcessesFilteredRow, error) {
	processes, err := p.querier.ListProcessesFiltered(ctx, repo.ListProcessesFilteredParams{
		ProjectID: util.ToNullInt32(&id),
	})
	if err != nil {
		return nil, err
	}
	return processes, nil
}
