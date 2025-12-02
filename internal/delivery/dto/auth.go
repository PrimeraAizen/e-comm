package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/PrimeraAizen/e-comm/internal/domain"
)

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required,min=8"`
}

func (r *RegisterRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return domain.ErrValidation
	}
	return nil
}

func (r *RegisterRequest) ToDomain() (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	return &domain.User{
		Password: string(hashedPassword),
		Email:    r.Email,
	}, nil
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (l *LoginRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(l); err != nil {
		return domain.ErrValidation
	}
	return nil
}

func (l *LoginRequest) ToDomain() *domain.User {
	return &domain.User{
		Email:    l.Email,
		Password: l.Password,
	}
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ProfileResponse represents user profile information
type ProfileResponse struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	MiddleName  string `json:"middle_name,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	Country     string `json:"country,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Email       string `json:"email"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	FirstName   *string `json:"first_name" validate:"omitempty,min=1,max=100"`
	LastName    *string `json:"last_name" validate:"omitempty,min=1,max=100"`
	MiddleName  *string `json:"middle_name" validate:"omitempty,max=100"`
	DateOfBirth *string `json:"date_of_birth" validate:"omitempty"`
	Gender      *string `json:"gender" validate:"omitempty,max=20"`
	Phone       *string `json:"phone" validate:"omitempty,max=50"`
	Address     *string `json:"address" validate:"omitempty"`
	City        *string `json:"city" validate:"omitempty,max=100"`
	Country     *string `json:"country" validate:"omitempty,max=100"`
	PostalCode  *string `json:"postal_code" validate:"omitempty,max=20"`
}

func (u *UpdateProfileRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(u); err != nil {
		return domain.ErrValidation
	}
	return nil
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

func (c *ChangePasswordRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return domain.ErrValidation
	}
	if c.NewPassword != c.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}
	return nil
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
}
