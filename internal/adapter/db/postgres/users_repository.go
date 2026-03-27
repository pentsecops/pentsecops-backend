package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/adapter/db/postgres/sqlc"
	"github.com/pentsecops/backend/internal/core/domain"
)

// UsersRepository implements the UsersRepository interface
type UsersRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewUsersRepository creates a new UsersRepository
func NewUsersRepository(db *sql.DB) domain.UsersRepository {
	return &UsersRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// CreateUser creates a new user
func (r *UsersRepository) CreateUser(ctx context.Context, params *domain.CreateUserParams) (*domain.User, error) {
	userID, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	sqlcParams := sqlc.CreateUserParams{
		ID:           userID,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		FullName:     params.FullName,
		Role:         params.Role,
		IsActive:     sql.NullBool{Bool: params.IsActive, Valid: true},
		CreatedAt:    sql.NullTime{Time: params.CreatedAt, Valid: true},
		UpdatedAt:    sql.NullTime{Time: params.UpdatedAt, Valid: true},
	}

	sqlcUser, err := r.queries.CreateUser(ctx, sqlcParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &domain.User{
		ID:           sqlcUser.ID,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		FullName:     sqlcUser.FullName,
		Role:         sqlcUser.Role,
		IsActive:     sqlcUser.IsActive.Bool,
		LastLogin:    timeFromNullTime(sqlcUser.LastLogin),
		CreatedAt:    sqlcUser.CreatedAt.Time,
		UpdatedAt:    sqlcUser.UpdatedAt.Time,
	}, nil
}

// GetUserByID retrieves a user by ID
func (r *UsersRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	sqlcUser, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &domain.User{
		ID:           sqlcUser.ID,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		FullName:     sqlcUser.FullName,
		Role:         sqlcUser.Role,
		IsActive:     sqlcUser.IsActive.Bool,
		LastLogin:    timeFromNullTime(sqlcUser.LastLogin),
		CreatedAt:    sqlcUser.CreatedAt.Time,
		UpdatedAt:    sqlcUser.UpdatedAt.Time,
	}, nil
}

// GetUserByEmail retrieves a user by email
func (r *UsersRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	sqlcUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &domain.User{
		ID:           sqlcUser.ID,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		FullName:     sqlcUser.FullName,
		Role:         sqlcUser.Role,
		IsActive:     sqlcUser.IsActive.Bool,
		LastLogin:    timeFromNullTime(sqlcUser.LastLogin),
		CreatedAt:    sqlcUser.CreatedAt.Time,
		UpdatedAt:    sqlcUser.UpdatedAt.Time,
	}, nil
}

// DeleteUser deletes a user by ID
func (r *UsersRepository) DeleteUser(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = r.queries.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// UpdateUserLastLogin updates the last login timestamp for a user
func (r *UsersRepository) UpdateUserLastLogin(ctx context.Context, id string, lastLogin time.Time) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	params := sqlc.UpdateUserLastLoginParams{
		LastLogin: sql.NullTime{Time: lastLogin, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:        userID,
	}

	err = r.queries.UpdateUserLastLogin(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// CheckEmailExists checks if an email already exists
func (r *UsersRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := r.queries.CheckEmailExists(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// ListUsers retrieves a paginated list of users with project counts
func (r *UsersRepository) ListUsers(ctx context.Context, limit, offset int) ([]*domain.UserWithProjectCount, error) {
	params := sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	sqlcUsers, err := r.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*domain.UserWithProjectCount, 0, len(sqlcUsers))
	for _, u := range sqlcUsers {
		users = append(users, &domain.UserWithProjectCount{
			ID:           u.ID.String(),
			Email:        u.Email,
			FullName:     u.FullName,
			Role:         u.Role,
			IsActive:     u.IsActive.Bool,
			LastLogin:    timeFromNullTime(u.LastLogin),
			CreatedAt:    u.CreatedAt.Time,
			UpdatedAt:    u.UpdatedAt.Time,
			ProjectCount: u.ProjectCount,
		})
	}

	return users, nil
}

// CountUsers returns the total count of users
func (r *UsersRepository) CountUsers(ctx context.Context) (int64, error) {
	count, err := r.queries.CountUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// GetUserStats retrieves user statistics
func (r *UsersRepository) GetUserStats(ctx context.Context) (*domain.UserStats, error) {
	stats, err := r.queries.GetUserStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return &domain.UserStats{
		ActivePentesters:   stats.ActivePentesters,
		ActiveStakeholders: stats.ActiveStakeholders,
		InactiveUsers:      stats.InactiveUsers,
		TotalUsers:         stats.TotalUsers,
	}, nil
}

// CountUsersByRole counts users by role and active status
func (r *UsersRepository) CountUsersByRole(ctx context.Context, role string, isActive bool) (int64, error) {
	params := sqlc.CountUsersByRoleParams{
		Role:     role,
		IsActive: sql.NullBool{Bool: isActive, Valid: true},
	}

	count, err := r.queries.CountUsersByRole(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by role: %w", err)
	}

	return count, nil
}

// CountInactiveUsers counts inactive users
func (r *UsersRepository) CountInactiveUsers(ctx context.Context) (int64, error) {
	count, err := r.queries.CountInactiveUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count inactive users: %w", err)
	}

	return count, nil
}

// ListAllUsersForExport retrieves all users for CSV export
func (r *UsersRepository) ListAllUsersForExport(ctx context.Context) ([]*domain.UserWithProjectCount, error) {
	sqlcUsers, err := r.queries.ListAllUsersForExport(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users for export: %w", err)
	}

	users := make([]*domain.UserWithProjectCount, 0, len(sqlcUsers))
	for _, u := range sqlcUsers {
		users = append(users, &domain.UserWithProjectCount{
			ID:           u.ID.String(),
			Email:        u.Email,
			FullName:     u.FullName,
			Role:         u.Role,
			IsActive:     u.IsActive.Bool,
			LastLogin:    timeFromNullTime(u.LastLogin),
			CreatedAt:    u.CreatedAt.Time,
			UpdatedAt:    time.Time{}, // Export doesn't have UpdatedAt
			ProjectCount: u.ProjectCount,
		})
	}

	return users, nil
}

// UpdateUser updates an existing user
func (r *UsersRepository) UpdateUser(ctx context.Context, id string, params *domain.UpdateUserParams) (*domain.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Build the update query dynamically
	query := "UPDATE users SET "
	args := []interface{}{}
	argIndex := 1
	updateFields := []string{}

	// Add fields to update
	updateFields = append(updateFields, fmt.Sprintf("full_name = $%d", argIndex))
	args = append(args, params.FullName)
	argIndex++

	updateFields = append(updateFields, fmt.Sprintf("email = $%d", argIndex))
	args = append(args, params.Email)
	argIndex++

	updateFields = append(updateFields, fmt.Sprintf("role = $%d", argIndex))
	args = append(args, params.Role)
	argIndex++

	if params.IsActive != nil {
		updateFields = append(updateFields, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	query += strings.Join(updateFields, ", ")
	query += fmt.Sprintf(" WHERE id = $%d RETURNING *", argIndex)
	args = append(args, userID)

	// Execute the update
	var sqlcUser sqlc.User
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&sqlcUser.ID,
		&sqlcUser.Email,
		&sqlcUser.PasswordHash,
		&sqlcUser.FullName,
		&sqlcUser.Role,
		&sqlcUser.IsActive,
		&sqlcUser.ForcePasswordChange,
		&sqlcUser.FailedLoginAttempts,
		&sqlcUser.LastFailedLogin,
		&sqlcUser.AccountLockedUntil,
		&sqlcUser.LastLogin,
		&sqlcUser.CreatedAt,
		&sqlcUser.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &domain.User{
		ID:           sqlcUser.ID,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		FullName:     sqlcUser.FullName,
		Role:         sqlcUser.Role,
		IsActive:     sqlcUser.IsActive.Bool,
		LastLogin:    timeFromNullTime(sqlcUser.LastLogin),
		CreatedAt:    sqlcUser.CreatedAt.Time,
		UpdatedAt:    sqlcUser.UpdatedAt.Time,
	}, nil
}

// GetUsersByRole retrieves users by role
func (r *UsersRepository) GetUsersByRole(ctx context.Context, role string) ([]domain.User, error) {
	query := `SELECT id, email, password_hash, full_name, role, is_active, last_login, created_at, updated_at FROM users WHERE role = $1 AND is_active = true`
	rows, err := r.db.QueryContext(ctx, query, role)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		var lastLogin sql.NullTime
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&user.Role,
			&user.IsActive,
			&lastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		user.LastLogin = timeFromNullTime(lastLogin)
		users = append(users, user)
	}

	return users, nil
}

// Helper function to convert sql.NullTime to *time.Time
func timeFromNullTime(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}
