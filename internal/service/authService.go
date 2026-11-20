package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/PrimeraAizen/e-comm/config"
	"github.com/PrimeraAizen/e-comm/internal/domain"
	"github.com/PrimeraAizen/e-comm/internal/repository"
	"github.com/PrimeraAizen/e-comm/pkg/logger"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*domain.Token, error)
	Login(ctx context.Context, email, password string) (*domain.Token, error)
	ValidateToken(tokenString string) (*domain.TokenClaims, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error)
	ChangePassword(ctx context.Context, userID string, currentPassword, newPassword string) error
}

type authService struct {
	userRepo             repository.UserRepository
	jwtSecret            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	logger               *logger.Logger
}

func NewAuthService(userRepo repository.UserRepository, log *logger.Logger, cfg *config.JWT) (AuthService, error) {
	accessDuration, err := time.ParseDuration(cfg.AccessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("parse access token duration: %w", err)
	}

	refreshDuration, err := time.ParseDuration(cfg.RefreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token duration: %w", err)
	}

	return &authService{
		userRepo:             userRepo,
		jwtSecret:            cfg.Secret,
		logger:               log,
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
	}, nil
}

func (s *authService) Register(ctx context.Context, email, password string) (*domain.Token, error) {
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && err != domain.ErrNotFound {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, domain.ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return s.generateAuthTokens(user)
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.Token, error) {
	dbUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := s.userRepo.UpdateLastLogin(ctx, dbUser.ID); err != nil {
		s.logger.WithComponent("auth").WithError(err).Error("failed to update last login")
	}

	return s.generateAuthTokens(dbUser)
}

func (s *authService) ValidateToken(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	return &domain.TokenClaims{
		UserID: userID,
		Email:  email,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	// Validate refresh token
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidToken
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	// Generate new tokens
	return s.generateAuthTokens(user)
}

func (s *authService) generateAuthTokens(user *domain.User) (*domain.Token, error) {
	// Generate access token
	accessToken, err := s.generateToken(user, s.accessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateToken(user, s.refreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &domain.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTokenDuration.Seconds()),
	}, nil
}

func (s *authService) generateToken(user *domain.User, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return tokenString, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID string, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return domain.ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := s.userRepo.UpdatePasswordHash(ctx, user.ID, string(hashedPassword)); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}
