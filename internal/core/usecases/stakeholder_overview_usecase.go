package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/pkg/auth/logger"
)

type stakeholderOverviewUseCase struct {
	repo domain.StakeholderOverviewRepository
}

// NewStakeholderOverviewUseCase creates a new stakeholder overview use case
func NewStakeholderOverviewUseCase(repo domain.StakeholderOverviewRepository) domain.StakeholderOverviewUseCase {
	return &stakeholderOverviewUseCase{
		repo: repo,
	}
}

// ============================================================================
// UC1-UC6: Get all security metrics cards
// ============================================================================

func (uc *stakeholderOverviewUseCase) GetSecurityMetrics(ctx context.Context) (*dto.StakeholderSecurityMetricsResponse, error) {
	logger.Info("Fetching stakeholder security metrics")

	// UC1: Overall Security Score
	score, err := uc.repo.CalculateSecurityScore(ctx)
	if err != nil {
		logger.Error("Failed to calculate security score", "error", err)
		return nil, fmt.Errorf("failed to calculate security score: %w", err)
	}

	status := uc.getSecurityStatus(score)

	// UC2: Active Projects Count
	totalActive, inProgress, completed, err := uc.repo.GetActiveProjectsCount(ctx)
	if err != nil {
		logger.Error("Failed to get active projects count", "error", err)
		return nil, fmt.Errorf("failed to get active projects count: %w", err)
	}

	completedText := fmt.Sprintf("%d in progress, %d completed", inProgress, completed)

	// UC3: Critical Issues Count
	criticalCount, err := uc.repo.GetCriticalIssuesCount(ctx)
	if err != nil {
		logger.Error("Failed to get critical issues count", "error", err)
		return nil, fmt.Errorf("failed to get critical issues count: %w", err)
	}

	// UC4: Open Vulnerabilities Count with Trend
	currentOpen, lastMonthOpen, err := uc.repo.GetOpenVulnerabilitiesCount(ctx)
	if err != nil {
		logger.Error("Failed to get open vulnerabilities count", "error", err)
		return nil, fmt.Errorf("failed to get open vulnerabilities count: %w", err)
	}

	difference := currentOpen - lastMonthOpen
	trendInfo := uc.getTrendInfo(difference, lastMonthOpen)

	// UC5: Remediation Rate
	totalVulns, remediatedVulns, err := uc.repo.CalculateRemediationRate(ctx)
	if err != nil {
		logger.Error("Failed to calculate remediation rate", "error", err)
		return nil, fmt.Errorf("failed to calculate remediation rate: %w", err)
	}

	var remediationRate float64
	if totalVulns > 0 {
		remediationRate = (float64(remediatedVulns) / float64(totalVulns)) * 100
	}

	// UC6: SLA Compliance
	totalWithDueDate, remediatedOnTime, err := uc.repo.CalculateSLACompliance(ctx)
	if err != nil {
		logger.Error("Failed to calculate SLA compliance", "error", err)
		return nil, fmt.Errorf("failed to calculate SLA compliance: %w", err)
	}

	var slaCompliance float64
	if totalWithDueDate > 0 {
		slaCompliance = (float64(remediatedOnTime) / float64(totalWithDueDate)) * 100
	}

	response := &dto.StakeholderSecurityMetricsResponse{
		SecurityScore: dto.SecurityScoreResponse{
			Score:  score,
			Status: status,
		},
		ActiveProjects: dto.ActiveProjectsCountResponse{
			TotalActive:   totalActive,
			InProgress:    inProgress,
			Completed:     completed,
			CompletedText: completedText,
		},
		CriticalIssues: dto.CriticalIssuesResponse{
			Count:   criticalCount,
			Message: "Requires immediate attention",
		},
		OpenVulnerabilities: dto.OpenVulnerabilitiesResponse{
			Count:      currentOpen,
			TrendInfo:  trendInfo,
			LastMonth:  lastMonthOpen,
			Difference: difference,
		},
		RemediationRate: dto.RemediationRateResponse{
			Rate:            remediationRate,
			TotalVulns:      totalVulns,
			RemediatedVulns: remediatedVulns,
		},
		SLACompliance: dto.SLAComplianceStakeholderResponse{
			ComplianceRate:   slaCompliance,
			TotalWithDueDate: totalWithDueDate,
			RemediatedOnTime: remediatedOnTime,
		},
	}

	logger.Info("Successfully fetched stakeholder security metrics")
	return response, nil
}

// ============================================================================
// UC7: Get vulnerability trend chart data
// ============================================================================

func (uc *stakeholderOverviewUseCase) GetVulnerabilityTrend(ctx context.Context) (*dto.VulnerabilityTrendChartResponse, error) {
	logger.Info("Fetching vulnerability trend chart data")

	trends, err := uc.repo.GetVulnerabilityTrend(ctx, 5)
	if err != nil {
		logger.Error("Failed to get vulnerability trend", "error", err)
		return nil, fmt.Errorf("failed to get vulnerability trend: %w", err)
	}

	monthlyTrends := make([]dto.MonthlyVulnerabilityTrend, 0, len(trends))
	for _, trend := range trends {
		monthlyTrends = append(monthlyTrends, dto.MonthlyVulnerabilityTrend{
			Month:    trend.Month.Format("Jan"),
			Open:     trend.Open,
			Resolved: trend.Resolved,
		})
	}

	response := &dto.VulnerabilityTrendChartResponse{
		Trends: monthlyTrends,
	}

	logger.Info("Successfully fetched vulnerability trend chart data", "count", len(monthlyTrends))
	return response, nil
}

// ============================================================================
// UC8: Get asset status chart data
// ============================================================================

func (uc *stakeholderOverviewUseCase) GetAssetStatus(ctx context.Context) (*dto.AssetStatusChartResponse, error) {
	logger.Info("Fetching asset status chart data")

	assets, err := uc.repo.GetAssetStatus(ctx)
	if err != nil {
		logger.Error("Failed to get asset status", "error", err)
		return nil, fmt.Errorf("failed to get asset status: %w", err)
	}

	assetCounts := make([]dto.AssetVulnerabilityCounts, 0, len(assets))
	for _, asset := range assets {
		total := asset.Critical + asset.High + asset.Medium + asset.Low
		assetCounts = append(assetCounts, dto.AssetVulnerabilityCounts{
			AssetName: asset.AssetName,
			Critical:  asset.Critical,
			High:      asset.High,
			Medium:    asset.Medium,
			Low:       asset.Low,
			Total:     total,
		})
	}

	response := &dto.AssetStatusChartResponse{
		Assets: assetCounts,
	}

	logger.Info("Successfully fetched asset status chart data", "count", len(assetCounts))
	return response, nil
}

// ============================================================================
// UC9: Get recent security events
// ============================================================================

func (uc *stakeholderOverviewUseCase) GetRecentSecurityEvents(ctx context.Context, limit int) (*dto.RecentSecurityEventsResponse, error) {
	logger.Info("Fetching recent security events", "limit", limit)

	events, err := uc.repo.GetRecentSecurityEvents(ctx, limit)
	if err != nil {
		logger.Error("Failed to get recent security events", "error", err)
		return nil, fmt.Errorf("failed to get recent security events: %w", err)
	}

	recentEvents := make([]dto.RecentSecurityEvent, 0, len(events))
	for _, event := range events {
		recentEvents = append(recentEvents, dto.RecentSecurityEvent{
			ID:          event.ID,
			Description: event.Description,
			EventType:   event.EventType,
			Timestamp:   event.Timestamp.Format(time.RFC3339),
		})
	}

	response := &dto.RecentSecurityEventsResponse{
		Events: recentEvents,
		Total:  len(recentEvents),
	}

	logger.Info("Successfully fetched recent security events", "count", len(recentEvents))
	return response, nil
}

// ============================================================================
// UC10: Get remediation updates
// ============================================================================

func (uc *stakeholderOverviewUseCase) GetRemediationUpdates(ctx context.Context, limit int) (*dto.RemediationUpdatesResponse, error) {
	logger.Info("Fetching remediation updates", "limit", limit)

	updates, err := uc.repo.GetRemediationUpdates(ctx, limit)
	if err != nil {
		logger.Error("Failed to get remediation updates", "error", err)
		return nil, fmt.Errorf("failed to get remediation updates: %w", err)
	}

	remediationUpdates := make([]dto.RemediationUpdate, 0, len(updates))
	for _, update := range updates {
		remediationUpdates = append(remediationUpdates, dto.RemediationUpdate{
			ID:             update.ID,
			VulnTitle:      update.VulnTitle,
			PreviousStatus: update.PreviousStatus,
			NewStatus:      update.NewStatus,
			RemediatedDate: update.RemediatedDate.Format(time.RFC3339),
			AssignedTeam:   update.AssignedTeam,
		})
	}

	response := &dto.RemediationUpdatesResponse{
		Updates: remediationUpdates,
		Total:   len(remediationUpdates),
	}

	logger.Info("Successfully fetched remediation updates", "count", len(remediationUpdates))
	return response, nil
}

// ============================================================================
// Helper functions
// ============================================================================

func (uc *stakeholderOverviewUseCase) getSecurityStatus(score float64) string {
	if score >= 8.0 {
		return "Excellent - Improving"
	} else if score >= 6.0 {
		return "Good - Stable"
	} else if score >= 4.0 {
		return "Fair - Needs Attention"
	} else {
		return "Poor - Critical Action Required"
	}
}

func (uc *stakeholderOverviewUseCase) getTrendInfo(difference, lastMonth int) string {
	if difference < 0 {
		return fmt.Sprintf("Down from %d last month", lastMonth)
	} else if difference > 0 {
		return fmt.Sprintf("Up from %d last month", lastMonth)
	} else {
		return fmt.Sprintf("Same as last month (%d)", lastMonth)
	}
}

