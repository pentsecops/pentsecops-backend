package dto

// AdminOverviewStatsResponse represents the complete overview statistics
type AdminOverviewStatsResponse struct {
	TotalProjects       TotalProjectsStats       `json:"total_projects"`
	TotalVulnerabilities TotalVulnerabilitiesStats `json:"total_vulnerabilities"`
	ActiveUsers         ActiveUsersStats         `json:"active_users"`
	OpenIssues          OpenIssuesStats          `json:"open_issues"`
	CompletedProjects   CompletedProjectsStats   `json:"completed_projects"`
	CriticalIssues      CriticalIssuesStats      `json:"critical_issues"`
}

// TotalProjectsStats represents total projects statistics
type TotalProjectsStats struct {
	Total      int64  `json:"total"`
	Open       int64  `json:"open"`
	InProgress int64  `json:"in_progress"`
	Breakdown  string `json:"breakdown"`
}

// TotalVulnerabilitiesStats represents total vulnerabilities statistics
type TotalVulnerabilitiesStats struct {
	Total    int64  `json:"total"`
	Subtitle string `json:"subtitle"`
}

// ActiveUsersStats represents active users statistics
type ActiveUsersStats struct {
	Total        int64  `json:"total"`
	Pentesters   int64  `json:"pentesters"`
	Stakeholders int64  `json:"stakeholders"`
	Breakdown    string `json:"breakdown"`
}

// OpenIssuesStats represents open issues statistics
type OpenIssuesStats struct {
	Count    int64  `json:"count"`
	Subtitle string `json:"subtitle"`
}

// CompletedProjectsStats represents completed projects statistics
type CompletedProjectsStats struct {
	Count    int64  `json:"count"`
	Subtitle string `json:"subtitle"`
}

// CriticalIssuesStats represents critical issues statistics
type CriticalIssuesStats struct {
	Count    int64  `json:"count"`
	Subtitle string `json:"subtitle"`
}

// VulnerabilitiesBySeverityResponse represents vulnerabilities grouped by severity
type VulnerabilitiesBySeverityResponse struct {
	Data []SeverityCount `json:"data"`
}

// SeverityCount represents count for a specific severity
type SeverityCount struct {
	Severity string `json:"severity"`
	Count    int64  `json:"count"`
}

// Top5DomainsResponse represents top 5 domains by vulnerabilities
type Top5DomainsResponse struct {
	Data []DomainVulnerabilityCount `json:"data"`
}

// DomainVulnerabilityCount represents vulnerability count for a domain
type DomainVulnerabilityCount struct {
	Domain              string `json:"domain"`
	VulnerabilityCount  int64  `json:"vulnerability_count"`
}

// ProjectStatusDistributionResponse represents project status distribution
type ProjectStatusDistributionResponse struct {
	Data []StatusCount `json:"data"`
}

// StatusCount represents count for a specific status
type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// RecentActivityResponse represents recent activity logs
type RecentActivityResponse struct {
	Activities []ActivityLogDTO `json:"activities"`
	Pagination PaginationMeta   `json:"pagination"`
}

// ActivityLogDTO represents an activity log entry
type ActivityLogDTO struct {
	ID         string  `json:"id"`
	UserID     *string `json:"user_id,omitempty"`
	UserName   *string `json:"user_name,omitempty"`
	Action     string  `json:"action"`
	EntityType *string `json:"entity_type,omitempty"`
	EntityID   *string `json:"entity_id,omitempty"`
	IPAddress  *string `json:"ip_address,omitempty"`
	UserAgent  *string `json:"user_agent,omitempty"`
	CreatedAt  string  `json:"created_at"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

