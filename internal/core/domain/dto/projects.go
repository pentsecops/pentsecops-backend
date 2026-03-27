package dto

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// Projects Sub-tab DTOs
// ============================================================================

// CreateProjectRequest - UC28: Create New Project
type CreateProjectRequest struct {
	Name       string    `json:"name" validate:"required,min=2,max=255"`
	Type       string    `json:"type" validate:"required,oneof=web network api mobile"`
	AssignedTo *uuid.UUID `json:"assigned_to,omitempty"`
	Deadline   time.Time `json:"deadline" validate:"required"`
	Scope      string    `json:"scope" validate:"max=5000"`
}

// CreateProjectResponse - UC28: Response after creating project
type CreateProjectResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	AssignedTo   uuid.UUID `json:"assigned_to"`
	AssignedToName string  `json:"assigned_to_name"`
	Deadline     time.Time `json:"deadline"`
	Scope        string    `json:"scope"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// UpdateProjectRequest - Request to update project
type UpdateProjectRequest struct {
	Name       string     `json:"name" validate:"omitempty,min=2,max=255"`
	Type       string     `json:"type" validate:"omitempty,oneof=web network api mobile"`
	AssignedTo *uuid.UUID `json:"assigned_to,omitempty"`
	Deadline   *time.Time `json:"deadline,omitempty"`
	Scope      string     `json:"scope" validate:"max=5000"`
	Status     string     `json:"status" validate:"omitempty,oneof=open in_progress completed"`
}

// UpdateProjectResponse - Response after updating project
type UpdateProjectResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	AssignedTo   uuid.UUID `json:"assigned_to"`
	AssignedToName string  `json:"assigned_to_name"`
	Deadline     time.Time `json:"deadline"`
	Scope        string    `json:"scope"`
	Status       string    `json:"status"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ProjectListItem - UC31: Display project in table
type ProjectListItem struct {
	ID                  uuid.UUID  `json:"id"`
	Name                string     `json:"name"`
	Type                string     `json:"type"`
	Status              string     `json:"status"`
	AssignedTo          *uuid.UUID `json:"assigned_to"`
	AssignedToName      string     `json:"assigned_to_name"`
	Deadline            time.Time  `json:"deadline"`
	VulnerabilityCount  int64      `json:"vulnerability_count"`
	CreatedAt           time.Time  `json:"created_at"`
}

// ListProjectsRequest - UC31, UC32: List projects with pagination
type ListProjectsRequest struct {
	Page    int `json:"page" validate:"min=1"`
	PerPage int `json:"per_page" validate:"min=1,max=100"`
}

// ListProjectsResponse - UC31: Response with projects list
type ListProjectsResponse struct {
	Projects   []ProjectListItem `json:"projects"`
	Pagination PaginationInfo    `json:"pagination"`
}

// ProjectStatsResponse - UC27: Project statistics
type ProjectStatsResponse struct {
	OpenCount       int64 `json:"open_count"`
	InProgressCount int64 `json:"in_progress_count"`
	CompletedCount  int64 `json:"completed_count"`
}

// PentesterOption - For dropdown in Create Project dialog
type PentesterOption struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
}

// GetPentestersResponse - List of pentesters for dropdown
type GetPentestersResponse struct {
	Pentesters []PentesterOption `json:"pentesters"`
}

