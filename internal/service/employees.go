package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
)

type EmployeeService struct {
	querier repo.Querier
}

func NewEmployeeService(db *pgxpool.Pool) *EmployeeService {
	return &EmployeeService{
		querier: repo.New(db),
	}
}

func (e *EmployeeService) CreateEmployee(ctx context.Context, req model.CreateEmployeeRequest) (repo.Employee, error) {
	// 插入数据
	employee, err := e.querier.CreateEmployee(ctx, repo.CreateEmployeeParams{
		Name:   req.Name,
		TeamID: req.TeamID,
	})

	if err != nil {
		return repo.Employee{}, err
	}

	return employee, nil
}

func (e *EmployeeService) UpdateEmployee(ctx context.Context, id int32, req model.UpdateEmployeeRequest) (repo.Employee, error) {
	// 更新数据
	employee, err := e.querier.UpdateEmployee(ctx, repo.UpdateEmployeeParams{
		ID:       id,
		Name:     req.Name,
		TeamID:   req.TeamID,
		IsActive: req.IsActive,
	})

	if err != nil {
		return repo.Employee{}, err
	}

	return employee, nil
}

func (e *EmployeeService) DeleteEmployee(ctx context.Context, id int32) error {
	// 删除数据
	if err := e.querier.DeleteEmployee(ctx, id); err != nil {
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
