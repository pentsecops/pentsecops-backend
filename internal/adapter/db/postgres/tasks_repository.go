package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/adapter/db/postgres/sqlc"
	"github.com/pentsecops/backend/internal/core/domain"
)

type TasksRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewTasksRepository(db *sql.DB) *TasksRepository {
	return &TasksRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// CreateTask creates a new task
func (r *TasksRepository) CreateTask(ctx context.Context, task *domain.CreateTaskParams) (*domain.Task, error) {
	// Parse UUIDs
	id, err := uuid.Parse(task.ID)
	if err != nil {
		return nil, err
	}

	var projectID uuid.NullUUID
	if task.ProjectID != nil {
		projectUUID, err := uuid.Parse(*task.ProjectID)
		if err != nil {
			return nil, err
		}
		projectID = uuid.NullUUID{UUID: projectUUID, Valid: true}
	}

	var assignedTo uuid.NullUUID
	if task.AssignedTo != nil {
		assignedToUUID, err := uuid.Parse(*task.AssignedTo)
		if err != nil {
			return nil, err
		}
		assignedTo = uuid.NullUUID{UUID: assignedToUUID, Valid: true}
	}

	var description sql.NullString
	if task.Description != nil {
		description = sql.NullString{String: *task.Description, Valid: true}
	}

	var priority sql.NullString
	if task.Priority != nil {
		priority = sql.NullString{String: *task.Priority, Valid: true}
	}

	var deadline sql.NullTime
	if task.Deadline != nil {
		deadline = sql.NullTime{Time: *task.Deadline, Valid: true}
	}

	// Create task using sqlc
	t, err := r.queries.CreateTask(ctx, sqlc.CreateTaskParams{
		ID:          id,
		ProjectID:   projectID,
		Title:       task.Title,
		Description: description,
		Status:      task.Status,
		Priority:    priority,
		AssignedTo:  assignedTo,
		Deadline:    deadline,
		CreatedAt:   sql.NullTime{Time: task.CreatedAt, Valid: true},
		UpdatedAt:   sql.NullTime{Time: task.UpdatedAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// Convert to domain model
	result := &domain.Task{
		ID:        t.ID,
		Title:     t.Title,
		Status:    t.Status,
		CreatedAt: t.CreatedAt.Time,
		UpdatedAt: t.UpdatedAt.Time,
	}

	if t.ProjectID.Valid {
		result.ProjectID = &t.ProjectID.UUID
	}

	if t.Description.Valid {
		result.Description = &t.Description.String
	}

	if t.Priority.Valid {
		result.Priority = &t.Priority.String
	}

	if t.AssignedTo.Valid {
		result.AssignedTo = &t.AssignedTo.UUID
	}

	if t.Deadline.Valid {
		result.Deadline = &t.Deadline.Time
	}

	return result, nil
}

// GetTaskByID retrieves a task by ID
func (r *TasksRepository) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	taskID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	t, err := r.queries.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	result := &domain.Task{
		ID:        t.ID,
		Title:     t.Title,
		Status:    t.Status,
		CreatedAt: t.CreatedAt.Time,
		UpdatedAt: t.UpdatedAt.Time,
	}

	if t.ProjectID.Valid {
		result.ProjectID = &t.ProjectID.UUID
	}

	if t.Description.Valid {
		result.Description = &t.Description.String
	}

	if t.Priority.Valid {
		result.Priority = &t.Priority.String
	}

	if t.AssignedTo.Valid {
		result.AssignedTo = &t.AssignedTo.UUID
	}

	if t.Deadline.Valid {
		result.Deadline = &t.Deadline.Time
	}

	return result, nil
}

// DeleteTask deletes a task
func (r *TasksRepository) DeleteTask(ctx context.Context, id string) error {
	taskID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteTask(ctx, taskID)
}

// UpdateTaskStatus updates a task's status
func (r *TasksRepository) UpdateTaskStatus(ctx context.Context, id string, status string) error {
	taskID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.UpdateTaskStatus(ctx, sqlc.UpdateTaskStatusParams{
		ID:     taskID,
		Status: status,
	})
}

// ListTasksByProject retrieves tasks for a specific project
func (r *TasksRepository) ListTasksByProject(ctx context.Context, projectID string) ([]*domain.TaskWithDetails, error) {
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.ListTasksByProject(ctx, uuid.NullUUID{UUID: projectUUID, Valid: true})
	if err != nil {
		return nil, err
	}

	tasks := make([]*domain.TaskWithDetails, 0, len(rows))
	for _, row := range rows {
		task := &domain.TaskWithDetails{
			ID:             row.ID.String(),
			Title:          row.Title,
			Status:         row.Status,
			AssignedToName: row.AssignedToName.String,
			CreatedAt:      row.CreatedAt.Time,
		}

		if row.ProjectID.Valid {
			projectIDStr := row.ProjectID.UUID.String()
			task.ProjectID = &projectIDStr
		}

		if row.Description.Valid {
			task.Description = &row.Description.String
		}

		if row.Priority.Valid {
			task.Priority = &row.Priority.String
		}

		if row.AssignedTo.Valid {
			assignedToStr := row.AssignedTo.UUID.String()
			task.AssignedTo = &assignedToStr
		}

		if row.Deadline.Valid {
			task.Deadline = &row.Deadline.Time
		}

		if row.CompletedAt.Valid {
			task.CompletedAt = &row.CompletedAt.Time
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// ListAllTasks retrieves all tasks
func (r *TasksRepository) ListAllTasks(ctx context.Context) ([]*domain.TaskWithDetails, error) {
	rows, err := r.queries.ListAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	tasks := make([]*domain.TaskWithDetails, 0, len(rows))
	for _, row := range rows {
		task := &domain.TaskWithDetails{
			ID:             row.ID.String(),
			Title:          row.Title,
			Status:         row.Status,
			ProjectName:    row.ProjectName.String,
			AssignedToName: row.AssignedToName.String,
			CreatedAt:      row.CreatedAt.Time,
		}

		if row.ProjectID.Valid {
			projectIDStr := row.ProjectID.UUID.String()
			task.ProjectID = &projectIDStr
		}

		if row.Description.Valid {
			task.Description = &row.Description.String
		}

		if row.Priority.Valid {
			task.Priority = &row.Priority.String
		}

		if row.AssignedTo.Valid {
			assignedToStr := row.AssignedTo.UUID.String()
			task.AssignedTo = &assignedToStr
		}

		if row.Deadline.Valid {
			task.Deadline = &row.Deadline.Time
		}

		if row.CompletedAt.Valid {
			task.CompletedAt = &row.CompletedAt.Time
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTasksByStatus retrieves tasks by status
func (r *TasksRepository) GetTasksByStatus(ctx context.Context, status string) ([]*domain.TaskWithDetails, error) {
	rows, err := r.queries.GetTasksByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	tasks := make([]*domain.TaskWithDetails, 0, len(rows))
	for _, row := range rows {
		task := &domain.TaskWithDetails{
			ID:             row.ID.String(),
			Title:          row.Title,
			Status:         row.Status,
			ProjectName:    row.ProjectName.String,
			AssignedToName: row.AssignedToName.String,
			CreatedAt:      row.CreatedAt.Time,
		}

		if row.ProjectID.Valid {
			projectIDStr := row.ProjectID.UUID.String()
			task.ProjectID = &projectIDStr
		}

		if row.Description.Valid {
			task.Description = &row.Description.String
		}

		if row.Priority.Valid {
			task.Priority = &row.Priority.String
		}

		if row.AssignedTo.Valid {
			assignedToStr := row.AssignedTo.UUID.String()
			task.AssignedTo = &assignedToStr
		}

		if row.Deadline.Valid {
			task.Deadline = &row.Deadline.Time
		}

		if row.CompletedAt.Valid {
			task.CompletedAt = &row.CompletedAt.Time
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
