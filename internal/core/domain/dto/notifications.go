package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateNotificationRequest - Request to create a notification
type CreateNotificationRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Message     string     `json:"message" validate:"required,min=1"`
	SentTo      string     `json:"sent_to" validate:"required,oneof=all_pentesters all_stakeholders specific_user"`
	RecipientID *uuid.UUID `json:"recipient_id,omitempty" validate:"omitempty,uuid"`
}

// NotificationResponse - Response for a notification
type NotificationResponse struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	Message       string     `json:"message"`
	SentTo        string     `json:"sent_to"`
	RecipientID   *uuid.UUID `json:"recipient_id,omitempty"`
	RecipientName *string    `json:"recipient_name,omitempty"`
	CreatedBy     uuid.UUID  `json:"created_by"`
	CreatedByName string     `json:"created_by_name"`
	Status        string     `json:"status"`
	Type          *string    `json:"type,omitempty"`
	IsRead        bool       `json:"is_read"`
	CreatedAt     time.Time  `json:"created_at"`
}

// ListNotificationsRequest - Request to list notifications
type ListNotificationsRequest struct {
	Page    int `json:"page" validate:"omitempty,min=1"`
	PerPage int `json:"per_page" validate:"omitempty,min=1,max=100"`
}

// ListNotificationsResponse - Response for listing notifications
type ListNotificationsResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Pagination    PaginationInfo         `json:"pagination"`
}

// TotalNotificationsSentResponse - Response for total notifications sent
type TotalNotificationsSentResponse struct {
	TotalSent int64 `json:"total_sent"`
}

// CreateAlertRequest - Request to create an alert
type CreateAlertRequest struct {
	Title     string     `json:"title" validate:"required,min=1,max=255"`
	Message   string     `json:"message" validate:"required,min=1"`
	AlertType string     `json:"alert_type" validate:"required,oneof=important critical warning info"`
	Severity  string     `json:"severity" validate:"required,oneof=critical high medium low"`
	Source    string     `json:"source" validate:"omitempty,oneof=pentester system admin"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty" validate:"omitempty,uuid"`
}

// AlertResponse - Response for an alert
type AlertResponse struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	Message       string     `json:"message"`
	AlertType     *string    `json:"alert_type,omitempty"`
	Priority      *string    `json:"priority,omitempty"`
	Source        *string    `json:"source,omitempty"`
	SenderID      *uuid.UUID `json:"sender_id,omitempty"`
	PentesterName *string    `json:"pentester_name,omitempty"`
	IsResolved    bool       `json:"is_resolved"`
	ResolvedBy    *uuid.UUID `json:"resolved_by,omitempty"`
	ResolvedAt    *time.Time `json:"resolved_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ListAlertsRequest - Request to list alerts
type ListAlertsRequest struct {
	Page    int `json:"page" validate:"omitempty,min=1"`
	PerPage int `json:"per_page" validate:"omitempty,min=1,max=100"`
}

// ListAlertsResponse - Response for listing alerts
type ListAlertsResponse struct {
	Alerts     []AlertResponse `json:"alerts"`
	Pagination PaginationInfo  `json:"pagination"`
}
