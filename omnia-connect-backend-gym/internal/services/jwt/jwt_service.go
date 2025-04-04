package jwt

import (
	"api-gateway/internal/domain"
	"context"
)

type JwtService interface {
	CreateJwtTokens(ctx context.Context, userID int) (*domain.UserTokens, error)
	VerifyToken(ctx context.Context, token string, tokenType TokenType) (int, error)
}
