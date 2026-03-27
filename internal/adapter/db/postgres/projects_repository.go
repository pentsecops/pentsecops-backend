package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/adapter/db/postgres/sqlc"
	"github.com/pentsecops/backend/internal/core/domain"
)

type ProjectsRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewProjectsRepository(db *sql.DB) *ProjectsRepository {
	return &ProjectsRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// CreateProject creates a new project
func (r *ProjectsRepository) CreateProject(ctx context.Context, project *domain.CreateProjectParams) (*domain.Project, error) {
	// Parse UUIDs
	id, err := uuid.Parse(project.ID)
	if err != nil {
		return nil, err
	}

	var assignedTo uuid.NullUUID
	if project.AssignedTo != nil {
		assignedToUUID, err := uuid.Parse(*project.AssignedTo)
		if err != nil {
			return nil, err
		}
		assignedTo = uuid.NullUUID{UUID: assignedToUUID, Valid: true}
	}

	var createdBy uuid.NullUUID
	if project.CreatedBy != nil {
		createdByUUID, err := uuid.Parse(*project.CreatedBy)
		if err != nil {
			return nil, err
		}
		createdBy = uuid.NullUUID{UUID: createdByUUID, Valid: true}
	} else {
		createdBy = uuid.NullUUID{Valid: false}
	}

	var scope sql.NullString
	if project.Scope != nil {
		scope = sql.NullString{String: *project.Scope, Valid: true}
	}

	// Create project using sqlc
	p, err := r.queries.CreateProject(ctx, sqlc.CreateProjectParams{
		ID:           id,
		Name:         project.Name,
		Type:         project.Type,
		AssignedTo:   assignedTo,
		Deadline:     project.Deadline,
		Scope:        scope,
		Status:       project.Status,
		CurrentPhase: sql.NullString{String: "pending", Valid: true},
		CreatedBy:    createdBy,
		CreatedAt:    sql.NullTime{Time: project.CreatedAt, Valid: true},
		UpdatedAt:    sql.NullTime{Time: project.UpdatedAt, Valid: true},
	})
	if err != nil {
		fmt.Printf("SQL Create Project Error: %v\n", err)
		return nil, err
	}

	// Convert to domain model
	result := &domain.Project{
		ID:        p.ID,
		Name:      p.Name,
		Type:      p.Type,
		Deadline:  &p.Deadline,
		Status:    p.Status,
		CreatedAt: p.CreatedAt.Time,
		UpdatedAt: p.UpdatedAt.Time,
	}

	if p.AssignedTo.Valid {
		result.AssignedTo = &p.AssignedTo.UUID
	}

	if p.Scope.Valid {
		result.Scope = &p.Scope.String
	}

	return result, nil
}

// GetProjectByID retrieves a project by ID
func (r *ProjectsRepository) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	projectID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	p, err := r.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	result := &domain.Project{
		ID:        p.ID,
		Name:      p.Name,
		Type:      p.Type,
		Deadline:  &p.Deadline,
		Status:    p.Status,
		CreatedAt: p.CreatedAt.Time,
		UpdatedAt: p.UpdatedAt.Time,
	}

	if p.AssignedTo.Valid {
		result.AssignedTo = &p.AssignedTo.UUID
	}

	if p.Scope.Valid {
		result.Scope = &p.Scope.String
	}

	return result, nil
}

// GetProjectByName retrieves a project by name
func (r *ProjectsRepository) GetProjectByName(ctx context.Context, name string) (*domain.Project, error) {
	p, err := r.queries.GetProjectByName(ctx, name)
	if err != nil {
		return nil, err
	}

	result := &domain.Project{
		ID:        p.ID,
		Name:      p.Name,
		Type:      p.Type,
		Deadline:  &p.Deadline,
		Status:    p.Status,
		CreatedAt: p.CreatedAt.Time,
		UpdatedAt: p.UpdatedAt.Time,
	}

	if p.AssignedTo.Valid {
		result.AssignedTo = &p.AssignedTo.UUID
	}

	if p.Scope.Valid {
		result.Scope = &p.Scope.String
	}

	return result, nil
}

// DeleteProject deletes a project
func (r *ProjectsRepository) DeleteProject(ctx context.Context, id string) error {
	projectID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteProject(ctx, projectID)
}

// UpdateProjectStatus updates a project's status
func (r *ProjectsRepository) UpdateProjectStatus(ctx context.Context, id string, status string) error {
	projectID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.UpdateProjectStatus(ctx, sqlc.UpdateProjectStatusParams{
		ID:     projectID,
		Status: status,
	})
}

// ListProjects retrieves projects with pagination
func (r *ProjectsRepository) ListProjects(ctx context.Context, limit, offset int) ([]*domain.ProjectWithDetails, error) {
	rows, err := r.queries.ListProjects(ctx, sqlc.ListProjectsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	projects := make([]*domain.ProjectWithDetails, 0, len(rows))
	for _, row := range rows {
		project := &domain.ProjectWithDetails{
			ID:                 row.ID.String(),
			Name:               row.Name,
			Type:               row.Type,
			Status:             row.Status,
			AssignedToName:     row.AssignedToName.String,
			Deadline:           row.Deadline,
			VulnerabilityCount: row.VulnerabilityCount,
			CreatedAt:          row.CreatedAt.Time,
		}

		if row.AssignedTo.Valid {
			assignedToStr := row.AssignedTo.UUID.String()
			project.AssignedTo = &assignedToStr
		}

		projects = append(projects, project)
	}

	return projects, nil
}

// CountProjects returns the total number of projects
func (r *ProjectsRepository) CountProjects(ctx context.Context) (int64, error) {
	return r.queries.CountProjects(ctx)
}

// GetProjectStats retrieves project statistics
func (r *ProjectsRepository) GetProjectStats(ctx context.Context) (*domain.ProjectStats, error) {
	stats, err := r.queries.GetProjectStats(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.ProjectStats{
		OpenCount:       stats.OpenCount,
		InProgressCount: stats.InProgressCount,
		CompletedCount:  stats.CompletedCount,
	}, nil
}

// GetPentesters retrieves all active pentesters
func (r *ProjectsRepository) GetPentesters(ctx context.Context) ([]*domain.Pentester, error) {
	rows, err := r.queries.GetPentesters(ctx)
	if err != nil {
		return nil, err
	}

	pentesters := make([]*domain.Pentester, 0, len(rows))
	for _, row := range rows {
		pentesters = append(pentesters, &domain.Pentester{
			ID:       row.ID.String(),
			FullName: row.FullName,
			Email:    row.Email,
		})
	}

	return pentesters, nil
}

// UpdateProject updates an existing project
func (r *ProjectsRepository) UpdateProject(ctx context.Context, id string, params *domain.UpdateProjectParams) (*domain.Project, error) {
	projectID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var assignedTo uuid.NullUUID
	if params.AssignedTo != nil {
		assignedTo = uuid.NullUUID{UUID: *params.AssignedTo, Valid: true}
	}

	var scope sql.NullString
	if params.Scope != nil {
		scope = sql.NullString{String: *params.Scope, Valid: true}
	}

	// Update project using raw SQL
	query := `
		UPDATE projects 
		SET name = $2, type = $3, assigned_to = $4, deadline = $5, scope = $6, status = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, name, type, assigned_to, deadline, scope, status, current_phase, created_by, created_at, updated_at`

	var p sqlc.Project
	err = r.db.QueryRowContext(ctx, query, projectID, params.Name, params.Type, assignedTo, params.Deadline, scope, params.Status).Scan(
		&p.ID, &p.Name, &p.Type, &p.AssignedTo, &p.Deadline, &p.Scope, &p.Status, &p.CurrentPhase, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		fmt.Printf("SQL Update Error: %v\n", err)
		return nil, err
	}

	// Convert to domain model
	result := &domain.Project{
		ID:        p.ID,
		Name:      p.Name,
		Type:      p.Type,
		Deadline:  &p.Deadline,
		Status:    p.Status,
		CreatedAt: p.CreatedAt.Time,
		UpdatedAt: p.UpdatedAt.Time,
	}

	if p.AssignedTo.Valid {
		result.AssignedTo = &p.AssignedTo.UUID
	}

	if p.Scope.Valid {
		result.Scope = &p.Scope.String
	}

	return result, nil
}
