package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/adapter/db/postgres/sqlc"
	"github.com/pentsecops/backend/internal/core/domain"
)

type VulnerabilitiesRepository struct {
	queries *sqlc.Queries
}

func NewVulnerabilitiesRepository(db *sql.DB) *VulnerabilitiesRepository {
	return &VulnerabilitiesRepository{
		queries: sqlc.New(db),
	}
}

func (r *VulnerabilitiesRepository) CreateVulnerability(ctx context.Context, vuln *domain.CreateVulnerabilityParams) (*domain.Vulnerability, error) {
	// Parse UUID
	id, err := uuid.Parse(vuln.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid vulnerability ID: %w", err)
	}

	// Parse optional UUIDs
	var domainID, projectID, discoveredBy uuid.NullUUID
	if vuln.DomainID != nil {
		domainUUID, err := uuid.Parse(*vuln.DomainID)
		if err != nil {
			return nil, fmt.Errorf("invalid domain ID: %w", err)
		}
		domainID = uuid.NullUUID{UUID: domainUUID, Valid: true}
	}
	if vuln.ProjectID != nil {
		projectUUID, err := uuid.Parse(*vuln.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("invalid project ID: %w", err)
		}
		projectID = uuid.NullUUID{UUID: projectUUID, Valid: true}
	}
	if vuln.DiscoveredBy != nil {
		discoveredByUUID, err := uuid.Parse(*vuln.DiscoveredBy)
		if err != nil {
			return nil, fmt.Errorf("invalid discovered_by ID: %w", err)
		}
		discoveredBy = uuid.NullUUID{UUID: discoveredByUUID, Valid: true}
	}

	// Convert optional fields
	var description, assignedTo, cweid, remediationNotes sql.NullString
	if vuln.Description != nil {
		description = sql.NullString{String: *vuln.Description, Valid: true}
	}
	if vuln.AssignedTo != nil {
		assignedTo = sql.NullString{String: *vuln.AssignedTo, Valid: true}
	}
	if vuln.CWEID != nil {
		cweid = sql.NullString{String: *vuln.CWEID, Valid: true}
	}
	if vuln.RemediationNotes != nil {
		remediationNotes = sql.NullString{String: *vuln.RemediationNotes, Valid: true}
	}

	var cvssScore sql.NullString
	if vuln.CVSSScore != nil {
		cvssScore = sql.NullString{String: fmt.Sprintf("%.1f", *vuln.CVSSScore), Valid: true}
	}

	var discoveredDate sql.NullTime
	if vuln.DiscoveredDate != nil {
		discoveredDate = sql.NullTime{Time: *vuln.DiscoveredDate, Valid: true}
	}

	params := sqlc.CreateVulnerabilityParams{
		ID:               id,
		Title:            vuln.Title,
		Description:      description,
		Severity:         vuln.Severity,
		Domain:           vuln.Domain,
		Status:           vuln.Status,
		DiscoveredDate:   discoveredDate,
		DueDate:          vuln.DueDate,
		AssignedTo:       assignedTo,
		CvssScore:        cvssScore,
		CweID:            cweid,
		DomainID:         domainID,
		ProjectID:        projectID,
		DiscoveredBy:     discoveredBy,
		RemediationNotes: remediationNotes,
		CreatedAt:        sql.NullTime{Time: vuln.CreatedAt, Valid: true},
		UpdatedAt:        sql.NullTime{Time: vuln.UpdatedAt, Valid: true},
	}

	v, err := r.queries.CreateVulnerability(ctx, params)
	if err != nil {
		fmt.Printf("SQL CreateVulnerability Error: %v\n", err)
		return nil, fmt.Errorf("failed to create vulnerability: %w", err)
	}

	return r.toVulnerability(&v), nil
}

func (r *VulnerabilitiesRepository) GetVulnerabilityByID(ctx context.Context, id string) (*domain.Vulnerability, error) {
	vulnID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid vulnerability ID: %w", err)
	}

	v, err := r.queries.GetVulnerabilityByID(ctx, vulnID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get vulnerability: %w", err)
	}

	return r.toVulnerability(&v), nil
}

func (r *VulnerabilitiesRepository) UpdateVulnerability(ctx context.Context, vuln *domain.UpdateVulnerabilityParams) (*domain.Vulnerability, error) {
	id, err := uuid.Parse(vuln.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid vulnerability ID: %w", err)
	}

	// Get existing vulnerability
	existing, err := r.queries.GetVulnerabilityByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing vulnerability: %w", err)
	}

	// Use existing values if not provided in update
	title := existing.Title
	if vuln.Title != nil {
		title = *vuln.Title
	}

	severity := existing.Severity
	if vuln.Severity != nil {
		severity = *vuln.Severity
	}

	domain := existing.Domain
	if vuln.Domain != nil {
		domain = *vuln.Domain
	}

	status := existing.Status
	if vuln.Status != nil {
		status = *vuln.Status
	}

	dueDate := existing.DueDate
	if vuln.DueDate != nil {
		dueDate = *vuln.DueDate
	}

	// Handle nullable fields
	description := existing.Description
	if vuln.Description != nil {
		description = sql.NullString{String: *vuln.Description, Valid: true}
	}

	assignedTo := existing.AssignedTo
	if vuln.AssignedTo != nil {
		assignedTo = sql.NullString{String: *vuln.AssignedTo, Valid: true}
	}

	cweid := existing.CweID
	if vuln.CWEID != nil {
		cweid = sql.NullString{String: *vuln.CWEID, Valid: true}
	}

	remediationNotes := existing.RemediationNotes
	if vuln.RemediationNotes != nil {
		remediationNotes = sql.NullString{String: *vuln.RemediationNotes, Valid: true}
	}

	cvssScore := existing.CvssScore
	if vuln.CVSSScore != nil {
		cvssScore = sql.NullString{String: fmt.Sprintf("%.1f", *vuln.CVSSScore), Valid: true}
	}

	discoveredDate := existing.DiscoveredDate
	if vuln.DiscoveredDate != nil {
		discoveredDate = sql.NullTime{Time: *vuln.DiscoveredDate, Valid: true}
	}

	params := sqlc.UpdateVulnerabilityParams{
		ID:               id,
		Title:            title,
		Description:      description,
		Severity:         severity,
		Domain:           domain,
		Status:           status,
		DiscoveredDate:   discoveredDate,
		DueDate:          dueDate,
		AssignedTo:       assignedTo,
		CvssScore:        cvssScore,
		CweID:            cweid,
		RemediationNotes: remediationNotes,
		UpdatedAt:        sql.NullTime{Time: vuln.UpdatedAt, Valid: true},
	}

	v, err := r.queries.UpdateVulnerability(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update vulnerability: %w", err)
	}

	return r.toVulnerability(&v), nil
}

func (r *VulnerabilitiesRepository) DeleteVulnerability(ctx context.Context, id string) error {
	vulnID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid vulnerability ID: %w", err)
	}

	err = r.queries.DeleteVulnerability(ctx, vulnID)
	if err != nil {
		return fmt.Errorf("failed to delete vulnerability: %w", err)
	}

	return nil
}

func (r *VulnerabilitiesRepository) ListVulnerabilities(ctx context.Context, limit, offset int) ([]*domain.Vulnerability, error) {
	params := sqlc.ListVulnerabilitiesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	vulns, err := r.queries.ListVulnerabilities(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list vulnerabilities: %w", err)
	}

	result := make([]*domain.Vulnerability, 0, len(vulns))
	for i := range vulns {
		result = append(result, r.toVulnerability(&vulns[i]))
	}

	return result, nil
}

func (r *VulnerabilitiesRepository) CountVulnerabilities(ctx context.Context) (int64, error) {
	count, err := r.queries.CountVulnerabilities(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count vulnerabilities: %w", err)
	}
	return count, nil
}

func (r *VulnerabilitiesRepository) SearchAndFilterVulnerabilities(ctx context.Context, search, severity, status string, limit, offset int) ([]*domain.Vulnerability, error) {
	// Prepare search pattern
	searchPattern := "%" + search + "%"
	if search == "" {
		searchPattern = "%%"
	}

	params := sqlc.SearchAndFilterVulnerabilitiesParams{
		Lower:   searchPattern,
		Column2: severity,
		Column3: status,
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	vulns, err := r.queries.SearchAndFilterVulnerabilities(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search and filter vulnerabilities: %w", err)
	}

	result := make([]*domain.Vulnerability, 0, len(vulns))
	for i := range vulns {
		result = append(result, r.toVulnerability(&vulns[i]))
	}

	return result, nil
}

func (r *VulnerabilitiesRepository) CountSearchAndFilterVulnerabilities(ctx context.Context, search, severity, status string) (int64, error) {
	searchPattern := "%" + search + "%"
	if search == "" {
		searchPattern = "%%"
	}

	params := sqlc.CountSearchAndFilterVulnerabilitiesParams{
		Lower:   searchPattern,
		Column2: severity,
		Column3: status,
	}

	count, err := r.queries.CountSearchAndFilterVulnerabilities(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to count filtered vulnerabilities: %w", err)
	}

	return count, nil
}

func (r *VulnerabilitiesRepository) GetVulnerabilityStats(ctx context.Context) (*domain.VulnerabilityStats, error) {
	stats, err := r.queries.GetVulnerabilityStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerability stats: %w", err)
	}

	return &domain.VulnerabilityStats{
		Total:      stats.Total,
		Critical:   stats.Critical,
		High:       stats.High,
		Medium:     stats.Medium,
		Low:        stats.Low,
		Open:       stats.Open,
		InProgress: stats.InProgress,
		Remediated: stats.Remediated,
		Verified:   stats.Verified,
	}, nil
}

func (r *VulnerabilitiesRepository) GetSLACompliance(ctx context.Context) (*domain.SLACompliance, error) {
	sla, err := r.queries.GetSLACompliance(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get SLA compliance: %w", err)
	}

	return &domain.SLACompliance{
		CriticalOverdue:  sla.CriticalOverdue,
		HighApproaching:  sla.HighApproaching,
		TotalWithDueDate: sla.TotalWithDueDate,
		RemediatedOnTime: sla.RemediatedOnTime,
	}, nil
}

func (r *VulnerabilitiesRepository) ExportVulnerabilities(ctx context.Context, search, severity, status string) ([]*domain.VulnerabilityExport, error) {
	searchPattern := "%" + search + "%"
	if search == "" {
		searchPattern = "%%"
	}

	params := sqlc.ExportVulnerabilitiesParams{
		Lower:   searchPattern,
		Column2: severity,
		Column3: status,
	}

	vulns, err := r.queries.ExportVulnerabilities(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to export vulnerabilities: %w", err)
	}

	result := make([]*domain.VulnerabilityExport, 0, len(vulns))
	for i := range vulns {
		export := &domain.VulnerabilityExport{
			ID:       vulns[i].ID.String(),
			Title:    vulns[i].Title,
			Severity: vulns[i].Severity,
			Domain:   vulns[i].Domain,
			Status:   vulns[i].Status,
			DueDate:  vulns[i].DueDate,
		}
		if vulns[i].DiscoveredDate.Valid {
			export.DiscoveredDate = &vulns[i].DiscoveredDate.Time
		}
		if vulns[i].AssignedTo.Valid {
			export.AssignedTo = &vulns[i].AssignedTo.String
		}
		result = append(result, export)
	}

	return result, nil
}

// Helper function to convert sqlc.Vulnerability to domain.Vulnerability
func (r *VulnerabilitiesRepository) toVulnerability(v *sqlc.Vulnerability) *domain.Vulnerability {
	result := &domain.Vulnerability{
		ID:        v.ID,
		Title:     v.Title,
		Severity:  v.Severity,
		Domain:    v.Domain,
		Status:    v.Status,
		DueDate:   &v.DueDate,
		CreatedAt: v.CreatedAt.Time,
		UpdatedAt: v.UpdatedAt.Time,
	}

	if v.Description.Valid {
		result.Description = &v.Description.String
	}
	if v.AssignedTo.Valid {
		result.AssignedTo = &v.AssignedTo.String
	}
	if v.DiscoveredDate.Valid {
		result.DiscoveredDate = &v.DiscoveredDate.Time
	}
	if v.CweID.Valid {
		result.CWEID = &v.CweID.String
	}
	if v.RemediationNotes.Valid {
		result.RemediationNotes = &v.RemediationNotes.String
	}
	if v.CvssScore.Valid {
		var score float64
		fmt.Sscanf(v.CvssScore.String, "%f", &score)
		result.CVSSScore = &score
	}

	return result
}
