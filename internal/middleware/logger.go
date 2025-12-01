package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lemonc7/project-manager/pkg/logger"
	"go.uber.org/zap"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)
		if err != nil {
			return err
		}

		req := c.Request()
		res := c.Response()

		fields := []zap.Field{
			zap.Int("status", res.Status),
			zap.String("request", req.Method+" "+req.RequestURI),
			zap.String("host", req.Host),
			zap.String("remote_ip", c.RealIP()),
			zap.String("latency", time.Since(start).String()),
			zap.Int64("size", res.Size),
			zap.String("user_agent", req.UserAgent()),
		}

		id := req.Header.Get(echo.HeaderXRequestID)
		if id != "" {
			fields = append(fields, zap.String("request_id", id))
		}

		n := res.Status
		switch {
		case n >= 500:
			logger.Error("Server Error", fields...)
		case n >= 400:
			logger.Warn("Client Error", fields...)
		case n >= 300:
			logger.Info("Redirection", fields...)
		default:
			logger.Info("Finished", fields...)
		}
		return nil
	}
}
