package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/adapter/db/postgres/sqlc"
	"github.com/pentsecops/backend/internal/core/domain"
)

type stakeholderVulnerabilitiesRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewStakeholderVulnerabilitiesRepository creates a new stakeholder vulnerabilities repository
func NewStakeholderVulnerabilitiesRepository(db *sql.DB) domain.StakeholderVulnerabilitiesRepository {
	return &stakeholderVulnerabilitiesRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// UC11: Get critical vulnerabilities count
func (r *stakeholderVulnerabilitiesRepository) GetCriticalVulnerabilitiesCount(ctx context.Context) (int, error) {
	count, err := r.queries.GetStakeholderCriticalVulnerabilitiesCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get critical vulnerabilities count: %w", err)
	}
	return int(count), nil
}

// UC12: Get high severity vulnerabilities count
func (r *stakeholderVulnerabilitiesRepository) GetHighSeverityVulnerabilitiesCount(ctx context.Context) (int, error) {
	count, err := r.queries.GetStakeholderHighSeverityVulnerabilitiesCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get high severity vulnerabilities count: %w", err)
	}
	return int(count), nil
}

// UC13: Get open issues count
func (r *stakeholderVulnerabilitiesRepository) GetOpenIssuesCount(ctx context.Context) (int, error) {
	count, err := r.queries.GetStakeholderOpenIssuesCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get open issues count: %w", err)
	}
	return int(count), nil
}

// UC14: Get remediation count
func (r *stakeholderVulnerabilitiesRepository) GetRemediationCount(ctx context.Context) (int, error) {
	count, err := r.queries.GetStakeholderRemediationCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get remediation count: %w", err)
	}
	return int(count), nil
}

// UC19: List vulnerabilities with search, filters, and pagination
func (r *stakeholderVulnerabilitiesRepository) ListVulnerabilities(ctx context.Context, search, severity, status string, limit, offset int) ([]domain.VulnerabilityListItem, error) {
	// Prepare nullable parameters
	var searchParam, severityParam, statusParam sql.NullString

	if search != "" {
		searchParam = sql.NullString{String: search, Valid: true}
	}
	if severity != "" && severity != "all" {
		severityParam = sql.NullString{String: severity, Valid: true}
	}
	if status != "" && status != "all" {
		statusParam = sql.NullString{String: status, Valid: true}
	}

	rows, err := r.queries.ListStakeholderVulnerabilities(ctx, sqlc.ListStakeholderVulnerabilitiesParams{
		Limit:    int32(limit),
		Offset:   int32(offset),
		Search:   searchParam,
		Severity: severityParam,
		Status:   statusParam,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list vulnerabilities: %w", err)
	}

	// Preallocate slice for performance
	vulnerabilities := make([]domain.VulnerabilityListItem, 0, len(rows))
	for _, row := range rows {
		// Convert discovered_date (sql.NullTime) to time.Time
		var discoveredDate time.Time
		if row.DiscoveredDate.Valid {
			discoveredDate = row.DiscoveredDate.Time
		}

		// Convert assigned_to (sql.NullString) to string
		assignedTo := ""
		if row.AssignedTo.Valid {
			assignedTo = row.AssignedTo.String
		}

		vulnerabilities = append(vulnerabilities, domain.VulnerabilityListItem{
			ID:             row.ID.String(),
			Title:          row.Title,
			Severity:       row.Severity,
			Domain:         row.Domain,
			Status:         row.Status,
			DiscoveredDate: discoveredDate,
			DueDate:        &row.DueDate,
			AssignedTo:     assignedTo,
		})
	}

	return vulnerabilities, nil
}

// UC19: Get total count for pagination
func (r *stakeholderVulnerabilitiesRepository) GetVulnerabilitiesCount(ctx context.Context, search, severity, status string) (int, error) {
	// Prepare nullable parameters
	var searchParam, severityParam, statusParam sql.NullString

	if search != "" {
		searchParam = sql.NullString{String: search, Valid: true}
	}
	if severity != "" && severity != "all" {
		severityParam = sql.NullString{String: severity, Valid: true}
	}
	if status != "" && status != "all" {
		statusParam = sql.NullString{String: status, Valid: true}
	}

	count, err := r.queries.GetStakeholderVulnerabilitiesCount(ctx, sqlc.GetStakeholderVulnerabilitiesCountParams{
		Search:   searchParam,
		Severity: severityParam,
		Status:   statusParam,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get vulnerabilities count: %w", err)
	}

	return int(count), nil
}

// UC24: Export vulnerabilities to CSV
func (r *stakeholderVulnerabilitiesRepository) ExportVulnerabilities(ctx context.Context, search, severity, status string) ([]domain.VulnerabilityExportItem, error) {
	// Prepare nullable parameters
	var searchParam, severityParam, statusParam sql.NullString

	if search != "" {
		searchParam = sql.NullString{String: search, Valid: true}
	}
	if severity != "" && severity != "all" {
		severityParam = sql.NullString{String: severity, Valid: true}
	}
	if status != "" && status != "all" {
		statusParam = sql.NullString{String: status, Valid: true}
	}

	rows, err := r.queries.ExportStakeholderVulnerabilities(ctx, sqlc.ExportStakeholderVulnerabilitiesParams{
		Search:   searchParam,
		Severity: severityParam,
		Status:   statusParam,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to export vulnerabilities: %w", err)
	}

	// Preallocate slice for performance
	vulnerabilities := make([]domain.VulnerabilityExportItem, 0, len(rows))
	for _, row := range rows {
		// Convert discovered_date (sql.NullTime) to time.Time
		var discoveredDate time.Time
		if row.DiscoveredDate.Valid {
			discoveredDate = row.DiscoveredDate.Time
		}

		// Convert assigned_to (sql.NullString) to string
		assignedTo := ""
		if row.AssignedTo.Valid {
			assignedTo = row.AssignedTo.String
		}

		vulnerabilities = append(vulnerabilities, domain.VulnerabilityExportItem{
			ID:             row.ID.String(),
			Title:          row.Title,
			Severity:       row.Severity,
			Domain:         row.Domain,
			Status:         row.Status,
			DiscoveredDate: discoveredDate,
			DueDate:        &row.DueDate,
			AssignedTo:     assignedTo,
		})
	}

	return vulnerabilities, nil
}

// UC25: Get critical vulnerabilities overdue count
func (r *stakeholderVulnerabilitiesRepository) GetCriticalOverdueCount(ctx context.Context) (int, error) {
	count, err := r.queries.GetStakeholderCriticalOverdueCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get critical overdue count: %w", err)
	}
	return int(count), nil
}

// UC26: Get high severity approaching deadline count
func (r *stakeholderVulnerabilitiesRepository) GetHighApproachingDeadlineCount(ctx context.Context) (int, error) {
	count, err := r.queries.GetStakeholderHighApproachingDeadlineCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get high approaching deadline count: %w", err)
	}
	return int(count), nil
}

// UC27: Get overall SLA compliance data
func (r *stakeholderVulnerabilitiesRepository) GetSLAComplianceData(ctx context.Context) (totalWithDueDate, remediatedOnTime int, err error) {
	result, err := r.queries.GetStakeholderSLAComplianceData(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get SLA compliance data: %w", err)
	}

	return int(result.TotalWithDueDate), int(result.RemediatedOnTime), nil
}

// Helper function to convert sql.NullTime to *time.Time
func convertNullTimeToTimePtr(nt *sql.NullTime) *time.Time {
	if nt == nil || !nt.Valid {
		return nil
	}
	return &nt.Time
}

// Helper to convert UUID to nullable UUID
func toNullUUID(id string) uuid.NullUUID {
	if id == "" {
		return uuid.NullUUID{Valid: false}
	}
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.NullUUID{Valid: false}
	}
	return uuid.NullUUID{UUID: parsedUUID, Valid: true}
}
