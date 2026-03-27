package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/adapter/db/postgres/sqlc"
	"github.com/pentsecops/backend/internal/core/domain"
)

// AuthRepository implements the domain.AuthRepository interface
type AuthRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewAuthRepository creates a new AuthRepository
func NewAuthRepository(db *sql.DB) domain.AuthRepository {
	return &AuthRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// GetUserByEmail retrieves a user by email
func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.queries.GetUserForAuth(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return r.sqlcUserToDomain(row), nil
}

// GetUserByID retrieves a user by ID
func (r *AuthRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	row, err := r.queries.GetUserByID(ctx, userUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return r.sqlcUserToDomainBasic(row), nil
}

// UpdateUserPassword updates a user's password
func (r *AuthRepository) UpdateUserPassword(ctx context.Context, userID, passwordHash string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = r.queries.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		ID:           userUUID,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// UpdateUserLastLogin updates a user's last login time
func (r *AuthRepository) UpdateUserLastLogin(ctx context.Context, userID string, lastLogin time.Time) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = r.queries.UpdateLastLoginTime(ctx, sqlc.UpdateLastLoginTimeParams{
		ID: userUUID,
		LastLogin: sql.NullTime{
			Time:  lastLogin,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// UpdateUserForcePasswordChange updates the force password change flag
func (r *AuthRepository) UpdateUserForcePasswordChange(ctx context.Context, userID string, forceChange bool) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = r.queries.UpdateUserForcePasswordChange(ctx, sqlc.UpdateUserForcePasswordChangeParams{
		ID: userUUID,
		ForcePasswordChange: sql.NullBool{
			Bool:  forceChange,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update force password change: %w", err)
	}

	return nil
}

// IncrementFailedLoginAttempts increments the failed login attempts counter
func (r *AuthRepository) IncrementFailedLoginAttempts(ctx context.Context, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = r.queries.IncrementFailedLoginAttempts(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("failed to increment failed login attempts: %w", err)
	}

	return nil
}

// ResetFailedLoginAttempts resets the failed login attempts counter
func (r *AuthRepository) ResetFailedLoginAttempts(ctx context.Context, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = r.queries.ResetFailedLoginAttempts(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("failed to reset failed login attempts: %w", err)
	}

	return nil
}

// LockAccount locks a user account until the specified time
func (r *AuthRepository) LockAccount(ctx context.Context, userID string, lockUntil time.Time) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = r.queries.LockAccount(ctx, sqlc.LockAccountParams{
		ID: userUUID,
		AccountLockedUntil: sql.NullTime{
			Time:  lockUntil,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to lock account: %w", err)
	}

	return nil
}

// UnlockAccount unlocks a user account
func (r *AuthRepository) UnlockAccount(ctx context.Context, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = r.queries.UnlockAccount(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("failed to unlock account: %w", err)
	}

	return nil
}

// CreateRefreshToken creates a new refresh token
func (r *AuthRepository) CreateRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	tokenID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate token ID: %w", err)
	}

	err = r.queries.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		ID: tokenID,
		UserID: uuid.NullUUID{
			UUID:  userUUID,
			Valid: true,
		},
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	return nil
}

// GetRefreshToken retrieves a refresh token by hash
func (r *AuthRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	row, err := r.queries.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found or expired")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return &domain.RefreshToken{
		ID:        row.ID,
		UserID:    row.UserID.UUID,
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

// DeleteRefreshToken deletes a refresh token
func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	err := r.queries.DeleteRefreshToken(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}

// DeleteAllUserRefreshTokens deletes all refresh tokens for a user
func (r *AuthRepository) DeleteAllUserRefreshTokens(ctx context.Context, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = r.queries.DeleteAllUserRefreshTokens(ctx, uuid.NullUUID{
		UUID:  userUUID,
		Valid: true,
	})
	if err != nil {
		return fmt.Errorf("failed to delete user refresh tokens: %w", err)
	}

	return nil
}

// DeleteExpiredRefreshTokens deletes all expired refresh tokens
func (r *AuthRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	err := r.queries.DeleteExpiredRefreshTokens(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}

	return nil
}

// Helper function to convert sqlc User to domain User (with auth fields)
func (r *AuthRepository) sqlcUserToDomain(row sqlc.User) *domain.User {
	user := &domain.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		FullName:     row.FullName,
		Role:         row.Role,
		IsActive:     row.IsActive.Bool,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}

	if row.ForcePasswordChange.Valid {
		user.ForcePasswordChange = row.ForcePasswordChange.Bool
	}

	if row.FailedLoginAttempts.Valid {
		user.FailedLoginAttempts = int(row.FailedLoginAttempts.Int32)
	}

	if row.LastFailedLogin.Valid {
		user.LastFailedLogin = &row.LastFailedLogin.Time
	}

	if row.AccountLockedUntil.Valid {
		user.AccountLockedUntil = &row.AccountLockedUntil.Time
	}

	if row.LastLogin.Valid {
		user.LastLogin = &row.LastLogin.Time
	}

	return user
}

// Helper function to convert sqlc User to domain User (basic fields)
func (r *AuthRepository) sqlcUserToDomainBasic(row sqlc.User) *domain.User {
	user := &domain.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		FullName:     row.FullName,
		Role:         row.Role,
		IsActive:     row.IsActive.Bool,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}

	if row.ForcePasswordChange.Valid {
		user.ForcePasswordChange = row.ForcePasswordChange.Bool
	}

	if row.FailedLoginAttempts.Valid {
		user.FailedLoginAttempts = int(row.FailedLoginAttempts.Int32)
	}

	if row.LastFailedLogin.Valid {
		user.LastFailedLogin = &row.LastFailedLogin.Time
	}

	if row.AccountLockedUntil.Valid {
		user.AccountLockedUntil = &row.AccountLockedUntil.Time
	}

	if row.LastLogin.Valid {
		user.LastLogin = &row.LastLogin.Time
	}

	return user
}
