package util

import (
	"net/http"
	"strconv"

	"github.com/lemonc7/zest"
)

func GetParamID(c *zest.Context) (int32, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return 0, zest.NewHTTPError(http.StatusUnprocessableEntity, "请求ID错误").Wrap(err)
	}
	return int32(id), nil
}
