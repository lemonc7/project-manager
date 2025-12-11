package util

import "github.com/labstack/echo/v4"

// 角色等级，数字越大，权限越高
var RoleLevel = map[string]int{
	"employee": 1,
	"manager":  2,
	"admin":    3,
}

func GetRoles(c echo.Context) []string {
	roles, ok := c.Get("roles").([]string)
	if !ok {
		return []string{}
	}
	return roles
}

func HasAtLeastRole(userRoles []string, required string) bool {
	reqLevel, ok := RoleLevel[required]
	if !ok {
		return false
	}
	
	for _, role := range userRoles {
		level, ok := RoleLevel[role]
		if !ok {
			continue
		}
		if level >= reqLevel {
			return true
		}
	}
	return false
}
