package dto

import "time"

// ============================================================================
// UC28-UC30: Overview Statistics Response
// ============================================================================

type StakeholderReportsStatsResponse struct {
	TotalReports   ReportStatCard `json:"total_reports"`
	UnderReview    ReportStatCard `json:"under_review"`
	Remediated     ReportStatCard `json:"remediated"`
}

type ReportStatCard struct {
	Count   int    `json:"count"`
	Message string `json:"message,omitempty"`
	Color   string `json:"color,omitempty"`
}

// ============================================================================
// UC31-UC33: List Reports Request/Response
// ============================================================================

type ListStakeholderReportsRequest struct {
	Page    int    `query:"page" validate:"omitempty,min=1"`
	PerPage int    `query:"per_page" validate:"omitempty,min=1,max=100"`
	Status  string `query:"status" validate:"omitempty,oneof=all received under_review shared remediated"`
}

type ListStakeholderReportsResponse struct {
	Reports    []StakeholderReportItem `json:"reports"`
	Pagination PaginationInfo          `json:"pagination"`
}

type StakeholderReportItem struct {
	ID                 string    `json:"id"`
	Title              string    `json:"title"`
	SubmittedBy        string    `json:"submitted_by"`
	SubmittedDate      time.Time `json:"submitted_date"`
	Project            string    `json:"project"`
	Status             string    `json:"status"`
	StatusColor        string    `json:"status_color"`
	VulnerabilitiesCount int     `json:"vulnerabilities_count"`
	EvidenceCount      int       `json:"evidence_count"`
}

// ============================================================================
// UC37-UC38: Report Viewer Response
// ============================================================================

type ViewStakeholderReportResponse struct {
	Report          ReportDetails                `json:"report"`
	Vulnerabilities []ReportVulnerabilityItem    `json:"vulnerabilities"`
	Evidence        ViewReportEvidenceResponse   `json:"evidence"`
}

type ReportDetails struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	SubmittedBy      string    `json:"submitted_by"`
	SubmissionDate   time.Time `json:"submission_date"`
	ProjectName      string    `json:"project_name"`
	Status           string    `json:"status"`
	StatusColor      string    `json:"status_color"`
	ExecutiveSummary string    `json:"executive_summary"`
}

type ReportVulnerabilityItem struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Severity     string `json:"severity"`
	SeverityColor string `json:"severity_color"`
	Domain       string `json:"domain"`
	Status       string `json:"status"`
	StatusColor  string `json:"status_color"`
	Description  string `json:"description"`
	Remediation  string `json:"remediation"`
}

// ============================================================================
// UC39-UC41: Evidence Files Request/Response
// ============================================================================

type ViewReportEvidenceRequest struct {
	ReportID string `query:"report_id" validate:"required,uuid"`
	Page     int    `query:"page" validate:"omitempty,min=1"`
	PerPage  int    `query:"per_page" validate:"omitempty,min=1,max=100"`
}

type ViewReportEvidenceResponse struct {
	Files      []EvidenceFileItem `json:"files"`
	Pagination PaginationInfo     `json:"pagination"`
}

type EvidenceFileItem struct {
	ID         string    `json:"id"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	UploadDate time.Time `json:"upload_date"`
	FilePath   string    `json:"file_path"`
}

// ============================================================================
// UC42: Download Report Request
// ============================================================================

type DownloadReportRequest struct {
	ReportID string `query:"report_id" validate:"required,uuid"`
}

