package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/pkg/utils"
)

type StakeholderReportsHandler struct {
	useCase  domain.StakeholderReportsUseCase
	validate *validator.Validate
}

// NewStakeholderReportsHandler creates a new stakeholder reports handler
func NewStakeholderReportsHandler(useCase domain.StakeholderReportsUseCase) *StakeholderReportsHandler {
	return &StakeholderReportsHandler{
		useCase:  useCase,
		validate: validator.New(),
	}
}

// UC28-UC30: Get reports statistics
func (h *StakeholderReportsHandler) GetReportsStats(c *fiber.Ctx) error {
	ctx := c.Context()

	stats, err := h.useCase.GetReportsStats(ctx)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "STATS_FETCH_FAILED", err.Error())
	}

	return utils.Success(c, stats, "Reports statistics retrieved successfully")
}

// UC31-UC36: List reports with status filter and pagination
func (h *StakeholderReportsHandler) ListReports(c *fiber.Ctx) error {
	ctx := c.Context()

	var req dto.ListStakeholderReportsRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid query parameters")
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "VALIDATION_FAILED", err.Error())
	}

	response, err := h.useCase.ListReports(ctx, &req)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "LIST_FETCH_FAILED", err.Error())
	}

	return utils.Success(c, response, "Reports retrieved successfully")
}

// UC37-UC38: View report details with vulnerabilities
func (h *StakeholderReportsHandler) ViewReport(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get report ID from query
	reportID := c.Query("report_id")
	if reportID == "" {
		return utils.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Report ID is required")
	}

	// Get evidence pagination parameters
	evidencePage := c.QueryInt("evidence_page", 1)
	evidencePerPage := c.QueryInt("evidence_per_page", 3)

	response, err := h.useCase.ViewReport(ctx, reportID, evidencePage, evidencePerPage)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "REPORT_FETCH_FAILED", err.Error())
	}

	return utils.Success(c, response, "Report details retrieved successfully")
}

// UC39-UC40: Get evidence files for a report
func (h *StakeholderReportsHandler) GetReportEvidenceFiles(c *fiber.Ctx) error {
	ctx := c.Context()

	var req dto.ViewReportEvidenceRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid query parameters")
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "VALIDATION_FAILED", err.Error())
	}

	response, err := h.useCase.GetReportEvidenceFiles(ctx, req.ReportID, req.Page, req.PerPage)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "EVIDENCE_FETCH_FAILED", err.Error())
	}

	return utils.Success(c, response, "Evidence files retrieved successfully")
}

// UC41: Download evidence file
func (h *StakeholderReportsHandler) DownloadEvidenceFile(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get file ID from query
	fileID := c.Query("file_id")
	if fileID == "" {
		return utils.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "File ID is required")
	}

	file, err := h.useCase.DownloadEvidenceFile(ctx, fileID)
	if err != nil {
		return utils.Error(c, fiber.StatusNotFound, "FILE_NOT_FOUND", err.Error())
	}

	// Set headers for file download
	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
	c.Set("Content-Length", fmt.Sprintf("%d", file.FileSize))

	// Send file path as response (in production, you would read and send the actual file)
	return utils.Success(c, fiber.Map{
		"file_id":   file.ID,
		"file_name": file.FileName,
		"file_path": file.FilePath,
		"file_size": file.FileSize,
		"message":   "File ready for download",
	}, "Evidence file retrieved successfully")
}

// UC42: Download report
func (h *StakeholderReportsHandler) DownloadReport(c *fiber.Ctx) error {
	ctx := c.Context()

	var req dto.DownloadReportRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid query parameters")
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "VALIDATION_FAILED", err.Error())
	}

	fileContent, filename, err := h.useCase.DownloadReport(ctx, req.ReportID)
	if err != nil {
		return utils.Error(c, fiber.StatusNotFound, "REPORT_NOT_FOUND", err.Error())
	}

	// Set headers for file download
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Send file content
	return c.Send(fileContent)
}

