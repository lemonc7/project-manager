package util

import (
	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/zest"
)

func GetRoles(c *zest.Context) []int32 {
	roles, ok := c.Get("roles").([]int32)
	if !ok {
		return []int32{}
	}
	return roles
}

func HasAtLeastRole(userRoles []int32, required model.Role) bool {
	for _, role := range userRoles {
		if role <= int32(required) {
			return true
		}
	}
	return false
}
