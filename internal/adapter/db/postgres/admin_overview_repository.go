package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pentsecops/backend/internal/core/domain"
)

// AdminOverviewRepository implements the AdminOverviewRepository interface
type AdminOverviewRepository struct {
	db *sql.DB
}

// NewAdminOverviewRepository creates a new AdminOverviewRepository
func NewAdminOverviewRepository(db *sql.DB) domain.AdminOverviewRepository {
	return &AdminOverviewRepository{db: db}
}

// GetTotalProjectsCount returns the total count of projects
func (r *AdminOverviewRepository) GetTotalProjectsCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM projects`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total projects count: %w", err)
	}
	return count, nil
}

// GetOpenProjectsCount returns the count of open projects
func (r *AdminOverviewRepository) GetOpenProjectsCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM projects WHERE status = 'open'`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get open projects count: %w", err)
	}
	return count, nil
}

// GetInProgressProjectsCount returns the count of in-progress projects
func (r *AdminOverviewRepository) GetInProgressProjectsCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM projects WHERE status = 'in_progress'`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get in-progress projects count: %w", err)
	}
	return count, nil
}

// GetCompletedProjectsCount returns the count of completed projects
func (r *AdminOverviewRepository) GetCompletedProjectsCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM projects WHERE status = 'completed'`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get completed projects count: %w", err)
	}
	return count, nil
}

// GetTotalVulnerabilitiesCount returns the total count of vulnerabilities
func (r *AdminOverviewRepository) GetTotalVulnerabilitiesCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM vulnerabilities`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total vulnerabilities count: %w", err)
	}
	return count, nil
}

// GetActiveUsersCount returns the count of active users
func (r *AdminOverviewRepository) GetActiveUsersCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users WHERE is_active = true`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get active users count: %w", err)
	}
	return count, nil
}

// GetActivePentestersCount returns the count of active pentesters
func (r *AdminOverviewRepository) GetActivePentestersCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users WHERE is_active = true AND role = 'pentester'`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get active pentesters count: %w", err)
	}
	return count, nil
}

// GetActiveStakeholdersCount returns the count of active stakeholders
func (r *AdminOverviewRepository) GetActiveStakeholdersCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users WHERE is_active = true AND role = 'stakeholder'`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get active stakeholders count: %w", err)
	}
	return count, nil
}

// GetOpenIssuesCount returns the count of open issues
func (r *AdminOverviewRepository) GetOpenIssuesCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM vulnerabilities WHERE status = 'open'`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get open issues count: %w", err)
	}
	return count, nil
}

// GetCriticalIssuesCount returns the count of critical issues
func (r *AdminOverviewRepository) GetCriticalIssuesCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM vulnerabilities WHERE severity = 'critical'`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get critical issues count: %w", err)
	}
	return count, nil
}

// GetVulnerabilitiesBySeverity returns vulnerabilities grouped by severity
func (r *AdminOverviewRepository) GetVulnerabilitiesBySeverity(ctx context.Context) ([]domain.SeverityCount, error) {
	query := `
		SELECT
			severity,
			COUNT(*) as count
		FROM vulnerabilities
		GROUP BY severity
		ORDER BY
			CASE severity
				WHEN 'Critical' THEN 1
				WHEN 'High' THEN 2
				WHEN 'Medium' THEN 3
				WHEN 'Low' THEN 4
			END
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerabilities by severity: %w", err)
	}
	defer rows.Close()

	var results []domain.SeverityCount
	for rows.Next() {
		var sc domain.SeverityCount
		if err := rows.Scan(&sc.Severity, &sc.Count); err != nil {
			return nil, fmt.Errorf("failed to scan severity count: %w", err)
		}
		results = append(results, sc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating severity counts: %w", err)
	}

	return results, nil
}

// GetTop5DomainsByVulnerabilities returns top 5 domains by vulnerability count
func (r *AdminOverviewRepository) GetTop5DomainsByVulnerabilities(ctx context.Context) ([]domain.DomainVulnCount, error) {
	query := `
		SELECT
			COALESCE(domain, 'Unknown') as domain,
			COUNT(*) as vulnerability_count
		FROM vulnerabilities
		WHERE domain IS NOT NULL AND domain != ''
		GROUP BY domain
		ORDER BY vulnerability_count DESC
		LIMIT 5
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get top domains: %w", err)
	}
	defer rows.Close()

	var results []domain.DomainVulnCount
	for rows.Next() {
		var dvc domain.DomainVulnCount
		if err := rows.Scan(&dvc.Domain, &dvc.Count); err != nil {
			return nil, fmt.Errorf("failed to scan domain count: %w", err)
		}
		results = append(results, dvc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating domain counts: %w", err)
	}

	return results, nil
}

// GetProjectStatusDistribution returns project status distribution
func (r *AdminOverviewRepository) GetProjectStatusDistribution(ctx context.Context) ([]domain.StatusCount, error) {
	query := `
		SELECT
			status,
			COUNT(*) as count
		FROM projects
		GROUP BY status
		ORDER BY
			CASE status
				WHEN 'Open' THEN 1
				WHEN 'In Progress' THEN 2
				WHEN 'Completed' THEN 3
			END
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get project status distribution: %w", err)
	}
	defer rows.Close()

	var results []domain.StatusCount
	for rows.Next() {
		var sc domain.StatusCount
		if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
			return nil, fmt.Errorf("failed to scan status count: %w", err)
		}
		results = append(results, sc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating status counts: %w", err)
	}

	return results, nil
}

// GetRecentActivityLogs returns recent activity logs with pagination
func (r *AdminOverviewRepository) GetRecentActivityLogs(ctx context.Context, limit, offset int) ([]domain.ActivityLog, error) {
	query := `
		SELECT 
			id,
			user_id,
			user_email,
			action,
			entity_type,
			entity_id,
			ip_address,
			user_agent,
			created_at
		FROM activity_logs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.ActivityLog
	for rows.Next() {
		var log domain.ActivityLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.UserEmail,
			&log.Action,
			&log.EntityType,
			&log.EntityID,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity log: %w", err)
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating activity logs: %w", err)
	}

	return logs, nil
}

// GetActivityLogsCount returns the total count of activity logs
func (r *AdminOverviewRepository) GetActivityLogsCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM activity_logs`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get activity logs count: %w", err)
	}
	return count, nil
}

// CreateActivityLog creates a new activity log
func (r *AdminOverviewRepository) CreateActivityLog(ctx context.Context, log *domain.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (
			id, user_id, user_email, action, entity_type, entity_id, ip_address, user_agent, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.UserID,
		log.UserEmail,
		log.Action,
		log.EntityType,
		log.EntityID,
		log.IPAddress,
		log.UserAgent,
		log.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create activity log: %w", err)
	}

	return nil
}
