package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SeverityCount represents vulnerability count by severity
type SeverityCount struct {
	Severity string
	Count    int64
}

// DomainVulnCount represents vulnerability count by domain
type DomainVulnCount struct {
	Domain string
	Count  int64
}

// StatusCount represents project count by status
type StatusCount struct {
	Status string
	Count  int64
}

// AdminOverviewRepository defines the interface for admin overview data access
type AdminOverviewRepository interface {
	// Statistics queries
	GetTotalProjectsCount(ctx context.Context) (int64, error)
	GetOpenProjectsCount(ctx context.Context) (int64, error)
	GetInProgressProjectsCount(ctx context.Context) (int64, error)
	GetCompletedProjectsCount(ctx context.Context) (int64, error)
	GetTotalVulnerabilitiesCount(ctx context.Context) (int64, error)
	GetActiveUsersCount(ctx context.Context) (int64, error)
	GetActivePentestersCount(ctx context.Context) (int64, error)
	GetActiveStakeholdersCount(ctx context.Context) (int64, error)
	GetOpenIssuesCount(ctx context.Context) (int64, error)
	GetCriticalIssuesCount(ctx context.Context) (int64, error)

	// Chart queries
	GetVulnerabilitiesBySeverity(ctx context.Context) ([]SeverityCount, error)
	GetTop5DomainsByVulnerabilities(ctx context.Context) ([]DomainVulnCount, error)
	GetProjectStatusDistribution(ctx context.Context) ([]StatusCount, error)

	// Activity logs
	GetRecentActivityLogs(ctx context.Context, limit, offset int) ([]ActivityLog, error)
	GetActivityLogsCount(ctx context.Context) (int64, error)
	CreateActivityLog(ctx context.Context, log *ActivityLog) error
}

// UsersRepository defines the interface for users data access
type UsersRepository interface {
	// User CRUD operations
	CreateUser(ctx context.Context, user *CreateUserParams) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, id string, params *UpdateUserParams) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUserLastLogin(ctx context.Context, id string, lastLogin time.Time) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)

	// User listing and pagination
	ListUsers(ctx context.Context, limit, offset int) ([]*UserWithProjectCount, error)
	CountUsers(ctx context.Context) (int64, error)

	// User statistics
	GetUserStats(ctx context.Context) (*UserStats, error)
	CountUsersByRole(ctx context.Context, role string, isActive bool) (int64, error)
	CountInactiveUsers(ctx context.Context) (int64, error)

	// Export
	ListAllUsersForExport(ctx context.Context) ([]*UserWithProjectCount, error)

	// Role-based queries for notifications
	GetUsersByRole(ctx context.Context, role string) ([]User, error)
}

// CreateUserParams represents parameters for creating a user
type CreateUserParams struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	Role         string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UpdateUserParams represents parameters for updating a user
type UpdateUserParams struct {
	FullName string
	Email    string
	Role     string
	IsActive *bool
}

// UserWithProjectCount represents a user with project count
type UserWithProjectCount struct {
	ID           string
	Email        string
	FullName     string
	Role         string
	IsActive     bool
	LastLogin    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ProjectCount int64
}

// UserStats represents user statistics
type UserStats struct {
	ActivePentesters   int64
	ActiveStakeholders int64
	InactiveUsers      int64
	TotalUsers         int64
}

// ProjectsRepository defines the interface for projects data access
type ProjectsRepository interface {
	// Project CRUD operations
	CreateProject(ctx context.Context, project *CreateProjectParams) (*Project, error)
	GetProjectByID(ctx context.Context, id string) (*Project, error)
	GetProjectByName(ctx context.Context, name string) (*Project, error)
	UpdateProject(ctx context.Context, id string, params *UpdateProjectParams) (*Project, error)
	DeleteProject(ctx context.Context, id string) error
	UpdateProjectStatus(ctx context.Context, id string, status string) error

	// Project listing and pagination
	ListProjects(ctx context.Context, limit, offset int) ([]*ProjectWithDetails, error)
	CountProjects(ctx context.Context) (int64, error)

	// Project statistics
	GetProjectStats(ctx context.Context) (*ProjectStats, error)

	// Get pentesters for dropdown
	GetPentesters(ctx context.Context) ([]*Pentester, error)
}

// CreateProjectParams represents parameters for creating a project
type CreateProjectParams struct {
	ID         string
	Name       string
	Type       string
	AssignedTo *string
	Deadline   time.Time
	Scope      *string
	Status     string
	CreatedBy  *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// UpdateProjectParams represents parameters for updating a project
type UpdateProjectParams struct {
	Name       string
	Type       string
	AssignedTo *uuid.UUID
	Deadline   time.Time
	Scope      *string
	Status     string
}

// ProjectWithDetails represents a project with additional details
type ProjectWithDetails struct {
	ID                 string
	Name               string
	Type               string
	Status             string
	AssignedTo         *string
	AssignedToName     string
	Deadline           time.Time
	VulnerabilityCount int64
	CreatedAt          time.Time
}

// ProjectStats represents project statistics
type ProjectStats struct {
	OpenCount       int64
	InProgressCount int64
	CompletedCount  int64
}

// Pentester represents a pentester for dropdown
type Pentester struct {
	ID       string
	FullName string
	Email    string
}

// TasksRepository defines the interface for tasks data access
type TasksRepository interface {
	// Task CRUD operations
	CreateTask(ctx context.Context, task *CreateTaskParams) (*Task, error)
	GetTaskByID(ctx context.Context, id string) (*Task, error)
	DeleteTask(ctx context.Context, id string) error
	UpdateTaskStatus(ctx context.Context, id string, status string) error

	// Task listing
	ListTasksByProject(ctx context.Context, projectID string) ([]*TaskWithDetails, error)
	ListAllTasks(ctx context.Context) ([]*TaskWithDetails, error)
	GetTasksByStatus(ctx context.Context, status string) ([]*TaskWithDetails, error)
}

// CreateTaskParams represents parameters for creating a task
type CreateTaskParams struct {
	ID          string
	ProjectID   *string
	Title       string
	Description *string
	Status      string
	Priority    *string
	AssignedTo  *string
	Deadline    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TaskWithDetails represents a task with additional details
type TaskWithDetails struct {
	ID             string
	ProjectID      *string
	ProjectName    string
	Title          string
	Description    *string
	Status         string
	Priority       *string
	AssignedTo     *string
	AssignedToName string
	Deadline       *time.Time
	CompletedAt    *time.Time
	CreatedAt      time.Time
}

// VulnerabilitiesRepository defines the interface for vulnerabilities data access
type VulnerabilitiesRepository interface {
	CreateVulnerability(ctx context.Context, vuln *CreateVulnerabilityParams) (*Vulnerability, error)
	GetVulnerabilityByID(ctx context.Context, id string) (*Vulnerability, error)
	UpdateVulnerability(ctx context.Context, vuln *UpdateVulnerabilityParams) (*Vulnerability, error)
	DeleteVulnerability(ctx context.Context, id string) error
	ListVulnerabilities(ctx context.Context, limit, offset int) ([]*Vulnerability, error)
	CountVulnerabilities(ctx context.Context) (int64, error)
	SearchAndFilterVulnerabilities(ctx context.Context, search, severity, status string, limit, offset int) ([]*Vulnerability, error)
	CountSearchAndFilterVulnerabilities(ctx context.Context, search, severity, status string) (int64, error)
	GetVulnerabilityStats(ctx context.Context) (*VulnerabilityStats, error)
	GetSLACompliance(ctx context.Context) (*SLACompliance, error)
	ExportVulnerabilities(ctx context.Context, search, severity, status string) ([]*VulnerabilityExport, error)
}

// CreateVulnerabilityParams represents parameters for creating a vulnerability
type CreateVulnerabilityParams struct {
	ID               string
	Title            string
	Description      *string
	Severity         string
	Domain           string
	Status           string
	DiscoveredDate   *time.Time
	DueDate          time.Time
	AssignedTo       *string
	CVSSScore        *float64
	CWEID            *string
	DomainID         *string
	ProjectID        *string
	DiscoveredBy     *string
	RemediationNotes *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// UpdateVulnerabilityParams represents parameters for updating a vulnerability (all fields optional except ID)
type UpdateVulnerabilityParams struct {
	ID               string
	Title            *string
	Description      *string
	Severity         *string
	Domain           *string
	Status           *string
	DiscoveredDate   *time.Time
	DueDate          *time.Time
	AssignedTo       *string
	CVSSScore        *float64
	CWEID            *string
	RemediationNotes *string
	UpdatedAt        time.Time
}

// VulnerabilityStats represents vulnerability statistics
type VulnerabilityStats struct {
	Total      int64
	Critical   int64
	High       int64
	Medium     int64
	Low        int64
	Open       int64
	InProgress int64
	Remediated int64
	Verified   int64
}

// SLACompliance represents SLA compliance metrics
type SLACompliance struct {
	CriticalOverdue  int64
	HighApproaching  int64
	TotalWithDueDate int64
	RemediatedOnTime int64
}

// VulnerabilityExport represents a vulnerability for export
type VulnerabilityExport struct {
	ID             string
	Title          string
	Severity       string
	Domain         string
	Status         string
	DiscoveredDate *time.Time
	DueDate        time.Time
	AssignedTo     *string
}

// DomainsRepository defines the interface for domains data access
type DomainsRepository interface {
	// Statistics
	GetDomainsStats(ctx context.Context) (*DomainsStats, error)

	// Domains CRUD
	ListDomains(ctx context.Context, limit, offset int) ([]DomainWithStats, error)
	CountDomains(ctx context.Context) (int64, error)
	GetDomainByID(ctx context.Context, id string) (*Domain, error)
	CreateDomain(ctx context.Context, params *CreateDomainParams) (*Domain, error)
	UpdateDomain(ctx context.Context, params *UpdateDomainParams) (*Domain, error)
	DeleteDomain(ctx context.Context, id string) error

	// Security Metrics
	GetSecurityMetrics(ctx context.Context) ([]SecurityMetric, error)
	CreateSecurityMetric(ctx context.Context, params *CreateSecurityMetricParams) (*SecurityMetric, error)

	// SLA Breach Analysis
	GetSLABreachAnalysis(ctx context.Context) ([]SLABreachDomain, error)
}

// DomainsStats represents domains overview statistics
type DomainsStats struct {
	TotalDomains         int64
	AverageRiskScore     float64
	CriticalIssues       int64
	SLACompliancePercent float64
}

// DomainWithStats represents a domain with vulnerability statistics
type DomainWithStats struct {
	ID                   string
	DomainName           string
	IPAddress            *string
	Description          *string
	RiskScore            *float64
	IsActive             bool
	LastScanned          *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
	TotalVulnerabilities int64
	CriticalCount        int64
	HighCount            int64
	MediumCount          int64
	LowCount             int64
	OpenIssues           int64
	SLACompliance        float64
}

// CreateDomainParams represents parameters for creating a domain
type CreateDomainParams struct {
	ID          string
	DomainName  string
	IPAddress   *string
	Description *string
	RiskScore   *float64
	IsActive    bool
	LastScanned *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UpdateDomainParams represents parameters for updating a domain
type UpdateDomainParams struct {
	ID          string
	DomainName  *string
	IPAddress   *string
	Description *string
	RiskScore   *float64
	LastScanned *time.Time
	UpdatedAt   time.Time
}

// SecurityMetric represents a security metric
type SecurityMetric struct {
	MetricName string
	AvgValue   float64
}

// CreateSecurityMetricParams represents parameters for creating a security metric
type CreateSecurityMetricParams struct {
	ID          string
	DomainID    string
	MetricName  string
	MetricValue float64
	MeasuredAt  time.Time
	CreatedAt   time.Time
}

// SLABreachDomain represents a domain with SLA compliance percentage
type SLABreachDomain struct {
	DomainName           string
	SLACompliancePercent float64
}

// NotificationsRepository defines the interface for notifications data access
type NotificationsRepository interface {
	// Notifications
	GetTotalNotificationsSent(ctx context.Context, createdBy string) (int64, error)
	ListNotifications(ctx context.Context, createdBy string, limit, offset int) ([]Notification, error)
	CountNotifications(ctx context.Context, createdBy string) (int64, error)
	CreateNotification(ctx context.Context, notification *Notification) error
	GetNotificationByID(ctx context.Context, id string) (*Notification, error)

	// Enhanced notification methods for admin
	GetNotifications(ctx context.Context, page, perPage int, sendTo, notificationType, priority string) ([]Notification, int, error)
	GetRecipientCount(ctx context.Context, notificationID uuid.UUID) (int, error)
	GetNotificationStats(ctx context.Context) (*NotificationStatsResponse, error)

	// Alerts
	ListImportantAlerts(ctx context.Context, limit, offset int) ([]Alert, error)
	CountImportantAlerts(ctx context.Context) (int64, error)
	CreateAlert(ctx context.Context, alert *Alert) (*Alert, error)
	GetAlertByID(ctx context.Context, id string) (*Alert, error)

	// Helper methods
	GetUsersByRole(ctx context.Context, role string) ([]User, error)
}

// AuthRepository defines the interface for authentication data access
type AuthRepository interface {
	// User authentication
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
	UpdateUserPassword(ctx context.Context, userID, passwordHash string) error
	UpdateUserLastLogin(ctx context.Context, userID string, lastLogin time.Time) error
	UpdateUserForcePasswordChange(ctx context.Context, userID string, forceChange bool) error

	// Account locking
	IncrementFailedLoginAttempts(ctx context.Context, userID string) error
	ResetFailedLoginAttempts(ctx context.Context, userID string) error
	LockAccount(ctx context.Context, userID string, lockUntil time.Time) error
	UnlockAccount(ctx context.Context, userID string) error

	// Refresh token management
	CreateRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	DeleteAllUserRefreshTokens(ctx context.Context, userID string) error
	DeleteExpiredRefreshTokens(ctx context.Context) error
}

// PentesterOverviewRepository defines the interface for pentester overview data access
type PentesterOverviewRepository interface {
	// Statistics queries
	GetPentesterActiveProjectsCount(ctx context.Context, pentesterID string) (int64, error)
	GetPentesterProjectsDueThisWeek(ctx context.Context, pentesterID string) (int64, error)
	GetPentesterReportsSubmittedCount(ctx context.Context, pentesterID string) (int64, error)
	GetPentesterReportsPendingReview(ctx context.Context, pentesterID string) (int64, error)
	GetPentesterVulnerabilitiesFoundCount(ctx context.Context, pentesterID string) (int64, error)
	GetPentesterCriticalHighVulnsCount(ctx context.Context, pentesterID string) (int64, error)
	GetPentesterAverageAssessmentTime(ctx context.Context, pentesterID string) (float64, error)

	// List queries
	GetPentesterActiveProjects(ctx context.Context, pentesterID string, limit int) ([]*PentesterProject, error)
	GetPentesterRecentVulnerabilities(ctx context.Context, pentesterID string, limit int) ([]*PentesterVulnerability, error)
	GetPentesterUpcomingDeadlines(ctx context.Context, pentesterID string) ([]*PentesterProjectDeadline, error)
}

// PentesterProject represents a project for pentester
type PentesterProject struct {
	ID                 string
	Name               string
	Type               string
	Status             string
	Deadline           time.Time
	Scope              *string
	VulnerabilityCount int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// PentesterVulnerability represents a vulnerability for pentester
type PentesterVulnerability struct {
	ID             string
	Title          string
	Severity       string
	Domain         string
	Status         string
	ProjectName    *string
	DiscoveredDate *time.Time
	CreatedAt      time.Time
}

// PentesterProjectDeadline represents a project with deadline info
type PentesterProjectDeadline struct {
	ID                 string
	Name               string
	Type               string
	Status             string
	Deadline           time.Time
	Scope              *string
	DaysLeft           float64
	VulnerabilityCount int64
}

// ========== Pentester Projects Repository ==========

// PentesterProjectsRepository defines the interface for pentester projects data access
type PentesterProjectsRepository interface {
	// UC9: Fetch and Display Assigned Projects List
	GetPentesterAssignedProjects(ctx context.Context, pentesterID string) ([]*PentesterAssignedProject, error)
	GetPentesterProjectByID(ctx context.Context, projectID string, pentesterID string) (*PentesterAssignedProject, error)
	CountPentesterAssignedProjects(ctx context.Context, pentesterID string) (int64, error)

	// UC12: Display Project Assets List
	GetProjectAssets(ctx context.Context, projectID string) ([]*ProjectAsset, error)

	// UC13: Display Project Requirements List
	GetProjectRequirements(ctx context.Context, projectID string) ([]*ProjectRequirement, error)
}

// PentesterAssignedProject represents a project assigned to pentester with full details
type PentesterAssignedProject struct {
	ID                 string
	Name               string
	Type               string
	Status             string
	Deadline           time.Time
	Scope              *string
	CurrentPhase       *string
	ProgressPercentage float64
	TotalTasks         int64
	CompletedTasks     int64
	CreatedAt          time.Time
}

// ProjectAsset represents an asset (domain/IP) for a project
type ProjectAsset struct {
	ID          string
	DomainID    string
	DomainName  string
	IPAddress   *string
	Description *string
	RiskScore   *float64
	IsActive    bool
	LastScanned *time.Time
}

// ProjectRequirement represents a requirement for a project
type ProjectRequirement struct {
	ID              string
	ProjectID       string
	RequirementText string
	IsCompleted     bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ========== Pentester Tasks Repository ==========

// PentesterTasksRepository defines the interface for pentester tasks data access
type PentesterTasksRepository interface {
	// UC16: Fetch and Display Task Board
	GetPentesterTaskBoard(ctx context.Context, pentesterID string) ([]*PentesterTask, error)
	GetPentesterTasksByStatus(ctx context.Context, pentesterID string, status string) ([]*PentesterTask, error)

	// UC17: Create New Task
	CreatePentesterTask(ctx context.Context, task *CreatePentesterTaskParams) (*PentesterTask, error)

	// UC18: Update Task Status
	UpdatePentesterTaskStatus(ctx context.Context, taskID string, status string) error

	// UC19: Edit Task Details
	UpdatePentesterTask(ctx context.Context, taskID string, params *UpdatePentesterTaskParams) (*PentesterTask, error)

	// UC20: Delete Task
	DeletePentesterTask(ctx context.Context, taskID string) error

	// Get single task
	GetPentesterTaskByID(ctx context.Context, taskID string, pentesterID string) (*PentesterTask, error)

	// Get projects for dropdown
	GetPentesterProjectsForDropdown(ctx context.Context, pentesterID string) ([]*PentesterProjectDropdown, error)
}

// PentesterTask represents a task in the pentester's task board
type PentesterTask struct {
	ID          string
	ProjectID   string
	ProjectName string
	ProjectType string
	Title       string
	Description *string
	Status      string
	Priority    *string
	Deadline    *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreatePentesterTaskParams represents parameters for creating a task
type CreatePentesterTaskParams struct {
	ID          string
	ProjectID   string
	Title       string
	Description *string
	Status      string
	Priority    *string
	AssignedTo  string
	Deadline    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UpdatePentesterTaskParams represents parameters for updating a task
type UpdatePentesterTaskParams struct {
	Title       string
	Description *string
	Priority    *string
	ProjectID   string
	Deadline    *time.Time
}

// PentesterProjectDropdown represents a project option for dropdown
type PentesterProjectDropdown struct {
	ID     string
	Name   string
	Type   string
	Status string
}

// ============================================================================
// PENTESTER SUBMIT REPORT REPOSITORY
// ============================================================================

// PentesterSubmitReportRepository defines the interface for pentester submit report data access
type PentesterSubmitReportRepository interface {
	// UC22: Get projects for dropdown
	GetProjectsForReportDropdown(ctx context.Context, pentesterID uuid.UUID) ([]ProjectForReportDropdown, error)

	// UC23: Submit a new report
	CreateReport(ctx context.Context, report *CreateReportParams) (*ReportDetails, error)
	CreateReportVulnerability(ctx context.Context, vuln *CreateReportVulnerabilityParams) error

	// UC27: Get submitted reports history
	GetSubmittedReports(ctx context.Context, pentesterID uuid.UUID, limit, offset int) ([]SubmittedReport, error)
	CountSubmittedReports(ctx context.Context, pentesterID uuid.UUID) (int64, error)

	// UC31: Get report details
	GetReportDetails(ctx context.Context, reportID, pentesterID uuid.UUID) (*ReportDetails, error)
	GetReportVulnerabilities(ctx context.Context, reportID uuid.UUID) ([]ReportVulnerabilityDetails, error)

	// UC30: Get rejected report for resubmission
	GetRejectedReportForResubmit(ctx context.Context, reportID, pentesterID uuid.UUID) (*ReportDetails, error)
}

// ProjectForReportDropdown represents a project option for report dropdown
type ProjectForReportDropdown struct {
	ID   uuid.UUID
	Name string
	Type string
}

// CreateReportParams represents parameters for creating a report
type CreateReportParams struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Title       string
	SubmittedBy uuid.UUID
	Status      string
}

// CreateReportVulnerabilityParams represents parameters for creating a report vulnerability
type CreateReportVulnerabilityParams struct {
	ID                        uuid.UUID
	ReportID                  uuid.UUID
	AssetTarget               string
	VulnerabilityTitle        string
	Severity                  string
	AttackVector              *string
	VulnerabilityDescription  string
	EvidencePOC               *string
	RemediationRecommendation *string
}

// SubmittedReport represents a report in the submitted reports history
type SubmittedReport struct {
	ID                   uuid.UUID
	ProjectID            *uuid.UUID
	ProjectName          string
	ProjectType          string
	Title                string
	Status               string
	Feedback             *string
	VulnerabilitiesCount int64
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// ReportDetails represents detailed information about a report
type ReportDetails struct {
	ID          uuid.UUID
	ProjectID   *uuid.UUID
	ProjectName string
	ProjectType string
	Title       string
	Status      string
	Feedback    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ReportVulnerabilityDetails represents vulnerability details in a report
type ReportVulnerabilityDetails struct {
	ID                        uuid.UUID
	ReportID                  uuid.UUID
	AssetTarget               string
	VulnerabilityTitle        string
	Severity                  string
	AttackVector              *string
	VulnerabilityDescription  string
	EvidencePOC               *string
	RemediationRecommendation *string
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

// ============================================================================
// PENTESTER ALERTS REPOSITORY
// ============================================================================

// PentesterAlertsRepository defines the interface for pentester alerts data access
type PentesterAlertsRepository interface {
	// UC33: Get alert statistics
	GetAlertStats(ctx context.Context, pentesterID uuid.UUID) (*AlertStats, error)

	// UC34: Get all alerts with filters and pagination
	GetAlerts(ctx context.Context, pentesterID uuid.UUID, alertType, search *string, limit, offset int) ([]PentesterAlert, error)
	CountAlerts(ctx context.Context, pentesterID uuid.UUID, alertType, search *string) (int64, error)

	// UC36: Mark alert as read/unread
	MarkAlertAsRead(ctx context.Context, alertID, pentesterID uuid.UUID) error
	MarkAlertAsUnread(ctx context.Context, alertID, pentesterID uuid.UUID) error

	// UC37: Dismiss alert
	DismissAlert(ctx context.Context, alertID, pentesterID uuid.UUID) error

	// UC38: Create alert to admin
	CreateAlertToAdmin(ctx context.Context, alert *CreateAlertParams) (*PentesterAlert, error)
	GetAdminUserID(ctx context.Context) (uuid.UUID, error)

	// Get single alert details
	GetAlertDetails(ctx context.Context, alertID, pentesterID uuid.UUID) (*PentesterAlert, error)
}

// AlertStats represents alert statistics for pentester
type AlertStats struct {
	TotalAlerts        int64
	UnreadAlerts       int64
	HighPriorityAlerts int64
}

// PentesterAlert represents an alert for pentester
type PentesterAlert struct {
	ID          uuid.UUID
	Title       string
	Message     string
	AlertType   string
	Priority    string
	Source      string
	SenderID    *uuid.UUID
	SenderName  string
	IsRead      bool
	IsDismissed bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateAlertParams represents parameters for creating an alert
type CreateAlertParams struct {
	ID          uuid.UUID
	Title       string
	Message     string
	AlertType   string
	Priority    string
	Source      string
	RecipientID uuid.UUID
	SenderID    uuid.UUID
}
