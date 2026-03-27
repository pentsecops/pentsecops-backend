package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
)

// AdminOverviewUseCase implements the AdminOverviewUseCase interface
type AdminOverviewUseCase struct {
	repo domain.AdminOverviewRepository
}

// NewAdminOverviewUseCase creates a new AdminOverviewUseCase
func NewAdminOverviewUseCase(repo domain.AdminOverviewRepository) domain.AdminOverviewUseCase {
	return &AdminOverviewUseCase{
		repo: repo,
	}
}

// GetOverviewStats returns all overview statistics
func (uc *AdminOverviewUseCase) GetOverviewStats(ctx context.Context) (*dto.AdminOverviewStatsResponse, error) {
	// Fetch all statistics from database
	totalProjects, err := uc.repo.GetTotalProjectsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total projects: %w", err)
	}

	openProjects, err := uc.repo.GetOpenProjectsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get open projects: %w", err)
	}

	inProgressProjects, err := uc.repo.GetInProgressProjectsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get in-progress projects: %w", err)
	}

	completedProjects, err := uc.repo.GetCompletedProjectsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed projects: %w", err)
	}

	totalVulns, err := uc.repo.GetTotalVulnerabilitiesCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total vulnerabilities: %w", err)
	}

	activeUsers, err := uc.repo.GetActiveUsersCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}

	activePentesters, err := uc.repo.GetActivePentestersCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active pentesters: %w", err)
	}

	activeStakeholders, err := uc.repo.GetActiveStakeholdersCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active stakeholders: %w", err)
	}

	openIssues, err := uc.repo.GetOpenIssuesCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get open issues: %w", err)
	}

	criticalIssues, err := uc.repo.GetCriticalIssuesCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get critical issues: %w", err)
	}

	// Build response
	stats := &dto.AdminOverviewStatsResponse{
		TotalProjects: dto.TotalProjectsStats{
			Total:      totalProjects,
			Open:       openProjects,
			InProgress: inProgressProjects,
			Breakdown:  fmt.Sprintf("%d open, %d in progress", openProjects, inProgressProjects),
		},
		TotalVulnerabilities: dto.TotalVulnerabilitiesStats{
			Total:    totalVulns,
			Subtitle: "Across all domains",
		},
		ActiveUsers: dto.ActiveUsersStats{
			Total:        activeUsers,
			Pentesters:   activePentesters,
			Stakeholders: activeStakeholders,
			Breakdown:    fmt.Sprintf("%d pentesters, %d stakeholders", activePentesters, activeStakeholders),
		},
		OpenIssues: dto.OpenIssuesStats{
			Count:    openIssues,
			Subtitle: "Require attention",
		},
		CompletedProjects: dto.CompletedProjectsStats{
			Count:    completedProjects,
			Subtitle: "Successfully finished",
		},
		CriticalIssues: dto.CriticalIssuesStats{
			Count:    criticalIssues,
			Subtitle: "High priority",
		},
	}

	return stats, nil
}

// GetVulnerabilitiesBySeverity returns vulnerabilities grouped by severity
func (uc *AdminOverviewUseCase) GetVulnerabilitiesBySeverity(ctx context.Context) (*dto.VulnerabilitiesBySeverityResponse, error) {
	severityCounts, err := uc.repo.GetVulnerabilitiesBySeverity(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vulnerabilities by severity: %w", err)
	}

	// Convert to DTO
	data := make([]dto.SeverityCount, len(severityCounts))
	for i, sc := range severityCounts {
		data[i] = dto.SeverityCount{
			Severity: sc.Severity,
			Count:    sc.Count,
		}
	}

	response := &dto.VulnerabilitiesBySeverityResponse{
		Data: data,
	}

	return response, nil
}

// GetTop5Domains returns top 5 domains by vulnerability count
func (uc *AdminOverviewUseCase) GetTop5Domains(ctx context.Context) (*dto.Top5DomainsResponse, error) {
	domainCounts, err := uc.repo.GetTop5DomainsByVulnerabilities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch top 5 domains: %w", err)
	}

	// Convert to DTO
	data := make([]dto.DomainVulnerabilityCount, len(domainCounts))
	for i, dc := range domainCounts {
		data[i] = dto.DomainVulnerabilityCount{
			Domain:             dc.Domain,
			VulnerabilityCount: dc.Count,
		}
	}

	response := &dto.Top5DomainsResponse{
		Data: data,
	}

	return response, nil
}

// GetProjectStatusDistribution returns project status distribution
func (uc *AdminOverviewUseCase) GetProjectStatusDistribution(ctx context.Context) (*dto.ProjectStatusDistributionResponse, error) {
	statusCounts, err := uc.repo.GetProjectStatusDistribution(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project status distribution: %w", err)
	}

	// Convert to DTO
	data := make([]dto.StatusCount, len(statusCounts))
	for i, sc := range statusCounts {
		data[i] = dto.StatusCount{
			Status: sc.Status,
			Count:  sc.Count,
		}
	}

	response := &dto.ProjectStatusDistributionResponse{
		Data: data,
	}

	return response, nil
}

// GetRecentActivity returns recent activity logs with pagination
func (uc *AdminOverviewUseCase) GetRecentActivity(ctx context.Context, page, perPage int) (*dto.RecentActivityResponse, error) {
	// Calculate offset
	offset := (page - 1) * perPage

	// Fetch activity logs
	logs, err := uc.repo.GetRecentActivityLogs(ctx, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recent activity: %w", err)
	}

	// Get total count
	total, err := uc.repo.GetActivityLogsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch activity logs count: %w", err)
	}

	// Convert to DTOs
	activities := make([]dto.ActivityLogDTO, len(logs))
	for i, log := range logs {
		var userID, userName, entityType, entityID, ipAddress, userAgent *string

		if log.UserID != nil {
			id := log.UserID.String()
			userID = &id
		}
		if log.UserEmail != "" {
			userName = &log.UserEmail
		}
		if log.EntityType != nil {
			entityType = log.EntityType
		}
		if log.EntityID != nil {
			id := log.EntityID.String()
			entityID = &id
		}
		if log.IPAddress != "" {
			ipAddress = &log.IPAddress
		}
		if log.UserAgent != "" {
			userAgent = &log.UserAgent
		}

		activities[i] = dto.ActivityLogDTO{
			ID:         log.ID.String(),
			UserID:     userID,
			UserName:   userName,
			Action:     log.Action,
			EntityType: entityType,
			EntityID:   entityID,
			IPAddress:  ipAddress,
			UserAgent:  userAgent,
			CreatedAt:  log.CreatedAt.Format(time.RFC3339),
		}
	}

	// Calculate pagination metadata
	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	response := &dto.RecentActivityResponse{
		Activities: activities,
		Pagination: dto.PaginationMeta{
			CurrentPage: page,
			PerPage:     perPage,
			Total:       total,
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrev:     page > 1,
		},
	}

	return response, nil
}
