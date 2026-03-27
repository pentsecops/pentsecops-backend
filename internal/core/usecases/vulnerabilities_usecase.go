package usecases

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
)

type VulnerabilitiesUseCase struct {
	repo      domain.VulnerabilitiesRepository
	validator *validator.Validate
}

func NewVulnerabilitiesUseCase(repo domain.VulnerabilitiesRepository) *VulnerabilitiesUseCase {
	return &VulnerabilitiesUseCase{
		repo:      repo,
		validator: validator.New(),
	}
}

// CreateVulnerability - UC45: Create New Vulnerability with All Details
func (uc *VulnerabilitiesUseCase) CreateVulnerability(ctx context.Context, req *dto.CreateVulnerabilityRequest) (*dto.VulnerabilityResponse, error) {
	// UC46: Validate required fields
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// UC47: Validate Due Date is After Discovered Date
	if req.DiscoveredDate != nil && req.DueDate.Before(*req.DiscoveredDate) {
		return nil, fmt.Errorf("due date must be after discovered date")
	}

	// Generate UUIDv7 for vulnerability
	vulnID := uuid.Must(uuid.NewV7())
	now := time.Now()

	params := &domain.CreateVulnerabilityParams{
		ID:               vulnID.String(),
		Title:            req.Title,
		Description:      req.Description,
		Severity:         req.Severity,
		Domain:           req.Domain,
		Status:           req.Status,
		DiscoveredDate:   req.DiscoveredDate,
		DueDate:          req.DueDate,
		AssignedTo:       req.AssignedTo,
		CVSSScore:        req.CVSSScore,
		CWEID:            req.CWEID,
		RemediationNotes: req.RemediationNotes,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	vuln, err := uc.repo.CreateVulnerability(ctx, params)
	if err != nil {
		fmt.Printf("CreateVulnerability Repository Error: %v\n", err)
		return nil, fmt.Errorf("failed to create vulnerability: %w", err)
	}

	return uc.toVulnerabilityResponse(vuln), nil
}

// GetVulnerabilityByID - Get vulnerability by ID
func (uc *VulnerabilitiesUseCase) GetVulnerabilityByID(ctx context.Context, id string) (*dto.VulnerabilityResponse, error) {
	vuln, err := uc.repo.GetVulnerabilityByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerability: %w", err)
	}
	if vuln == nil {
		return nil, fmt.Errorf("vulnerability not found")
	}

	return uc.toVulnerabilityResponse(vuln), nil
}

// UpdateVulnerability - UC52: Edit Vulnerability from Table
func (uc *VulnerabilitiesUseCase) UpdateVulnerability(ctx context.Context, id string, req *dto.UpdateVulnerabilityRequest) (*dto.VulnerabilityResponse, error) {
	// Validate required fields
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Validate Due Date is After Discovered Date (only if both are provided)
	if req.DiscoveredDate != nil && req.DueDate != nil && req.DueDate.Before(*req.DiscoveredDate) {
		return nil, fmt.Errorf("due date must be after discovered date")
	}

	params := &domain.UpdateVulnerabilityParams{
		ID:               id,
		Title:            req.Title,
		Description:      req.Description,
		Severity:         req.Severity,
		Domain:           req.Domain,
		Status:           req.Status,
		DiscoveredDate:   req.DiscoveredDate,
		DueDate:          req.DueDate,
		AssignedTo:       req.AssignedTo,
		CVSSScore:        req.CVSSScore,
		CWEID:            req.CWEID,
		RemediationNotes: req.RemediationNotes,
		UpdatedAt:        time.Now(),
	}

	vuln, err := uc.repo.UpdateVulnerability(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update vulnerability: %w", err)
	}

	return uc.toVulnerabilityResponse(vuln), nil
}

// DeleteVulnerability - Delete vulnerability
func (uc *VulnerabilitiesUseCase) DeleteVulnerability(ctx context.Context, id string) error {
	err := uc.repo.DeleteVulnerability(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete vulnerability: %w", err)
	}
	return nil
}

// ListVulnerabilities - UC40, UC41, UC42, UC43, UC48, UC49: List vulnerabilities with search and filters
func (uc *VulnerabilitiesUseCase) ListVulnerabilities(ctx context.Context, req *dto.ListVulnerabilitiesRequest) (*dto.ListVulnerabilitiesResponse, error) {
	// Validate request
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Set defaults
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 5
	}

	offset := (page - 1) * perPage

	// Get vulnerabilities with search and filters
	var vulns []*domain.Vulnerability
	var total int64
	var err error

	if req.Search != "" || req.Severity != "" || req.Status != "" {
		// UC43: Combine Search and Multiple Filters
		vulns, err = uc.repo.SearchAndFilterVulnerabilities(ctx, req.Search, req.Severity, req.Status, perPage, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to search and filter vulnerabilities: %w", err)
		}
		total, err = uc.repo.CountSearchAndFilterVulnerabilities(ctx, req.Search, req.Severity, req.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to count filtered vulnerabilities: %w", err)
		}
	} else {
		// UC48: Display All Vulnerabilities with Pagination
		vulns, err = uc.repo.ListVulnerabilities(ctx, perPage, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to list vulnerabilities: %w", err)
		}
		total, err = uc.repo.CountVulnerabilities(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to count vulnerabilities: %w", err)
		}
	}

	// Convert to response
	vulnerabilities := make([]dto.VulnerabilityResponse, 0, len(vulns))
	for _, v := range vulns {
		vulnerabilities = append(vulnerabilities, *uc.toVulnerabilityResponse(v))
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return &dto.ListVulnerabilitiesResponse{
		Vulnerabilities: vulnerabilities,
		Pagination: dto.PaginationInfo{
			CurrentPage: page,
			PerPage:     perPage,
			Total:       total,
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrev:     page > 1,
		},
	}, nil
}

// GetVulnerabilityStats - UC39: Fetch and Display Vulnerability Statistics
func (uc *VulnerabilitiesUseCase) GetVulnerabilityStats(ctx context.Context) (*dto.VulnerabilityStatsResponse, error) {
	stats, err := uc.repo.GetVulnerabilityStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerability stats: %w", err)
	}

	return &dto.VulnerabilityStatsResponse{
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

// GetSLACompliance - UC54: Display SLA Compliance Information
func (uc *VulnerabilitiesUseCase) GetSLACompliance(ctx context.Context) (*dto.SLAComplianceResponse, error) {
	sla, err := uc.repo.GetSLACompliance(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get SLA compliance: %w", err)
	}

	// Calculate compliance percentage
	var compliancePercent float64
	if sla.TotalWithDueDate > 0 {
		compliancePercent = (float64(sla.RemediatedOnTime) / float64(sla.TotalWithDueDate)) * 100
	}

	return &dto.SLAComplianceResponse{
		CriticalOverdue:   sla.CriticalOverdue,
		HighApproaching:   sla.HighApproaching,
		TotalWithDueDate:  sla.TotalWithDueDate,
		RemediatedOnTime:  sla.RemediatedOnTime,
		CompliancePercent: compliancePercent,
	}, nil
}

// ExportVulnerabilitiesToCSV - UC53: Export Vulnerabilities to CSV
func (uc *VulnerabilitiesUseCase) ExportVulnerabilitiesToCSV(ctx context.Context, req *dto.ListVulnerabilitiesRequest) ([]byte, error) {
	// Get vulnerabilities with filters
	vulns, err := uc.repo.ExportVulnerabilities(ctx, req.Search, req.Severity, req.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to export vulnerabilities: %w", err)
	}

	// Create CSV
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "Title", "Severity", "Domain", "Status", "Discovered Date", "Due Date", "Assigned To"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, v := range vulns {
		discoveredDate := ""
		if v.DiscoveredDate != nil {
			discoveredDate = v.DiscoveredDate.Format("2006-01-02")
		}
		dueDate := v.DueDate.Format("2006-01-02")
		assignedTo := ""
		if v.AssignedTo != nil {
			assignedTo = *v.AssignedTo
		}

		row := []string{
			v.ID,
			v.Title,
			v.Severity,
			v.Domain,
			v.Status,
			discoveredDate,
			dueDate,
			assignedTo,
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return []byte(buf.String()), nil
}

// Helper function to convert domain.Vulnerability to dto.VulnerabilityResponse
func (uc *VulnerabilitiesUseCase) toVulnerabilityResponse(v *domain.Vulnerability) *dto.VulnerabilityResponse {
	return &dto.VulnerabilityResponse{
		ID:               v.ID.String(),
		Title:            v.Title,
		Description:      v.Description,
		Severity:         v.Severity,
		Domain:           v.Domain,
		Status:           v.Status,
		DiscoveredDate:   v.DiscoveredDate,
		DueDate:          *v.DueDate,
		AssignedTo:       v.AssignedTo,
		CVSSScore:        v.CVSSScore,
		CWEID:            v.CWEID,
		RemediationNotes: v.RemediationNotes,
		CreatedAt:        v.CreatedAt,
		UpdatedAt:        v.UpdatedAt,
	}
}
