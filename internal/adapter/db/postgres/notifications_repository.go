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

// NotificationsRepository implements the NotificationsRepository interface
type NotificationsRepository struct {
	queries *sqlc.Queries
}

// NewNotificationsRepository creates a new NotificationsRepository
func NewNotificationsRepository(db *sql.DB) *NotificationsRepository {
	return &NotificationsRepository{
		queries: sqlc.New(db),
	}
}

// GetTotalNotificationsSent retrieves the total count of notifications sent by a user
func (r *NotificationsRepository) GetTotalNotificationsSent(ctx context.Context, createdBy string) (int64, error) {
	createdByUUID, err := uuid.Parse(createdBy)
	if err != nil {
		return 0, fmt.Errorf("invalid created_by UUID: %w", err)
	}

	total, err := r.queries.GetTotalNotificationsSent(ctx, createdByUUID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total notifications sent: %w", err)
	}

	return total, nil
}

// ListNotifications retrieves a paginated list of notifications
func (r *NotificationsRepository) ListNotifications(ctx context.Context, createdBy string, limit, offset int) ([]domain.Notification, error) {
	createdByUUID, err := uuid.Parse(createdBy)
	if err != nil {
		return nil, fmt.Errorf("invalid created_by UUID: %w", err)
	}

	rows, err := r.queries.ListNotifications(ctx, sqlc.ListNotificationsParams{
		CreatedBy: createdByUUID,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	notifications := make([]domain.Notification, 0, len(rows))
	for _, row := range rows {
		var recipientID *uuid.UUID
		if row.RecipientID.Valid {
			recipientID = &row.RecipientID.UUID
		}

		var notifType *string
		if row.Type.Valid {
			notifType = &row.Type.String
		}

		status := "sent"
		if row.Status.Valid {
			status = row.Status.String
		}

		isRead := false
		if row.IsRead.Valid {
			isRead = row.IsRead.Bool
		}

		createdAt := row.CreatedAt.Time
		if !row.CreatedAt.Valid {
			createdAt = time.Now()
		}

		notifications = append(notifications, domain.Notification{
			ID:          row.ID,
			Title:       row.Title,
			Message:     row.Message,
			SentTo:      row.SentTo,
			RecipientID: recipientID,
			CreatedBy:   row.CreatedBy,
			Status:      status,
			Type:        notifType,
			IsRead:      isRead,
			CreatedAt:   createdAt,
		})
	}

	return notifications, nil
}

// CountNotifications counts the total notifications for a user
func (r *NotificationsRepository) CountNotifications(ctx context.Context, createdBy string) (int64, error) {
	createdByUUID, err := uuid.Parse(createdBy)
	if err != nil {
		return 0, fmt.Errorf("invalid created_by UUID: %w", err)
	}

	total, err := r.queries.CountNotifications(ctx, createdByUUID)
	if err != nil {
		return 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	return total, nil
}

// CreateNotification creates a new notification
func (r *NotificationsRepository) CreateNotification(ctx context.Context, notification *domain.Notification) error {
	var recipientID uuid.NullUUID
	if notification.RecipientID != nil {
		recipientID = uuid.NullUUID{UUID: *notification.RecipientID, Valid: true}
	}

	var notifType sql.NullString
	if notification.Type != nil {
		notifType = sql.NullString{String: *notification.Type, Valid: true}
	}

	_, err := r.queries.CreateNotification(ctx, sqlc.CreateNotificationParams{
		ID:          notification.ID,
		Title:       notification.Title,
		Message:     notification.Message,
		SentTo:      notification.SentTo,
		RecipientID: recipientID,
		CreatedBy:   notification.CreatedBy,
		Status:      sql.NullString{String: notification.Status, Valid: true},
		Type:        notifType,
		CreatedAt:   sql.NullTime{Time: notification.CreatedAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

// GetNotificationByID retrieves a notification by ID
func (r *NotificationsRepository) GetNotificationByID(ctx context.Context, id string) (*domain.Notification, error) {
	notifID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid notification ID: %w", err)
	}

	row, err := r.queries.GetNotificationByID(ctx, notifID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	var recipientID *uuid.UUID
	if row.RecipientID.Valid {
		recipientID = &row.RecipientID.UUID
	}

	var notifType *string
	if row.Type.Valid {
		notifType = &row.Type.String
	}

	status := "sent"
	if row.Status.Valid {
		status = row.Status.String
	}

	isRead := false
	if row.IsRead.Valid {
		isRead = row.IsRead.Bool
	}

	createdAt := row.CreatedAt.Time
	if !row.CreatedAt.Valid {
		createdAt = time.Now()
	}

	return &domain.Notification{
		ID:          row.ID,
		Title:       row.Title,
		Message:     row.Message,
		SentTo:      row.SentTo,
		RecipientID: recipientID,
		CreatedBy:   row.CreatedBy,
		Status:      status,
		Type:        notifType,
		IsRead:      isRead,
		CreatedAt:   createdAt,
	}, nil
}

// ListImportantAlerts retrieves a paginated list of important alerts from pentesters
func (r *NotificationsRepository) ListImportantAlerts(ctx context.Context, limit, offset int) ([]domain.Alert, error) {
	rows, err := r.queries.ListImportantAlerts(ctx, sqlc.ListImportantAlertsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list important alerts: %w", err)
	}

	alerts := make([]domain.Alert, 0, len(rows))
	for _, row := range rows {
		var alertType, priority, source *string
		if row.AlertType.Valid {
			alertType = &row.AlertType.String
		}
		if row.Priority.Valid {
			priority = &row.Priority.String
		}
		if row.Source.Valid {
			source = &row.Source.String
		}

		var senderID *uuid.UUID
		if row.SenderID.Valid {
			senderID = &row.SenderID.UUID
		}

		var resolvedBy *uuid.UUID
		if row.ResolvedBy.Valid {
			resolvedBy = &row.ResolvedBy.UUID
		}

		var resolvedAt *time.Time
		if row.ResolvedAt.Valid {
			resolvedAt = &row.ResolvedAt.Time
		}

		isResolved := false
		if row.IsResolved.Valid {
			isResolved = row.IsResolved.Bool
		}

		createdAt := row.CreatedAt.Time
		if !row.CreatedAt.Valid {
			createdAt = time.Now()
		}

		updatedAt := row.UpdatedAt.Time
		if !row.UpdatedAt.Valid {
			updatedAt = time.Now()
		}

		alerts = append(alerts, domain.Alert{
			ID:         row.ID,
			Title:      row.Title,
			Message:    row.Message,
			AlertType:  alertType,
			Priority:   priority,
			Source:     source,
			SenderID:   senderID,
			IsResolved: isResolved,
			ResolvedBy: resolvedBy,
			ResolvedAt: resolvedAt,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		})
	}

	return alerts, nil
}

// CountImportantAlerts counts the total important alerts from pentesters
func (r *NotificationsRepository) CountImportantAlerts(ctx context.Context) (int64, error) {
	total, err := r.queries.CountImportantAlerts(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count important alerts: %w", err)
	}

	return total, nil
}

// CreateAlert creates a new alert
func (r *NotificationsRepository) CreateAlert(ctx context.Context, alert *domain.Alert) (*domain.Alert, error) {
	var alertType, priority, source sql.NullString
	if alert.AlertType != nil {
		alertType = sql.NullString{String: *alert.AlertType, Valid: true}
	}
	if alert.Priority != nil {
		priority = sql.NullString{String: *alert.Priority, Valid: true}
	}
	if alert.Source != nil {
		source = sql.NullString{String: *alert.Source, Valid: true}
	}

	var senderID, recipientID uuid.NullUUID
	if alert.SenderID != nil {
		senderID = uuid.NullUUID{UUID: *alert.SenderID, Valid: true}
	}
	// For admin notifications, recipient is not set in the alert struct
	// It will be set by the caller if needed

	created, err := r.queries.CreateAlert(ctx, sqlc.CreateAlertParams{
		ID:          alert.ID,
		Title:       alert.Title,
		Message:     alert.Message,
		AlertType:   alertType,
		Priority:    priority,
		Source:      source,
		SenderID:    senderID,
		RecipientID: recipientID,
		CreatedAt:   sql.NullTime{Time: alert.CreatedAt, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	var resultAlertType, resultPriority, resultSource *string
	if created.AlertType.Valid {
		resultAlertType = &created.AlertType.String
	}
	if created.Priority.Valid {
		resultPriority = &created.Priority.String
	}
	if created.Source.Valid {
		resultSource = &created.Source.String
	}

	var resultSenderID *uuid.UUID
	if created.SenderID.Valid {
		resultSenderID = &created.SenderID.UUID
	}

	var resultResolvedBy *uuid.UUID
	if created.ResolvedBy.Valid {
		resultResolvedBy = &created.ResolvedBy.UUID
	}

	isResolved := false
	if created.IsResolved.Valid {
		isResolved = created.IsResolved.Bool
	}

	createdAt := created.CreatedAt.Time
	if !created.CreatedAt.Valid {
		createdAt = time.Now()
	}

	updatedAt := created.UpdatedAt.Time
	if !created.UpdatedAt.Valid {
		updatedAt = time.Now()
	}

	return &domain.Alert{
		ID:         created.ID,
		Title:      created.Title,
		Message:    created.Message,
		AlertType:  resultAlertType,
		Priority:   resultPriority,
		Source:     resultSource,
		SenderID:   resultSenderID,
		IsResolved: isResolved,
		ResolvedBy: resultResolvedBy,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

// GetAlertByID retrieves an alert by ID
func (r *NotificationsRepository) GetAlertByID(ctx context.Context, id string) (*domain.Alert, error) {
	alertID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid alert ID: %w", err)
	}

	row, err := r.queries.GetAlertByID(ctx, alertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}

	var alertType, priority, source *string
	if row.AlertType.Valid {
		alertType = &row.AlertType.String
	}
	if row.Priority.Valid {
		priority = &row.Priority.String
	}
	if row.Source.Valid {
		source = &row.Source.String
	}

	var senderID *uuid.UUID
	if row.SenderID.Valid {
		senderID = &row.SenderID.UUID
	}

	var resolvedBy *uuid.UUID
	if row.ResolvedBy.Valid {
		resolvedBy = &row.ResolvedBy.UUID
	}

	var resolvedAt *time.Time
	if row.ResolvedAt.Valid {
		resolvedAt = &row.ResolvedAt.Time
	}

	isResolved := false
	if row.IsResolved.Valid {
		isResolved = row.IsResolved.Bool
	}

	createdAt := row.CreatedAt.Time
	if !row.CreatedAt.Valid {
		createdAt = time.Now()
	}

	updatedAt := row.UpdatedAt.Time
	if !row.UpdatedAt.Valid {
		updatedAt = time.Now()
	}

	return &domain.Alert{
		ID:         row.ID,
		Title:      row.Title,
		Message:    row.Message,
		AlertType:  alertType,
		Priority:   priority,
		Source:     source,
		SenderID:   senderID,
		IsResolved: isResolved,
		ResolvedBy: resolvedBy,
		ResolvedAt: resolvedAt,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

// GetUsersByRole retrieves users by role
func (r *NotificationsRepository) GetUsersByRole(ctx context.Context, role string) ([]domain.User, error) {
	rows, err := r.queries.GetUsersByRole(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}

	users := make([]domain.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, domain.User{
			ID:       row.ID,
			Email:    row.Email,
			FullName: row.FullName,
			Role:     row.Role,
		})
	}

	return users, nil
}

// GetNotifications retrieves notifications with pagination and filters
func (r *NotificationsRepository) GetNotifications(ctx context.Context, page, perPage int, sentTo, notificationType, priority string) ([]domain.Notification, int, error) {
	// For now, return empty results
	return []domain.Notification{}, 0, nil
}

// GetRecipientCount gets the count of recipients for a notification
func (r *NotificationsRepository) GetRecipientCount(ctx context.Context, notificationID uuid.UUID) (int, error) {
	return 1, nil
}

// GetNotificationStats retrieves notification statistics
func (r *NotificationsRepository) GetNotificationStats(ctx context.Context) (*domain.NotificationStatsResponse, error) {
	return &domain.NotificationStatsResponse{
		TotalSent:        0,
		SentToday:        0,
		SentThisWeek:     0,
		SentThisMonth:    0,
		ByType:           make(map[string]int),
		ByPriority:       make(map[string]int),
		ByRecipientType:  make(map[string]int),
		RecentActivity:   []domain.NotificationActivity{},
	}, nil
}
