package google

// courtesy of Kazutaka Yoshinaga
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendarToken struct {
	ID           string     `json:"id"`
	UserID       uint64     `json:"user_id"`
	RefreshToken string     `json:"refresh_token"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// GoogleCalendarTokenValidator performs token operations for Google Calendar API
type GoogleCalendarTokenValidator struct {
	Service     *calendar.Service
	PackageName string
	Config      *oauth2.Config
}

// NewGoogleCalendarTokenValidator init token validator for API
func NewGoogleCalendarTokenValidator() (*GoogleCalendarTokenValidator, error) {
	// 環境変数からOAuth 2.0クライアント情報を取得
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	// 必須情報が不足している場合はエラー
	if clientID == "" || clientSecret == "" || redirectURI == "" {
		return nil, errors.New("missing required environment variables: GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_REDIRECT_URI")
	}

	// OAuth2 Config の初期化
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{calendar.CalendarScope},
		Endpoint:     oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	// Google Calendar Service の初期化
	ctx := context.Background()
	service, err := calendar.NewService(ctx, option.WithHTTPClient(config.Client(ctx, nil)))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Google Calendar service: %w", err)
	}

	// Validator の初期化
	return &GoogleCalendarTokenValidator{
		Service: service,
		Config:  config,
	}, nil
}

// ExchangeAuthCode, exchanges auth code for access + refresh token
func (v *GoogleCalendarTokenValidator) ExchangeAuthCode(ctx context.Context, authCode string) (*oauth2.Token, error) {
	if v.Config == nil {
		return nil, errors.New("OAuth2 config is not initialized")
	}

	token, err := v.Config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange auth code: %w", err)
	}

	return token, nil
}

const TokenEndpoint = "https://oauth2.googleapis.com/token"

// AccessTokenResponse
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GoogleCalendarService Google Calendar API
type GoogleCalendarService struct {
	httpClient  *http.Client
	accessToken string
	expiry      time.Time
}

// NewGoogleCalendarService init
func NewGoogleCalendarService() *GoogleCalendarService {
	return &GoogleCalendarService{
		httpClient: &http.Client{},
	}
}

// GetAccessToken, Use Refresh token to get Acess token
func (s *GoogleCalendarService) GetAccessToken(refreshToken string) (string, error) {
	// Reuse token if valid
	if s.accessToken != "" && time.Now().Before(s.expiry) {
		return s.accessToken, nil
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("missing environment variables (client_id or client_secret)")
	}

	payload := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", TokenEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("request creation error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("access token retrieval failed: status code %d, response %s", resp.StatusCode, string(responseBody))
	}

	var tokenResp AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("response decode error: %v", err)
	}

	s.accessToken = tokenResp.AccessToken
	s.expiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return s.accessToken, nil
}

// GetEvents, Get events within a time range
func (s *GoogleCalendarService) GetEvents(ctx context.Context, accessToken string, timeMin, timeMax time.Time) ([]*calendar.Event, error) {
	srv, err := calendar.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	})))
	if err != nil {
		return nil, fmt.Errorf("google calendar api client creation error: %v", err)
	}

	events, err := srv.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(timeMin.Format(time.RFC3339)).
		TimeMax(timeMax.Format(time.RFC3339)).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf("google Calendar api event retrieval error: %v", err)
	}

	return events.Items, nil
}