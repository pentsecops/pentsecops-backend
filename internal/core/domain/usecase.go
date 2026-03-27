package domain

import (
	"context"

	"github.com/pentsecops/backend/internal/core/domain/dto"
)

// AdminOverviewUseCase defines the interface for admin overview business logic
type AdminOverviewUseCase interface {
	GetOverviewStats(ctx context.Context) (*dto.AdminOverviewStatsResponse, error)
	GetVulnerabilitiesBySeverity(ctx context.Context) (*dto.VulnerabilitiesBySeverityResponse, error)
	GetTop5Domains(ctx context.Context) (*dto.Top5DomainsResponse, error)
	GetProjectStatusDistribution(ctx context.Context) (*dto.ProjectStatusDistributionResponse, error)
	GetRecentActivity(ctx context.Context, page, perPage int) (*dto.RecentActivityResponse, error)
}

// UsersUseCase defines the interface for users management business logic
type UsersUseCase interface {
	// User management
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	ListUsers(ctx context.Context, page, perPage int) (*dto.ListUsersResponse, error)
	UpdateUser(ctx context.Context, userID string, req *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, userID string) error
	RefreshUsers(ctx context.Context, page, perPage int) (*dto.ListUsersResponse, error)

	// Statistics
	GetUserStats(ctx context.Context) (*dto.UserStatsResponse, error)

	// Export
	ExportUsersToCSV(ctx context.Context) ([]byte, error)
}

// ProjectsUseCase defines the interface for projects management business logic
type ProjectsUseCase interface {
	// Project management
	CreateProject(ctx context.Context, req *dto.CreateProjectRequest, createdBy string) (*dto.CreateProjectResponse, error)
	ListProjects(ctx context.Context, page, perPage int) (*dto.ListProjectsResponse, error)
	UpdateProject(ctx context.Context, projectID string, req *dto.UpdateProjectRequest) (*dto.UpdateProjectResponse, error)
	DeleteProject(ctx context.Context, projectID string) error

	// Statistics
	GetProjectStats(ctx context.Context) (*dto.ProjectStatsResponse, error)

	// Helpers
	GetPentesters(ctx context.Context) (*dto.GetPentestersResponse, error)
}

// TasksUseCase defines the interface for tasks management business logic
type TasksUseCase interface {
	// Task management
	CreateTask(ctx context.Context, req *dto.CreateTaskRequest) (*dto.CreateTaskResponse, error)
	ListTasksByProject(ctx context.Context, projectID string) (*dto.ListTasksResponse, error)
	ListAllTasks(ctx context.Context) (*dto.TaskBoardResponse, error)
	UpdateTaskStatus(ctx context.Context, taskID string, req *dto.UpdateTaskStatusRequest) (*dto.UpdateTaskStatusResponse, error)
	DeleteTask(ctx context.Context, taskID string) error
}

// VulnerabilitiesUseCase defines the interface for vulnerabilities management business logic
type VulnerabilitiesUseCase interface {
	// Vulnerability management
	CreateVulnerability(ctx context.Context, req *dto.CreateVulnerabilityRequest) (*dto.VulnerabilityResponse, error)
	GetVulnerabilityByID(ctx context.Context, id string) (*dto.VulnerabilityResponse, error)
	UpdateVulnerability(ctx context.Context, id string, req *dto.UpdateVulnerabilityRequest) (*dto.VulnerabilityResponse, error)
	DeleteVulnerability(ctx context.Context, id string) error
	ListVulnerabilities(ctx context.Context, req *dto.ListVulnerabilitiesRequest) (*dto.ListVulnerabilitiesResponse, error)

	// Statistics
	GetVulnerabilityStats(ctx context.Context) (*dto.VulnerabilityStatsResponse, error)
	GetSLACompliance(ctx context.Context) (*dto.SLAComplianceResponse, error)

	// Export
	ExportVulnerabilitiesToCSV(ctx context.Context, req *dto.ListVulnerabilitiesRequest) ([]byte, error)
}

// DomainsUseCase defines the interface for domains management business logic
type DomainsUseCase interface {
	// Statistics
	GetDomainsStats(ctx context.Context) (*dto.DomainsStatsResponse, error)

	// Domains management
	ListDomains(ctx context.Context, req *dto.ListDomainsRequest) (*dto.ListDomainsResponse, error)
	GetDomainByID(ctx context.Context, id string) (*dto.DomainResponse, error)
	CreateDomain(ctx context.Context, req *dto.CreateDomainRequest) (*dto.DomainResponse, error)
	UpdateDomain(ctx context.Context, id string, req *dto.UpdateDomainRequest) (*dto.DomainResponse, error)
	DeleteDomain(ctx context.Context, id string) error

	// Security Metrics
	GetSecurityMetrics(ctx context.Context) (*dto.SecurityMetricsResponse, error)
	CreateSecurityMetric(ctx context.Context, req *dto.CreateSecurityMetricRequest) error

	// SLA Breach Analysis
	GetSLABreachAnalysis(ctx context.Context) (*dto.SLABreachAnalysisResponse, error)
}

// NotificationsUseCase defines the interface for notifications management business logic
type NotificationsUseCase interface {
	// Notifications
	GetTotalNotificationsSent(ctx context.Context, createdBy string) (*dto.TotalNotificationsSentResponse, error)
	ListNotifications(ctx context.Context, createdBy string, req *dto.ListNotificationsRequest) (*dto.ListNotificationsResponse, error)
	CreateNotification(ctx context.Context, createdBy string, req *dto.CreateNotificationRequest) (*dto.NotificationResponse, error)

	// Alerts
	ListImportantAlerts(ctx context.Context, req *dto.ListAlertsRequest) (*dto.ListAlertsResponse, error)
}

// AuthUseCase defines the interface for authentication business logic
type AuthUseCase interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) (*dto.ChangePasswordResponse, error)
	Logout(ctx context.Context, req *dto.LogoutRequest) (*dto.LogoutResponse, error)
}

// PentesterOverviewUseCase defines the interface for pentester overview business logic
type PentesterOverviewUseCase interface {
	GetOverviewStats(ctx context.Context, pentesterID string) (*dto.PentesterOverviewStatsResponse, error)
	GetActiveProjects(ctx context.Context, pentesterID string, limit int) (*dto.PentesterActiveProjectsResponse, error)
	GetRecentVulnerabilities(ctx context.Context, pentesterID string, limit int) (*dto.PentesterRecentVulnerabilitiesResponse, error)
	GetUpcomingDeadlines(ctx context.Context, pentesterID string) (*dto.PentesterUpcomingDeadlinesResponse, error)
}

// PentesterProjectsUseCase defines the interface for pentester projects business logic
type PentesterProjectsUseCase interface {
	// UC9: Fetch and Display Assigned Projects List
	GetAssignedProjects(ctx context.Context, pentesterID string) (*dto.GetPentesterProjectsResponse, error)
	GetProjectDetails(ctx context.Context, projectID string, pentesterID string) (*dto.GetPentesterProjectDetailsResponse, error)

	// UC12: Display Project Assets List
	GetProjectAssets(ctx context.Context, projectID string) (*dto.GetProjectAssetsResponse, error)

	// UC13: Display Project Requirements List
	GetProjectRequirements(ctx context.Context, projectID string) (*dto.GetProjectRequirementsResponse, error)
}

// PentesterTasksUseCase defines the interface for pentester tasks business logic
type PentesterTasksUseCase interface {
	// UC16: Fetch and Display Task Board
	GetTaskBoard(ctx context.Context, pentesterID string) (*dto.GetPentesterTaskBoardResponse, error)

	// UC17: Create New Task
	CreateTask(ctx context.Context, pentesterID string, req *dto.CreatePentesterTaskRequest) (*dto.CreatePentesterTaskResponse, error)

	// UC18: Update Task Status
	UpdateTaskStatus(ctx context.Context, taskID string, pentesterID string, req *dto.UpdatePentesterTaskStatusRequest) (*dto.UpdatePentesterTaskStatusResponse, error)

	// UC19: Edit Task Details
	UpdateTask(ctx context.Context, taskID string, pentesterID string, req *dto.UpdatePentesterTaskRequest) (*dto.UpdatePentesterTaskResponse, error)

	// UC20: Delete Task
	DeleteTask(ctx context.Context, taskID string, pentesterID string) (*dto.DeletePentesterTaskResponse, error)

	// Get single task details
	GetTaskDetails(ctx context.Context, taskID string, pentesterID string) (*dto.GetPentesterTaskDetailsResponse, error)

	// Get projects for dropdown
	GetProjectsForDropdown(ctx context.Context, pentesterID string) (*dto.GetPentesterProjectsDropdownResponse, error)
}

// PentesterSubmitReportUseCase defines the interface for pentester submit report business logic
type PentesterSubmitReportUseCase interface {
	// UC22: Get projects for report dropdown
	GetProjectsForReportDropdown(ctx context.Context, pentesterID string) (*dto.ProjectsForReportDropdownResponse, error)

	// UC23: Submit vulnerability report
	SubmitReport(ctx context.Context, pentesterID string, req *dto.SubmitReportRequest) (*dto.SubmitReportResponse, error)

	// UC27: Get submitted reports history
	GetSubmittedReportsHistory(ctx context.Context, pentesterID string, page, perPage int) (*dto.SubmittedReportsHistoryResponse, error)

	// UC31: Get report details
	GetReportDetails(ctx context.Context, reportID string, pentesterID string) (*dto.ReportDetailsResponse, error)

	// UC30: Resubmit rejected report
	ResubmitReport(ctx context.Context, reportID string, pentesterID string, req *dto.ResubmitReportRequest) (*dto.ResubmitReportResponse, error)
}

// PentesterAlertsUseCase defines the interface for pentester alerts business logic
type PentesterAlertsUseCase interface {
	// UC33: Get alert statistics
	GetAlertStats(ctx context.Context, pentesterID string) (*dto.AlertStatsResponse, error)

	// UC34: Get all alerts with filters and pagination
	GetAlerts(ctx context.Context, pentesterID string, alertType, search *string, page, perPage int) (*dto.AlertsListResponse, error)

	// UC35: Get alert types
	GetAlertTypes(ctx context.Context) (*dto.AlertTypesResponse, error)

	// UC36: Mark alert as read/unread
	MarkAlertAsRead(ctx context.Context, alertID string, pentesterID string, isRead bool) (*dto.MarkAlertReadResponse, error)

	// UC37: Dismiss alert
	DismissAlert(ctx context.Context, alertID string, pentesterID string) (*dto.DismissAlertResponse, error)

	// UC38: Create alert to admin
	CreateAlertToAdmin(ctx context.Context, pentesterID string, req *dto.CreateAlertToAdminRequest) (*dto.CreateAlertToAdminResponse, error)

	// UC39: Get alert guidelines
	GetAlertGuidelines(ctx context.Context) (*dto.AlertGuidelinesResponse, error)
}

// StakeholderOverviewUseCase defines the interface for stakeholder overview business logic
type StakeholderOverviewUseCase interface {
	// UC1-UC6: Get all security metrics cards
	GetSecurityMetrics(ctx context.Context) (*dto.StakeholderSecurityMetricsResponse, error)

	// UC7: Get vulnerability trend chart data
	GetVulnerabilityTrend(ctx context.Context) (*dto.VulnerabilityTrendChartResponse, error)

	// UC8: Get asset status chart data
	GetAssetStatus(ctx context.Context) (*dto.AssetStatusChartResponse, error)

	// UC9: Get recent security events
	GetRecentSecurityEvents(ctx context.Context, limit int) (*dto.RecentSecurityEventsResponse, error)

	// UC10: Get remediation updates
	GetRemediationUpdates(ctx context.Context, limit int) (*dto.RemediationUpdatesResponse, error)
}

// StakeholderVulnerabilitiesUseCase defines business logic for stakeholder vulnerabilities tab
type StakeholderVulnerabilitiesUseCase interface {
	// UC11-UC14: Get vulnerabilities statistics
	GetVulnerabilitiesStats(ctx context.Context) (*dto.StakeholderVulnerabilitiesStatsResponse, error)

	// UC15-UC23: List vulnerabilities with search, filters, and pagination
	ListVulnerabilities(ctx context.Context, req *dto.ListStakeholderVulnerabilitiesRequest) (*dto.ListStakeholderVulnerabilitiesResponse, error)

	// UC24: Export vulnerabilities to CSV
	ExportVulnerabilitiesToCSV(ctx context.Context, req *dto.ExportStakeholderVulnerabilitiesRequest) ([]byte, error)

	// UC25-UC27: Get SLA compliance data
	GetSLACompliance(ctx context.Context) (*dto.StakeholderSLAComplianceResponse, error)
}

// StakeholderReportsUseCase defines business logic for stakeholder reports tab
type StakeholderReportsUseCase interface {
	// UC28-UC30: Get reports statistics
	GetReportsStats(ctx context.Context) (*dto.StakeholderReportsStatsResponse, error)

	// UC31-UC36: List reports with status filter and pagination
	ListReports(ctx context.Context, req *dto.ListStakeholderReportsRequest) (*dto.ListStakeholderReportsResponse, error)

	// UC37-UC38: View report details with vulnerabilities
	ViewReport(ctx context.Context, reportID string, evidencePage, evidencePerPage int) (*dto.ViewStakeholderReportResponse, error)

	// UC39-UC40: Get evidence files for a report
	GetReportEvidenceFiles(ctx context.Context, reportID string, page, perPage int) (*dto.ViewReportEvidenceResponse, error)

	// UC41: Download evidence file
	DownloadEvidenceFile(ctx context.Context, fileID string) (*EvidenceFile, error)

	// UC42: Download report
	DownloadReport(ctx context.Context, reportID string) ([]byte, string, error)
}
