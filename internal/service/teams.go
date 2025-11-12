package service

import (
	"context"
	"database/sql"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/pkg/util"
)

type TeamService struct {
	querier repo.Querier
}

func NewTeamService(db *sql.DB) *TeamService {
	return &TeamService{
		querier: repo.New(db),
	}
	
}

func (t *TeamService) CreateTeam(ctx context.Context, req model.CreateTeamRequest) error {
	// 插入数据
	err := t.querier.CreateTeam(ctx, repo.CreateTeamParams{
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		return err
	}

	return nil
}

func (t *TeamService) UpdateTeam(ctx context.Context, req model.UpdateTeamRequest) error {
	// 更新数据
	err := t.querier.UpdateTeam(ctx, repo.UpdateTeamParams{
		ID:          req.ID,
		Name:        util.ToNullString(req.Name),
		Description: util.ToNullString(req.Description),
	})

	if err != nil {
		return err
	}

	return nil
}

func (t *TeamService) DeleteTeam(ctx context.Context, req model.DeleteTeamRequest) error {
	// 删除数据
	if err := t.querier.DeleteTeam(ctx, req.ID); err != nil {
		return err
	}

	return nil
}

func (t *TeamService) ListTeams(ctx context.Context) ([]repo.Team, error) {
	// 获取数据
	teams, err := t.querier.ListTeams(ctx)
	if err != nil {
		return nil, err
	}

	return teams, nil
}
