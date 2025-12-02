package domain

import (
	"time"
)

type User struct {
	ID          int        `json:"id" bson:"_id"`
	Email       string     `json:"email" bson:"email"`
	Password     string     `json:"-" bson:"password_hash"`
	Status      string     `json:"status" bson:"status"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" bson:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" bson:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	User         *User  `json:"user"`
}
