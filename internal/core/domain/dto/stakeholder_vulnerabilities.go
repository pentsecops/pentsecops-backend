package dto

import "time"

// ============================================================================
// UC11-UC14: Overview Statistics Response
// ============================================================================

type StakeholderVulnerabilitiesStatsResponse struct {
	Critical    VulnStatCard `json:"critical"`
	High        VulnStatCard `json:"high"`
	OpenIssues  VulnStatCard `json:"open_issues"`
	Remediation VulnStatCard `json:"remediation"`
}

type VulnStatCard struct {
	Count   int    `json:"count"`
	Message string `json:"message,omitempty"`
	Color   string `json:"color,omitempty"`
}

// ============================================================================
// UC15-UC18: Search & Filter Request
// ============================================================================

type ListStakeholderVulnerabilitiesRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	PerPage  int    `query:"per_page" validate:"omitempty,min=1,max=100"`
	Search   string `query:"search"`
	Severity string `query:"severity" validate:"omitempty,oneof=all critical high medium low"`
	Status   string `query:"status" validate:"omitempty,oneof=all open in_progress remediated verified"`
}

// ============================================================================
// UC19-UC23: Vulnerabilities Table Response
// ============================================================================

type ListStakeholderVulnerabilitiesResponse struct {
	Vulnerabilities []StakeholderVulnerabilityItem `json:"vulnerabilities"`
	Pagination      PaginationInfo                 `json:"pagination"`
}

type StakeholderVulnerabilityItem struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Severity       string    `json:"severity"`
	SeverityColor  string    `json:"severity_color"`
	Domain         string    `json:"domain"`
	Status         string    `json:"status"`
	StatusColor    string    `json:"status_color"`
	DiscoveredDate time.Time `json:"discovered_date"`
	DueDate        time.Time `json:"due_date,omitempty"`
	AssignedTo     string    `json:"assigned_to"`
	IsOverdue      bool      `json:"is_overdue"`
}

// ============================================================================
// UC24: Export CSV Response
// ============================================================================

type ExportStakeholderVulnerabilitiesRequest struct {
	Search   string `query:"search"`
	Severity string `query:"severity" validate:"omitempty,oneof=all critical high medium low"`
	Status   string `query:"status" validate:"omitempty,oneof=all open in_progress remediated verified"`
}

// ============================================================================
// UC25-UC27: SLA Compliance Response
// ============================================================================

type StakeholderSLAComplianceResponse struct {
	CriticalOverdue         SLACard `json:"critical_overdue"`
	HighApproachingDeadline SLACard `json:"high_approaching_deadline"`
	OverallSLACompliance    SLACard `json:"overall_sla_compliance"`
}

type SLACard struct {
	Count      int     `json:"count,omitempty"`
	Percentage float64 `json:"percentage,omitempty"`
	Message    string  `json:"message"`
	Status     string  `json:"status"`
	Color      string  `json:"color"`
}
