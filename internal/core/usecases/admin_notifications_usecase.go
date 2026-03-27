package usecases

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
)

type AdminNotificationsUseCase struct {
	notificationsRepo domain.NotificationsRepository
	usersRepo         domain.UsersRepository
}

func NewAdminNotificationsUseCase(
	notificationsRepo domain.NotificationsRepository,
	usersRepo domain.UsersRepository,
) *AdminNotificationsUseCase {
	return &AdminNotificationsUseCase{
		notificationsRepo: notificationsRepo,
		usersRepo:         usersRepo,
	}
}

// SendNotification sends a notification to specified recipients
func (uc *AdminNotificationsUseCase) SendNotification(ctx context.Context, req dto.SendNotificationRequest, adminID uuid.UUID) (*dto.SendNotificationResponse, error) {
	log.Printf("Admin %s sending notification: %s", adminID.String(), req.Title)

	// Validate specific user requirement
	if req.SentTo == "specific_user" && req.UserID == nil {
		return nil, fmt.Errorf("user_id is required when sent_to is 'specific_user'")
	}

	// Set defaults
	if req.Type == "" {
		req.Type = "info"
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	var recipientCount int
	var err error

	switch req.SentTo {
	case "all_pentesters":
		recipientCount, err = uc.sendToAllPentesters(ctx, req, adminID)
	case "all_stakeholders":
		recipientCount, err = uc.sendToAllStakeholders(ctx, req, adminID)
	case "specific_user":
		recipientCount, err = uc.sendToSpecificUser(ctx, req, adminID, *req.UserID)
	default:
		return nil, fmt.Errorf("invalid sent_to value: %s", req.SentTo)
	}

	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}

	response := &dto.SendNotificationResponse{
		ID:             uuid.New(),
		Title:          req.Title,
		Message:        req.Message,
		SentTo:         req.SentTo,
		RecipientCount: recipientCount,
		Type:           req.Type,
		Priority:       req.Priority,
		CreatedAt:      time.Now(),
		DeliveryStatus: "sent",
	}

	log.Printf("Notification sent successfully to %d recipients", recipientCount)
	return response, nil
}

func (uc *AdminNotificationsUseCase) sendToAllPentesters(ctx context.Context, req dto.SendNotificationRequest, adminID uuid.UUID) (int, error) {
	pentesters, err := uc.usersRepo.GetUsersByRole(ctx, "pentester")
	if err != nil {
		return 0, fmt.Errorf("failed to get pentesters: %w", err)
	}

	count := 0
	for _, pentester := range pentesters {
		notification := &domain.Notification{
			ID:          uuid.New(),
			Title:       req.Title,
			Message:     req.Message,
			SentTo:      req.SentTo,
			RecipientID: &pentester.ID,
			CreatedBy:   adminID,
			Status:      "sent",
			Type:        &req.Type,
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		if err := uc.notificationsRepo.CreateNotification(ctx, notification); err != nil {
			log.Printf("Failed to create notification for pentester %s: %v", pentester.ID.String(), err)
			continue
		}
		count++
	}

	return count, nil
}

func (uc *AdminNotificationsUseCase) sendToAllStakeholders(ctx context.Context, req dto.SendNotificationRequest, adminID uuid.UUID) (int, error) {
	stakeholders, err := uc.usersRepo.GetUsersByRole(ctx, "stakeholder")
	if err != nil {
		return 0, fmt.Errorf("failed to get stakeholders: %w", err)
	}

	count := 0
	for _, stakeholder := range stakeholders {
		notification := &domain.Notification{
			ID:          uuid.New(),
			Title:       req.Title,
			Message:     req.Message,
			SentTo:      req.SentTo,
			RecipientID: &stakeholder.ID,
			CreatedBy:   adminID,
			Status:      "sent",
			Type:        &req.Type,
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		if err := uc.notificationsRepo.CreateNotification(ctx, notification); err != nil {
			log.Printf("Failed to create notification for stakeholder %s: %v", stakeholder.ID.String(), err)
			continue
		}
		count++
	}

	return count, nil
}

func (uc *AdminNotificationsUseCase) sendToSpecificUser(ctx context.Context, req dto.SendNotificationRequest, adminID, userID uuid.UUID) (int, error) {
	// Verify user exists
	user, err := uc.usersRepo.GetUserByID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}

	notification := &domain.Notification{
		ID:          uuid.New(),
		Title:       req.Title,
		Message:     req.Message,
		SentTo:      "specific_user",
		RecipientID: &user.ID,
		CreatedBy:   adminID,
		Status:      "sent",
		Type:        &req.Type,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	if err := uc.notificationsRepo.CreateNotification(ctx, notification); err != nil {
		return 0, fmt.Errorf("failed to create notification: %w", err)
	}

	return 1, nil
}

// GetNotifications retrieves notifications with pagination and filters
func (uc *AdminNotificationsUseCase) GetNotifications(ctx context.Context, req dto.NotificationListRequest) (*dto.NotificationListResponse, error) {
	log.Printf("Admin retrieving notifications list")

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 20
	}

	notifications, total, err := uc.notificationsRepo.GetNotifications(ctx, req.Page, req.PerPage, req.SentTo, req.Type, req.Priority)
	if err != nil {
		log.Printf("Failed to get notifications: %v", err)
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	// Convert to DTOs
	notificationItems := make([]dto.NotificationItem, len(notifications))
	for i, notification := range notifications {
		recipientCount := 1
		if notification.SentTo == "all_pentesters" || notification.SentTo == "all_stakeholders" {
			// Get actual count from database
			recipientCount, _ = uc.notificationsRepo.GetRecipientCount(ctx, notification.ID)
		}

		createdBy := "Admin"
		if creator, err := uc.usersRepo.GetUserByID(ctx, notification.CreatedBy); err == nil {
			createdBy = creator.FullName
		}

		notificationItems[i] = dto.NotificationItem{
			ID:             notification.ID,
			Title:          notification.Title,
			Message:        notification.Message,
			SentTo:         notification.SentTo,
			RecipientCount: recipientCount,
			Type:           *notification.Type,
			Priority:       "medium", // Default since not stored in current model
			Status:         notification.Status,
			CreatedBy:      createdBy,
			CreatedAt:      notification.CreatedAt,
		}
	}

	// Calculate pagination
	totalPages := (total + req.PerPage - 1) / req.PerPage
	hasNext := req.Page < totalPages
	hasPrev := req.Page > 1

	pagination := dto.PaginationInfo{
		CurrentPage: req.Page,
		PerPage:     req.PerPage,
		Total:       int64(total),
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}

	return &dto.NotificationListResponse{
		Notifications: notificationItems,
		Pagination:    pagination,
	}, nil
}

// GetNotificationStats retrieves notification statistics
func (uc *AdminNotificationsUseCase) GetNotificationStats(ctx context.Context) (*dto.NotificationStatsResponse, error) {
	log.Printf("Admin retrieving notification statistics")

	stats, err := uc.notificationsRepo.GetNotificationStats(ctx)
	if err != nil {
		log.Printf("Failed to get notification stats: %v", err)
		return nil, fmt.Errorf("failed to get notification stats: %w", err)
	}

	return &dto.NotificationStatsResponse{
		TotalSent:        int(stats.TotalSent),
		SentToday:        int(stats.SentToday),
		SentThisWeek:     int(stats.SentThisWeek),
		SentThisMonth:    int(stats.SentThisMonth),
		ByType:           stats.ByType,
		ByPriority:       stats.ByPriority,
		ByRecipientType:  stats.ByRecipientType,
		RecentActivity:   []dto.NotificationActivity{},
	}, nil
}