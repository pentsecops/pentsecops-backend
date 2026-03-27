package usecases

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/pkg/auth"
	"golang.org/x/crypto/argon2"
)

const (
	// Account locking settings
	MaxFailedLoginAttempts = 5
	AccountLockDuration    = 15 * time.Minute

	// Password hashing parameters (Argon2id)
	Argon2Time    = 1
	Argon2Memory  = 64 * 1024
	Argon2Threads = 4
	Argon2KeyLen  = 32
	Argon2SaltLen = 16
)

// AuthUseCase implements authentication business logic
type AuthUseCase struct {
	repo          domain.AuthRepository
	pasetoService *auth.PasetoService
	validator     *validator.Validate
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewAuthUseCase creates a new AuthUseCase
func NewAuthUseCase(
	repo domain.AuthRepository,
	pasetoService *auth.PasetoService,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) domain.AuthUseCase {
	return &AuthUseCase{
		repo:                 repo,
		pasetoService:        pasetoService,
		validator:            validator.New(),
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// Login authenticates a user and returns tokens
func (uc *AuthUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Validate request
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Get user by email
	user, err := uc.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if account is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	// Check if account is locked
	if user.AccountLockedUntil != nil && time.Now().Before(*user.AccountLockedUntil) {
		remainingTime := time.Until(*user.AccountLockedUntil).Round(time.Minute)
		return nil, fmt.Errorf("account is locked. Try again in %v", remainingTime)
	}

	// Unlock account if lock period has expired
	if user.AccountLockedUntil != nil && time.Now().After(*user.AccountLockedUntil) {
		if err := uc.repo.UnlockAccount(ctx, user.ID.String()); err != nil {
			// Log error but continue with login
			fmt.Printf("Failed to unlock account: %v\n", err)
		}
		user.FailedLoginAttempts = 0
		user.AccountLockedUntil = nil
	}

	// Verify password
	if !uc.verifyPassword(req.Password, user.PasswordHash) {
		// Increment failed login attempts
		if err := uc.repo.IncrementFailedLoginAttempts(ctx, user.ID.String()); err != nil {
			fmt.Printf("Failed to increment login attempts: %v\n", err)
		}

		// Lock account if max attempts reached
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= MaxFailedLoginAttempts {
			lockUntil := time.Now().Add(AccountLockDuration)
			if err := uc.repo.LockAccount(ctx, user.ID.String(), lockUntil); err != nil {
				fmt.Printf("Failed to lock account: %v\n", err)
			}
			return nil, fmt.Errorf("account locked due to too many failed login attempts. Try again in %v", AccountLockDuration)
		}

		return nil, fmt.Errorf("invalid email or password")
	}

	// Reset failed login attempts on successful login
	if user.FailedLoginAttempts > 0 {
		if err := uc.repo.ResetFailedLoginAttempts(ctx, user.ID.String()); err != nil {
			fmt.Printf("Failed to reset login attempts: %v\n", err)
		}
	}

	// Generate access token
	accessToken, err := uc.pasetoService.GenerateAccessToken(
		user.ID.String(),
		user.Email,
		user.Role,
		uc.accessTokenDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := uc.pasetoService.GenerateRefreshToken(
		user.ID.String(),
		user.Email,
		user.Role,
		uc.refreshTokenDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token hash in database
	refreshTokenHash := uc.hashToken(refreshToken)
	expiresAt := time.Now().Add(uc.refreshTokenDuration)
	if err := uc.repo.CreateRefreshToken(ctx, user.ID.String(), refreshTokenHash, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Update last login time
	if err := uc.repo.UpdateUserLastLogin(ctx, user.ID.String(), time.Now()); err != nil {
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	// Clean up expired tokens (async, don't wait)
	go func() {
		if err := uc.repo.DeleteExpiredRefreshTokens(context.Background()); err != nil {
			fmt.Printf("Failed to delete expired tokens: %v\n", err)
		}
	}()

	return &dto.LoginResponse{
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
		TokenType:           "Bearer",
		ExpiresIn:           int(uc.accessTokenDuration.Seconds()),
		UserID:              user.ID.String(),
		Email:               user.Email,
		FullName:            user.FullName,
		Role:                user.Role,
		ForcePasswordChange: user.ForcePasswordChange,
		LastLogin:           user.LastLogin,
	}, nil
}

// RefreshToken generates new tokens using a refresh token
func (uc *AuthUseCase) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	// Validate request
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Validate refresh token
	claims, err := uc.pasetoService.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check token type
	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Check if refresh token exists in database
	refreshTokenHash := uc.hashToken(req.RefreshToken)
	storedToken, err := uc.repo.GetRefreshToken(ctx, refreshTokenHash)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found or expired")
	}

	// Verify user ID matches
	if storedToken.UserID.String() != claims.UserID {
		return nil, fmt.Errorf("token user mismatch")
	}

	// Get user to verify account status
	user, err := uc.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	// Delete old refresh token
	if err := uc.repo.DeleteRefreshToken(ctx, refreshTokenHash); err != nil {
		fmt.Printf("Failed to delete old refresh token: %v\n", err)
	}

	// Generate new access token
	newAccessToken, err := uc.pasetoService.GenerateAccessToken(
		user.ID.String(),
		user.Email,
		user.Role,
		uc.accessTokenDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := uc.pasetoService.GenerateRefreshToken(
		user.ID.String(),
		user.Email,
		user.Role,
		uc.refreshTokenDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store new refresh token hash
	newRefreshTokenHash := uc.hashToken(newRefreshToken)
	expiresAt := time.Now().Add(uc.refreshTokenDuration)
	if err := uc.repo.CreateRefreshToken(ctx, user.ID.String(), newRefreshTokenHash, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(uc.accessTokenDuration.Seconds()),
	}, nil
}

// ChangePassword changes a user's password
func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) (*dto.ChangePasswordResponse, error) {
	// Validate request
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Get user
	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Verify old password
	if !uc.verifyPassword(req.OldPassword, user.PasswordHash) {
		return nil, fmt.Errorf("invalid old password")
	}

	// Hash new password
	newPasswordHash := uc.hashPassword(req.NewPassword)

	// Update password
	if err := uc.repo.UpdateUserPassword(ctx, userID, newPasswordHash); err != nil {
		return nil, fmt.Errorf("failed to update password: %w", err)
	}

	// Disable force password change flag
	if user.ForcePasswordChange {
		if err := uc.repo.UpdateUserForcePasswordChange(ctx, userID, false); err != nil {
			fmt.Printf("Failed to update force password change: %v\n", err)
		}
	}

	// Invalidate all refresh tokens for security
	if err := uc.repo.DeleteAllUserRefreshTokens(ctx, userID); err != nil {
		fmt.Printf("Failed to delete user refresh tokens: %v\n", err)
	}

	return &dto.ChangePasswordResponse{
		Message: "Password changed successfully. Please login again.",
	}, nil
}

// Logout invalidates a refresh token
func (uc *AuthUseCase) Logout(ctx context.Context, req *dto.LogoutRequest) (*dto.LogoutResponse, error) {
	// Validate request
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Hash and delete refresh token
	refreshTokenHash := uc.hashToken(req.RefreshToken)
	if err := uc.repo.DeleteRefreshToken(ctx, refreshTokenHash); err != nil {
		// Don't fail if token doesn't exist
		fmt.Printf("Failed to delete refresh token: %v\n", err)
	}

	return &dto.LogoutResponse{
		Message: "Logged out successfully",
	}, nil
}

// hashPassword hashes a password using Argon2id
func (uc *AuthUseCase) hashPassword(password string) string {
	// In production, you should generate a random salt for each password
	// For now, using a simple approach
	salt := []byte("pentsecops-salt-2025") // TODO: Use random salt per user
	hash := argon2.IDKey([]byte(password), salt, Argon2Time, Argon2Memory, Argon2Threads, Argon2KeyLen)
	return hex.EncodeToString(hash)
}

// verifyPassword verifies a password against a hash
func (uc *AuthUseCase) verifyPassword(password, hash string) bool {
	passwordHash := uc.hashPassword(password)
	return passwordHash == hash
}

// hashToken creates a SHA-256 hash of a token for storage
func (uc *AuthUseCase) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

