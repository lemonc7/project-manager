package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
)

type EmployeeService struct {
	db             *pgxpool.Pool
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

	tx, err := e.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := repo.New(e.db).WithTx(tx)

	// 创建employee
	id, err := qtx.CreateEmployee(ctx, repo.CreateEmployeeParams{
		Name:         req.Name,
		Username:     req.Username,
		PasswordHash: passwordHash,
		TeamID:       req.TeamID,
	})
	if err != nil {
		return fmt.Errorf("创建employee: %w", err)
	}

	// 先删除对应employee的所有权限
	if err := qtx.DeleteEmployeeRoles(ctx, id); err != nil {
		return fmt.Errorf("删除employee权限: %w", err)
	}

	// 添加employee权限
	if err := qtx.AddEmployeeRoles(ctx, repo.AddEmployeeRolesParams{
		EmployeeID: id,
		RoleIds:    req.Roles,
	}); err != nil {
		return fmt.Errorf("添加employee权限: %w", err)
	}

	return tx.Commit(ctx)
}

func (e *EmployeeService) UpdateEmployee(ctx context.Context, req model.UpdateEmployeeRequest) error {
	if len(req.Roles) == 0 {
		return e.querier.UpdateEmployee(ctx, repo.UpdateEmployeeParams{
			ID:       req.ID,
			Name:     req.Name,
			Username: req.Username,
			TeamID:   req.TeamID,
			IsActive: req.IsActive,
		})
	}

	tx, err := e.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := repo.New(e.db).WithTx(tx)

	if err := qtx.UpdateEmployee(ctx, repo.UpdateEmployeeParams{
		ID:       req.ID,
		Name:     req.Name,
		Username: req.Username,
		TeamID:   req.TeamID,
		IsActive: req.IsActive,
	}); err != nil {
		return fmt.Errorf("更新employee: %w", err)
	}

	if err := qtx.DeleteEmployeeRoles(ctx, req.ID); err != nil {
		return fmt.Errorf("删除employee权限: %w", err)
	}

	if err := qtx.AddEmployeeRoles(ctx, repo.AddEmployeeRolesParams{
		EmployeeID: req.ID,
		RoleIds:    req.Roles,
	}); err != nil {
		return fmt.Errorf("添加employee权限: %w", err)
	}

	return tx.Commit(ctx)
}

func (e *EmployeeService) DeleteEmployee(ctx context.Context, id int32) error {
	return e.querier.UpdateEmployee(ctx, repo.UpdateEmployeeParams{
		ID:       id,
		IsActive: new(false),
	})
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
		return "", err
	}

	return employee.Name, nil
}
