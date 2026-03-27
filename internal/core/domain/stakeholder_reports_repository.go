package domain

import (
	"context"
	"time"
)

// StakeholderReportsRepository defines data access methods for stakeholder reports tab
type StakeholderReportsRepository interface {
	// UC28: Get total reports count
	GetTotalReportsCount(ctx context.Context) (int, error)

	// UC29: Get under review reports count
	GetUnderReviewReportsCount(ctx context.Context) (int, error)

	// UC30: Get remediated reports count
	GetRemediatedReportsCount(ctx context.Context) (int, error)

	// UC32: List reports with status filter and pagination
	ListReports(ctx context.Context, status string, limit, offset int) ([]ReportListItem, error)

	// UC32: Get total count for pagination
	GetReportsCount(ctx context.Context, status string) (int, error)

	// UC37: Get report details by ID
	GetReportByID(ctx context.Context, reportID string) (*ReportDetail, error)

	// UC38: Get vulnerabilities for a report
	GetReportVulnerabilities(ctx context.Context, reportID string) ([]ReportVulnerability, error)

	// UC39: Get evidence files for a report with pagination
	GetReportEvidenceFiles(ctx context.Context, reportID string, limit, offset int) ([]EvidenceFile, error)

	// UC39: Get evidence files count for pagination
	GetReportEvidenceFilesCount(ctx context.Context, reportID string) (int, error)

	// UC41: Get evidence file by ID for download
	GetEvidenceFileByID(ctx context.Context, fileID string) (*EvidenceFile, error)

	// UC42: Get report file path for download
	GetReportFilePath(ctx context.Context, reportID string) (string, error)
}

// ReportListItem represents a report in the list
type ReportListItem struct {
	ID                   string
	Title                string
	SubmittedBy          string
	SubmittedDate        time.Time
	ProjectName          string
	Status               string
	VulnerabilitiesCount int
	EvidenceCount        int
}

// ReportDetail represents full report details
type ReportDetail struct {
	ID               string
	Title            string
	SubmittedBy      string
	SubmissionDate   time.Time
	ProjectName      string
	Status           string
	ExecutiveSummary string
}

// ReportVulnerability represents a vulnerability in a report
type ReportVulnerability struct {
	ID          string
	Title       string
	Severity    string
	Domain      string
	Status      string
	Description string
	Remediation string
}

// EvidenceFile represents an evidence file
type EvidenceFile struct {
	ID         string
	FileName   string
	FileSize   int64
	UploadDate time.Time
	FilePath   string
}

