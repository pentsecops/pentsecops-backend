package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/core/domain"
)

type AuditHandler struct {
	auditRepo domain.AuditRepository
}

func NewAuditHandler(auditRepo domain.AuditRepository) *AuditHandler {
	return &AuditHandler{
		auditRepo: auditRepo,
	}
}

// GetActivityLogs retrieves activity logs with filters and pagination
func (h *AuditHandler) GetActivityLogs(c *fiber.Ctx) error {
	log.Printf("Admin requesting activity logs")

	// Parse query parameters
	filter := domain.AuditLogFilter{
		Page:    1,
		PerPage: 20,
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.Page = p
		}
	}

	if perPage := c.Query("per_page"); perPage != "" {
		if pp, err := strconv.Atoi(perPage); err == nil && pp > 0 && pp <= 100 {
			filter.PerPage = pp
		}
	}

	if userID := c.Query("user_id"); userID != "" {
		if id, err := uuid.Parse(userID); err == nil {
			filter.UserID = &id
		}
	}

	if action := c.Query("action"); action != "" {
		filter.Action = action
	}

	if entityType := c.Query("entity_type"); entityType != "" {
		filter.EntityType = entityType
	}

	if entityID := c.Query("entity_id"); entityID != "" {
		if id, err := uuid.Parse(entityID); err == nil {
			filter.EntityID = &id
		}
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &date
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			// Set to end of day
			endOfDay := date.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filter.EndDate = &endOfDay
		}
	}

	// Get activity logs
	logs, total, err := h.auditRepo.GetActivityLogs(c.Context(), filter)
	if err != nil {
		log.Printf("Failed to get activity logs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve activity logs",
			},
		})
	}

	// Calculate pagination
	totalPages := (total + filter.PerPage - 1) / filter.PerPage
	hasNext := filter.Page < totalPages
	hasPrev := filter.Page > 1

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"logs": logs,
			"pagination": fiber.Map{
				"current_page": filter.Page,
				"per_page":     filter.PerPage,
				"total":        total,
				"total_pages":  totalPages,
				"has_next":     hasNext,
				"has_prev":     hasPrev,
			},
		},
		"message": "Activity logs retrieved successfully",
	})
}

// GetActivityStats retrieves activity statistics
func (h *AuditHandler) GetActivityStats(c *fiber.Ctx) error {
	log.Printf("Admin requesting activity statistics")

	// Get logs for the last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	filter := domain.AuditLogFilter{
		StartDate: &startDate,
		EndDate:   &endDate,
		Page:      1,
		PerPage:   10000, // Get all for stats
	}

	logs, total, err := h.auditRepo.GetActivityLogs(c.Context(), filter)
	if err != nil {
		log.Printf("Failed to get activity logs for stats: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve activity statistics",
			},
		})
	}

	// Calculate statistics
	stats := map[string]interface{}{
		"total_activities": total,
		"activities_by_action": make(map[string]int),
		"activities_by_user": make(map[string]int),
		"activities_by_entity": make(map[string]int),
		"activities_by_day": make(map[string]int),
		"error_count": 0,
		"success_count": 0,
	}

	actionStats := make(map[string]int)
	userStats := make(map[string]int)
	entityStats := make(map[string]int)
	dayStats := make(map[string]int)
	errorCount := 0
	successCount := 0

	for _, log := range logs {
		// Action statistics
		actionStats[log.Action]++

		// User statistics
		userStats[log.UserEmail]++

		// Entity statistics
		if log.EntityType != nil {
			entityStats[*log.EntityType]++
		}

		// Daily statistics
		day := log.CreatedAt.Format("2006-01-02")
		dayStats[day]++

		// Status statistics
		if log.StatusCode >= 400 {
			errorCount++
		} else {
			successCount++
		}
	}

	stats["activities_by_action"] = actionStats
	stats["activities_by_user"] = userStats
	stats["activities_by_entity"] = entityStats
	stats["activities_by_day"] = dayStats
	stats["error_count"] = errorCount
	stats["success_count"] = successCount

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
		"message": "Activity statistics retrieved successfully",
	})
}

// ExportActivityLogs exports activity logs to CSV
func (h *AuditHandler) ExportActivityLogs(c *fiber.Ctx) error {
	log.Printf("Admin exporting activity logs to CSV")

	// Parse filters (similar to GetActivityLogs)
	filter := domain.AuditLogFilter{
		Page:    1,
		PerPage: 10000, // Export all matching records
	}

	// Apply same filters as GetActivityLogs
	if userID := c.Query("user_id"); userID != "" {
		if id, err := uuid.Parse(userID); err == nil {
			filter.UserID = &id
		}
	}

	if action := c.Query("action"); action != "" {
		filter.Action = action
	}

	if entityType := c.Query("entity_type"); entityType != "" {
		filter.EntityType = entityType
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &date
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			endOfDay := date.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filter.EndDate = &endOfDay
		}
	}

	// Get activity logs
	logs, _, err := h.auditRepo.GetActivityLogs(c.Context(), filter)
	if err != nil {
		log.Printf("Failed to get activity logs for export: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to export activity logs",
			},
		})
	}

	// Generate CSV content
	csvContent := "ID,User Email,User Role,Action,Entity Type,Entity ID,IP Address,Endpoint,Method,Status Code,Created At\n"
	
	for _, log := range logs {
		entityType := ""
		if log.EntityType != nil {
			entityType = *log.EntityType
		}
		
		entityID := ""
		if log.EntityID != nil {
			entityID = log.EntityID.String()
		}

		csvContent += fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%d,%s\n",
			log.ID.String(),
			log.UserEmail,
			log.UserRole,
			log.Action,
			entityType,
			entityID,
			log.IPAddress,
			log.Endpoint,
			log.Method,
			log.StatusCode,
			log.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}

	// Set headers for file download
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=activity_logs_"+time.Now().Format("20060102_150405")+".csv")

	return c.SendString(csvContent)
}