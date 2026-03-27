package dto

import "time"

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	TokenType             string    `json:"token_type"` // "Bearer"
	ExpiresIn             int       `json:"expires_in"` // seconds
	UserID                string    `json:"user_id"`
	Email                 string    `json:"email"`
	FullName              string    `json:"full_name"`
	Role                  string    `json:"role"`
	ForcePasswordChange   bool      `json:"force_password_change"`
	LastLogin             *time.Time `json:"last_login,omitempty"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents a token refresh response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"` // "Bearer"
	ExpiresIn    int    `json:"expires_in"` // seconds
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"required,min=8,nefield=OldPassword"`
}

// ChangePasswordResponse represents a password change response
type ChangePasswordResponse struct {
	Message string `json:"message"`
}

// LogoutRequest represents a logout request
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutResponse represents a logout response
type LogoutResponse struct {
	Message string `json:"message"`
}

