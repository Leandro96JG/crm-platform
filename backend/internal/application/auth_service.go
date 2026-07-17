package application

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type authService struct {
	userRepo  domain.UserRepository
	secretKey string
}

type AuthService interface {
	Login(ctx context.Context, username string, password string) (string, *domain.User, error)
	Logout(ctx context.Context, userID string, sessionToken string) error
	VerifyToken(tokenString string) (map[string]string, error)
	VerifyUserSession(ctx context.Context, userID string, sessionToken string) error
}

type Claims struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	SessionToken string `json:"session_token"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo domain.UserRepository, secretKey string) AuthService {
	return &authService{
		userRepo:  userRepo,
		secretKey: secretKey,
	}
}

func (s *authService) Login(ctx context.Context, login string, password string) (string, *domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, login)
	if err != nil {
		return "", nil, err
	}

	if user == nil {
		user, err = s.userRepo.FindByUsername(ctx, login)
		if err != nil {
			return "", nil, err
		}
	}

	if user == nil {
		return "", nil, domain.NewUnauthorizedError("invalid credentials")
	}

	sessionToken := uuid.New().String()
	claims := &Claims{
		UserID:       user.UserID,
		Username:     user.Username,
		SessionToken: sessionToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

func (s *authService) Logout(ctx context.Context, userID string, sessionToken string) error {
	return nil
}

func (s *authService) VerifyToken(tokenString string) (map[string]string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, domain.NewUnauthorizedError("invalid token")
	}

	return map[string]string{
		"user_id":       claims.UserID,
		"session_token": claims.SessionToken,
	}, nil
}

func (s *authService) VerifyUserSession(ctx context.Context, userID string, sessionToken string) error {
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
