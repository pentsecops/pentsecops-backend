package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/pkg/utils"
)

// DomainsHandler handles HTTP requests for domains
type DomainsHandler struct {
	usecase domain.DomainsUseCase
}

// NewDomainsHandler creates a new DomainsHandler
func NewDomainsHandler(usecase domain.DomainsUseCase) *DomainsHandler {
	return &DomainsHandler{
		usecase: usecase,
	}
}

// GetDomainsStats handles GET /api/admin/domains/stats
func (h *DomainsHandler) GetDomainsStats(c *fiber.Ctx) error {
	stats, err := h.usecase.GetDomainsStats(c.Context())
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, stats, "Domains statistics retrieved successfully")
}

// ListDomains handles GET /api/admin/domains
func (h *DomainsHandler) ListDomains(c *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "5"))

	req := &dto.ListDomainsRequest{
		Page:    page,
		PerPage: perPage,
	}

	domains, err := h.usecase.ListDomains(c.Context(), req)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, domains, "Domains retrieved successfully")
}

// GetDomainByID handles GET /api/admin/domains/:id
func (h *DomainsHandler) GetDomainByID(c *fiber.Ctx) error {
	id := c.Params("id")

	domain, err := h.usecase.GetDomainByID(c.Context(), id)
	if err != nil {
		return utils.Error(c, fiber.StatusNotFound, "DOMAIN_NOT_FOUND", err.Error())
	}

	return utils.Success(c, domain, "Domain retrieved successfully")
}

// CreateDomain handles POST /api/admin/domains
func (h *DomainsHandler) CreateDomain(c *fiber.Ctx) error {
	var req dto.CreateDomainRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST_BODY", err.Error())
	}

	domain, err := h.usecase.CreateDomain(c.Context(), &req)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, domain, "Domain created successfully")
}

// UpdateDomain handles PUT /api/admin/domains/:id
func (h *DomainsHandler) UpdateDomain(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateDomainRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST_BODY", err.Error())
	}

	domain, err := h.usecase.UpdateDomain(c.Context(), id, &req)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, domain, "Domain updated successfully")
}

// DeleteDomain handles DELETE /api/admin/domains/:id
func (h *DomainsHandler) DeleteDomain(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.usecase.DeleteDomain(c.Context(), id)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, nil, "Domain deleted successfully")
}

// GetSecurityMetrics handles GET /api/admin/domains/security-metrics
func (h *DomainsHandler) GetSecurityMetrics(c *fiber.Ctx) error {
	metrics, err := h.usecase.GetSecurityMetrics(c.Context())
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, metrics, "Security metrics retrieved successfully")
}

// CreateSecurityMetric handles POST /api/admin/domains/security-metrics
func (h *DomainsHandler) CreateSecurityMetric(c *fiber.Ctx) error {
	var req dto.CreateSecurityMetricRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST_BODY", err.Error())
	}

	err := h.usecase.CreateSecurityMetric(c.Context(), &req)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, nil, "Security metric created successfully")
}

// GetSLABreachAnalysis handles GET /api/admin/domains/sla-breach
func (h *DomainsHandler) GetSLABreachAnalysis(c *fiber.Ctx) error {
	analysis, err := h.usecase.GetSLABreachAnalysis(c.Context())
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	return utils.Success(c, analysis, "SLA breach analysis retrieved successfully")
}

