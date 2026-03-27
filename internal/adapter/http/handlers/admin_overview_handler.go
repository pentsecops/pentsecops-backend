package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/pkg/utils"
)

// AdminOverviewHandler handles admin overview HTTP requests
type AdminOverviewHandler struct {
	useCase domain.AdminOverviewUseCase
}

// NewAdminOverviewHandler creates a new AdminOverviewHandler
func NewAdminOverviewHandler(useCase domain.AdminOverviewUseCase) *AdminOverviewHandler {
	return &AdminOverviewHandler{
		useCase: useCase,
	}
}

// GetOverviewStats handles GET /api/admin/overview/stats
func (h *AdminOverviewHandler) GetOverviewStats(c *fiber.Ctx) error {
	ctx := c.Context()

	stats, err := h.useCase.GetOverviewStats(ctx)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch overview statistics")
	}

	return utils.Success(c, stats, "Overview statistics retrieved successfully")
}

// GetVulnerabilitiesBySeverity handles GET /api/admin/overview/vulnerabilities-by-severity
func (h *AdminOverviewHandler) GetVulnerabilitiesBySeverity(c *fiber.Ctx) error {
	ctx := c.Context()

	data, err := h.useCase.GetVulnerabilitiesBySeverity(ctx)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch vulnerabilities by severity")
	}

	return utils.Success(c, data, "Vulnerabilities by severity retrieved successfully")
}

// GetTop5Domains handles GET /api/admin/overview/top-domains
func (h *AdminOverviewHandler) GetTop5Domains(c *fiber.Ctx) error {
	ctx := c.Context()

	data, err := h.useCase.GetTop5Domains(ctx)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch top domains")
	}

	return utils.Success(c, data, "Top 5 domains retrieved successfully")
}

// GetProjectStatusDistribution handles GET /api/admin/overview/project-status
func (h *AdminOverviewHandler) GetProjectStatusDistribution(c *fiber.Ctx) error {
	ctx := c.Context()

	data, err := h.useCase.GetProjectStatusDistribution(ctx)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch project status distribution")
	}

	return utils.Success(c, data, "Project status distribution retrieved successfully")
}

// GetRecentActivity handles GET /api/admin/overview/recent-activity
func (h *AdminOverviewHandler) GetRecentActivity(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 3)

	// Validate pagination parameters
	if page < 1 {
		return utils.BadRequest(c, "Page must be greater than 0", nil)
	}
	if perPage < 1 || perPage > 100 {
		return utils.BadRequest(c, "Per page must be between 1 and 100", nil)
	}

	data, err := h.useCase.GetRecentActivity(ctx, page, perPage)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch recent activity")
	}

	return utils.Success(c, data, "Recent activity retrieved successfully")
}
