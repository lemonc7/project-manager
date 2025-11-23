package model

import "errors"

var (
	ErrUsernameOrPasswordInvalid = errors.New("用户名或密码错误")
	ErrOldPasswordInvalid        = errors.New("旧密码错误")
	ErrNoPermission              = errors.New("没有权限")
)
