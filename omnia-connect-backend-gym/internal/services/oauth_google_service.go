package services

import (
	"api-gateway/internal/domain"
	"api-gateway/internal/repositories"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type UserInfoGoogle struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}
type OauthGoogleService interface {
	OauthGoogleLogin(ctx context.Context) (state, url string)
	GetUserDataFromGoogle(ctx context.Context, code string) (*domain.User, error)
}

type OAuthGoogleServiceImpl struct {
	userRepository    repositories.UserRepository
	conf              *oauth2.Config
	oauthGoogleUrlAPI string
}

func NewOauthGoogleService(userRepository repositories.UserRepository, clientID, clientSecret string) *OAuthGoogleServiceImpl {
	return &OAuthGoogleServiceImpl{
		userRepository: userRepository,
		conf: &oauth2.Config{
			RedirectURL:  "https://theomnia.ru/api/auth/google/callback",
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		},
		oauthGoogleUrlAPI: "https://www.googleapis.com/oauth2/v2/userinfo",
	}
}

func (gs *OAuthGoogleServiceImpl) OauthGoogleLogin(ctx context.Context) (state, url string) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie()

	u := gs.conf.AuthCodeURL(oauthState)
	return oauthState, u
}

func (gs *OAuthGoogleServiceImpl) GetUserDataFromGoogle(ctx context.Context, code string) (*domain.User, error) {

	token, err := gs.conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	client := gs.conf.Client(ctx, token)
	response, err := client.Get(gs.oauthGoogleUrlAPI)
	if err != nil {
		return nil, fmt.Errorf("failed getting userInfo info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	var userInfo UserInfoGoogle
	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal userInfo info: %s", err.Error())
	}
	fmt.Printf("%+v\n", userInfo)
	user, err := gs.userRepository.GetUserByGoogleID(ctx, userInfo.ID)
	slog.ErrorContext(ctx, "get user for google", slog.Any("error", err))
	if err != nil {
		if errors.Is(err, repositories.ErrorUserNotFound) {
			newUser := &domain.User{
				GoogleID:  &userInfo.ID,
				LastName:  userInfo.FamilyName,
				FirstName: userInfo.GivenName,
				Email:     userInfo.Email,
				Roles:     []string{"user"},
			}
			err := gs.userRepository.CreateUser(ctx, newUser)
			slog.InfoContext(ctx, "create user for google", slog.Any("user", newUser))
			if err != nil {
				slog.ErrorContext(ctx, "error create user", slog.String("error", err.Error()))
				return nil, err
			}
			return newUser, nil
		}
		return nil, err
	}
	fmt.Println(user)
	return user, nil
}

func generateStateOauthCookie() string {

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return state
}
