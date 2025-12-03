package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword string, password string) bool
}

type JWTer interface {
	Generate(username string, id int32, roles []string) (string, error)
}

type AuthService struct {
	db             *pgxpool.Pool
	querier        repo.Querier
	passwordHasher PasswordHasher
	jwter          JWTer
}

func NewAuthService(db *pgxpool.Pool, ph PasswordHasher, j JWTer) *AuthService {
	return &AuthService{
		db:             db,
		querier:        repo.New(db),
		passwordHasher: ph,
		jwter:          j,
	}
}

const ADMIN_INITIALIZED = "admin_initialized"

func (a *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	// 从数据库中获取数据
	user, err := a.querier.GetEmployeeByUsername(ctx, req.Username)
	if err != nil {
		return nil, model.ErrUsernameOrPasswordInvalid
	}

	// 验证密码
	if !a.passwordHasher.ComparePassword(user.PasswordHash, req.Password) {
		return nil, model.ErrUsernameOrPasswordInvalid
	}

	// 验证权限
	if !slices.Contains(user.Roles, "admin") &&
		!slices.Contains(user.Roles, "manager") &&
		!slices.Contains(user.Roles, "employee") {
		return nil, model.ErrNoPermission
	}

	// 生成token
	accessToken, err := a.jwter.Generate(req.Username, user.ID, user.Roles)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &model.LoginResponse{
		AccessToken: accessToken,
		UserID:      user.ID,
		Username:    user.Username,
		Roles:       user.Roles,
	}, nil
}

func (a *AuthService) ResetPassword(ctx context.Context, req model.ResetPasswordRequest) error {
	passwordHash, err := a.passwordHasher.HashPassword(req.Password)
	if err != nil {
		return err
	}

	err = a.querier.UpdatePassword(ctx, repo.UpdatePasswordParams{
		ID:           req.ID,
		PasswordHash: passwordHash,
	})

	return err
}

func (a *AuthService) ChangePassword(ctx context.Context, username string, req model.ChangePasswordRequest) error {
	// 获取用户
	user, err := a.querier.GetEmployeeByUsername(ctx, username)
	if err != nil {
		return err
	}

	// 验证旧密码
	if ok := a.passwordHasher.ComparePassword(user.PasswordHash, req.OldPassword); !ok {
		return model.ErrOldPasswordInvalid
	}

	// 生成新密码
	newHashed, err := a.passwordHasher.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// 更新密码
	err = a.querier.UpdatePassword(ctx, repo.UpdatePasswordParams{
		ID:           user.ID,
		PasswordHash: newHashed,
	})

	return err
}

func (a *AuthService) InitAdmin(ctx context.Context, req model.InitAdminRequest) error {
	// 判断管理员是否可以初始化
	allowInit, err := a.querier.AllowInitAdmin(ctx)
	if err != nil {
		return fmt.Errorf("获取管理员初始状态失败: %w", err)
	}
	if !allowInit {
		return fmt.Errorf("管理员已存在,不允许再次初始化")
	}

	// 密码加密
	passwordHash, err := a.passwordHasher.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("生成密码hash失败: %w", err)
	}

	// 生成事务
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 使用事务
	qtx := repo.New(a.db).WithTx(tx)

	// 注册管理员
	err = qtx.CreateEmployee(ctx, repo.CreateEmployeeParams{
		Name:         req.Name,
		Username:     req.Username,
		PasswordHash: passwordHash,
		Roles:        []string{"admin"},
	})
	if err != nil {
		return err
	}

	// 添加admin_initialized设置为true
	value := "true"
	description := "管理员初始化状态"
	err = qtx.CreateSetting(ctx, repo.CreateSettingParams{
		Key:         ADMIN_INITIALIZED,
		Value:       value,
		Description: description,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
