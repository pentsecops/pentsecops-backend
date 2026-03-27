package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/adapter/db/postgres/sqlc"
	"github.com/pentsecops/backend/internal/core/domain"
)

type stakeholderReportsRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewStakeholderReportsRepository creates a new stakeholder reports repository
func NewStakeholderReportsRepository(db *sql.DB) domain.StakeholderReportsRepository {
	return &stakeholderReportsRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// UC28: Get total reports count
func (r *stakeholderReportsRepository) GetTotalReportsCount(ctx context.Context) (int, error) {
	count, err := r.queries.GetStakeholderTotalReportsCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get total reports count: %w", err)
	}
	return int(count), nil
}

// UC29: Get under review reports count
func (r *stakeholderReportsRepository) GetUnderReviewReportsCount(ctx context.Context) (int, error) {
	count, err := r.queries.GetStakeholderUnderReviewReportsCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get under review reports count: %w", err)
	}
	return int(count), nil
}

// UC30: Get remediated reports count
func (r *stakeholderReportsRepository) GetRemediatedReportsCount(ctx context.Context) (int, error) {
	count, err := r.queries.GetStakeholderRemediatedReportsCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get remediated reports count: %w", err)
	}
	return int(count), nil
}

// UC32: List reports with status filter and pagination
func (r *stakeholderReportsRepository) ListReports(ctx context.Context, status string, limit, offset int) ([]domain.ReportListItem, error) {
	// Prepare nullable status parameter
	var statusParam sql.NullString
	if status != "" && status != "all" {
		statusParam = sql.NullString{String: status, Valid: true}
	}

	rows, err := r.queries.ListStakeholderReports(ctx, sqlc.ListStakeholderReportsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		Status: statusParam,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list reports: %w", err)
	}

	// Preallocate slice for performance
	reports := make([]domain.ReportListItem, 0, len(rows))
	for _, row := range rows {
		submittedBy := ""
		if row.SubmittedBy.Valid {
			submittedBy = row.SubmittedBy.String
		}

		projectName := ""
		if row.ProjectName.Valid {
			projectName = row.ProjectName.String
		}

		var submittedDate time.Time
		if row.SubmittedDate.Valid {
			submittedDate = row.SubmittedDate.Time
		}

		status := ""
		if row.Status.Valid {
			status = row.Status.String
		}

		reports = append(reports, domain.ReportListItem{
			ID:                   row.ID.String(),
			Title:                row.Title,
			SubmittedBy:          submittedBy,
			SubmittedDate:        submittedDate,
			ProjectName:          projectName,
			Status:               status,
			VulnerabilitiesCount: int(row.VulnerabilitiesCount),
			EvidenceCount:        int(row.EvidenceCount),
		})
	}

	return reports, nil
}

// UC32: Get total count for pagination
func (r *stakeholderReportsRepository) GetReportsCount(ctx context.Context, status string) (int, error) {
	// Prepare nullable status parameter
	var statusParam sql.NullString
	if status != "" && status != "all" {
		statusParam = sql.NullString{String: status, Valid: true}
	}

	count, err := r.queries.GetStakeholderReportsCount(ctx, statusParam)
	if err != nil {
		return 0, fmt.Errorf("failed to get reports count: %w", err)
	}

	return int(count), nil
}

// UC37: Get report details by ID
func (r *stakeholderReportsRepository) GetReportByID(ctx context.Context, reportID string) (*domain.ReportDetail, error) {
	reportUUID, err := uuid.Parse(reportID)
	if err != nil {
		return nil, fmt.Errorf("invalid report ID: %w", err)
	}

	row, err := r.queries.GetStakeholderReportByID(ctx, reportUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("report not found")
		}
		return nil, fmt.Errorf("failed to get report: %w", err)
	}

	submittedBy := ""
	if row.SubmittedBy.Valid {
		submittedBy = row.SubmittedBy.String
	}

	projectName := ""
	if row.ProjectName.Valid {
		projectName = row.ProjectName.String
	}

	executiveSummary := ""
	if row.ExecutiveSummary.Valid {
		executiveSummary = row.ExecutiveSummary.String
	}

	var submissionDate time.Time
	if row.SubmissionDate.Valid {
		submissionDate = row.SubmissionDate.Time
	}

	status := ""
	if row.Status.Valid {
		status = row.Status.String
	}

	return &domain.ReportDetail{
		ID:               row.ID.String(),
		Title:            row.Title,
		SubmittedBy:      submittedBy,
		SubmissionDate:   submissionDate,
		ProjectName:      projectName,
		Status:           status,
		ExecutiveSummary: executiveSummary,
	}, nil
}

// UC38: Get vulnerabilities for a report
func (r *stakeholderReportsRepository) GetReportVulnerabilities(ctx context.Context, reportID string) ([]domain.ReportVulnerability, error) {
	reportUUID, err := uuid.Parse(reportID)
	if err != nil {
		return nil, fmt.Errorf("invalid report ID: %w", err)
	}

	// Convert to NullUUID
	reportNullUUID := uuid.NullUUID{UUID: reportUUID, Valid: true}

	rows, err := r.queries.GetStakeholderReportVulnerabilities(ctx, reportNullUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report vulnerabilities: %w", err)
	}

	// Preallocate slice for performance
	vulnerabilities := make([]domain.ReportVulnerability, 0, len(rows))
	for _, row := range rows {
		remediation := ""
		if row.Remediation.Valid {
			remediation = row.Remediation.String
		}

		vulnerabilities = append(vulnerabilities, domain.ReportVulnerability{
			ID:          row.ID.String(),
			Title:       row.Title,
			Severity:    row.Severity,
			Domain:      row.Domain,
			Status:      row.Status,
			Description: row.Description,
			Remediation: remediation,
		})
	}

	return vulnerabilities, nil
}

// UC39: Get evidence files for a report with pagination
func (r *stakeholderReportsRepository) GetReportEvidenceFiles(ctx context.Context, reportID string, limit, offset int) ([]domain.EvidenceFile, error) {
	reportUUID, err := uuid.Parse(reportID)
	if err != nil {
		return nil, fmt.Errorf("invalid report ID: %w", err)
	}

	// Convert to NullUUID
	reportNullUUID := uuid.NullUUID{UUID: reportUUID, Valid: true}

	rows, err := r.queries.GetStakeholderReportEvidenceFiles(ctx, sqlc.GetStakeholderReportEvidenceFilesParams{
		ReportID: reportNullUUID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get evidence files: %w", err)
	}

	// Preallocate slice for performance
	files := make([]domain.EvidenceFile, 0, len(rows))
	for _, row := range rows {
		var uploadDate time.Time
		if row.UploadDate.Valid {
			uploadDate = row.UploadDate.Time
		}

		files = append(files, domain.EvidenceFile{
			ID:         row.ID.String(),
			FileName:   row.FileName,
			FileSize:   row.FileSize,
			UploadDate: uploadDate,
			FilePath:   row.FilePath,
		})
	}

	return files, nil
}

// UC39: Get evidence files count for pagination
func (r *stakeholderReportsRepository) GetReportEvidenceFilesCount(ctx context.Context, reportID string) (int, error) {
	reportUUID, err := uuid.Parse(reportID)
	if err != nil {
		return 0, fmt.Errorf("invalid report ID: %w", err)
	}

	// Convert to NullUUID
	reportNullUUID := uuid.NullUUID{UUID: reportUUID, Valid: true}

	count, err := r.queries.GetStakeholderReportEvidenceFilesCount(ctx, reportNullUUID)
	if err != nil {
		return 0, fmt.Errorf("failed to get evidence files count: %w", err)
	}

	return int(count), nil
}

// UC41: Get evidence file by ID for download
func (r *stakeholderReportsRepository) GetEvidenceFileByID(ctx context.Context, fileID string) (*domain.EvidenceFile, error) {
	fileUUID, err := uuid.Parse(fileID)
	if err != nil {
		return nil, fmt.Errorf("invalid file ID: %w", err)
	}

	row, err := r.queries.GetStakeholderEvidenceFileByID(ctx, fileUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("evidence file not found")
		}
		return nil, fmt.Errorf("failed to get evidence file: %w", err)
	}

	var uploadDate time.Time
	if row.UploadDate.Valid {
		uploadDate = row.UploadDate.Time
	}

	return &domain.EvidenceFile{
		ID:         row.ID.String(),
		FileName:   row.FileName,
		FileSize:   row.FileSize,
		UploadDate: uploadDate,
		FilePath:   row.FilePath,
	}, nil
}

// UC42: Get report file path for download
func (r *stakeholderReportsRepository) GetReportFilePath(ctx context.Context, reportID string) (string, error) {
	reportUUID, err := uuid.Parse(reportID)
	if err != nil {
		return "", fmt.Errorf("invalid report ID: %w", err)
	}

	filePath, err := r.queries.GetStakeholderReportFilePath(ctx, reportUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("report not found")
		}
		return "", fmt.Errorf("failed to get report file path: %w", err)
	}

	if !filePath.Valid || filePath.String == "" {
		return "", fmt.Errorf("report file not available")
	}

	return filePath.String, nil
}
