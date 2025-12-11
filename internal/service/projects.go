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

type ProjectService struct {
	querier repo.Querier
}

func NewProjectService(db *pgxpool.Pool) *ProjectService {
	return &ProjectService{
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
	return p.querier.UpdateProject(ctx, repo.UpdateProjectParams{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Remark:      req.Remark,
		Status:      req.Status,
		ManagerID:   req.ManagerID,
		Priority:    req.Priority,
		StartedAt:   req.StartedAt,
		Duration:    req.Duration,
		FinishedAt:  req.FinishedAt,
	})
}

func (p *ProjectService) DeleteProject(ctx context.Context, id int32) error {
	return p.querier.DeleteProject(ctx, id)
}

func (p *ProjectService) ListProjectsFiltered(ctx context.Context, req model.ListProjectsRequest) ([]repo.ListProjectsFilteredRow, error) {
	// 获取数据
	projects, err := p.querier.ListProjectsFiltered(ctx, repo.ListProjectsFilteredParams{
		StartAt: req.StartAt,
		EndAt:   req.EndAt,
	})
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *ProjectService) StartProject(ctx context.Context, id int32) error {
	_, err := p.querier.StartProject(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("项目状态必须为'排单'")
		}
		return err
	}

	return nil
}

func (p *ProjectService) FinishProject(ctx context.Context, id int32) error {
	return p.querier.FinishProject(ctx, id)
}

func (p *ProjectService) GetProjectByID(ctx context.Context, id int32) (repo.GetProjectByIDRow, error) {
	return p.querier.GetProjectByID(ctx, id)
}
