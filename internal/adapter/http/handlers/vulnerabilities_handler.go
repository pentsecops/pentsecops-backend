package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/pkg/utils"
)

type VulnerabilitiesHandler struct {
	usecase domain.VulnerabilitiesUseCase
}

func NewVulnerabilitiesHandler(usecase domain.VulnerabilitiesUseCase) *VulnerabilitiesHandler {
	return &VulnerabilitiesHandler{
		usecase: usecase,
	}
}

// CreateVulnerability - UC44, UC45: Open Add Vulnerability Dialog and Create New Vulnerability
func (h *VulnerabilitiesHandler) CreateVulnerability(c *fiber.Ctx) error {
	fmt.Printf("CreateVulnerability handler called\n")
	var req dto.CreateVulnerabilityRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Printf("Body parser error: %v\n", err)
		return utils.BadRequest(c, "Invalid request body", nil)
	}
	fmt.Printf("Request parsed successfully: %+v\n", req)

	response, err := h.usecase.CreateVulnerability(c.Context(), &req)
	if err != nil {
		fmt.Printf("Create vulnerability error: %v\n", err)
		// Check for specific errors
		if err.Error() == "due date must be after discovered date" {
			return utils.BadRequest(c, err.Error(), nil)
		}
		if strings.Contains(err.Error(), "validation failed") {
			return utils.BadRequest(c, err.Error(), nil)
		}
		return utils.InternalServerError(c, "Failed to create vulnerability")
	}

	return utils.Success(c, response, "Vulnerability added successfully")
}

// GetVulnerabilityByID - Get vulnerability by ID
func (h *VulnerabilitiesHandler) GetVulnerabilityByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Vulnerability ID is required", nil)
	}

	response, err := h.usecase.GetVulnerabilityByID(c.Context(), id)
	if err != nil {
		if err.Error() == "vulnerability not found" {
			return utils.NotFound(c, "Vulnerability not found")
		}
		return utils.InternalServerError(c, "Failed to get vulnerability")
	}

	return utils.Success(c, response, "Vulnerability retrieved successfully")
}

// UpdateVulnerability - UC52: Edit Vulnerability from Table
func (h *VulnerabilitiesHandler) UpdateVulnerability(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Vulnerability ID is required", nil)
	}

	var req dto.UpdateVulnerabilityRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	response, err := h.usecase.UpdateVulnerability(c.Context(), id, &req)
	if err != nil {
		// Check for specific errors
		if err.Error() == "due date must be after discovered date" {
			return utils.BadRequest(c, err.Error(), nil)
		}
		return utils.InternalServerError(c, "Failed to update vulnerability")
	}

	return utils.Success(c, response, "Vulnerability updated successfully")
}

// DeleteVulnerability - Delete vulnerability
func (h *VulnerabilitiesHandler) DeleteVulnerability(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Vulnerability ID is required", nil)
	}

	err := h.usecase.DeleteVulnerability(c.Context(), id)
	if err != nil {
		return utils.InternalServerError(c, "Failed to delete vulnerability")
	}

	return utils.Success(c, nil, "Vulnerability deleted successfully")
}

// ListVulnerabilities - UC40, UC41, UC42, UC43, UC48, UC49: List vulnerabilities with search and filters
func (h *VulnerabilitiesHandler) ListVulnerabilities(c *fiber.Ctx) error {
	var req dto.ListVulnerabilitiesRequest

	// Parse query parameters
	if err := c.QueryParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid query parameters", nil)
	}

	response, err := h.usecase.ListVulnerabilities(c.Context(), &req)
	if err != nil {
		return utils.InternalServerError(c, "Failed to list vulnerabilities")
	}

	return utils.Success(c, response, "Vulnerabilities retrieved successfully")
}

// GetVulnerabilityStats - UC39: Fetch and Display Vulnerability Statistics
func (h *VulnerabilitiesHandler) GetVulnerabilityStats(c *fiber.Ctx) error {
	response, err := h.usecase.GetVulnerabilityStats(c.Context())
	if err != nil {
		return utils.InternalServerError(c, "Failed to get vulnerability statistics")
	}

	return utils.Success(c, response, "Vulnerability statistics retrieved successfully")
}

// GetSLACompliance - UC54: Display SLA Compliance Information
func (h *VulnerabilitiesHandler) GetSLACompliance(c *fiber.Ctx) error {
	response, err := h.usecase.GetSLACompliance(c.Context())
	if err != nil {
		return utils.InternalServerError(c, "Failed to get SLA compliance")
	}

	return utils.Success(c, response, "SLA compliance retrieved successfully")
}

// ExportVulnerabilitiesToCSV - UC53: Export Vulnerabilities to CSV
func (h *VulnerabilitiesHandler) ExportVulnerabilitiesToCSV(c *fiber.Ctx) error {
	var req dto.ListVulnerabilitiesRequest

	// Parse query parameters
	if err := c.QueryParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid query parameters", nil)
	}

	csvData, err := h.usecase.ExportVulnerabilitiesToCSV(c.Context(), &req)
	if err != nil {
		return utils.InternalServerError(c, "Failed to export vulnerabilities")
	}

	// Set headers for CSV download
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("vulnerabilities_export_%s.csv", timestamp)

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.Send(csvData)
}

