package service

import (
	"context"
	"database/sql"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
)

type EmployeeService struct {
	querier repo.Querier
}

func (e *EmployeeService) CreateEmployee(ctx context.Context, req model.CreateEmployeeRequest) error {
	var teamID sql.NullInt32
	if req.TeamID != nil {
		teamID.Valid = true
		teamID.Int32 = *req.TeamID
	}

	// 插入数据
	_, err := e.querier.CreateEmployee(ctx, repo.CreateEmployeeParams{
		Name:   req.Name,
		TeamID: teamID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (e *EmployeeService) UpdateEmployee(ctx context.Context, req model.UpdateEmployeeRequest) error {
	var teamID sql.NullInt32
	if req.TeamID != nil {
		teamID.Valid = true
		teamID.Int32 = *req.TeamID
	}

	// 更新数据
	_, err := e.querier.UpdateEmployee(ctx, repo.UpdateEmployeeParams{
		ID:       req.ID,
		Name:     req.Name,
		TeamID:   teamID,
		IsActive: req.IsActive,
	})

	if err != nil {
		return err
	}

	return nil
}

func (e *EmployeeService) DeleteEmployee(ctx context.Context, req model.DeleteEmployeeRequest) error {
	// 删除数据
	if err := e.querier.DeleteEmployee(ctx, req.ID); err != nil {
		return err
	}

	return nil
}

func (e *EmployeeService) ListEmployees(ctx context.Context) ([]repo.ListEmployeesRow, error) {
	// 获取数据
	employees, err := e.querier.ListEmployees(ctx)
	if err != nil {
		return nil, err
	}

	return employees, nil
}
