package dto

// ============================================================================
// UC1: Overall Security Score
// ============================================================================

type SecurityScoreResponse struct {
	Score  float64 `json:"score"`  // 0-10 scale
	Status string  `json:"status"` // "Good - Improving", "Fair - Stable", "Poor - Declining"
}

// ============================================================================
// UC2: Active Projects Count with Breakdown
// ============================================================================

type ActiveProjectsCountResponse struct {
	TotalActive   int `json:"total_active"`
	InProgress    int `json:"in_progress"`
	Completed     int `json:"completed"`
	CompletedText string `json:"completed_text"` // "X in progress, Y completed"
}

// ============================================================================
// UC3: Critical Issues Count
// ============================================================================

type CriticalIssuesResponse struct {
	Count   int    `json:"count"`
	Message string `json:"message"` // "Requires immediate attention"
}

// ============================================================================
// UC4: Open Vulnerabilities Count with Trend
// ============================================================================

type OpenVulnerabilitiesResponse struct {
	Count      int    `json:"count"`
	TrendInfo  string `json:"trend_info"`  // "Down from X last month" or "Up from X last month"
	LastMonth  int    `json:"last_month"`
	Difference int    `json:"difference"` // Positive = increase, Negative = decrease
}

// ============================================================================
// UC5: Remediation Rate
// ============================================================================

type RemediationRateResponse struct {
	Rate              float64 `json:"rate"`               // Percentage 0-100
	TotalVulns        int     `json:"total_vulns"`
	RemediatedVulns   int     `json:"remediated_vulns"`
}

// ============================================================================
// UC6: SLA Compliance Percentage
// ============================================================================

type SLAComplianceStakeholderResponse struct {
	ComplianceRate    float64 `json:"compliance_rate"`     // Percentage 0-100
	TotalWithDueDate  int     `json:"total_with_due_date"`
	RemediatedOnTime  int     `json:"remediated_on_time"`
}

// ============================================================================
// Combined Security Metrics Cards Response (UC1-UC6)
// ============================================================================

type StakeholderSecurityMetricsResponse struct {
	SecurityScore        SecurityScoreResponse            `json:"security_score"`
	ActiveProjects       ActiveProjectsCountResponse      `json:"active_projects"`
	CriticalIssues       CriticalIssuesResponse           `json:"critical_issues"`
	OpenVulnerabilities  OpenVulnerabilitiesResponse      `json:"open_vulnerabilities"`
	RemediationRate      RemediationRateResponse          `json:"remediation_rate"`
	SLACompliance        SLAComplianceStakeholderResponse `json:"sla_compliance"`
}

// ============================================================================
// UC7: Vulnerability Trend Line Chart
// ============================================================================

type MonthlyVulnerabilityTrend struct {
	Month    string `json:"month"`     // "Jun", "Jul", "Aug", "Sep", "Oct"
	Open     int    `json:"open"`      // Count of open vulnerabilities
	Resolved int    `json:"resolved"`  // Count of resolved vulnerabilities
}

type VulnerabilityTrendChartResponse struct {
	Trends []MonthlyVulnerabilityTrend `json:"trends"`
}

// ============================================================================
// UC8: Asset Status Bar Chart
// ============================================================================

type AssetVulnerabilityCounts struct {
	AssetName string `json:"asset_name"` // Domain name
	Critical  int    `json:"critical"`
	High      int    `json:"high"`
	Medium    int    `json:"medium"`
	Low       int    `json:"low"`
	Total     int    `json:"total"`
}

type AssetStatusChartResponse struct {
	Assets []AssetVulnerabilityCounts `json:"assets"`
}

// ============================================================================
// UC9: Recent Security Events
// ============================================================================

type RecentSecurityEvent struct {
	ID          string `json:"id"`
	Description string `json:"description"` // Event description
	EventType   string `json:"event_type"`  // "vulnerability_discovered", "vulnerability_remediated", "report_submitted", "project_status_changed"
	Timestamp   string `json:"timestamp"`   // ISO 8601 format
}

type RecentSecurityEventsResponse struct {
	Events []RecentSecurityEvent `json:"events"`
	Total  int                   `json:"total"`
}

// ============================================================================
// UC10: Remediation Updates
// ============================================================================

type RemediationUpdate struct {
	ID             string `json:"id"`
	VulnTitle      string `json:"vuln_title"`
	PreviousStatus string `json:"previous_status"`
	NewStatus      string `json:"new_status"` // "Remediated" or "Verified"
	RemediatedDate string `json:"remediated_date"`
	AssignedTeam   string `json:"assigned_team"` // Team/person who remediated
}

type RemediationUpdatesResponse struct {
	Updates []RemediationUpdate `json:"updates"`
	Total   int                 `json:"total"`
}

