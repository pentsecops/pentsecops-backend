package domain

import (
	"context"
	"time"
)

// StakeholderOverviewRepository defines the interface for stakeholder overview data access
type StakeholderOverviewRepository interface {
	// UC1: Calculate overall security score
	CalculateSecurityScore(ctx context.Context) (float64, error)

	// UC2: Get active projects count with breakdown
	GetActiveProjectsCount(ctx context.Context) (totalActive, inProgress, completed int, err error)

	// UC3: Get critical issues count
	GetCriticalIssuesCount(ctx context.Context) (int, error)

	// UC4: Get open vulnerabilities count with trend
	GetOpenVulnerabilitiesCount(ctx context.Context) (current, lastMonth int, err error)

	// UC5: Calculate remediation rate
	CalculateRemediationRate(ctx context.Context) (total, remediated int, err error)

	// UC6: Calculate SLA compliance
	CalculateSLACompliance(ctx context.Context) (totalWithDueDate, remediatedOnTime int, err error)

	// UC7: Get vulnerability trend data for past 5 months
	GetVulnerabilityTrend(ctx context.Context, months int) ([]MonthlyTrendData, error)

	// UC8: Get asset status with vulnerability counts by severity
	GetAssetStatus(ctx context.Context) ([]AssetVulnCounts, error)

	// UC9: Get recent security events
	GetRecentSecurityEvents(ctx context.Context, limit int) ([]SecurityEvent, error)

	// UC10: Get remediation updates
	GetRemediationUpdates(ctx context.Context, limit int) ([]RemediationUpdateData, error)
}

// MonthlyTrendData represents vulnerability trend data for a month
type MonthlyTrendData struct {
	Month    time.Time
	Open     int
	Resolved int
}

// AssetVulnCounts represents vulnerability counts by severity for an asset
type AssetVulnCounts struct {
	AssetName string
	Critical  int
	High      int
	Medium    int
	Low       int
}

// SecurityEvent represents a security event in the activity log
type SecurityEvent struct {
	ID          string
	Description string
	EventType   string
	Timestamp   time.Time
}

// RemediationUpdateData represents a remediation update
type RemediationUpdateData struct {
	ID             string
	VulnTitle      string
	PreviousStatus string
	NewStatus      string
	RemediatedDate time.Time
	AssignedTeam   string
}

