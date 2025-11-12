package service

import (
	"context"
	"database/sql"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/pkg/util"
)

type EmployeeService struct {
	querier repo.Querier
}

func NewEmployeeService(db *sql.DB) *EmployeeService {
	return &EmployeeService{
		querier: repo.New(db),
	}
}

func (e *EmployeeService) CreateEmployee(ctx context.Context, req model.CreateEmployeeRequest) error {
	// 插入数据
	err := e.querier.CreateEmployee(ctx, repo.CreateEmployeeParams{
		Name:   req.Name,
		TeamID: util.ToNullInt32(req.TeamID),
	})

	if err != nil {
		return err
	}

	return nil
}

func (e *EmployeeService) UpdateEmployee(ctx context.Context, req model.UpdateEmployeeRequest) error {
	// 更新数据
	err := e.querier.UpdateEmployee(ctx, repo.UpdateEmployeeParams{
		ID:       req.ID,
		Name:     util.ToNullString(req.Name),
		TeamID:   util.ToNullInt32(req.TeamID),
		IsActive: util.ToNullBool(req.IsActive),
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
