package utils

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Предопределенные ошибки
var (
	ErrUserIdHeaderRequired = errors.New("X-User-Id header is required")
	ErrInvalidUserIdFormat  = errors.New("Invalid user ID format")
)

// ExtractUserIdFromHeader извлекает ID пользователя из заголовка X-User-Id
// Возвращает ID пользователя и ошибку, если ID не удалось извлечь или преобразовать
func ExtractUserIdFromHeader(c *gin.Context) (int, error) {
	userIdStr := c.Request.Header.Get("X-User-Id")
	if userIdStr == "" {
		return 0, ErrUserIdHeaderRequired
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		slog.Error("Failed to parse user ID", "error", err)
		return 0, ErrInvalidUserIdFormat
	}

	return userId, nil
}
