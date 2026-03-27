package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID                  uuid.UUID  `json:"id"`
	Email               string     `json:"email"`
	PasswordHash        string     `json:"-"`
	FullName            string     `json:"full_name"`
	Role                string     `json:"role"`
	IsActive            bool       `json:"is_active"`
	ForcePasswordChange bool       `json:"force_password_change"`
	FailedLoginAttempts int        `json:"-"`
	LastFailedLogin     *time.Time `json:"-"`
	AccountLockedUntil  *time.Time `json:"-"`
	LastLogin           *time.Time `json:"last_login,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// Project represents a penetration testing project
type Project struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	Status       string     `json:"status"`
	AssignedTo   *uuid.UUID `json:"assigned_to,omitempty"`
	Deadline     *time.Time `json:"deadline,omitempty"`
	Scope        *string    `json:"scope,omitempty"`
	CurrentPhase *string    `json:"current_phase,omitempty"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Vulnerability represents a security vulnerability
type Vulnerability struct {
	ID               uuid.UUID  `json:"id"`
	Title            string     `json:"title"`
	Description      *string    `json:"description,omitempty"`
	Severity         string     `json:"severity"`
	Domain           string     `json:"domain"`
	Status           string     `json:"status"`
	DiscoveredDate   *time.Time `json:"discovered_date,omitempty"`
	DueDate          *time.Time `json:"due_date,omitempty"`
	AssignedTo       *string    `json:"assigned_to,omitempty"`
	CVSSScore        *float64   `json:"cvss_score,omitempty"`
	CWEID            *string    `json:"cwe_id,omitempty"`
	DomainID         *uuid.UUID `json:"domain_id,omitempty"`
	ProjectID        *uuid.UUID `json:"project_id,omitempty"`
	DiscoveredBy     *uuid.UUID `json:"discovered_by,omitempty"`
	RemediationNotes *string    `json:"remediation_notes,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ActivityLog represents a system activity log
type ActivityLog struct {
	ID           uuid.UUID  `json:"id"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	UserEmail    string     `json:"user_email"`
	UserRole     string     `json:"user_role"`
	Action       string     `json:"action"`
	EntityType   *string    `json:"entity_type,omitempty"`
	EntityID     *uuid.UUID `json:"entity_id,omitempty"`
	OldValues    *string    `json:"old_values,omitempty"`
	NewValues    *string    `json:"new_values,omitempty"`
	IPAddress    string     `json:"ip_address"`
	UserAgent    string     `json:"user_agent"`
	Endpoint     string     `json:"endpoint"`
	Method       string     `json:"method"`
	StatusCode   int        `json:"status_code"`
	ErrorMessage *string    `json:"error_message,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// AuditLogFilter represents filters for audit log queries
type AuditLogFilter struct {
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Action     string     `json:"action,omitempty"`
	EntityType string     `json:"entity_type,omitempty"`
	EntityID   *uuid.UUID `json:"entity_id,omitempty"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
}

// Domain represents a domain being tested
type Domain struct {
	ID            uuid.UUID  `json:"id"`
	DomainName    string     `json:"domain_name"`
	IPAddress     *string    `json:"ip_address,omitempty"`
	RiskScore     *float64   `json:"risk_score,omitempty"`
	LastScannedAt *time.Time `json:"last_scanned_at,omitempty"`
	SLACompliance *float64   `json:"sla_compliance,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Report represents a penetration testing report
type Report struct {
	ID               uuid.UUID  `json:"id"`
	Title            string     `json:"title"`
	SubmittedBy      uuid.UUID  `json:"submitted_by"`
	SubmissionDate   time.Time  `json:"submission_date"`
	ProjectID        *uuid.UUID `json:"project_id,omitempty"`
	Status           string     `json:"status"`
	ExecutiveSummary *string    `json:"executive_summary,omitempty"`
	Feedback         *string    `json:"feedback,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Task represents a task in the system
type Task struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	ProjectID   *uuid.UUID `json:"project_id,omitempty"`
	AssignedTo  *uuid.UUID `json:"assigned_to,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
	Status      string     `json:"status"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Constants for roles
const (
	RoleAdmin       = "admin"
	RolePentester   = "pentester"
	RoleStakeholder = "stakeholder"
)

// Constants for project status
const (
	ProjectStatusOpen       = "open"
	ProjectStatusInProgress = "in_progress"
	ProjectStatusCompleted  = "completed"
	ProjectStatusOnHold     = "on_hold"
	ProjectStatusCancelled  = "cancelled"
)

// Constants for project types
const (
	ProjectTypeWeb     = "web"
	ProjectTypeNetwork = "network"
	ProjectTypeAPI     = "api"
	ProjectTypeMobile  = "mobile"
)

// Constants for vulnerability severity
const (
	SeverityCritical = "critical"
	SeverityHigh     = "high"
	SeverityMedium   = "medium"
	SeverityLow      = "low"
)

// Constants for vulnerability status
const (
	VulnStatusOpen       = "open"
	VulnStatusInProgress = "in_progress"
	VulnStatusRemediated = "remediated"
	VulnStatusVerified   = "verified"
)

// Constants for report status
const (
	ReportStatusReceived    = "Received"
	ReportStatusUnderReview = "Under Review"
	ReportStatusApproved    = "Approved"
	ReportStatusRejected    = "Rejected"
	ReportStatusShared      = "Shared"
	ReportStatusRemediated  = "Remediated"
)

// Constants for task status
const (
	TaskStatusToDo       = "to_do"
	TaskStatusInProgress = "in_progress"
	TaskStatusDone       = "done"
)

// Constants for task priority
const (
	PriorityLow      = "low"
	PriorityMedium   = "medium"
	PriorityHigh     = "high"
	PriorityCritical = "critical"
)

// Notification represents a notification sent to users
type Notification struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Message     string     `json:"message"`
	SentTo      string     `json:"sent_to"`
	RecipientID *uuid.UUID `json:"recipient_id,omitempty"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	Status      string     `json:"status"`
	Type        *string    `json:"type,omitempty"`
	IsRead      bool       `json:"is_read"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Alert represents an alert in the system
type Alert struct {
	ID         uuid.UUID  `json:"id"`
	Title      string     `json:"title"`
	Message    string     `json:"message"`
	AlertType  *string    `json:"alert_type,omitempty"`
	Priority   *string    `json:"priority,omitempty"`
	Source     *string    `json:"source,omitempty"`
	SenderID   *uuid.UUID `json:"sender_id,omitempty"`
	IsResolved bool       `json:"is_resolved"`
	ResolvedBy *uuid.UUID `json:"resolved_by,omitempty"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Constants for notification sent_to
const (
	SentToAllPentesters   = "all_pentesters"
	SentToAllStakeholders = "all_stakeholders"
	SentToSpecificUser    = "specific_user"
)

// Constants for notification status
const (
	NotificationStatusSent    = "sent"
	NotificationStatusFailed  = "failed"
	NotificationStatusPending = "pending"
)

// Constants for alert type
const (
	AlertTypeImportant = "important"
	AlertTypeCritical  = "critical"
	AlertTypeWarning   = "warning"
	AlertTypeInfo      = "info"
)

// Constants for alert source
const (
	AlertSourcePentester = "pentester"
	AlertSourceSystem    = "system"
	AlertSourceAdmin     = "admin"
)

// RefreshToken represents a refresh token in the system
type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// NotificationStatsResponse represents notification statistics
type NotificationStatsResponse struct {
	TotalSent        int64                    `json:"total_sent"`
	SentToday        int64                    `json:"sent_today"`
	SentThisWeek     int64                    `json:"sent_this_week"`
	SentThisMonth    int64                    `json:"sent_this_month"`
	ByType           map[string]int           `json:"by_type"`
	ByPriority       map[string]int           `json:"by_priority"`
	ByRecipientType  map[string]int           `json:"by_recipient_type"`
	RecentActivity   []NotificationActivity   `json:"recent_activity"`
}

// NotificationActivity represents recent notification activity
type NotificationActivity struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	SentTo    string    `json:"sent_to"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

// AuditRepository interface for audit logging operations
type AuditRepository interface {
	LogActivity(ctx context.Context, log *ActivityLog) error
	GetActivityLogs(ctx context.Context, filter AuditLogFilter) ([]ActivityLog, int, error)
	LogDatabaseOperation(ctx context.Context, userID *uuid.UUID, userEmail, userRole, action, entityType string, entityID *uuid.UUID, oldValues, newValues interface{})
}
