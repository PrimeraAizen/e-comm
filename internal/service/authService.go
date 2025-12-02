package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/PrimeraAizen/e-comm/config"
	"github.com/PrimeraAizen/e-comm/internal/domain"
	"github.com/PrimeraAizen/e-comm/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, req *domain.User) (*domain.Token, error)
	Login(ctx context.Context, req *domain.User) (*domain.Token, error)
	ValidateToken(tokenString string) (*domain.TokenClaims, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error)
}

type authService struct {
	userRepo             repository.UserRepository
	jwtSecret            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) (AuthService, error) {
	accessDuration, err := time.ParseDuration(cfg.JWT.AccessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("parse access token duration: %w", err)
	}

	refreshDuration, err := time.ParseDuration(cfg.JWT.RefreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token duration: %w", err)
	}

	return &authService{
		userRepo:             userRepo,
		jwtSecret:            cfg.JWT.Secret,
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
	}, nil
}

func (s *authService) Register(ctx context.Context, user *domain.User) (*domain.Token, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil && err != domain.ErrNotFound {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, domain.ErrAlreadyExists
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Generate tokens
	return s.generateAuthTokens(user)
}

func (s *authService) Login(ctx context.Context, user *domain.User) (*domain.Token, error) {
	// Get user by email
	dbUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	fmt.Println("DB user:", dbUser)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Check user status
	if dbUser.Status != "active" {
		return nil, domain.ErrUserInactive
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, dbUser.ID); err != nil {
		// Log error but don't fail the login
		fmt.Printf("failed to update last login: %v\n", err)
	}

	// Generate tokens
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
		return nil, fmt.Errorf("parse token: %w", err)
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

	// Get user
	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrInvalidToken
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	// Check user status
	if user.Status != "active" {
		return nil, domain.ErrUserInactive
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

	// Remove password hash from response
	user.Password = ""

	return &domain.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTokenDuration.Seconds()),
		User:         user,
	}, nil
}

func (s *authService) generateToken(user *domain.User, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": strconv.Itoa(user.ID),
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
