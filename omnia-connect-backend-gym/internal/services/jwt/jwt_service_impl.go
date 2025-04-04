package jwt

import (
	"api-gateway/internal/config"
	"api-gateway/internal/domain"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"time"
)

type TokenType string

var (
	AccessTokenType  TokenType = "access_token"
	RefreshTokenType TokenType = "refresh_token"

	TokenTypeKey = "token_type"
)

type JwtServiceImpl struct {
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	secretKey       string
}

func NewJwtService(config config.JwtConfig) *JwtServiceImpl {
	return &JwtServiceImpl{
		accessTokenTTL:  config.AccessTokenTTL,
		refreshTokenTTL: config.RefreshTokenTTL,
		secretKey:       config.Secret,
	}
}

func (j *JwtServiceImpl) VerifyToken(ctx context.Context, token string, tokenType TokenType) (int, error) {

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil || !parsedToken.Valid {
		slog.DebugContext(ctx, "Token is invalid", slog.String("token", token))
		return 0, ErrorInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		slog.DebugContext(ctx, "Invalid token", slog.String("token", token))
		return 0, ErrorInvalidToken
	}

	if val, ok := claims[TokenTypeKey]; !ok || val != string(tokenType) {
		slog.DebugContext(ctx, "Token type mismatch", slog.String("token", token))
		return 0, ErrorInvalidToken
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID in token")
	}

	return int(userID), nil

}

func (j *JwtServiceImpl) CreateJwtTokens(ctx context.Context, userID int) (*domain.UserTokens, error) {
	accessToken, err := j.createAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.createRefreshToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JwtServiceImpl) createAccessToken(ctx context.Context, userID int) (string, error) {
	return j.createToken(ctx, userID, j.accessTokenTTL, AccessTokenType)
}

func (j *JwtServiceImpl) createRefreshToken(ctx context.Context, userID int) (string, error) {
	return j.createToken(ctx, userID, j.refreshTokenTTL, RefreshTokenType)
}

func (j *JwtServiceImpl) createToken(ctx context.Context, userID int, tokenTTL time.Duration, tokenType TokenType) (string, error) {
	claims := jwt.MapClaims{
		"id":         userID,
		"exp":        time.Now().Add(tokenTTL).Unix(),
		TokenTypeKey: tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		slog.ErrorContext(ctx, "Failed signed token", slog.String("error", err.Error()))
		return "", err
	}
	return signedToken, nil
}
