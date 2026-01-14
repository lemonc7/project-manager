package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/service"
	"github.com/lemonc7/project-manager/pkg/util"
	"github.com/lemonc7/zest"
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

func (a *AuthHandler) Login(c *zest.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	res, err := a.authService.Login(c.Context(), req)
	if err != nil {
		var code int
		if errors.Is(err, model.ErrUsernameOrPasswordInvalid) {
			code = http.StatusBadRequest
		} else if errors.Is(err, model.ErrNoPermission) {
			code = http.StatusForbidden
		} else {
			code = http.StatusInternalServerError
		}
		return zest.NewHTTPError(code, err.Error()).Wrap(err)
	}
	return c.JSON(http.StatusOK, res)
}

func (a *AuthHandler) ResetPassword(c *zest.Context) error {
	id, err := util.GetParamID(c)
	if err != nil {
		return err
	}
	var req model.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := a.authService.ResetPassword(c.Context(), id, req.Password); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "重置密码失败").Wrap(err)
	}
	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功重置密码",
	})
}

func (a *AuthHandler) ChangePassword(c *zest.Context) error {
	// 获取jwt中的user
	username, ok := c.Get("username").(string)
	if !ok {
		return zest.NewHTTPError(http.StatusUnauthorized)
	}

	var req model.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := a.authService.ChangePassword(c.Context(), username, req); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, model.ErrOldPasswordInvalid) {
			code = http.StatusBadRequest
		}
		return zest.NewHTTPError(code, "更改密码失败").Wrap(err)
	}
	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功修改密码",
	})
}

func (a *AuthHandler) InitAdmin(c *zest.Context) error {
	var req model.InitAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := a.authService.InitAdmin(c.Context(), req); err != nil {
		return zest.NewHTTPError(http.StatusInternalServerError, "管理员初始化失败").Wrap(err)
	}
	return c.JSON(http.StatusOK, zest.Map{
		"message": "成功初始化管理员",
	})
}
