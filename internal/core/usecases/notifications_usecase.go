package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
)

// NotificationsUseCase implements the NotificationsUseCase interface
type NotificationsUseCase struct {
	repo domain.NotificationsRepository
}

// NewNotificationsUseCase creates a new NotificationsUseCase
func NewNotificationsUseCase(repo domain.NotificationsRepository) *NotificationsUseCase {
	return &NotificationsUseCase{
		repo: repo,
	}
}

// GetTotalNotificationsSent retrieves the total count of notifications sent
func (uc *NotificationsUseCase) GetTotalNotificationsSent(ctx context.Context, createdBy string) (*dto.TotalNotificationsSentResponse, error) {
	total, err := uc.repo.GetTotalNotificationsSent(ctx, createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to get total notifications sent: %w", err)
	}

	return &dto.TotalNotificationsSentResponse{
		TotalSent: total,
	}, nil
}

// ListNotifications retrieves a paginated list of notifications
func (uc *NotificationsUseCase) ListNotifications(ctx context.Context, createdBy string, req *dto.ListNotificationsRequest) (*dto.ListNotificationsResponse, error) {
	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Set defaults
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 5
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Get notifications
	notifications, err := uc.repo.ListNotifications(ctx, createdBy, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	// Get total count
	total, err := uc.repo.CountNotifications(ctx, createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to count notifications: %w", err)
	}

	// Calculate total pages
	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	// Build response
	notificationResponses := make([]dto.NotificationResponse, 0, len(notifications))
	for _, notif := range notifications {
		notificationResponses = append(notificationResponses, dto.NotificationResponse{
			ID:          notif.ID,
			Title:       notif.Title,
			Message:     notif.Message,
			SentTo:      notif.SentTo,
			RecipientID: notif.RecipientID,
			CreatedBy:   notif.CreatedBy,
			Status:      notif.Status,
			Type:        notif.Type,
			IsRead:      notif.IsRead,
			CreatedAt:   notif.CreatedAt,
		})
	}

	return &dto.ListNotificationsResponse{
		Notifications: notificationResponses,
		Pagination: dto.PaginationInfo{
			CurrentPage: page,
			PerPage:     perPage,
			Total:       total,
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrev:     page > 1,
		},
	}, nil
}

// CreateNotification creates a new notification
func (uc *NotificationsUseCase) CreateNotification(ctx context.Context, createdBy string, req *dto.CreateNotificationRequest) (*dto.NotificationResponse, error) {
	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Validate recipient_id if sent_to is specific_user
	if req.SentTo == domain.SentToSpecificUser && req.RecipientID == nil {
		return nil, fmt.Errorf("recipient_id is required when sent_to is specific_user")
	}

	// Parse created_by UUID
	createdByUUID, err := uuid.Parse(createdBy)
	if err != nil {
		return nil, fmt.Errorf("invalid created_by UUID: %w", err)
	}

	// Create notification
	notification := &domain.Notification{
		ID:          uuid.Must(uuid.NewV7()),
		Title:       req.Title,
		Message:     req.Message,
		SentTo:      req.SentTo,
		RecipientID: req.RecipientID,
		CreatedBy:   createdByUUID,
		Status:      domain.NotificationStatusSent,
		Type:        nil, // Can be set based on business logic
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	err = uc.repo.CreateNotification(ctx, notification)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return &dto.NotificationResponse{
		ID:          notification.ID,
		Title:       notification.Title,
		Message:     notification.Message,
		SentTo:      notification.SentTo,
		RecipientID: notification.RecipientID,
		CreatedBy:   notification.CreatedBy,
		Status:      notification.Status,
		Type:        notification.Type,
		IsRead:      notification.IsRead,
		CreatedAt:   notification.CreatedAt,
	}, nil
}

// ListImportantAlerts retrieves a paginated list of important alerts from pentesters
func (uc *NotificationsUseCase) ListImportantAlerts(ctx context.Context, req *dto.ListAlertsRequest) (*dto.ListAlertsResponse, error) {
	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Set defaults
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 5
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Get alerts
	alerts, err := uc.repo.ListImportantAlerts(ctx, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list important alerts: %w", err)
	}

	// Get total count
	total, err := uc.repo.CountImportantAlerts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count important alerts: %w", err)
	}

	// Calculate total pages
	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	// Build response
	alertResponses := make([]dto.AlertResponse, 0, len(alerts))
	for _, alert := range alerts {
		alertResponses = append(alertResponses, dto.AlertResponse{
			ID:         alert.ID,
			Title:      alert.Title,
			Message:    alert.Message,
			AlertType:  alert.AlertType,
			Priority:   alert.Priority,
			Source:     alert.Source,
			SenderID:   alert.SenderID,
			IsResolved: alert.IsResolved,
			ResolvedBy: alert.ResolvedBy,
			ResolvedAt: alert.ResolvedAt,
			CreatedAt:  alert.CreatedAt,
			UpdatedAt:  alert.UpdatedAt,
		})
	}

	return &dto.ListAlertsResponse{
		Alerts: alertResponses,
		Pagination: dto.PaginationInfo{
			CurrentPage: page,
			PerPage:     perPage,
			Total:       total,
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrev:     page > 1,
		},
	}, nil
}
