package service

import (
	"context"
	"database/sql"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
)

type ProjectService struct {
	querier repo.Querier
}

func (p *ProjectService) CreateProject(ctx context.Context, req model.CreateProjectRequest) error {
	var managerID sql.NullInt32
	if req.ManagerID != nil {
		managerID.Valid = true
		managerID.Int32 = *req.ManagerID
	}
	// 插入数据
	_, err := p.querier.CreateProject(ctx, repo.CreateProjectParams{
		Name:        req.Name,
		Description: req.Description,
		ManagerID:   managerID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) UpdateProject(ctx context.Context, req model.UpdateProjectRequest) error {
	var managerID sql.NullInt32
	var startTime sql.NullTime
	if req.ManagerID != nil {
		managerID.Valid = true
		managerID.Int32 = *req.ManagerID
	}

	if req.StartTime != nil {
		startTime.Valid = true
		startTime.Time = *req.StartTime
	}

	// 更新数据
	_, err := p.querier.UpdateProject(ctx, repo.UpdateProjectParams{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Status:      repo.ProjectStatus(req.Status),
		ManagerID:   managerID,
		StartTime:   startTime,
		Duration:    req.Duration,
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) DeleteProject(ctx context.Context, req model.DeleteProjectRequest) error {
	// 删除数据
	if err := p.querier.DeleteProject(ctx, req.ID); err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) ListProjectsFiltered(ctx context.Context, req model.ListProjectsRequest) ([]repo.ListProjectsFilteredRow, error) {
	var status string
	var managerID sql.NullInt32
	var startTime, endTime sql.NullTime
	if req.Status != nil {
		status = string(*req.Status)
	}
	if req.ManagerID != nil {
		managerID.Valid = true
		managerID.Int32 = *req.ManagerID
	}
	if req.StartTime != nil {
		startTime.Valid = true
		startTime.Time = *req.StartTime
	}
	if req.EndTime != nil {
		endTime.Valid = true
		endTime.Time = *req.EndTime
	}

	// 获取数据
	projects, err := p.querier.ListProjectsFiltered(ctx, repo.ListProjectsFilteredParams{
		Column1:     status,
		ManagerID:   managerID,
		StartTime:   startTime,
		StartTime_2: endTime,
	})
	if err != nil {
		return nil, err
	}

	return projects, nil
}
