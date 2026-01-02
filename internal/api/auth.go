package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/util"
)

type AuthService interface {
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
	ResetPassword(ctx context.Context, userID int32, password string) error
	ChangePassword(ctx context.Context, username string, req model.ChangePasswordRequest) error
	InitAdmin(ctx context.Context, req model.InitAdminRequest) error
}

var _ AuthService = (*service.AuthService)(nil)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{
		authService: service,
	}
}

func (a *AuthHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := util.BindRequest(
		c.Request().Body,
		model.LoginRequestSchema,
		&req,
	); err != nil {
		return err
	}

	res, err := a.authService.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, model.ErrUsernameOrPasswordInvalid) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("登录: %w", err))
		}
		if errors.Is(err, model.ErrNoPermission) {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("登录: %w", err))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("登录: %w", err))
	}
	return c.JSON(http.StatusOK, res)
}

func (a *AuthHandler) ResetPassword(c echo.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	var req model.ResetPasswordRequest
	if err := util.BindRequest(
		c.Request().Body,
		model.ResetPasswordRequestSchema,
		&req,
	); err != nil {
		return err
	}

	if err := a.authService.ResetPassword(c.Request().Context(), id, req.Password); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("重置密码: %w", err))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功重置密码",
	})
}

func (a *AuthHandler) ChangePassword(c echo.Context) error {
	// 获取jwt中的user
	username, ok := c.Get("username").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("未授权"))
	}

	var req model.ChangePasswordRequest
	if err := util.BindRequest(
		c.Request().Body,
		model.ChangePasswordRequestSchema,
		&req,
	); err != nil {
		return err
	}

	if err := a.authService.ChangePassword(c.Request().Context(), username, req); err != nil {
		if errors.Is(err, model.ErrOldPasswordInvalid) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("更改密码: %w", err))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("更改密码: %w", err))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功修改密码",
	})
}

func (a *AuthHandler) InitAdmin(c echo.Context) error {
	var req model.InitAdminRequest
	if err := util.BindRequest(
		c.Request().Body,
		model.InitAdminRequestSchema,
		&req,
	); err != nil {
		return err
	}

	if err := a.authService.InitAdmin(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("管理员初始化: %w", err))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "成功初始化管理员",
	})
}
