package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ContextKey はコンテキストキーの型
type ContextKey string

const (
	UserIDKey ContextKey = "user_id"
)

// AuthMiddleware はJWT認証ミドルウェア
func AuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Authorization header is required"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid authorization header format"})
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
			}

			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user_id in token"})
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user_id format"})
			}

			// コンテキストにuser_idを設定
			c.Set(string(UserIDKey), userID)
			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), UserIDKey, userID)))

			return next(c)
		}
	}
}

// GetUserID はコンテキストからuser_idを取得
func GetUserID(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get(string(UserIDKey)).(uuid.UUID)
	if !ok {
		return uuid.Nil, echo.ErrUnauthorized
	}
	return userID, nil
}

