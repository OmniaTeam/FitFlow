package handlers

import (
	"api-gateway/internal/services"
	"api-gateway/internal/services/jwt"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AuthHandler struct {
	accessTokenMaxAge  int
	refreshTokenMaxAge int
	authService        services.AuthService
	oauthGoogleService services.OauthGoogleService
	oauthVkidService   services.OAuthVkidServiceImpl
	redirect           string
}

func NewAuthHandler(authService services.AuthService, oauthGoogleService services.OauthGoogleService,
	oAuthVkidServiceImpl services.OAuthVkidServiceImpl,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration,
	redirect string) *AuthHandler {
	return &AuthHandler{
		authService:        authService,
		oauthGoogleService: oauthGoogleService,
		oauthVkidService:   oAuthVkidServiceImpl,
		accessTokenMaxAge:  int(accessTokenTTL.Seconds()),
		refreshTokenMaxAge: int(refreshTokenTTL.Seconds()),
		redirect:           redirect,
	}
}

const (
	vkAuthURL           = "https://id.vk.com/authorize"
	vkTokenURL          = "https://id.vk.com/access_token"
	vkAPIURL            = "https://api.vk.ru/method/users.get"
	defaultScope        = "vkid.personal_info"
	codeChallengeMethod = "S256"
)

type PKCEParams struct {
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

type VKTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	UserID       int    `json:"user_id"`
	Email        string `json:"email,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type VKUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type VKUserResponse struct {
	Response []VKUser `json:"response"`
}

// CreateUser godoc
// @Summary      Create user
// @Tags         auth
// @Accept       json
// @Param        userCreate body models.UserRegisterRequest true "Data for register"
// @Success      201
// @Failure 401
// @Router       /auth/create [post]
//func (ah *AuthHandler) CreateUser(c *gin.Context) {
//	ctx := c.Request.Context()
//
//	var userRegister models.UserRegisterRequest
//	if err := c.ShouldBindJSON(&userRegister); err != nil {
//		slog.InfoContext(ctx, "error binding", slog.String("error", err.Error()))
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	user, err := ah.authService.CreateUser(ctx, &userRegister)
//	if err != nil {
//		if errors.Is(err, repositories.ErrorUserAlreadyExist) {
//			c.JSON(http.StatusBadRequest, gin.H{"user": "already exists"})
//			return
//		}
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusCreated, gin.H{"user": user})
//}

func (ah *AuthHandler) GetUser(c *gin.Context) {

	ctx := c.Request.Context()

	accessToken, err := c.Cookie("access_token")
	if err != nil {
		fmt.Println(accessToken, err)
		c.Status(http.StatusUnauthorized)
		return
	}

	user, err := ah.authService.Authenticate(ctx, accessToken)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, user)

}

func (ah *AuthHandler) Validate(c *gin.Context) {
	ctx := c.Request.Context()

	originalPath := c.GetHeader("X-Original-URI")

	if ah.authService.IsPublicRoute(ctx, c.Request.Method, originalPath) {
		c.Status(http.StatusOK)
		return
	}

	accessToken, err := c.Cookie(string(jwt.AccessTokenType))
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	user, err := ah.authService.Validate(ctx, c.Request.Method, originalPath, accessToken)
	if err != nil {
		if errors.Is(err, services.ErrorForbidden) {
			c.Status(http.StatusForbidden)
			return
		} else {
			c.Status(http.StatusUnauthorized)
			return
		}

	}

	c.Header("X-User-Id", strconv.Itoa(user.ID))
	c.Status(http.StatusOK)
}

// Login godoc
// @Summary      Login user
// @Tags         auth
// @Accept       json
// @Param        userLogin body models.UserLoginRequest true "Data for login"
// @Success      200
// @Failure 401
// @Router       /auth/login [post]
//func (ah *AuthHandler) Login(c *gin.Context) {
//
//	ctx := c.Request.Context()
//
//	var userLogin models.UserLoginRequest
//	if err := c.ShouldBindJSON(&userLogin); err != nil {
//		slog.InfoContext(ctx, "error binding", slog.String("error", err.Error()))
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	tokens, err := ah.authService.Login(ctx, &userLogin)
//	if err != nil {
//		c.Status(http.StatusUnauthorized)
//		return
//	}
//
//	c.SetCookie(string(jwt.AccessTokenType), tokens.AccessToken, ah.accessTokenMaxAge, "/", "", true, true)
//	c.SetCookie(string(jwt.RefreshTokenType), tokens.RefreshToken, ah.refreshTokenMaxAge, "/auth/refresh", "", true, true)
//	c.Status(http.StatusOK)
//}

// Refresh godoc
// @Summary      Refresh tokens
// @Tags         auth
// @Success      200
// @Failure 401
// @Failure 403
// @Router       /auth/refresh [get]
func (ah *AuthHandler) Refresh(c *gin.Context) {

	ctx := c.Request.Context()

	refreshToken, err := c.Cookie(string(jwt.RefreshTokenType))
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	tokens, err := ah.authService.Refresh(ctx, refreshToken)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.SetCookie(string(jwt.AccessTokenType), tokens.AccessToken, ah.accessTokenMaxAge, "/", "", true, true)
	c.SetCookie(string(jwt.RefreshTokenType), tokens.RefreshToken, ah.refreshTokenMaxAge, "/auth/refresh", "", true, true)

	c.Status(http.StatusOK)

}

// ChangePassword godoc
// @Summary      Change user password
// @Tags         auth
// @Accept       json
// @Param        changePassword body models.UserChangePassword true "Data for change password"
// @Success      200
// @Failure 401
// @Router       /auth/change_password [post]
//func (ah *AuthHandler) ChangePassword(c *gin.Context) {
//	ctx := c.Request.Context()
//
//	var userChangePassword models.UserChangePassword
//	if err := c.ShouldBindJSON(&userChangePassword); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	accessToken, err := c.Cookie(string(jwt.AccessTokenType))
//	if err != nil {
//		c.Status(http.StatusUnauthorized)
//		return
//	}
//
//	_, err = ah.authService.ChangePassword(ctx, accessToken, &userChangePassword)
//	if err != nil {
//		if errors.Is(err, services.ErrorUnauthorized) {
//			c.Status(http.StatusUnauthorized)
//			return
//		}
//		c.Status(http.StatusInternalServerError)
//		return
//	}
//
//	c.Status(200)
//}

// Logout godoc
// @Summary      Logout user
// @Tags         auth
// @Success      200
// @Router       /auth/logout [get]
func (ah *AuthHandler) Logout(c *gin.Context) {

	c.SetCookie(string(jwt.AccessTokenType), "", -1, "/", "", true, true)
	c.SetCookie(string(jwt.RefreshTokenType), "", -1, "/auth/refresh", "", true, true)

	c.Status(http.StatusOK)

}

// AuthGoogleLogin godoc
// @Summary      AuthGoogleLogin
// @Tags         auth
// @Success      301
// @Failure 401
// @Failure 403
// @Router       /auth/google/login [get]
func (ah *AuthHandler) AuthGoogleLogin(c *gin.Context) {
	ctx := c.Request.Context()
	state, url := ah.oauthGoogleService.OauthGoogleLogin(ctx)

	c.SetCookie("oauth_state", state, 60, "/api/auth/google/callback", "", true, true)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// AuthVkidLogin godoc
// @Summary      AuthVkidLogin
// @Tags         auth
// @Success      301
// @Failure 401
// @Failure 403
// @Router       /auth/vkid/login [get]
func (ah *AuthHandler) AuthVkidLogin(c *gin.Context) {
	//pkce, err := generatePKCEParams()
	//if err != nil {
	//	c.AbortWithError(http.StatusInternalServerError, err)
	//	return
	//}
	//
	//// Генерация state
	//state, err := generateRandomString(32)
	//if err != nil {
	//	c.AbortWithError(http.StatusInternalServerError, err)
	//	return
	//}
	//
	//// Сохранение в сессию (в production используйте защищенное хранилище)
	//
	//// Формирование URL для аутентификации
	//authURL := fmt.Sprintf(
	//	"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&code_challenge=%s&code_challenge_method=%s",
	//	vkAuthURL,
	//	"53326175",
	//	"https://theomnia.ru/auth/vkid/callback",
	//	strings.Join([]string{"vkid.personal_info", "email"}, ","),
	//	state,
	//	pkce.CodeChallenge,
	//	pkce.CodeChallengeMethod,
	//)

	authURL := ah.oauthVkidService.VkAuthUrl()

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (ah *AuthHandler) AuthVkidCallback(c *gin.Context) {

	ctx := c.Request.Context()

	code := c.Query("code")
	state := c.Query("state")
	device := c.Query("device_id")

	fmt.Println("state", state)
	fmt.Println("code", code)
	fmt.Println("device", device)

	if code == "" || state == "" {
		c.Redirect(http.StatusPermanentRedirect, "https://theomnia.ru")
		return
	}

	user, err := ah.oauthVkidService.GetUser(code, state, device)
	if err != nil {
		fmt.Println("get user", err)
		return
	}

	tokens, err := ah.authService.CreateTokens(ctx, user.ID)
	if err != nil {
		fmt.Println("create tokens", err)
		return
	}

	c.SetCookie(string(jwt.AccessTokenType), tokens.AccessToken, ah.accessTokenMaxAge, "/", "", true, true)
	c.SetCookie(string(jwt.RefreshTokenType), tokens.RefreshToken, ah.refreshTokenMaxAge, "/auth/refresh", "", true, true)

	c.Redirect(http.StatusPermanentRedirect, "https://theomnia.ru")

	fmt.Println(user)
}

func (ah *AuthHandler) AuthGoogleCallback(c *gin.Context) {

	ctx := c.Request.Context()

	state, err := c.Cookie("oauth_state")
	if err != nil {
		slog.Error("no state", slog.String("error", err.Error()))
		c.Redirect(http.StatusPermanentRedirect, "https://theomnia.ru")
		return
	}

	value := c.Query("state")

	if value != state {
		slog.Error("invalid state")
		c.Redirect(http.StatusPermanentRedirect, "https://theomnia.ru")
		return
	}

	code := c.Query("code")
	user, err := ah.oauthGoogleService.GetUserDataFromGoogle(ctx, code)
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusPermanentRedirect, "https://theomnia.ru")
		return
	}

	slog.InfoContext(ctx, "user from google", slog.Any("user", user))

	tokens, err := ah.authService.CreateTokens(ctx, user.ID)
	if err != nil {
		return
	}

	c.SetCookie(string(jwt.AccessTokenType), tokens.AccessToken, ah.accessTokenMaxAge, "/", "", true, true)
	c.SetCookie(string(jwt.RefreshTokenType), tokens.RefreshToken, ah.refreshTokenMaxAge, "/auth/refresh", "", true, true)

	c.Redirect(http.StatusPermanentRedirect, "https://theomnia.ru")

	fmt.Println(user)

}

func generatePKCEParams() (PKCEParams, error) {
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return PKCEParams{}, err
	}

	codeChallenge := generateCodeChallengeS256(codeVerifier)

	return PKCEParams{
		CodeVerifier:        codeVerifier,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}, nil
}

func generateCodeVerifier() (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	b := make([]byte, 64) // 64 символа - хороший баланс безопасности
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := range b {
		b[i] = chars[int(b[i])%len(chars)]
	}
	return string(b), nil
}

// Генерация code_challenge по RFC 7636: BASE64URL-ENCODE(SHA256(ASCII(code_verifier)))
func generateCodeChallengeS256(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// Генерация случайного state
func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func exchangeCodeForToken(clientID, clientSecret, redirectURI, code, codeVerifier string) (*VKTokenResponse, error) {
	reqBody := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s&code_verifier=%s",
		code, redirectURI, clientID, clientSecret, codeVerifier,
	)

	resp, err := http.Post(vkTokenURL, "application/x-www-form-urlencoded", strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token VKTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

// Получение информации о пользователе
func getVKUserInfo(accessToken string, userID int) (*VKUser, error) {
	url := fmt.Sprintf("%s?user_ids=%d&access_token=%s&v=5.131&fields=first_name,last_name",
		vkAPIURL, userID, accessToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Response []VKUser `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Response) == 0 {
		return nil, fmt.Errorf("no user data in response")
	}

	return &response.Response[0], nil
}
