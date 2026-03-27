package domain

import (
	"context"
	"time"
)

// StakeholderVulnerabilitiesRepository defines data access methods for stakeholder vulnerabilities tab
type StakeholderVulnerabilitiesRepository interface {
	// UC11: Get critical vulnerabilities count
	GetCriticalVulnerabilitiesCount(ctx context.Context) (int, error)

	// UC12: Get high severity vulnerabilities count
	GetHighSeverityVulnerabilitiesCount(ctx context.Context) (int, error)

	// UC13: Get open issues count (open + in_progress)
	GetOpenIssuesCount(ctx context.Context) (int, error)

	// UC14: Get remediation count (remediated + verified)
	GetRemediationCount(ctx context.Context) (int, error)

	// UC19: List vulnerabilities with search, filters, and pagination
	ListVulnerabilities(ctx context.Context, search, severity, status string, limit, offset int) ([]VulnerabilityListItem, error)

	// UC19: Get total count for pagination
	GetVulnerabilitiesCount(ctx context.Context, search, severity, status string) (int, error)

	// UC24: Export vulnerabilities to CSV (all matching filters)
	ExportVulnerabilities(ctx context.Context, search, severity, status string) ([]VulnerabilityExportItem, error)

	// UC25: Get critical vulnerabilities overdue count
	GetCriticalOverdueCount(ctx context.Context) (int, error)

	// UC26: Get high severity approaching deadline count (within 3 days)
	GetHighApproachingDeadlineCount(ctx context.Context) (int, error)

	// UC27: Get overall SLA compliance data
	GetSLAComplianceData(ctx context.Context) (totalWithDueDate, remediatedOnTime int, err error)
}

// VulnerabilityListItem represents a vulnerability in the list
type VulnerabilityListItem struct {
	ID             string
	Title          string
	Severity       string
	Domain         string
	Status         string
	DiscoveredDate time.Time
	DueDate        *time.Time
	AssignedTo     string
}

// VulnerabilityExportItem represents a vulnerability for CSV export
type VulnerabilityExportItem struct {
	ID             string
	Title          string
	Severity       string
	Domain         string
	Status         string
	DiscoveredDate time.Time
	DueDate        *time.Time
	AssignedTo     string
}

