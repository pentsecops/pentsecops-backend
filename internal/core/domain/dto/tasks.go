package dto

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// Task Board Sub-tab DTOs
// ============================================================================

// CreateTaskRequest - UC37: Create New Task
type CreateTaskRequest struct {
	ProjectID   uuid.UUID  `json:"project_id" validate:"required"`
	Title       string     `json:"title" validate:"required,min=2,max=255"`
	Description string     `json:"description" validate:"max=5000"`
	Priority    string     `json:"priority" validate:"required,oneof=low medium high critical"`
	AssignedTo  uuid.UUID  `json:"assigned_to" validate:"required"`
	Deadline    *time.Time `json:"deadline"`
}

// CreateTaskResponse - UC37: Response after creating task
type CreateTaskResponse struct {
	ID             uuid.UUID  `json:"id"`
	ProjectID      uuid.UUID  `json:"project_id"`
	ProjectName    string     `json:"project_name"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	AssignedTo     uuid.UUID  `json:"assigned_to"`
	AssignedToName string     `json:"assigned_to_name"`
	Deadline       *time.Time `json:"deadline"`
	CreatedAt      time.Time  `json:"created_at"`
}

// TaskListItem - UC36: Display task in task board
type TaskListItem struct {
	ID             uuid.UUID  `json:"id"`
	ProjectID      *uuid.UUID `json:"project_id"`
	ProjectName    string     `json:"project_name"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	AssignedTo     *uuid.UUID `json:"assigned_to"`
	AssignedToName string     `json:"assigned_to_name"`
	Deadline       *time.Time `json:"deadline"`
	CompletedAt    *time.Time `json:"completed_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

// ListTasksResponse - UC36: Response with tasks list
type ListTasksResponse struct {
	Tasks []TaskListItem `json:"tasks"`
}

// ListTasksByProjectRequest - UC36: Get tasks for specific project
type ListTasksByProjectRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
}

// UpdateTaskStatusRequest - UC38: Update task status
type UpdateTaskStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=to_do in_progress done"`
}

// UpdateTaskStatusResponse - UC38: Response after updating status
type UpdateTaskStatusResponse struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskBoardResponse - UC36: Task board with tasks grouped by status
type TaskBoardResponse struct {
	ToDo       []TaskListItem `json:"to_do"`
	InProgress []TaskListItem `json:"in_progress"`
	Done       []TaskListItem `json:"done"`
}

