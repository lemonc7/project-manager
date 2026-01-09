package util

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/lemonc7/zest"
)

func GetParamID(c *zest.Context) (int32, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return 0, zest.NewHTTPError(http.StatusBadRequest, fmt.Errorf("请求ID错误: %w", err))
	}
	return int32(id), nil
}
