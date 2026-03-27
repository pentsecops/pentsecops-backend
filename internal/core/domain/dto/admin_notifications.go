package dto

import (
	"time"

	"github.com/google/uuid"
)

// SendNotificationRequest represents the request to send a notification
type SendNotificationRequest struct {
	Title     string     `json:"title" validate:"required,min=1,max=255"`
	Message   string     `json:"message" validate:"required,min=1"`
	SentTo    string     `json:"sent_to" validate:"required,oneof=all_pentesters all_stakeholders specific_user"`
	UserID    *uuid.UUID `json:"user_id,omitempty"` // Required when sent_to is "specific_user"
	Type      string     `json:"type,omitempty" validate:"omitempty,oneof=info warning error success"`
	Priority  string     `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
}

// SendNotificationResponse represents the response after sending a notification
type SendNotificationResponse struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Message          string    `json:"message"`
	SentTo           string    `json:"sent_to"`
	RecipientCount   int       `json:"recipient_count"`
	Type             string    `json:"type"`
	Priority         string    `json:"priority"`
	CreatedAt        time.Time `json:"created_at"`
	DeliveryStatus   string    `json:"delivery_status"`
}

// NotificationListRequest represents the request to list notifications
type NotificationListRequest struct {
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
	SentTo   string `json:"sent_to,omitempty"`
	Type     string `json:"type,omitempty"`
	Priority string `json:"priority,omitempty"`
}

// NotificationListResponse represents the response for listing notifications
type NotificationListResponse struct {
	Notifications []NotificationItem `json:"notifications"`
	Pagination    PaginationInfo     `json:"pagination"`
}

// NotificationItem represents a single notification in the list
type NotificationItem struct {
	ID             uuid.UUID  `json:"id"`
	Title          string     `json:"title"`
	Message        string     `json:"message"`
	SentTo         string     `json:"sent_to"`
	RecipientCount int        `json:"recipient_count"`
	Type           string     `json:"type"`
	Priority       string     `json:"priority"`
	Status         string     `json:"status"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
}

// NotificationStatsResponse represents notification statistics
type NotificationStatsResponse struct {
	TotalSent        int                    `json:"total_sent"`
	SentToday        int                    `json:"sent_today"`
	SentThisWeek     int                    `json:"sent_this_week"`
	SentThisMonth    int                    `json:"sent_this_month"`
	ByType           map[string]int         `json:"by_type"`
	ByPriority       map[string]int         `json:"by_priority"`
	ByRecipientType  map[string]int         `json:"by_recipient_type"`
	RecentActivity   []NotificationActivity `json:"recent_activity"`
}

// NotificationActivity represents recent notification activity
type NotificationActivity struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	SentTo    string    `json:"sent_to"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

// UserNotificationResponse represents notifications for a specific user
type UserNotificationResponse struct {
	Notifications []UserNotificationItem `json:"notifications"`
	UnreadCount   int                    `json:"unread_count"`
	Pagination    PaginationInfo         `json:"pagination"`
}

// UserNotificationItem represents a notification for a user
type UserNotificationItem struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Priority  string    `json:"priority"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// MarkNotificationReadRequest represents request to mark notification as read
type MarkNotificationReadRequest struct {
	NotificationID uuid.UUID `json:"notification_id" validate:"required"`
	IsRead         bool      `json:"is_read"`
}