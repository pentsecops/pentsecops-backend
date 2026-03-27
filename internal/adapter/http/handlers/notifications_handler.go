package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/pkg/auth/logger"
	"github.com/pentsecops/backend/pkg/utils"
)

// NotificationsHandler handles notifications-related HTTP requests
type NotificationsHandler struct {
	usecase domain.NotificationsUseCase
}

// NewNotificationsHandler creates a new NotificationsHandler
func NewNotificationsHandler(usecase domain.NotificationsUseCase) *NotificationsHandler {
	return &NotificationsHandler{
		usecase: usecase,
	}
}

// GetTotalNotificationsSent handles GET /api/admin/notifications/total
func (h *NotificationsHandler) GetTotalNotificationsSent(c *fiber.Ctx) error {
	// Get created_by from context (set by auth middleware)
	createdBy := c.Locals("user_id").(string)

	result, err := h.usecase.GetTotalNotificationsSent(c.Context(), createdBy)
	if err != nil {
		logger.Error("Failed to get total notifications sent", "error", err)
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to get total notifications sent")
	}

	return utils.Success(c, result, "Total notifications sent retrieved successfully")
}

// ListNotifications handles GET /api/admin/notifications
func (h *NotificationsHandler) ListNotifications(c *fiber.Ctx) error {
	// Get created_by from context (set by auth middleware)
	createdBy := c.Locals("user_id").(string)

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "5"))

	req := &dto.ListNotificationsRequest{
		Page:    page,
		PerPage: perPage,
	}

	result, err := h.usecase.ListNotifications(c.Context(), createdBy, req)
	if err != nil {
		logger.Error("Failed to list notifications", "error", err)
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to list notifications")
	}

	return utils.Success(c, result, "Notifications retrieved successfully")
}

// CreateNotification handles POST /api/admin/notifications
func (h *NotificationsHandler) CreateNotification(c *fiber.Ctx) error {
	// Get created_by from context (set by auth middleware)
	createdBy := c.Locals("user_id").(string)

	// Parse request body
	var req dto.CreateNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", "error", err)
		return utils.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	result, err := h.usecase.CreateNotification(c.Context(), createdBy, &req)
	if err != nil {
		logger.Error("Failed to create notification", "error", err)
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, result, "Notification sent successfully")
}

// ListImportantAlerts handles GET /api/admin/notifications/alerts
func (h *NotificationsHandler) ListImportantAlerts(c *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "5"))

	req := &dto.ListAlertsRequest{
		Page:    page,
		PerPage: perPage,
	}

	result, err := h.usecase.ListImportantAlerts(c.Context(), req)
	if err != nil {
		logger.Error("Failed to list important alerts", "error", err)
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to list important alerts")
	}

	return utils.Success(c, result, "Important alerts retrieved successfully")
}

