package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
)

type ProjectService struct {
	querier repo.Querier
}

func NewProjectService(db *pgxpool.Pool) *ProjectService {
	return &ProjectService{
		querier: repo.New(db),
	}
}

func (p *ProjectService) CreateProject(ctx context.Context, req model.CreateProjectRequest) (repo.Project, error) {
	// 插入数据
	project, err := p.querier.CreateProject(ctx, repo.CreateProjectParams{
		Name:        req.Name,
		Description: req.Description,
		ManagerID:   req.ManagerID,
		Remark:      req.Remark,
		Duration:    req.Duration,
	})

	if err != nil {
		return repo.Project{}, err
	}

	return project, nil
}

func (p *ProjectService) UpdateProject(ctx context.Context, id int32, req model.UpdateProjectRequest) (repo.Project, error) {
	// 更新数据
	project, err := p.querier.UpdateProject(ctx, repo.UpdateProjectParams{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Remark:      req.Remark,
		Status:      req.Status,
		ManagerID:   req.ManagerID,
		StartedAt:   req.StartedAt,
		Duration:    req.Duration,
		FinishedAt:  req.FinishedAt,
	})

	if err != nil {
		return repo.Project{}, err
	}

	return project, nil
}

func (p *ProjectService) DeleteProject(ctx context.Context, id int32) error {
	// 删除数据
	if err := p.querier.DeleteProject(ctx, id); err != nil {
		return err
	}

	return nil
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
