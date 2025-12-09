package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
)

type TeamService struct {
	querier repo.Querier
}

func NewTeamService(db *pgxpool.Pool) *TeamService {
	return &TeamService{
		querier: repo.New(db),
	}

}

func (t *TeamService) CreateTeam(ctx context.Context, req model.CreateTeamRequest) error {
	return t.querier.CreateTeam(ctx, repo.CreateTeamParams{
		Name:        req.Name,
		Description: req.Description,
	})
}

func (t *TeamService) UpdateTeam(ctx context.Context, id int32, req model.UpdateTeamRequest) error {
	return t.querier.UpdateTeam(ctx, repo.UpdateTeamParams{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})
}

func (t *TeamService) DeleteTeam(ctx context.Context, id int32) error {
	return t.querier.DeleteTeam(ctx, id)
}

func (t *TeamService) ListTeams(ctx context.Context) ([]repo.ListTeamsRow, error) {
	// 获取数据
	teams, err := t.querier.ListTeams(ctx)
	if err != nil {
		return nil, err
	}

	return teams, nil
}
