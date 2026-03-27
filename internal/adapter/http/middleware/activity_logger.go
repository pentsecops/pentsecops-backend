package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/core/domain"
)

type ActivityLog struct {
	Timestamp    time.Time `json:"timestamp"`
	UserID       string    `json:"user_id,omitempty"`
	UserEmail    string    `json:"user_email,omitempty"`
	UserRole     string    `json:"user_role,omitempty"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	StatusCode   int       `json:"status_code"`
	RequestBody  string    `json:"request_body,omitempty"`
	ResponseBody string    `json:"response_body,omitempty"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	Duration     string    `json:"duration"`
	Action       string    `json:"action"`
	EntityType   string    `json:"entity_type,omitempty"`
	EntityID     string    `json:"entity_id,omitempty"`
}

func ActivityLogger(auditRepo domain.AuditRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Read request body
		var requestBody string
		if c.Body() != nil {
			// Don't log sensitive data
			if !isSensitiveEndpoint(c.Path()) {
				requestBody = string(c.Body())
				if len(requestBody) > 500 {
					requestBody = requestBody[:500] + "...[TRUNCATED]"
				}
			} else {
				requestBody = "[SENSITIVE_DATA_HIDDEN]"
			}
		}

		// Continue to next handler
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get user info from context
		userID := getUserID(c)
		userEmail := getUserEmail(c)
		userRole := getUserRole(c)

		// Get response body (limited)
		responseBody := string(c.Response().Body())
		if len(responseBody) > 1000 {
			responseBody = responseBody[:1000] + "...[TRUNCATED]"
		}

		// Determine action and entity
		action, entityType, entityID := parseAction(c.Method(), c.Path(), requestBody)

		// Create activity log
		activityLog := ActivityLog{
			Timestamp:    start,
			UserID:       userID,
			UserEmail:    userEmail,
			UserRole:     userRole,
			Method:       c.Method(),
			Path:         c.Path(),
			StatusCode:   c.Response().StatusCode(),
			RequestBody:  requestBody,
			ResponseBody: responseBody,
			IPAddress:    c.IP(),
			UserAgent:    c.Get("User-Agent"),
			Duration:     duration.String(),
			Action:       action,
			EntityType:   entityType,
			EntityID:     entityID,
		}

		// Log the activity to file and database
		logActivity(activityLog)
		logActivityToDatabase(auditRepo, activityLog)

		return err
	}
}

func isSensitiveEndpoint(path string) bool {
	sensitiveEndpoints := []string{
		"/auth/login",
		"/auth/change-password",
		"/auth/refresh",
	}
	
	for _, endpoint := range sensitiveEndpoints {
		if strings.Contains(path, endpoint) {
			return true
		}
	}
	return false
}

func getUserID(c *fiber.Ctx) string {
	if userID := c.Locals("user_id"); userID != nil {
		return userID.(string)
	}
	return ""
}

func getUserEmail(c *fiber.Ctx) string {
	if userEmail := c.Locals("user_email"); userEmail != nil {
		return userEmail.(string)
	}
	return ""
}

func getUserRole(c *fiber.Ctx) string {
	if userRole := c.Locals("user_role"); userRole != nil {
		return userRole.(string)
	}
	return ""
}

func parseAction(method, path, requestBody string) (action, entityType, entityID string) {
	// Extract entity ID from path
	pathParts := strings.Split(path, "/")
	if len(pathParts) > 0 {
		lastPart := pathParts[len(pathParts)-1]
		if isUUID(lastPart) {
			entityID = lastPart
		}
	}

	// Determine entity type and action
	switch {
	case strings.Contains(path, "/users"):
		entityType = "user"
		action = getActionFromMethod(method, "user")
	case strings.Contains(path, "/projects"):
		entityType = "project"
		action = getActionFromMethod(method, "project")
	case strings.Contains(path, "/vulnerabilities"):
		entityType = "vulnerability"
		action = getActionFromMethod(method, "vulnerability")
	case strings.Contains(path, "/tasks"):
		entityType = "task"
		action = getActionFromMethod(method, "task")
	case strings.Contains(path, "/auth/login"):
		entityType = "auth"
		action = "login"
	case strings.Contains(path, "/auth/logout"):
		entityType = "auth"
		action = "logout"
	case strings.Contains(path, "/auth/change-password"):
		entityType = "auth"
		action = "change_password"
	case strings.Contains(path, "/overview"):
		entityType = "overview"
		action = "view_dashboard"
	default:
		entityType = "system"
		action = fmt.Sprintf("%s_%s", strings.ToLower(method), strings.ReplaceAll(path, "/", "_"))
	}

	return action, entityType, entityID
}

func getActionFromMethod(method, entityType string) string {
	switch method {
	case "GET":
		return fmt.Sprintf("view_%s", entityType)
	case "POST":
		return fmt.Sprintf("create_%s", entityType)
	case "PUT", "PATCH":
		return fmt.Sprintf("update_%s", entityType)
	case "DELETE":
		return fmt.Sprintf("delete_%s", entityType)
	default:
		return fmt.Sprintf("unknown_%s", entityType)
	}
}

func isUUID(s string) bool {
	return len(s) == 36 && strings.Count(s, "-") == 4
}

func logActivity(activity ActivityLog) {
	// Log as JSON for structured logging
	activityJSON, err := json.Marshal(activity)
	if err != nil {
		log.Printf("Failed to marshal activity log: %v", err)
		return
	}

	// Log with different levels based on status code
	switch {
	case activity.StatusCode >= 500:
		log.Printf("ERROR ACTIVITY: %s", string(activityJSON))
	case activity.StatusCode >= 400:
		log.Printf("WARN ACTIVITY: %s", string(activityJSON))
	default:
		log.Printf("INFO ACTIVITY: %s", string(activityJSON))
	}

	// Also log human-readable format
	logMessage := fmt.Sprintf("User[%s|%s|%s] %s %s -> %d (%s) [%s]",
		activity.UserID,
		activity.UserEmail,
		activity.UserRole,
		activity.Method,
		activity.Path,
		activity.StatusCode,
		activity.Duration,
		activity.Action,
	)

	log.Printf("USER_ACTIVITY: %s", logMessage)
}

func logActivityToDatabase(auditRepo domain.AuditRepository, activity ActivityLog) {
	// Convert to domain ActivityLog
	var userID *uuid.UUID
	if activity.UserID != "" {
		if id, err := uuid.Parse(activity.UserID); err == nil {
			userID = &id
		}
	}

	var entityID *uuid.UUID
	if activity.EntityID != "" {
		if id, err := uuid.Parse(activity.EntityID); err == nil {
			entityID = &id
		}
	}

	var entityType *string
	if activity.EntityType != "" {
		entityType = &activity.EntityType
	}

	var errorMessage *string
	if activity.StatusCode >= 400 {
		errorMessage = &activity.ResponseBody
	}

	dbLog := &domain.ActivityLog{
		UserID:       userID,
		UserEmail:    activity.UserEmail,
		UserRole:     activity.UserRole,
		Action:       activity.Action,
		EntityType:   entityType,
		EntityID:     entityID,
		IPAddress:    activity.IPAddress,
		UserAgent:    activity.UserAgent,
		Endpoint:     activity.Path,
		Method:       activity.Method,
		StatusCode:   activity.StatusCode,
		ErrorMessage: errorMessage,
	}

	// Log asynchronously to avoid blocking requests
	go func() {
		if err := auditRepo.LogActivity(context.Background(), dbLog); err != nil {
			log.Printf("Failed to log activity to database: %v", err)
		}
	}()
}