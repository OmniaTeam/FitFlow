package middlewares

import (
	"api-gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

func LoggingMiddleware(c *gin.Context) {
	start := time.Now().UTC()

	traceID := uuid.New().String()

	ctx := logger.WithTraceID(c.Request.Context(), traceID)
	c.Request = c.Request.WithContext(ctx)

	c.Next()

	end := time.Now().UTC()
	latency := end.Sub(start).Milliseconds()
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	fields := []any{
		slog.Int("status", c.Writer.Status()),
		slog.String("method", c.Request.Method),
		slog.String("path", path),
		slog.String("query", query),
		slog.String("ip", c.ClientIP()),
		slog.String("user-agent", c.Request.UserAgent()),
		slog.Int("latency", int(latency)),
	}
	if len(c.Errors) > 0 {
		for _, e := range c.Errors.Errors() {
			slog.ErrorContext(ctx, e, fields...)
		}
	} else {
		slog.InfoContext(ctx, path, fields...)
	}
}
