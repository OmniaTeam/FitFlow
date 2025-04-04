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
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const lengthCodeVerifier = 43
const lengthState = 20

type OAuthVkidServiceImpl struct {
	userRepository  repositories.UserRepository
	conf            *oauth2.Config
	oauthVkidUrlAPI string
	storage         *PKCEStorage
}

func NewVkidService(clientID, clientSecret string, userRepository repositories.UserRepository) *OAuthVkidServiceImpl {
	return &OAuthVkidServiceImpl{conf: &oauth2.Config{
		ClientID:     clientID,     // Из переменных окружения
		ClientSecret: clientSecret, // Из переменных окружения
		//RedirectURL:  "https://theomnia.ru/auth/vkid/callback",
		RedirectURL: "https://theomnia.ru/api/auth/vkid/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://id.vk.com/authorize",
			TokenURL: "https://id.vk.com/oauth2/auth",
		},
		Scopes: []string{"email", "phone"},
	},
		storage:        NewPKCEStorage(5 * time.Minute),
		userRepository: userRepository,
	}
}

type PKCEStorage struct {
	store      sync.Map // хранит map[string]*pkceEntry
	cleanupInt time.Duration
}

type pkceEntry struct {
	codeVerifier string
	expiresAt    time.Time
}

func NewPKCEStorage(cleanupInterval time.Duration) *PKCEStorage {
	s := &PKCEStorage{
		cleanupInt: cleanupInterval,
	}

	// Запускаем горутину для периодической очистки
	if cleanupInterval > 0 {
		go s.cleanupExpired()
	}

	return s
}

// Store сохраняет code_verifier с указанным TTL
func (s *PKCEStorage) Store(state, codeVerifier string) {
	s.store.Store(state, &pkceEntry{
		codeVerifier: codeVerifier,
		expiresAt:    time.Now().Add(5 * time.Minute),
	})
}

// GetAndDelete извлекает и удаляет code_verifier
func (s *PKCEStorage) GetAndDelete(state string) (string, bool) {
	value, loaded := s.store.LoadAndDelete(state)
	if !loaded {
		return "", false
	}

	entry := value.(*pkceEntry)
	if time.Now().After(entry.expiresAt) {
		return "", false
	}

	return entry.codeVerifier, true
}

func (s *PKCEStorage) cleanupExpired() {
	if s.cleanupInt <= 0 {
		return
	}

	for {
		time.Sleep(s.cleanupInt)

		now := time.Now()
		s.store.Range(func(key, value interface{}) bool {
			entry := value.(*pkceEntry)
			if now.After(entry.expiresAt) {
				s.store.Delete(key)
			}
			return true
		})
	}
}

type VKUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type VKUserResponse struct {
	Response []VKUser `json:"response"`
}

type VKTokenResponse struct {
	AccessToken string `json:"access_token"`
	Email       string `json:"email"`
	UserID      int    `json:"user_id"`
}

func (o *OAuthVkidServiceImpl) GetUser(code, state, device string) (*domain.User, error) {
	codeVerifier, _ := o.storage.GetAndDelete(state)
	//if !ok {
	//	return nil, fmt.Errorf("code verifier not found in storage")
	//}

	token, err := o.exchangeCodeForToken(code, codeVerifier, device)
	if err != nil {
		return nil, err
	}

	info, err := o.GetUserInfo(token.AccessToken)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(info.UserID)
	user, err := o.userRepository.GetUserByVkID(context.Background(), info.UserID)
	if err != nil {
		if errors.Is(err, repositories.ErrorUserNotFound) {
			newUser := &domain.User{
				VkID:      &info.UserID,
				LastName:  info.LastName,
				FirstName: info.FirstName,
				Email:     info.Email,
				Roles:     []string{"user"},
			}
			err := o.userRepository.CreateUser(context.Background(), newUser)
			slog.InfoContext(context.Background(), "create user for vkid", slog.Any("user", newUser))
			if err != nil {
				slog.ErrorContext(context.Background(), "error create user", slog.String("error", err.Error()))
				return nil, err
			}
			return newUser, nil
		}
		return nil, err
	}

	fmt.Println("return", user, err)
	return user, nil

}

func (o *OAuthVkidServiceImpl) VkAuthUrl() string {
	verifier, _ := GenerateCodeVerifier()
	state, _ := GenerateStateToken()
	o.storage.Store(state, verifier)
	url := o.conf.AuthCodeURL(state, oauth2.AccessTypeOnline, oauth2.S256ChallengeOption(verifier))
	fmt.Println("url ", url)
	return url
}

func GenerateCodeVerifier() (string, error) {

	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
	b := make([]byte, lengthCodeVerifier)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("ошибка при генерации случайного числа: %v", err)
		}
		b[i] = chars[num.Int64()]
	}

	return string(b), nil
}

func GenerateStateToken() (string, error) {

	// Генерируем случайные байты
	b := make([]byte, lengthState)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Кодируем в URL-safe base64 без паддинга
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (o *OAuthVkidServiceImpl) exchangeCodeForToken(code, codeVerifier, device string) (*oauth2.Token, error) {
	token, err := o.conf.Exchange(context.Background(), code, oauth2.VerifierOption(codeVerifier), oauth2.SetAuthURLParam("device_id", device))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("token", token)
	fmt.Println(err)

	return token, nil
}

type UserInfoVkid struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Email     string `json:"email,omitempty"`
	Sex       int    `json:"sex,omitempty"` // 1 - женский, 2 - мужской
	Verified  bool   `json:"verified"`
	Birthday  string `json:"birthday,omitempty"`
}

func (o *OAuthVkidServiceImpl) GetUserInfo(accessToken string) (*UserInfoVkid, error) {
	fmt.Println(accessToken)
	// Формируем URL запроса
	apiURL := "https://id.vk.com/oauth2/user_info"

	// Подготавливаем данные для тела запроса
	formData := url.Values{}
	formData.Set("access_token", accessToken)
	formData.Set("client_id", o.conf.ClientID)

	// Создаем тело запроса
	requestBody := strings.NewReader(formData.Encode())

	// Создаем HTTP-запрос с телом
	req, err := http.NewRequest("POST", apiURL, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Выводим тело ответа
	fmt.Println("Response body:", string(body))
	fmt.Println("Response Status:", resp.Status)

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		var apiError struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if err := json.Unmarshal(body, &apiError); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("vk api error: %s - %s", apiError.Error, apiError.ErrorDescription)
	}

	// Парсим ответ
	var response struct {
		User UserInfoVkid `json:"user"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.User, nil
}
