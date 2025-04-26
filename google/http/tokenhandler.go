package http

// courtesy of Kazutaka Yoshinaga
import (
	"context"
	"net/http"
	"time"

	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/nahnhh/points-hunter/google"
)


type GoogleCalendarTokenRequest struct {
	UserID   uint64 `json:"userId" validate:"required"`
	AuthCode string `json:"authCode" validate:"required"`
}


type GoogleCalendarTokenService interface {
	SaveToken(ctx context.Context, token google.GoogleCalendarToken) error
	GetTokenByUserID(ctx context.Context, userID uint64) (google.GoogleCalendarToken, error)
}


type GoogleCalendarTokenHandler struct {
	TokenValidator *google.GoogleCalendarTokenValidator
	TokenService   GoogleCalendarTokenService
}


func NewGoogleCalendarTokenHandler(
	tokenValidator *google.GoogleCalendarTokenValidator,
	tokenService GoogleCalendarTokenService,
) *GoogleCalendarTokenHandler {
	return &GoogleCalendarTokenHandler{
		TokenValidator: tokenValidator,
		TokenService:   tokenService,
	}
}

// 認証コードをリフレッシュトークンに交換しデータベースに保存
func (h *GoogleCalendarTokenHandler) HandleExchangeAuthCode(c echo.Context) error {
	var req GoogleCalendarTokenRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind request", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// リクエスト内容のデバッグログ
	slog.Info("Received request", "userID", req.UserID, "authCode", req.AuthCode)


	// リクエストのバリデーション
	if req.UserID == 0 || req.AuthCode == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "userId and authCode are required"})
	}

	ctx := context.Background()

	// AuthCodeを使用してアクセストークンとリフレッシュトークンを取得
	token, err := h.TokenValidator.ExchangeAuthCode(ctx, req.AuthCode)
	if err != nil {
		slog.Error("Failed to exchange auth code", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to exchange auth code"})
	}

	// トークン内容のデバッグログ
	slog.Info("Token exchanged successfully", "accessToken", token.AccessToken, "refreshToken", token.RefreshToken, "expiry", token.Expiry)


	// トークンデータを構築
	googleToken := google.GoogleCalendarToken{
		UserID:       req.UserID,
		RefreshToken: token.RefreshToken,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// トークンをデータベースに保存
	if err := h.TokenService.SaveToken(ctx, googleToken); err != nil {
		slog.Error("Failed to save Google Calendar token", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Token saved successfully"})
}