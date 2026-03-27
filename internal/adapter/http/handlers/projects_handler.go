package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/pkg/utils"
)

type ProjectsHandler struct {
	usecase domain.ProjectsUseCase
}

func NewProjectsHandler(usecase domain.ProjectsUseCase) *ProjectsHandler {
	return &ProjectsHandler{
		usecase: usecase,
	}
}

// CreateProject - UC28: Create New Project
func (h *ProjectsHandler) CreateProject(c *fiber.Ctx) error {
	var req dto.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	// Get created_by from context (from auth middleware)
	createdBy, ok := c.Locals("user_id").(string)
	if !ok || createdBy == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	response, err := h.usecase.CreateProject(c.Context(), &req, createdBy)
	if err != nil {
		// Log the actual error for debugging
		fmt.Printf("Create project error: %v\n", err)
		// Check for specific errors
		if err.Error() == "a project with this name already exists" {
			return utils.BadRequest(c, err.Error(), nil)
		}
		if err.Error() == "validation failed" {
			return utils.BadRequest(c, err.Error(), nil)
		}
		return utils.InternalServerError(c, "Failed to create project")
	}

	return utils.Success(c, response, "Project created successfully")
}

// UpdateProject - Update an existing project
func (h *ProjectsHandler) UpdateProject(c *fiber.Ctx) error {
	projectID := c.Params("id")
	if projectID == "" {
		return utils.BadRequest(c, "Project ID is required", nil)
	}

	var req dto.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	response, err := h.usecase.UpdateProject(c.Context(), projectID, &req)
	if err != nil {
		fmt.Printf("Update project error: %v\n", err)
		if err.Error() == "project not found" {
			return utils.NotFound(c, "Project not found")
		}
		if err.Error() == "validation failed" {
			return utils.BadRequest(c, err.Error(), nil)
		}
		return utils.InternalServerError(c, "Failed to update project")
	}

	return utils.Success(c, response, "Project updated successfully")
}

// ListProjects - UC31, UC32: List projects with pagination
func (h *ProjectsHandler) ListProjects(c *fiber.Ctx) error {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "5"))

	response, err := h.usecase.ListProjects(c.Context(), page, perPage)
	if err != nil {
		return utils.InternalServerError(c, "Failed to retrieve projects")
	}

	return utils.Success(c, response, "Projects retrieved successfully")
}

// DeleteProject - Delete a project
func (h *ProjectsHandler) DeleteProject(c *fiber.Ctx) error {
	projectID := c.Params("id")
	if projectID == "" {
		return utils.BadRequest(c, "Project ID is required", nil)
	}

	if err := h.usecase.DeleteProject(c.Context(), projectID); err != nil {
		return utils.InternalServerError(c, "Failed to delete project")
	}

	return utils.Success(c, nil, "Project deleted successfully")
}

// GetProjectStats - UC27: Get project statistics
func (h *ProjectsHandler) GetProjectStats(c *fiber.Ctx) error {
	response, err := h.usecase.GetProjectStats(c.Context())
	if err != nil {
		return utils.InternalServerError(c, "Failed to retrieve project statistics")
	}

	return utils.Success(c, response, "Project statistics retrieved successfully")
}

// GetPentesters - Get list of pentesters for dropdown
func (h *ProjectsHandler) GetPentesters(c *fiber.Ctx) error {
	response, err := h.usecase.GetPentesters(c.Context())
	if err != nil {
		return utils.InternalServerError(c, "Failed to retrieve pentesters")
	}

	return utils.Success(c, response, "Pentesters retrieved successfully")
}
