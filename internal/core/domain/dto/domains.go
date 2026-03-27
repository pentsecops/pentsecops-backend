package dto

import "time"

// DomainsStatsResponse - Response for domains overview statistics (UC55-58)
type DomainsStatsResponse struct {
	TotalDomains        int64   `json:"total_domains"`
	AverageRiskScore    float64 `json:"average_risk_score"`
	CriticalIssues      int64   `json:"critical_issues"`
	SLACompliancePercent float64 `json:"sla_compliance_percent"`
}

// DomainResponse - Response for a single domain with vulnerability breakdown (UC59-64)
type DomainResponse struct {
	ID                   string   `json:"id"`
	DomainName           string   `json:"domain_name"`
	IPAddress            *string  `json:"ip_address,omitempty"`
	Description          *string  `json:"description,omitempty"`
	RiskScore            *float64 `json:"risk_score,omitempty"`
	RiskLevel            string   `json:"risk_level"` // High Risk / Medium Risk / Low Risk / Minimal Risk
	TotalVulnerabilities int64    `json:"total_vulnerabilities"`
	CriticalCount        int64    `json:"critical_count"`
	HighCount            int64    `json:"high_count"`
	MediumCount          int64    `json:"medium_count"`
	LowCount             int64    `json:"low_count"`
	SLACompliance        float64  `json:"sla_compliance"`
	OpenIssues           int64    `json:"open_issues"`
	LastScanned          *time.Time `json:"last_scanned,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// ListDomainsRequest - Request for listing domains with pagination (UC59-60)
type ListDomainsRequest struct {
	Page    int `json:"page" validate:"min=1"`
	PerPage int `json:"per_page" validate:"min=1,max=100"`
}

// ListDomainsResponse - Response for listing domains with pagination
type ListDomainsResponse struct {
	Domains    []DomainResponse `json:"domains"`
	Pagination PaginationInfo   `json:"pagination"`
}

// SecurityMetricsResponse - Response for security metrics radar chart (UC65)
type SecurityMetricsResponse struct {
	Authentication   float64 `json:"authentication"`
	Authorization    float64 `json:"authorization"`
	InputValidation  float64 `json:"input_validation"`
	Encryption       float64 `json:"encryption"`
	Configuration    float64 `json:"configuration"`
	NetworkSecurity  float64 `json:"network_security"`
}

// SLABreachItem - Single item in SLA breach analysis
type SLABreachItem struct {
	DomainName           string  `json:"domain_name"`
	SLACompliancePercent float64 `json:"sla_compliance_percent"`
}

// SLABreachAnalysisResponse - Response for SLA breach analysis (UC66)
type SLABreachAnalysisResponse struct {
	Domains []SLABreachItem `json:"domains"`
}

// CreateDomainRequest - Request to create a new domain
type CreateDomainRequest struct {
	DomainName  string   `json:"domain_name" validate:"required,max=255"`
	IPAddress   *string  `json:"ip_address,omitempty" validate:"omitempty,ip"`
	Description *string  `json:"description,omitempty"`
	RiskScore   *float64 `json:"risk_score,omitempty" validate:"omitempty,min=0,max=10"`
}

// UpdateDomainRequest - Request to update a domain
type UpdateDomainRequest struct {
	DomainName  *string  `json:"domain_name,omitempty" validate:"omitempty,max=255"`
	IPAddress   *string  `json:"ip_address,omitempty" validate:"omitempty,ip"`
	Description *string  `json:"description,omitempty"`
	RiskScore   *float64 `json:"risk_score,omitempty" validate:"omitempty,min=0,max=10"`
}

// CreateSecurityMetricRequest - Request to create a security metric
type CreateSecurityMetricRequest struct {
	DomainID    string   `json:"domain_id" validate:"required,uuid"`
	MetricName  string   `json:"metric_name" validate:"required,oneof=authentication authorization input_validation encryption configuration network_security"`
	MetricValue float64  `json:"metric_value" validate:"required,min=0,max=10"`
}

