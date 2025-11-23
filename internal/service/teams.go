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

func (t *TeamService) CreateTeam(ctx context.Context, req model.CreateTeamRequest) (repo.CreateTeamRow, error) {
	// 插入数据
	team, err := t.querier.CreateTeam(ctx, repo.CreateTeamParams{
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		return repo.CreateTeamRow{}, err
	}

	return team, nil
}

func (t *TeamService) UpdateTeam(ctx context.Context, id int32, req model.UpdateTeamRequest) (repo.UpdateTeamRow, error) {
	// 更新数据
	team, err := t.querier.UpdateTeam(ctx, repo.UpdateTeamParams{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		return repo.UpdateTeamRow{}, err
	}

	return team, err
}

func (t *TeamService) DeleteTeam(ctx context.Context, id int32) error {
	// 删除数据
	if err := t.querier.DeleteTeam(ctx, id); err != nil {
		return err
	}

	return nil
}

func (t *TeamService) ListTeams(ctx context.Context) ([]repo.ListTeamsRow, error) {
	// 获取数据
	teams, err := t.querier.ListTeams(ctx)
	if err != nil {
		return nil, err
	}

	return teams, nil
}
