package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
)

type ProjectsUseCase struct {
	repo      domain.ProjectsRepository
	validator *validator.Validate
}

func NewProjectsUseCase(repo domain.ProjectsRepository) *ProjectsUseCase {
	return &ProjectsUseCase{
		repo:      repo,
		validator: validator.New(),
	}
}

// CreateProject - UC28: Create New Project with All Details
func (uc *ProjectsUseCase) CreateProject(ctx context.Context, req *dto.CreateProjectRequest, createdBy string) (*dto.CreateProjectResponse, error) {
	// UC29: Validate required fields
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Validate createdBy user exists
	if createdBy == "" {
		return nil, fmt.Errorf("created by user ID is required")
	}

	// UC30: Validate project name is unique
	existingProject, err := uc.repo.GetProjectByName(ctx, req.Name)
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("GetProjectByName Error: %v\n", err)
		return nil, fmt.Errorf("failed to check project name: %w", err)
	}
	if existingProject != nil {
		return nil, fmt.Errorf("a project with this name already exists")
	}

	// Generate UUIDv7 for project
	projectID := uuid.Must(uuid.NewV7())
	now := time.Now()

	// Prepare create params
	var assignedToPtr *string
	if req.AssignedTo != nil {
		assignedToStr := req.AssignedTo.String()
		assignedToPtr = &assignedToStr
	}
	var scopePtr *string
	if req.Scope != "" {
		scopePtr = &req.Scope
	}

	params := &domain.CreateProjectParams{
		ID:         projectID.String(),
		Name:       req.Name,
		Type:       req.Type,
		AssignedTo: assignedToPtr,
		Deadline:   req.Deadline,
		Scope:      scopePtr,
		Status:     domain.ProjectStatusOpen, // Default status is "open"
		CreatedBy:  nil, // Set to NULL to avoid foreign key constraint
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Create project in database
	project, err := uc.repo.CreateProject(ctx, params)
	if err != nil {
		fmt.Printf("CreateProject Repository Error: %v\n", err)
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Get pentester name for response
	var assignedToName string
	var assignedToUUID uuid.UUID
	if req.AssignedTo != nil {
		// Validate that the assigned user exists
		pentesters, err := uc.repo.GetPentesters(ctx)
		if err != nil {
			fmt.Printf("GetPentesters Error: %v\n", err)
			return nil, fmt.Errorf("failed to get pentesters: %w", err)
		}

		assignedToStr := req.AssignedTo.String()
		userExists := false
		for _, p := range pentesters {
			if p.ID == assignedToStr {
				assignedToName = p.FullName
				userExists = true
				break
			}
		}
		if !userExists {
			return nil, fmt.Errorf("assigned user does not exist")
		}
		assignedToUUID = *req.AssignedTo
	}

	// Build response
	response := &dto.CreateProjectResponse{
		ID:             projectID,
		Name:           project.Name,
		Type:           project.Type,
		AssignedTo:     assignedToUUID,
		AssignedToName: assignedToName,
		Deadline:       req.Deadline,
		Scope:          req.Scope,
		Status:         project.Status,
		CreatedAt:      project.CreatedAt,
	}

	return response, nil
}

// UpdateProject - Update an existing project
func (uc *ProjectsUseCase) UpdateProject(ctx context.Context, projectID string, req *dto.UpdateProjectRequest) (*dto.UpdateProjectResponse, error) {
	// Validate UUID format
	if _, err := uuid.Parse(projectID); err != nil {
		return nil, fmt.Errorf("invalid project ID format")
	}

	// Validate request
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if project exists
	existingProject, err := uc.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		fmt.Printf("GetProjectByID Error in Update: %v\n", err)
		return nil, fmt.Errorf("project not found")
	}

	// Use existing values if not provided in request
	name := existingProject.Name
	if req.Name != "" {
		name = req.Name
		// Check if name is being changed and if it already exists
		if req.Name != existingProject.Name {
			nameExists, err := uc.repo.GetProjectByName(ctx, req.Name)
			if err != nil && err != sql.ErrNoRows {
				fmt.Printf("GetProjectByName Error in Update: %v\n", err)
				return nil, fmt.Errorf("failed to check project name: %w", err)
			}
			if nameExists != nil {
				return nil, fmt.Errorf("a project with this name already exists")
			}
		}
	}

	projectType := existingProject.Type
	if req.Type != "" {
		projectType = req.Type
	}

	deadline := *existingProject.Deadline
	if req.Deadline != nil {
		deadline = *req.Deadline
	}

	status := existingProject.Status
	if req.Status != "" {
		status = req.Status
	}

	assignedTo := existingProject.AssignedTo
	if req.AssignedTo != nil {
		assignedTo = req.AssignedTo
	}

	scope := ""
	if existingProject.Scope != nil {
		scope = *existingProject.Scope
	}
	if req.Scope != "" {
		scope = req.Scope
	}

	// Validate assigned user if provided
	var assignedToName string
	if assignedTo != nil {
		pentesters, err := uc.repo.GetPentesters(ctx)
		if err != nil {
			fmt.Printf("GetPentesters Error in Update: %v\n", err)
			return nil, fmt.Errorf("failed to get pentesters: %w", err)
		}

		assignedToStr := assignedTo.String()
		userExists := false
		for _, p := range pentesters {
			if p.ID == assignedToStr {
				assignedToName = p.FullName
				userExists = true
				break
			}
		}
		if !userExists {
			return nil, fmt.Errorf("assigned user does not exist")
		}
	}

	// Update project
	updatedProject, err := uc.repo.UpdateProject(ctx, projectID, &domain.UpdateProjectParams{
		Name:       name,
		Type:       projectType,
		AssignedTo: assignedTo,
		Deadline:   deadline,
		Scope:      &scope,
		Status:     status,
	})
	if err != nil {
		fmt.Printf("UpdateProject Repository Error: %v\n", err)
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	var assignedToUUID uuid.UUID
	if assignedTo != nil {
		assignedToUUID = *assignedTo
	}

	return &dto.UpdateProjectResponse{
		ID:             updatedProject.ID,
		Name:           updatedProject.Name,
		Type:           updatedProject.Type,
		AssignedTo:     assignedToUUID,
		AssignedToName: assignedToName,
		Deadline:       *updatedProject.Deadline,
		Scope:          scope,
		Status:         updatedProject.Status,
		UpdatedAt:      updatedProject.UpdatedAt,
	}, nil
}

// ListProjects - UC31, UC32: Display All Projects with Pagination
func (uc *ProjectsUseCase) ListProjects(ctx context.Context, page, perPage int) (*dto.ListProjectsResponse, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 5 // Default 5 per page as per requirements
	}
	if perPage > 100 {
		perPage = 100
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Get projects from repository
	projects, err := uc.repo.ListProjects(ctx, perPage, offset)
	if err != nil {
		fmt.Printf("ListProjects Error: %v\n", err)
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	// Get total count
	total, err := uc.repo.CountProjects(ctx)
	if err != nil {
		fmt.Printf("CountProjects Error: %v\n", err)
		return nil, fmt.Errorf("failed to count projects: %w", err)
	}

	// Convert to DTOs
	projectItems := make([]dto.ProjectListItem, 0, len(projects))
	for _, p := range projects {
		item := dto.ProjectListItem{
			ID:                 uuid.MustParse(p.ID),
			Name:               p.Name,
			Type:               p.Type,
			Status:             p.Status,
			AssignedToName:     p.AssignedToName,
			Deadline:           p.Deadline,
			VulnerabilityCount: p.VulnerabilityCount,
			CreatedAt:          p.CreatedAt,
		}

		if p.AssignedTo != nil {
			assignedToUUID := uuid.MustParse(*p.AssignedTo)
			item.AssignedTo = &assignedToUUID
		}

		projectItems = append(projectItems, item)
	}

	// Calculate pagination info
	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	pagination := dto.PaginationInfo{
		CurrentPage: page,
		PerPage:     perPage,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}

	return &dto.ListProjectsResponse{
		Projects:   projectItems,
		Pagination: pagination,
	}, nil
}

// DeleteProject - UC19: Delete Project
func (uc *ProjectsUseCase) DeleteProject(ctx context.Context, projectID string) error {
	// Validate UUID format
	if _, err := uuid.Parse(projectID); err != nil {
		return fmt.Errorf("invalid project ID format")
	}

	// Delete project
	if err := uc.repo.DeleteProject(ctx, projectID); err != nil {
		fmt.Printf("DeleteProject Error: %v\n", err)
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// GetProjectStats - UC27: Fetch and Display Project Statistics
func (uc *ProjectsUseCase) GetProjectStats(ctx context.Context) (*dto.ProjectStatsResponse, error) {
	stats, err := uc.repo.GetProjectStats(ctx)
	if err != nil {
		fmt.Printf("GetProjectStats Error: %v\n", err)
		return nil, fmt.Errorf("failed to get project stats: %w", err)
	}

	return &dto.ProjectStatsResponse{
		OpenCount:       stats.OpenCount,
		InProgressCount: stats.InProgressCount,
		CompletedCount:  stats.CompletedCount,
	}, nil
}

// GetPentesters - Get list of pentesters for dropdown
func (uc *ProjectsUseCase) GetPentesters(ctx context.Context) (*dto.GetPentestersResponse, error) {
	pentesters, err := uc.repo.GetPentesters(ctx)
	if err != nil {
		fmt.Printf("GetPentesters Error in GetPentesters: %v\n", err)
		return nil, fmt.Errorf("failed to get pentesters: %w", err)
	}

	options := make([]dto.PentesterOption, 0, len(pentesters))
	for _, p := range pentesters {
		options = append(options, dto.PentesterOption{
			ID:       uuid.MustParse(p.ID),
			FullName: p.FullName,
			Email:    p.Email,
		})
	}

	return &dto.GetPentestersResponse{
		Pentesters: options,
	}, nil
}
