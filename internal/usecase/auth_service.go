package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"subscription-manager-api/internal/domain"
	"subscription-manager-api/internal/infrastructure/persistence"
)

// AuthService は認証ユースケース
type AuthService struct {
	userRepo *persistence.UserRepository
	jwtSecret string
}

// NewAuthService は新しい認証ユースケースを作成
func NewAuthService(userRepo *persistence.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// SignUp はユーザー登録
func (s *AuthService) SignUp(ctx context.Context, email, password, name string) (*domain.User, string, error) {
	// 既存ユーザーをチェック
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", errors.New("user already exists")
	}

	// ユーザー作成
	user := &domain.User{
		Email: email,
		Name:  name,
	}
	if err := user.SetPassword(password); err != nil {
		return nil, "", err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	// JWTトークン生成
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login はログイン
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}

	if !user.CheckPassword(password) {
		return "", errors.New("invalid email or password")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateToken はJWTトークンを生成
func (s *AuthService) generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7日間有効
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

