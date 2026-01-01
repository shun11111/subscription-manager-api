package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"subscription-manager-api/internal/usecase"
)

// AuthHandler は認証ハンドラ
type AuthHandler struct {
	authService *usecase.AuthService
}

// NewAuthHandler は新しい認証ハンドラを作成
func NewAuthHandler(authService *usecase.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SignUpRequest はサインアップリクエスト
type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// SignUpResponse はサインアップレスポンス
type SignUpResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

// SignUp はユーザー登録
func (h *AuthHandler) SignUp(c echo.Context) error {
	var req SignUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	user, token, err := h.authService.SignUp(c.Request().Context(), req.Email, req.Password, req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, SignUpResponse{
		UserID: user.ID.String(),
		Token:  token,
	})
}

// LoginRequest はログインリクエスト
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse はログインレスポンス
type LoginResponse struct {
	Token string `json:"token"`
}

// Login はログイン
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	token, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, LoginResponse{Token: token})
}
