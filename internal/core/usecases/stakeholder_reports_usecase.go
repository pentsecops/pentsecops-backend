package usecases

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
)

type stakeholderReportsUseCase struct {
	repo domain.StakeholderReportsRepository
}

// NewStakeholderReportsUseCase creates a new stakeholder reports use case
func NewStakeholderReportsUseCase(repo domain.StakeholderReportsRepository) domain.StakeholderReportsUseCase {
	return &stakeholderReportsUseCase{
		repo: repo,
	}
}

// UC28-UC30: Get reports statistics
func (uc *stakeholderReportsUseCase) GetReportsStats(ctx context.Context) (*dto.StakeholderReportsStatsResponse, error) {
	// UC28: Get total reports count
	totalReports, err := uc.repo.GetTotalReportsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total reports count: %w", err)
	}

	// UC29: Get under review reports count
	underReview, err := uc.repo.GetUnderReviewReportsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get under review reports count: %w", err)
	}

	// UC30: Get remediated reports count
	remediated, err := uc.repo.GetRemediatedReportsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get remediated reports count: %w", err)
	}

	return &dto.StakeholderReportsStatsResponse{
		TotalReports: dto.ReportStatCard{
			Count:   totalReports,
			Message: "All submitted reports",
			Color:   "blue",
		},
		UnderReview: dto.ReportStatCard{
			Count:   underReview,
			Message: "Currently being reviewed",
			Color:   "orange",
		},
		Remediated: dto.ReportStatCard{
			Count:   remediated,
			Message: "Completed and remediated",
			Color:   "green",
		},
	}, nil
}

// UC31-UC36: List reports with status filter and pagination
func (uc *stakeholderReportsUseCase) ListReports(ctx context.Context, req *dto.ListStakeholderReportsRequest) (*dto.ListStakeholderReportsResponse, error) {
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
	totalCount, err := uc.repo.GetReportsCount(ctx, req.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports count: %w", err)
	}

	// Get reports
	reports, err := uc.repo.ListReports(ctx, req.Status, req.PerPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list reports: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(req.PerPage)))

	// Preallocate response slice
	items := make([]dto.StakeholderReportItem, 0, len(reports))

	for _, report := range reports {
		// UC34: Determine status color
		statusColor := getReportStatusColor(report.Status)

		items = append(items, dto.StakeholderReportItem{
			ID:                   report.ID,
			Title:                report.Title,
			SubmittedBy:          report.SubmittedBy,
			SubmittedDate:        report.SubmittedDate,
			Project:              report.ProjectName,
			Status:               report.Status,
			StatusColor:          statusColor,
			VulnerabilitiesCount: report.VulnerabilitiesCount, // UC35
			EvidenceCount:        report.EvidenceCount,        // UC36
		})
	}

	return &dto.ListStakeholderReportsResponse{
		Reports: items,
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

// UC37-UC38: View report details with vulnerabilities
func (uc *stakeholderReportsUseCase) ViewReport(ctx context.Context, reportID string, evidencePage, evidencePerPage int) (*dto.ViewStakeholderReportResponse, error) {
	// UC37: Get report details
	report, err := uc.repo.GetReportByID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}

	// UC38: Get vulnerabilities for the report
	vulnerabilities, err := uc.repo.GetReportVulnerabilities(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report vulnerabilities: %w", err)
	}

	// Convert vulnerabilities to DTO
	vulnItems := make([]dto.ReportVulnerabilityItem, 0, len(vulnerabilities))
	for _, vuln := range vulnerabilities {
		vulnItems = append(vulnItems, dto.ReportVulnerabilityItem{
			ID:            vuln.ID,
			Title:         vuln.Title,
			Severity:      vuln.Severity,
			SeverityColor: getReportSeverityColor(vuln.Severity),
			Domain:        vuln.Domain,
			Status:        vuln.Status,
			StatusColor:   getReportVulnStatusColor(vuln.Status),
			Description:   vuln.Description,
			Remediation:   vuln.Remediation,
		})
	}

	// UC39-UC40: Get evidence files with pagination
	evidenceResponse, err := uc.GetReportEvidenceFiles(ctx, reportID, evidencePage, evidencePerPage)
	if err != nil {
		return nil, fmt.Errorf("failed to get evidence files: %w", err)
	}

	return &dto.ViewStakeholderReportResponse{
		Report: dto.ReportDetails{
			ID:               report.ID,
			Title:            report.Title,
			SubmittedBy:      report.SubmittedBy,
			SubmissionDate:   report.SubmissionDate,
			ProjectName:      report.ProjectName,
			Status:           report.Status,
			StatusColor:      getReportStatusColor(report.Status),
			ExecutiveSummary: report.ExecutiveSummary,
		},
		Vulnerabilities: vulnItems,
		Evidence:        *evidenceResponse,
	}, nil
}

// UC39-UC40: Get evidence files for a report
func (uc *stakeholderReportsUseCase) GetReportEvidenceFiles(ctx context.Context, reportID string, page, perPage int) (*dto.ViewReportEvidenceResponse, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 3
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Get total count for pagination
	totalCount, err := uc.repo.GetReportEvidenceFilesCount(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get evidence files count: %w", err)
	}

	// Get evidence files
	files, err := uc.repo.GetReportEvidenceFiles(ctx, reportID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get evidence files: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))

	// Convert to DTO
	fileItems := make([]dto.EvidenceFileItem, 0, len(files))
	for _, file := range files {
		fileItems = append(fileItems, dto.EvidenceFileItem{
			ID:         file.ID,
			FileName:   file.FileName,
			FileSize:   file.FileSize,
			UploadDate: file.UploadDate,
			FilePath:   file.FilePath,
		})
	}

	return &dto.ViewReportEvidenceResponse{
		Files: fileItems,
		Pagination: dto.PaginationInfo{
			CurrentPage: page,
			PerPage:     perPage,
			Total:       int64(totalCount),
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrev:     page > 1,
		},
	}, nil
}

// UC41: Download evidence file
func (uc *stakeholderReportsUseCase) DownloadEvidenceFile(ctx context.Context, fileID string) (*domain.EvidenceFile, error) {
	file, err := uc.repo.GetEvidenceFileByID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get evidence file: %w", err)
	}

	return file, nil
}

// UC42: Download report
func (uc *stakeholderReportsUseCase) DownloadReport(ctx context.Context, reportID string) ([]byte, string, error) {
	// Get report file path
	filePath, err := uc.repo.GetReportFilePath(ctx, reportID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get report file path: %w", err)
	}

	// Read file content
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read report file: %w", err)
	}

	// Extract filename from path
	filename := fmt.Sprintf("report_%s.pdf", reportID)

	return fileContent, filename, nil
}

// Helper: Get report status color (UC34)
func getReportStatusColor(status string) string {
	switch status {
	case "received":
		return "gray"
	case "under_review":
		return "blue"
	case "shared", "remediated":
		return "green"
	default:
		return "gray"
	}
}

// Helper: Get severity color for reports
func getReportSeverityColor(severity string) string {
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

// Helper: Get vulnerability status color for reports
func getReportVulnStatusColor(status string) string {
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
