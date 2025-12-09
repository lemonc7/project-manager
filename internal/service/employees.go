package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
)

type EmployeeService struct {
	querier        repo.Querier
	passwordHasher PasswordHasher
}

func NewEmployeeService(db *pgxpool.Pool, ph PasswordHasher) *EmployeeService {
	return &EmployeeService{
		querier:        repo.New(db),
		passwordHasher: ph,
	}
}

func (e *EmployeeService) CreateEmployee(ctx context.Context, req model.CreateEmployeeRequest) error {
	passwordHash, err := e.passwordHasher.HashPassword(req.Password)
	if err != nil {
		return err
	}
	// 插入数据
	err = e.querier.CreateEmployee(ctx, repo.CreateEmployeeParams{
		Name:         req.Name,
		Username:     req.Username,
		PasswordHash: passwordHash,
		Roles:        req.Roles,
		TeamID:       req.TeamID,
	})

	return err
}

func (e *EmployeeService) UpdateEmployee(ctx context.Context, id int32, req model.UpdateEmployeeRequest) error {
	return e.querier.UpdateEmployee(ctx, repo.UpdateEmployeeParams{
		ID:       id,
		Name:     req.Name,
		Roles:    req.Roles,
		TeamID:   req.TeamID,
		IsActive: req.IsActive,
	})
}

func (e *EmployeeService) DeleteEmployee(ctx context.Context, id int32) error {
	return e.querier.DeleteEmployee(ctx, id)
}

func (e *EmployeeService) ListEmployees(ctx context.Context) ([]repo.ListEmployeesRow, error) {
	// 获取数据
	data, err := e.querier.ListEmployees(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *EmployeeService) GetEmployee(ctx context.Context, username string) (string, error) {
	// 获取数据
	employee, err := e.querier.GetEmployeeByUsername(ctx, username)
	if err != nil {
		return "", nil
	}

	return employee.Name, nil
}
