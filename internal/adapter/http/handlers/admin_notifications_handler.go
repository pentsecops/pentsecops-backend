package handlers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/internal/core/usecases"
)

type AdminNotificationsHandler struct {
	useCase *usecases.AdminNotificationsUseCase
}

func NewAdminNotificationsHandler(useCase *usecases.AdminNotificationsUseCase) *AdminNotificationsHandler {
	return &AdminNotificationsHandler{
		useCase: useCase,
	}
}

// SendNotification sends a notification to specified recipients
func (h *AdminNotificationsHandler) SendNotification(c *fiber.Ctx) error {
	log.Printf("Admin sending notification")

	var req dto.SendNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse send notification request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
			},
		})
	}

	// Get admin ID from context
	adminID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		log.Printf("Invalid admin ID: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "Invalid user ID",
			},
		})
	}

	// Send notification
	response, err := h.useCase.SendNotification(c.Context(), req, adminID)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "SEND_FAILED",
				"message": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Notification sent successfully",
	})
}

// GetNotifications retrieves notifications with pagination and filters
func (h *AdminNotificationsHandler) GetNotifications(c *fiber.Ctx) error {
	log.Printf("Admin retrieving notifications")

	// Parse query parameters
	req := dto.NotificationListRequest{
		Page:    1,
		PerPage: 20,
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			req.Page = p
		}
	}

	if perPage := c.Query("per_page"); perPage != "" {
		if pp, err := strconv.Atoi(perPage); err == nil && pp > 0 && pp <= 100 {
			req.PerPage = pp
		}
	}

	req.SentTo = c.Query("sent_to")
	req.Type = c.Query("type")
	req.Priority = c.Query("priority")

	// Get notifications
	response, err := h.useCase.GetNotifications(c.Context(), req)
	if err != nil {
		log.Printf("Failed to get notifications: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "FETCH_FAILED",
				"message": "Failed to retrieve notifications",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Notifications retrieved successfully",
	})
}

// GetNotificationStats retrieves notification statistics
func (h *AdminNotificationsHandler) GetNotificationStats(c *fiber.Ctx) error {
	log.Printf("Admin retrieving notification statistics")

	stats, err := h.useCase.GetNotificationStats(c.Context())
	if err != nil {
		log.Printf("Failed to get notification stats: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "STATS_FAILED",
				"message": "Failed to retrieve notification statistics",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
		"message": "Notification statistics retrieved successfully",
	})
}

// GetAvailableUsers retrieves users available for specific notifications
func (h *AdminNotificationsHandler) GetAvailableUsers(c *fiber.Ctx) error {
	log.Printf("Admin retrieving available users for notifications")

	role := c.Query("role") // "pentester" or "stakeholder"
	if role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "MISSING_ROLE",
				"message": "Role parameter is required",
			},
		})
	}

	if role != "pentester" && role != "stakeholder" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INVALID_ROLE",
				"message": "Role must be 'pentester' or 'stakeholder'",
			},
		})
	}

	// This would be implemented in the use case if needed
	// For now, return a simple response
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"users": []fiber.Map{},
		},
		"message": "Available users retrieved successfully",
	})
}