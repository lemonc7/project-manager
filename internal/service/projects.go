package service

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/pkg/util"
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

func (p *ProjectService) ListProjectsFiltered(ctx context.Context, req model.ListProjectsRequest) ([]model.Project, error) {
	// 获取数据
	data, err := p.querier.ListProjectsFiltered(ctx, repo.ListProjectsFilteredParams{
		Status:    req.Status,
		ManagerID: req.ManagerID,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
	})
	if err != nil {
		return nil, err
	}

	projects := make([]model.Project, len(data))
	for i, p := range data {
		projects[i] = model.Project{
			ID:            p.ID,
			Name:          p.Name,
			Description:   p.Description,
			Remark:        p.Remark,
			Status:        p.Status,
			ManagerID:     p.ManagerID,
			Priority:      p.Priority,
			StartedAt:     p.StartedAt,
			RemainingDays: util.CalculateRemainingDays(p.StartedAt, p.Duration),
			Duration:      p.Duration,
			FinishedAt:    p.FinishedAt,
		}
	}

	// 排序
	slices.SortFunc(projects, func(a, b model.Project) int {
		return util.SortByRemaining(a, b)
	})

	return projects, nil
}

func (p *ProjectService) ListProjectsTodo(ctx context.Context, req model.ListProjectsRequest) ([]model.Project, error) {
	// 获取数据
	data, err := p.querier.ListProjectsTodo(ctx, req.ManagerID)
	if err != nil {
		return nil, err
	}

	projects := make([]model.Project, len(data))
	for i, p := range data {
		projects[i] = model.Project{
			ID:            p.ID,
			Name:          p.Name,
			Description:   p.Description,
			Remark:        p.Remark,
			Status:        p.Status,
			ManagerID:     p.ManagerID,
			Priority:      p.Priority,
			StartedAt:     p.StartedAt,
			RemainingDays: util.CalculateRemainingDays(p.StartedAt, p.Duration),
			Duration:      p.Duration,
			FinishedAt:    p.FinishedAt,
		}
	}
	return projects, nil

}

func (p *ProjectService) ListProjectsDoing(ctx context.Context, req model.ListProjectsRequest) ([]model.Project, error) {
	// 获取数据
	data, err := p.querier.ListProjectsDoing(ctx, repo.ListProjectsDoingParams{
		ManagerID: req.ManagerID,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
	})

	if err != nil {
		return nil, err
	}

	projects := make([]model.Project, len(data))
	for i, p := range data {
		projects[i] = model.Project{
			ID:            p.ID,
			Name:          p.Name,
			Description:   p.Description,
			Remark:        p.Remark,
			Status:        p.Status,
			ManagerID:     p.ManagerID,
			Priority:      p.Priority,
			StartedAt:     p.StartedAt,
			RemainingDays: util.CalculateRemainingDays(p.StartedAt, p.Duration),
			Duration:      p.Duration,
			FinishedAt:    p.FinishedAt,
		}
	}

	return projects, nil
}

func (p *ProjectService) ListProjectsDone(ctx context.Context, req model.ListProjectsRequest) ([]model.Project, error) {
	// 获取数据
	data, err := p.querier.ListProjectsDone(ctx, repo.ListProjectsDoneParams{
		ManagerID: req.ManagerID,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
	})

	if err != nil {
		return nil, err
	}

	projects := make([]model.Project, len(data))
	for i, p := range data {
		projects[i] = model.Project{
			ID:            p.ID,
			Name:          p.Name,
			Description:   p.Description,
			Remark:        p.Remark,
			Status:        p.Status,
			ManagerID:     p.ManagerID,
			Priority:      p.Priority,
			StartedAt:     p.StartedAt,
			RemainingDays: util.CalculateRemainingDays(p.StartedAt, p.Duration),
			Duration:      p.Duration,
			FinishedAt:    p.FinishedAt,
		}
	}

	return projects, nil
}
