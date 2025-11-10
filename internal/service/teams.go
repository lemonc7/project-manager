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

	// 插入数据库
	_, err = t.querier.CreateTeam(ctx, repo.CreateTeamParams{
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		return err
	}

	return nil
}
