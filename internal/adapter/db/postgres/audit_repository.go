package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/core/domain"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// LogActivity creates a new audit log entry
func (r *AuditRepository) LogActivity(ctx context.Context, activityLog *domain.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (
			id, user_id, user_email, user_role, action, entity_type, entity_id,
			old_values, new_values, ip_address, user_agent, endpoint, method,
			status_code, error_message, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)`

	activityLog.ID = uuid.New()
	activityLog.CreatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		activityLog.ID,
		activityLog.UserID,
		activityLog.UserEmail,
		activityLog.UserRole,
		activityLog.Action,
		activityLog.EntityType,
		activityLog.EntityID,
		activityLog.OldValues,
		activityLog.NewValues,
		activityLog.IPAddress,
		activityLog.UserAgent,
		activityLog.Endpoint,
		activityLog.Method,
		activityLog.StatusCode,
		activityLog.ErrorMessage,
		activityLog.CreatedAt,
	)

	if err != nil {
		log.Printf("Failed to log activity: %v", err)
		return fmt.Errorf("failed to log activity: %w", err)
	}

	return nil
}

// GetActivityLogs retrieves activity logs with filters and pagination
func (r *AuditRepository) GetActivityLogs(ctx context.Context, filter domain.AuditLogFilter) ([]domain.ActivityLog, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if filter.UserID != nil {
		whereClause += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.Action != "" {
		whereClause += fmt.Sprintf(" AND action ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Action+"%")
		argIndex++
	}

	if filter.EntityType != "" {
		whereClause += fmt.Sprintf(" AND entity_type = $%d", argIndex)
		args = append(args, filter.EntityType)
		argIndex++
	}

	if filter.EntityID != nil {
		whereClause += fmt.Sprintf(" AND entity_id = $%d", argIndex)
		args = append(args, *filter.EntityID)
		argIndex++
	}

	if filter.StartDate != nil {
		whereClause += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		whereClause += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM activity_logs %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("Failed to count activity logs: %v", err)
		return nil, 0, fmt.Errorf("failed to count activity logs: %w", err)
	}

	// Get paginated results
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PerPage <= 0 {
		filter.PerPage = 20
	}

	offset := (filter.Page - 1) * filter.PerPage
	query := fmt.Sprintf(`
		SELECT id, user_id, user_email, user_role, action, entity_type, entity_id,
			   old_values, new_values, ip_address, user_agent, endpoint, method,
			   status_code, error_message, created_at
		FROM activity_logs %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to query activity logs: %v", err)
		return nil, 0, fmt.Errorf("failed to query activity logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.ActivityLog
	for rows.Next() {
		var activityLog domain.ActivityLog
		err := rows.Scan(
			&activityLog.ID,
			&activityLog.UserID,
			&activityLog.UserEmail,
			&activityLog.UserRole,
			&activityLog.Action,
			&activityLog.EntityType,
			&activityLog.EntityID,
			&activityLog.OldValues,
			&activityLog.NewValues,
			&activityLog.IPAddress,
			&activityLog.UserAgent,
			&activityLog.Endpoint,
			&activityLog.Method,
			&activityLog.StatusCode,
			&activityLog.ErrorMessage,
			&activityLog.CreatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan activity log: %v", err)
			continue
		}
		logs = append(logs, activityLog)
	}

	return logs, total, nil
}

// LogDatabaseOperation logs database CRUD operations
func (r *AuditRepository) LogDatabaseOperation(ctx context.Context, userID *uuid.UUID, userEmail, userRole, action, entityType string, entityID *uuid.UUID, oldValues, newValues interface{}) {
	var oldJSON, newJSON *string

	if oldValues != nil {
		if data, err := json.Marshal(oldValues); err == nil {
			oldStr := string(data)
			oldJSON = &oldStr
		}
	}

	if newValues != nil {
		if data, err := json.Marshal(newValues); err == nil {
			newStr := string(data)
			newJSON = &newStr
		}
	}

	activityLog := &domain.ActivityLog{
		UserID:     userID,
		UserEmail:  userEmail,
		UserRole:   userRole,
		Action:     action,
		EntityType: &entityType,
		EntityID:   entityID,
		OldValues:  oldJSON,
		NewValues:  newJSON,
		IPAddress:  "database",
		UserAgent:  "system",
		Endpoint:   "database_operation",
		Method:     "DB",
		StatusCode: 200,
	}

	// Log asynchronously to avoid blocking database operations
	go func() {
		if err := r.LogActivity(context.Background(), activityLog); err != nil {
			log.Printf("Failed to log database operation: %v", err)
		}
	}()
}