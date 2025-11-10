package service

import (
	"context"
	"fmt"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
)

type TeamService struct {
	querier repo.Querier
}

func (t *TeamService) CreateTeam(ctx context.Context, req model.CreateTeamRequest) error {
	// 判断小组名是否存在
	exist, err := t.querier.IsTeamExists(ctx, req.Name)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("小组已存在")
	}

	// 插入数据
	_, err = t.querier.CreateTeam(ctx, repo.CreateTeamParams{
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		return err
	}

	return nil
}

func (t *TeamService) UpdateTeam(ctx context.Context, req model.UpdateTeamRequest) error {
	// 判断小组名是否存在
	team, err := t.querier.GetTeamByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if req.NewName != team.Name {
		exist, err := t.querier.IsTeamExists(ctx, req.NewName)
		if err != nil {
			return err
		}
		if exist {
			return fmt.Errorf("小组名已被占用")
		}
	}

	// 更新数据
	_, err = t.querier.UpdateTeam(ctx, repo.UpdateTeamParams{
		ID:          req.ID,
		Name:        req.NewName,
		Description: req.Description,
	})

	if err != nil {
		return err
	}

	return nil
}

func (t *TeamService) DeleteTeam(ctx context.Context, req model.DeleteTeamRequest) error {
	// 判断小组是否存在
	_, err := t.querier.GetTeamByID(ctx, req.ID)
	if err != nil {
		return err
	}

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
