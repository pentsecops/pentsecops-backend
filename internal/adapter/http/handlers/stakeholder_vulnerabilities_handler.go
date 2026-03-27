package handlers

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/pkg/utils"
)

type StakeholderVulnerabilitiesHandler struct {
	useCase  domain.StakeholderVulnerabilitiesUseCase
	validate *validator.Validate
}

// NewStakeholderVulnerabilitiesHandler creates a new stakeholder vulnerabilities handler
func NewStakeholderVulnerabilitiesHandler(useCase domain.StakeholderVulnerabilitiesUseCase) *StakeholderVulnerabilitiesHandler {
	return &StakeholderVulnerabilitiesHandler{
		useCase:  useCase,
		validate: validator.New(),
	}
}

// UC11-UC14: Get vulnerabilities statistics
func (h *StakeholderVulnerabilitiesHandler) GetVulnerabilitiesStats(c *fiber.Ctx) error {
	ctx := c.Context()

	stats, err := h.useCase.GetVulnerabilitiesStats(ctx)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "STATS_FETCH_ERROR", "Failed to fetch vulnerabilities statistics")
	}

	return utils.Success(c, stats, "Vulnerabilities statistics fetched successfully")
}

// UC15-UC23: List vulnerabilities with search, filters, and pagination
func (h *StakeholderVulnerabilitiesHandler) ListVulnerabilities(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse query parameters
	var req dto.ListStakeholderVulnerabilitiesRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "INVALID_QUERY_PARAMS", "Invalid query parameters")
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PerPage < 1 {
		req.PerPage = 5
	}

	// Get vulnerabilities
	response, err := h.useCase.ListVulnerabilities(ctx, &req)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "LIST_ERROR", "Failed to fetch vulnerabilities")
	}

	return utils.Success(c, response, "Vulnerabilities fetched successfully")
}

// UC24: Export vulnerabilities to CSV
func (h *StakeholderVulnerabilitiesHandler) ExportVulnerabilitiesToCSV(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse query parameters
	var req dto.ExportStakeholderVulnerabilitiesRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "INVALID_QUERY_PARAMS", "Invalid query parameters")
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	// Export to CSV
	csvData, err := h.useCase.ExportVulnerabilitiesToCSV(ctx, &req)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "EXPORT_ERROR", "Failed to export vulnerabilities")
	}

	// Generate filename with timestamp
	filename := fmt.Sprintf("vulnerabilities_export_%s.csv", time.Now().Format("20060102_150405"))

	// Set headers for CSV download
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.Send(csvData)
}

// UC25-UC27: Get SLA compliance data
func (h *StakeholderVulnerabilitiesHandler) GetSLACompliance(c *fiber.Ctx) error {
	ctx := c.Context()

	slaData, err := h.useCase.GetSLACompliance(ctx)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "SLA_FETCH_ERROR", "Failed to fetch SLA compliance data")
	}

	return utils.Success(c, slaData, "SLA compliance data fetched successfully")
}

