package usecases

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"time"

	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
)

type stakeholderVulnerabilitiesUseCase struct {
	repo domain.StakeholderVulnerabilitiesRepository
}

// NewStakeholderVulnerabilitiesUseCase creates a new stakeholder vulnerabilities use case
func NewStakeholderVulnerabilitiesUseCase(repo domain.StakeholderVulnerabilitiesRepository) domain.StakeholderVulnerabilitiesUseCase {
	return &stakeholderVulnerabilitiesUseCase{
		repo: repo,
	}
}

// UC11-UC14: Get vulnerabilities statistics
func (uc *stakeholderVulnerabilitiesUseCase) GetVulnerabilitiesStats(ctx context.Context) (*dto.StakeholderVulnerabilitiesStatsResponse, error) {
	// UC11: Get critical count
	criticalCount, err := uc.repo.GetCriticalVulnerabilitiesCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get critical count: %w", err)
	}

	// UC12: Get high count
	highCount, err := uc.repo.GetHighSeverityVulnerabilitiesCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get high count: %w", err)
	}

	// UC13: Get open issues count
	openIssuesCount, err := uc.repo.GetOpenIssuesCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get open issues count: %w", err)
	}

	// UC14: Get remediation count
	remediationCount, err := uc.repo.GetRemediationCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get remediation count: %w", err)
	}

	// Calculate remediation percentage
	totalVulns := criticalCount + highCount + openIssuesCount + remediationCount
	var remediationPercentage float64
	if totalVulns > 0 {
		remediationPercentage = (float64(remediationCount) / float64(totalVulns)) * 100
	}

	return &dto.StakeholderVulnerabilitiesStatsResponse{
		Critical: dto.VulnStatCard{
			Count:   criticalCount,
			Message: "Requires immediate attention",
			Color:   "red",
		},
		High: dto.VulnStatCard{
			Count:   highCount,
			Message: "High priority issues",
			Color:   "orange",
		},
		OpenIssues: dto.VulnStatCard{
			Count:   openIssuesCount,
			Message: "Currently being addressed",
			Color:   "blue",
		},
		Remediation: dto.VulnStatCard{
			Count:   remediationCount,
			Message: fmt.Sprintf("%.1f%% remediation rate", remediationPercentage),
			Color:   "green",
		},
	}, nil
}

// UC15-UC23: List vulnerabilities with search, filters, and pagination
func (uc *stakeholderVulnerabilitiesUseCase) ListVulnerabilities(ctx context.Context, req *dto.ListStakeholderVulnerabilitiesRequest) (*dto.ListStakeholderVulnerabilitiesResponse, error) {
	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PerPage < 1 {
		req.PerPage = 5
	}

	// Calculate offset
	offset := (req.Page - 1) * req.PerPage

	// Get total count for pagination
	totalCount, err := uc.repo.GetVulnerabilitiesCount(ctx, req.Search, req.Severity, req.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerabilities count: %w", err)
	}

	// Get vulnerabilities
	vulnerabilities, err := uc.repo.ListVulnerabilities(ctx, req.Search, req.Severity, req.Status, req.PerPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list vulnerabilities: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(req.PerPage)))

	// Preallocate response slice
	items := make([]dto.StakeholderVulnerabilityItem, 0, len(vulnerabilities))
	now := time.Now()

	for _, vuln := range vulnerabilities {
		// UC21: Determine severity color
		severityColor := getSeverityColor(vuln.Severity)

		// UC22: Determine status color
		statusColor := getStatusColor(vuln.Status)

		// UC23: Check if overdue
		isOverdue := false
		var dueDate time.Time
		if vuln.DueDate != nil {
			dueDate = *vuln.DueDate
			if vuln.Status == "open" || vuln.Status == "in_progress" {
				isOverdue = vuln.DueDate.Before(now)
			}
		}

		items = append(items, dto.StakeholderVulnerabilityItem{
			ID:             vuln.ID,
			Title:          vuln.Title,
			Severity:       vuln.Severity,
			SeverityColor:  severityColor,
			Domain:         vuln.Domain,
			Status:         vuln.Status,
			StatusColor:    statusColor,
			DiscoveredDate: vuln.DiscoveredDate,
			DueDate:        dueDate,
			AssignedTo:     vuln.AssignedTo,
			IsOverdue:      isOverdue,
		})
	}

	return &dto.ListStakeholderVulnerabilitiesResponse{
		Vulnerabilities: items,
		Pagination: dto.PaginationInfo{
			CurrentPage: req.Page,
			PerPage:     req.PerPage,
			Total:       int64(totalCount),
			TotalPages:  totalPages,
			HasNext:     req.Page < totalPages,
			HasPrev:     req.Page > 1,
		},
	}, nil
}

// UC24: Export vulnerabilities to CSV
func (uc *stakeholderVulnerabilitiesUseCase) ExportVulnerabilitiesToCSV(ctx context.Context, req *dto.ExportStakeholderVulnerabilitiesRequest) ([]byte, error) {
	// Get all vulnerabilities matching filters
	vulnerabilities, err := uc.repo.ExportVulnerabilities(ctx, req.Search, req.Severity, req.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to export vulnerabilities: %w", err)
	}

	// Create CSV buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "Title", "Severity", "Domain", "Status", "Discovered Date", "Due Date", "Assigned To"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, vuln := range vulnerabilities {
		dueDateStr := ""
		if vuln.DueDate != nil {
			dueDateStr = vuln.DueDate.Format("2006-01-02")
		}

		row := []string{
			vuln.ID,
			vuln.Title,
			vuln.Severity,
			vuln.Domain,
			vuln.Status,
			vuln.DiscoveredDate.Format("2006-01-02"),
			dueDateStr,
			vuln.AssignedTo,
		}

		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}

// UC25-UC27: Get SLA compliance data
func (uc *stakeholderVulnerabilitiesUseCase) GetSLACompliance(ctx context.Context) (*dto.StakeholderSLAComplianceResponse, error) {
	// UC25: Get critical overdue count
	criticalOverdue, err := uc.repo.GetCriticalOverdueCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get critical overdue count: %w", err)
	}

	// UC26: Get high approaching deadline count
	highApproaching, err := uc.repo.GetHighApproachingDeadlineCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get high approaching deadline count: %w", err)
	}

	// UC27: Get overall SLA compliance
	totalWithDueDate, remediatedOnTime, err := uc.repo.GetSLAComplianceData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get SLA compliance data: %w", err)
	}

	// Calculate compliance percentage
	var compliancePercentage float64
	if totalWithDueDate > 0 {
		compliancePercentage = (float64(remediatedOnTime) / float64(totalWithDueDate)) * 100
	}

	return &dto.StakeholderSLAComplianceResponse{
		CriticalOverdue: dto.SLACard{
			Count:   criticalOverdue,
			Message: "Critical vulnerabilities overdue",
			Status:  "Requires immediate attention",
			Color:   "red",
		},
		HighApproachingDeadline: dto.SLACard{
			Count:   highApproaching,
			Message: "High severity approaching deadline",
			Status:  "Due within 3 days",
			Color:   "orange",
		},
		OverallSLACompliance: dto.SLACard{
			Percentage: compliancePercentage,
			Message:    "Overall SLA compliance",
			Status:     "This month",
			Color:      "green",
		},
	}, nil
}

// Helper: Get severity color (UC21)
func getSeverityColor(severity string) string {
	switch severity {
	case "critical":
		return "red"
	case "high":
		return "orange"
	case "medium":
		return "yellow"
	case "low":
		return "blue"
	default:
		return "gray"
	}
}

// Helper: Get status color (UC22)
func getStatusColor(status string) string {
	switch status {
	case "open":
		return "red"
	case "in_progress":
		return "blue"
	case "remediated":
		return "green"
	case "verified":
		return "gray"
	default:
		return "gray"
	}
}
