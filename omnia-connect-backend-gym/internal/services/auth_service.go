package services

import (
	"api-gateway/internal/config"
	"api-gateway/internal/domain"
	_ "api-gateway/internal/models"
	"api-gateway/internal/repositories"
	"api-gateway/internal/services/jwt"
	"api-gateway/pkg/logger"
	"context"
	"errors"
	"fmt"
	_ "golang.org/x/crypto/bcrypt"
	"log/slog"
	"regexp"
)

var (
	ErrorForbidden    = errors.New("forbidden")
	ErrorUnauthorized = errors.New("unauthorized")
)

type AuthService interface {
	Validate(ctx context.Context, method, path, accessToken string) (*domain.User, error)
	//Email(ctx context.Context, request *models.UserLoginRequest) (*domain.UserTokens, error)
	//ChangePassword(ctx context.Context, accessToken string, request *models.UserChangePassword) (*domain.User, error)
	Refresh(ctx context.Context, refreshToken string) (*domain.UserTokens, error)
	IsPublicRoute(ctx context.Context, method, path string) bool
	CreateTokens(ctx context.Context, id int) (*domain.UserTokens, error)
	Authenticate(ctx context.Context, accessToken string) (*domain.User, error)
	//CreateUser(ctx context.Context, user *models.UserRegisterRequest) (*domain.User, error)
}

type authService struct {
	jwtService     jwt.JwtService
	userRepository repositories.UserRepository
	publicRoutes   []config.RouteConfig
	privateRoutes  []config.RouteConfig
}

func NewAuthService(jwtService jwt.JwtService, userRepository repositories.UserRepository,
	publicRoutes []config.RouteConfig, privateRoutes []config.RouteConfig) AuthService {
	return &authService{
		jwtService:     jwtService,
		userRepository: userRepository,
		publicRoutes:   publicRoutes,
		privateRoutes:  privateRoutes,
	}
}

func (as *authService) IsPublicRoute(ctx context.Context, method, path string) bool {
	fmt.Println(method, path, as.publicRoutes)
	for _, public := range as.publicRoutes {
		if matchPattern(public.Path, path) && public.Method == method {
			slog.DebugContext(ctx, "public route", slog.String("method", method), slog.String("path", path))
			return true
		}
	}
	return false
}

func (as *authService) Validate(ctx context.Context, method, path, accessToken string) (*domain.User, error) {

	userID, err := as.jwtService.VerifyToken(ctx, accessToken, jwt.AccessTokenType)
	if err != nil {
		return nil, ErrorUnauthorized
	}

	user, err := as.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrorUnauthorized
	}

	for _, private := range as.privateRoutes {
		if matchPattern(private.Path, path) && private.Method == method {

			if len(private.Roles) == 0 {
				return user, nil
			}

			for _, role := range private.Roles {
				for _, userRole := range user.Roles {
					if role == userRole {
						return user, nil
					}
				}
			}

			return nil, ErrorForbidden
		}
	}

	return user, nil
}

//func (as *authService) ChangePassword(ctx context.Context, accessToken string, request *models.UserChangePassword) (*domain.User, error) {
//
//	userID, err := as.jwtService.VerifyToken(ctx, accessToken, jwt.AccessTokenType)
//	if err != nil {
//		return nil, ErrorUnauthorized
//	}
//
//	user, err := as.userRepository.GetUserByID(ctx, userID)
//	if err != nil {
//		return nil, ErrorUnauthorized
//	}
//
//	if user.Password == nil {
//		return nil, ErrorUnauthorized
//	}
//
//	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(request.OldPassword))
//	if err != nil {
//		return nil, ErrorUnauthorized
//	}
//
//	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
//	if err != nil {
//		return nil, ErrorUnauthorized
//	}
//
//	hp := string(hashPassword)
//
//	user.Password = &hp
//
//	if err = as.userRepository.UpdateUser(ctx, user); err != nil {
//		return nil, ErrorUnauthorized
//	}
//
//	return user, nil
//}

//func (as *authService) CreateUser(ctx context.Context, userRegister *models.UserRegisterRequest) (*domain.User, error) {
//	//pass, _ := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)
//	//p := string(pass)
//
//	user := &domain.User{}
//
//	err := as.userRepository.CreateUser(ctx, user)
//	if err != nil {
//		return nil, err
//	}
//	return user, nil
//}

func (as *authService) Authenticate(ctx context.Context, accessToken string) (*domain.User, error) {
	userID, err := as.jwtService.VerifyToken(ctx, accessToken, jwt.AccessTokenType)
	if err != nil {
		return nil, err
	}

	user, err := as.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (as *authService) CreateTokens(ctx context.Context, id int) (*domain.UserTokens, error) {
	userTokens, err := as.jwtService.CreateJwtTokens(ctx, id)
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, "Generated tokens for user", slog.Any("tokens", userTokens))
	return userTokens, nil
}

//func (as *authService) Email(ctx context.Context, request *models.UserLoginRequest) (*domain.UserTokens, error) {
//	user, err := as.userRepository.GetUserByLogin(ctx, request.Email)
//	if err != nil {
//		return nil, err
//	}
//
//	if user.Password == nil {
//		return nil, ErrorUnauthorized
//	}
//
//	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(request.Password))
//	if err != nil {
//		return nil, ErrorUnauthorized
//	}
//
//	logger.WithUserID(ctx, user.ID)
//
//	userTokens, err := as.CreateTokens(ctx, user.ID)
//	return userTokens, nil
//}

func (as *authService) Refresh(ctx context.Context, refreshToken string) (*domain.UserTokens, error) {
	userID, err := as.jwtService.VerifyToken(ctx, refreshToken, jwt.RefreshTokenType)
	if err != nil {
		return nil, err
	}

	_, err = as.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	logger.WithUserID(ctx, userID)

	userTokens, err := as.jwtService.CreateJwtTokens(ctx, userID)
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, "Generated tokens for user", slog.Any("tokens", userTokens))

	return userTokens, nil

}

func matchPattern(pattern, path string) bool {
	regexPattern := "^" + regexp.QuoteMeta(pattern) + "$"
	regexPattern = regexp.MustCompile(`\\\*`).ReplaceAllString(regexPattern, ".*") // Заменяем * на .*

	// Компилируем регулярное выражение
	re := regexp.MustCompile(regexPattern)

	// Проверяем соответствие
	return re.MatchString(path)
}
