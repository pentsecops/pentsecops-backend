package dto

import "time"

// CreateVulnerabilityRequest - Request to create a new vulnerability
type CreateVulnerabilityRequest struct {
	Title            string     `json:"title" validate:"required"`
	Description      *string    `json:"description,omitempty"`
	Severity         string     `json:"severity" validate:"required,oneof=critical high medium low"`
	Domain           string     `json:"domain" validate:"required"`
	Status           string     `json:"status" validate:"required,oneof=open in_progress remediated verified"`
	DiscoveredDate   *time.Time `json:"discovered_date,omitempty"`
	DueDate          time.Time  `json:"due_date" validate:"required"`
	AssignedTo       *string    `json:"assigned_to,omitempty"`
	CVSSScore        *float64   `json:"cvss_score,omitempty" validate:"omitempty,min=0,max=10"`
	CWEID            *string    `json:"cwe_id,omitempty"`
	RemediationNotes *string    `json:"remediation_notes,omitempty"`
}

// UpdateVulnerabilityRequest - Request to update a vulnerability (all fields optional for partial updates)
type UpdateVulnerabilityRequest struct {
	Title            *string    `json:"title,omitempty"`
	Description      *string    `json:"description,omitempty"`
	Severity         *string    `json:"severity,omitempty" validate:"omitempty,oneof=critical high medium low"`
	Domain           *string    `json:"domain,omitempty"`
	Status           *string    `json:"status,omitempty" validate:"omitempty,oneof=open in_progress remediated verified"`
	DiscoveredDate   *time.Time `json:"discovered_date,omitempty"`
	DueDate          *time.Time `json:"due_date,omitempty"`
	AssignedTo       *string    `json:"assigned_to,omitempty"`
	CVSSScore        *float64   `json:"cvss_score,omitempty" validate:"omitempty,min=0,max=10"`
	CWEID            *string    `json:"cwe_id,omitempty"`
	RemediationNotes *string    `json:"remediation_notes,omitempty"`
}

// VulnerabilityResponse - Response for a single vulnerability
type VulnerabilityResponse struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Description      *string    `json:"description,omitempty"`
	Severity         string     `json:"severity"`
	Domain           string     `json:"domain"`
	Status           string     `json:"status"`
	DiscoveredDate   *time.Time `json:"discovered_date,omitempty"`
	DueDate          time.Time  `json:"due_date"`
	AssignedTo       *string    `json:"assigned_to,omitempty"`
	CVSSScore        *float64   `json:"cvss_score,omitempty"`
	CWEID            *string    `json:"cwe_id,omitempty"`
	RemediationNotes *string    `json:"remediation_notes,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ListVulnerabilitiesRequest - Request to list vulnerabilities with filters
type ListVulnerabilitiesRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	PerPage  int    `query:"per_page" validate:"omitempty,min=1,max=100"`
	Search   string `query:"search"`
	Severity string `query:"severity" validate:"omitempty,oneof=critical high medium low"`
	Status   string `query:"status" validate:"omitempty,oneof=open in_progress remediated verified"`
}

// ListVulnerabilitiesResponse - Response for listing vulnerabilities
type ListVulnerabilitiesResponse struct {
	Vulnerabilities []VulnerabilityResponse `json:"vulnerabilities"`
	Pagination      PaginationInfo          `json:"pagination"`
}

// VulnerabilityStatsResponse - Response for vulnerability statistics
type VulnerabilityStatsResponse struct {
	Total      int64 `json:"total"`
	Critical   int64 `json:"critical"`
	High       int64 `json:"high"`
	Medium     int64 `json:"medium"`
	Low        int64 `json:"low"`
	Open       int64 `json:"open"`
	InProgress int64 `json:"in_progress"`
	Remediated int64 `json:"remediated"`
	Verified   int64 `json:"verified"`
}

// SLAComplianceResponse - Response for SLA compliance metrics
type SLAComplianceResponse struct {
	CriticalOverdue   int64   `json:"critical_overdue"`
	HighApproaching   int64   `json:"high_approaching"`
	TotalWithDueDate  int64   `json:"total_with_due_date"`
	RemediatedOnTime  int64   `json:"remediated_on_time"`
	CompliancePercent float64 `json:"compliance_percent"`
}

// ExportVulnerabilityResponse - Response for exporting vulnerabilities
type ExportVulnerabilityResponse struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Severity       string     `json:"severity"`
	Domain         string     `json:"domain"`
	Status         string     `json:"status"`
	DiscoveredDate *time.Time `json:"discovered_date,omitempty"`
	DueDate        time.Time  `json:"due_date"`
	AssignedTo     *string    `json:"assigned_to,omitempty"`
}
