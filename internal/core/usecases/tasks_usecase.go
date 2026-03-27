package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
)

type TasksUseCase struct {
	repo      domain.TasksRepository
	validator *validator.Validate
}

func NewTasksUseCase(repo domain.TasksRepository) *TasksUseCase {
	return &TasksUseCase{
		repo:      repo,
		validator: validator.New(),
	}
}

// CreateTask - UC37: Create New Task
func (uc *TasksUseCase) CreateTask(ctx context.Context, req *dto.CreateTaskRequest) (*dto.CreateTaskResponse, error) {
	// Validate required fields
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate UUIDv7 for task
	taskID := uuid.Must(uuid.NewV7())
	now := time.Now()

	// Prepare create params
	projectIDStr := req.ProjectID.String()
	assignedToStr := req.AssignedTo.String()
	var descriptionPtr *string
	if req.Description != "" {
		descriptionPtr = &req.Description
	}
	priorityStr := req.Priority

	params := &domain.CreateTaskParams{
		ID:          taskID.String(),
		ProjectID:   &projectIDStr,
		Title:       req.Title,
		Description: descriptionPtr,
		Status:      domain.TaskStatusToDo, // Default status is "to_do"
		Priority:    &priorityStr,
		AssignedTo:  &assignedToStr,
		Deadline:    req.Deadline,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Create task in database
	task, err := uc.repo.CreateTask(ctx, params)
	if err != nil {
		fmt.Printf("CreateTask Repository Error: %v\n", err)
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Build response (we'll need to get project name and assigned to name)
	// For now, return basic info
	response := &dto.CreateTaskResponse{
		ID:          taskID,
		ProjectID:   req.ProjectID,
		Title:       task.Title,
		Description: req.Description,
		Status:      task.Status,
		Priority:    req.Priority,
		AssignedTo:  req.AssignedTo,
		Deadline:    req.Deadline,
		CreatedAt:   task.CreatedAt,
	}

	return response, nil
}

// ListTasksByProject - UC36: Display tasks for a specific project
func (uc *TasksUseCase) ListTasksByProject(ctx context.Context, projectID string) (*dto.ListTasksResponse, error) {
	// Validate UUID format
	if _, err := uuid.Parse(projectID); err != nil {
		return nil, fmt.Errorf("invalid project ID format")
	}

	// Get tasks from repository
	tasks, err := uc.repo.ListTasksByProject(ctx, projectID)
	if err != nil {
		fmt.Printf("ListTasksByProject Error: %v\n", err)
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	// Convert to DTOs
	taskItems := make([]dto.TaskListItem, 0, len(tasks))
	for _, t := range tasks {
		item := dto.TaskListItem{
			ID:             uuid.MustParse(t.ID),
			Title:          t.Title,
			Status:         t.Status,
			AssignedToName: t.AssignedToName,
			CreatedAt:      t.CreatedAt,
		}

		if t.ProjectID != nil {
			projectUUID := uuid.MustParse(*t.ProjectID)
			item.ProjectID = &projectUUID
		}

		if t.Description != nil {
			item.Description = *t.Description
		}

		if t.Priority != nil {
			item.Priority = *t.Priority
		}

		if t.AssignedTo != nil {
			assignedToUUID := uuid.MustParse(*t.AssignedTo)
			item.AssignedTo = &assignedToUUID
		}

		item.Deadline = t.Deadline
		item.CompletedAt = t.CompletedAt

		taskItems = append(taskItems, item)
	}

	return &dto.ListTasksResponse{
		Tasks: taskItems,
	}, nil
}

// ListAllTasks - UC36: Display Task Board Interface
func (uc *TasksUseCase) ListAllTasks(ctx context.Context) (*dto.TaskBoardResponse, error) {
	// Get all tasks
	allTasks, err := uc.repo.ListAllTasks(ctx)
	if err != nil {
		fmt.Printf("ListAllTasks Error: %v\n", err)
		return nil, fmt.Errorf("failed to list all tasks: %w", err)
	}

	// Group tasks by status
	toDo := make([]dto.TaskListItem, 0)
	inProgress := make([]dto.TaskListItem, 0)
	done := make([]dto.TaskListItem, 0)

	for _, t := range allTasks {
		item := dto.TaskListItem{
			ID:             uuid.MustParse(t.ID),
			ProjectName:    t.ProjectName,
			Title:          t.Title,
			Status:         t.Status,
			AssignedToName: t.AssignedToName,
			CreatedAt:      t.CreatedAt,
		}

		if t.ProjectID != nil {
			projectUUID := uuid.MustParse(*t.ProjectID)
			item.ProjectID = &projectUUID
		}

		if t.Description != nil {
			item.Description = *t.Description
		}

		if t.Priority != nil {
			item.Priority = *t.Priority
		}

		if t.AssignedTo != nil {
			assignedToUUID := uuid.MustParse(*t.AssignedTo)
			item.AssignedTo = &assignedToUUID
		}

		item.Deadline = t.Deadline
		item.CompletedAt = t.CompletedAt

		// Group by status
		switch t.Status {
		case domain.TaskStatusToDo:
			toDo = append(toDo, item)
		case domain.TaskStatusInProgress:
			inProgress = append(inProgress, item)
		case domain.TaskStatusDone:
			done = append(done, item)
		}
	}

	return &dto.TaskBoardResponse{
		ToDo:       toDo,
		InProgress: inProgress,
		Done:       done,
	}, nil
}

// UpdateTaskStatus - UC38: Update Task Status
func (uc *TasksUseCase) UpdateTaskStatus(ctx context.Context, taskID string, req *dto.UpdateTaskStatusRequest) (*dto.UpdateTaskStatusResponse, error) {
	// Validate UUID format
	if _, err := uuid.Parse(taskID); err != nil {
		return nil, fmt.Errorf("invalid task ID format")
	}

	// Validate request
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update task status
	if err := uc.repo.UpdateTaskStatus(ctx, taskID, req.Status); err != nil {
		fmt.Printf("UpdateTaskStatus Error: %v\n", err)
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	return &dto.UpdateTaskStatusResponse{
		ID:        uuid.MustParse(taskID),
		Status:    req.Status,
		UpdatedAt: time.Now(),
	}, nil
}

// DeleteTask - Delete a task
func (uc *TasksUseCase) DeleteTask(ctx context.Context, taskID string) error {
	// Validate UUID format
	if _, err := uuid.Parse(taskID); err != nil {
		return fmt.Errorf("invalid task ID format")
	}

	// Delete task
	if err := uc.repo.DeleteTask(ctx, taskID); err != nil {
		fmt.Printf("DeleteTask Error: %v\n", err)
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

