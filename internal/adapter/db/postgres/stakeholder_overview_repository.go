package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pentsecops/backend/internal/adapter/db/postgres/sqlc"
	"github.com/pentsecops/backend/internal/core/domain"
)

type stakeholderOverviewRepository struct {
	queries *sqlc.Queries
}

// NewStakeholderOverviewRepository creates a new stakeholder overview repository
func NewStakeholderOverviewRepository(db *sql.DB) domain.StakeholderOverviewRepository {
	return &stakeholderOverviewRepository{
		queries: sqlc.New(db),
	}
}

// ============================================================================
// UC1: Calculate overall security score
// ============================================================================

func (r *stakeholderOverviewRepository) CalculateSecurityScore(ctx context.Context) (float64, error) {
	// Get all vulnerabilities with severity and status
	vulns, err := r.queries.GetVulnerabilitiesForSecurityScore(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get vulnerabilities for security score: %w", err)
	}

	if len(vulns) == 0 {
		return 10.0, nil // Perfect score if no vulnerabilities
	}

	// Calculate score based on severity distribution and remediation
	var weights = map[string]float64{
		"critical": 10.0,
		"high":     5.0,
		"medium":   2.0,
		"low":      1.0,
	}

	var maxPossibleScore float64
	var actualScore float64

	for _, vuln := range vulns {
		severity := vuln.Severity
		status := vuln.Status

		weight := weights[severity]
		maxPossibleScore += weight

		// If remediated or verified, add full weight to actual score
		if status == "remediated" || status == "verified" {
			actualScore += weight
		} else if status == "in_progress" {
			// Partial credit for in-progress
			actualScore += weight * 0.5
		}
		// Open vulnerabilities get 0 points
	}

	if maxPossibleScore == 0 {
		return 10.0, nil
	}

	// Calculate score on 0-10 scale
	score := (actualScore / maxPossibleScore) * 10.0

	return score, nil
}

// ============================================================================
// UC2: Get active projects count with breakdown
// ============================================================================

func (r *stakeholderOverviewRepository) GetActiveProjectsCount(ctx context.Context) (totalActive, inProgress, completed int, err error) {
	result, err := r.queries.GetActiveProjectsCounts(ctx)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get active projects count: %w", err)
	}

	return int(result.TotalActive), int(result.InProgress), int(result.Completed), nil
}

// ============================================================================
// UC3: Get critical issues count
// ============================================================================

func (r *stakeholderOverviewRepository) GetCriticalIssuesCount(ctx context.Context) (int, error) {
	count, err := r.queries.GetStakeholderCriticalIssuesCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get critical issues count: %w", err)
	}

	return int(count), nil
}

// ============================================================================
// UC4: Get open vulnerabilities count with trend
// ============================================================================

func (r *stakeholderOverviewRepository) GetOpenVulnerabilitiesCount(ctx context.Context) (current, lastMonth int, err error) {
	// Get current month count
	currentCount, err := r.queries.GetOpenVulnerabilitiesCurrentMonth(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get current open vulnerabilities count: %w", err)
	}

	// Get last month count
	lastMonthCount, err := r.queries.GetOpenVulnerabilitiesLastMonth(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get last month open vulnerabilities count: %w", err)
	}

	return int(currentCount), int(lastMonthCount), nil
}

// ============================================================================
// UC5: Calculate remediation rate
// ============================================================================

func (r *stakeholderOverviewRepository) CalculateRemediationRate(ctx context.Context) (total, remediated int, err error) {
	result, err := r.queries.GetRemediationRateCounts(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get remediation rate counts: %w", err)
	}

	return int(result.Total), int(result.Remediated), nil
}

// ============================================================================
// UC6: Calculate SLA compliance
// ============================================================================

func (r *stakeholderOverviewRepository) CalculateSLACompliance(ctx context.Context) (totalWithDueDate, remediatedOnTime int, err error) {
	result, err := r.queries.GetSLAComplianceCounts(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get SLA compliance counts: %w", err)
	}

	return int(result.TotalWithDueDate), int(result.RemediatedOnTime), nil
}

// ============================================================================
// UC7: Get vulnerability trend data for past 5 months
// ============================================================================

func (r *stakeholderOverviewRepository) GetVulnerabilityTrend(ctx context.Context, months int) ([]domain.MonthlyTrendData, error) {
	rows, err := r.queries.GetVulnerabilityTrendByMonth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerability trend: %w", err)
	}

	trends := make([]domain.MonthlyTrendData, 0, len(rows))
	for _, row := range rows {
		// Convert Unix timestamp to time.Time
		monthTime := time.Unix(row.Month, 0)
		trends = append(trends, domain.MonthlyTrendData{
			Month:    monthTime,
			Open:     int(row.OpenCount),
			Resolved: int(row.ResolvedCount),
		})
	}

	return trends, nil
}

// ============================================================================
// UC8: Get asset status with vulnerability counts by severity
// ============================================================================

func (r *stakeholderOverviewRepository) GetAssetStatus(ctx context.Context) ([]domain.AssetVulnCounts, error) {
	rows, err := r.queries.GetAssetVulnerabilityCounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset vulnerability counts: %w", err)
	}

	assets := make([]domain.AssetVulnCounts, 0, len(rows))
	for _, row := range rows {
		assets = append(assets, domain.AssetVulnCounts{
			AssetName: row.AssetName,
			Critical:  int(row.Critical),
			High:      int(row.High),
			Medium:    int(row.Medium),
			Low:       int(row.Low),
		})
	}

	return assets, nil
}

// ============================================================================
// UC9: Get recent security events
// ============================================================================

func (r *stakeholderOverviewRepository) GetRecentSecurityEvents(ctx context.Context, limit int) ([]domain.SecurityEvent, error) {
	rows, err := r.queries.GetRecentSecurityEvents(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get recent security events: %w", err)
	}

	events := make([]domain.SecurityEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, domain.SecurityEvent{
			ID:          row.ID.String(),
			Description: row.Description,
			EventType:   row.EventType.String,
			Timestamp:   row.CreatedAt.Time,
		})
	}

	return events, nil
}

// ============================================================================
// UC10: Get remediation updates
// ============================================================================

func (r *stakeholderOverviewRepository) GetRemediationUpdates(ctx context.Context, limit int) ([]domain.RemediationUpdateData, error) {
	rows, err := r.queries.GetRemediationUpdates(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get remediation updates: %w", err)
	}

	updates := make([]domain.RemediationUpdateData, 0, len(rows))
	for _, row := range rows {
		var remediatedDate time.Time
		if row.RemediatedDate.Valid {
			remediatedDate = row.RemediatedDate.Time
		}

		updates = append(updates, domain.RemediationUpdateData{
			ID:             row.ID.String(),
			VulnTitle:      row.VulnTitle,
			PreviousStatus: row.PreviousStatus,
			NewStatus:      row.NewStatus,
			RemediatedDate: remediatedDate,
			AssignedTeam:   row.AssignedTeam,
		})
	}

	return updates, nil
}
