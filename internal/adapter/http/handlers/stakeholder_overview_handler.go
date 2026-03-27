package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/pkg/auth/logger"
	"github.com/pentsecops/backend/pkg/utils"
)

// StakeholderOverviewHandler handles stakeholder overview HTTP requests
type StakeholderOverviewHandler struct {
	useCase domain.StakeholderOverviewUseCase
}

// NewStakeholderOverviewHandler creates a new StakeholderOverviewHandler
func NewStakeholderOverviewHandler(useCase domain.StakeholderOverviewUseCase) *StakeholderOverviewHandler {
	return &StakeholderOverviewHandler{
		useCase: useCase,
	}
}

// ============================================================================
// UC1-UC6: Get Security Metrics Cards
// ============================================================================

// GetSecurityMetrics handles GET /api/stakeholder/overview/security-metrics
func (h *StakeholderOverviewHandler) GetSecurityMetrics(c *fiber.Ctx) error {
	logger.Info("Handler: GetSecurityMetrics called")

	response, err := h.useCase.GetSecurityMetrics(c.Context())
	if err != nil {
		logger.Error("Failed to get security metrics", "error", err)
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, response, "Security metrics fetched successfully")
}

// ============================================================================
// UC7: Get Vulnerability Trend Chart
// ============================================================================

// GetVulnerabilityTrend handles GET /api/stakeholder/overview/vulnerability-trend
func (h *StakeholderOverviewHandler) GetVulnerabilityTrend(c *fiber.Ctx) error {
	logger.Info("Handler: GetVulnerabilityTrend called")

	response, err := h.useCase.GetVulnerabilityTrend(c.Context())
	if err != nil {
		logger.Error("Failed to get vulnerability trend", "error", err)
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, response, "Vulnerability trend fetched successfully")
}

// ============================================================================
// UC8: Get Asset Status Chart
// ============================================================================

// GetAssetStatus handles GET /api/stakeholder/overview/asset-status
func (h *StakeholderOverviewHandler) GetAssetStatus(c *fiber.Ctx) error {
	logger.Info("Handler: GetAssetStatus called")

	response, err := h.useCase.GetAssetStatus(c.Context())
	if err != nil {
		logger.Error("Failed to get asset status", "error", err)
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, response, "Asset status fetched successfully")
}

// ============================================================================
// UC9: Get Recent Security Events
// ============================================================================

// GetRecentSecurityEvents handles GET /api/stakeholder/overview/recent-events
func (h *StakeholderOverviewHandler) GetRecentSecurityEvents(c *fiber.Ctx) error {
	logger.Info("Handler: GetRecentSecurityEvents called")

	// Get limit from query params (default: 10)
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	response, err := h.useCase.GetRecentSecurityEvents(c.Context(), limit)
	if err != nil {
		logger.Error("Failed to get recent security events", "error", err)
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, response, "Recent security events fetched successfully")
}

// ============================================================================
// UC10: Get Remediation Updates
// ============================================================================

// GetRemediationUpdates handles GET /api/stakeholder/overview/remediation-updates
func (h *StakeholderOverviewHandler) GetRemediationUpdates(c *fiber.Ctx) error {
	logger.Info("Handler: GetRemediationUpdates called")

	// Get limit from query params (default: 10)
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	response, err := h.useCase.GetRemediationUpdates(c.Context(), limit)
	if err != nil {
		logger.Error("Failed to get remediation updates", "error", err)
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, response, "Remediation updates fetched successfully")
}

