package util

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetParamID(c echo.Context) (int32, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("请求ID错误: %w", err))
	}
	return int32(id), nil
}
