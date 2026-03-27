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

// DomainsRepository implements the domain.DomainsRepository interface
type DomainsRepository struct {
	queries *sqlc.Queries
}

// NewDomainsRepository creates a new DomainsRepository
func NewDomainsRepository(db *sql.DB) *DomainsRepository {
	return &DomainsRepository{
		queries: sqlc.New(db),
	}
}

// GetDomainsStats retrieves domains overview statistics
func (r *DomainsRepository) GetDomainsStats(ctx context.Context) (*domain.DomainsStats, error) {
	stats, err := r.queries.GetDomainsStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get domains stats: %w", err)
	}

	// Convert interface{} to proper types
	avgRiskScore := 0.0
	if stats.AvgRiskScore != nil {
		switch v := stats.AvgRiskScore.(type) {
		case float64:
			avgRiskScore = v
		case string:
			fmt.Sscanf(v, "%f", &avgRiskScore)
		case []byte:
			fmt.Sscanf(string(v), "%f", &avgRiskScore)
		}
	}

	criticalIssues := int64(0)
	if stats.CriticalIssues != nil {
		if val, ok := stats.CriticalIssues.(int64); ok {
			criticalIssues = val
		}
	}

	slaCompliance := 0.0
	if stats.SlaCompliancePercent != nil {
		switch v := stats.SlaCompliancePercent.(type) {
		case float64:
			slaCompliance = v
		case string:
			fmt.Sscanf(v, "%f", &slaCompliance)
		case []byte:
			fmt.Sscanf(string(v), "%f", &slaCompliance)
		}
	}

	return &domain.DomainsStats{
		TotalDomains:         stats.TotalDomains,
		AverageRiskScore:     avgRiskScore,
		CriticalIssues:       criticalIssues,
		SLACompliancePercent: slaCompliance,
	}, nil
}

// ListDomains retrieves a paginated list of domains with statistics
func (r *DomainsRepository) ListDomains(ctx context.Context, limit, offset int) ([]domain.DomainWithStats, error) {
	rows, err := r.queries.ListDomains(ctx, sqlc.ListDomainsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}

	domains := make([]domain.DomainWithStats, 0, len(rows))
	for _, row := range rows {
		d := domain.DomainWithStats{
			ID:         row.ID.String(),
			DomainName: row.DomainName,
		}

		// Convert interface{} counts to int64
		if row.TotalVulnerabilities != nil {
			if val, ok := row.TotalVulnerabilities.(int64); ok {
				d.TotalVulnerabilities = val
			}
		}
		if row.CriticalCount != nil {
			if val, ok := row.CriticalCount.(int64); ok {
				d.CriticalCount = val
			}
		}
		if row.HighCount != nil {
			if val, ok := row.HighCount.(int64); ok {
				d.HighCount = val
			}
		}
		if row.MediumCount != nil {
			if val, ok := row.MediumCount.(int64); ok {
				d.MediumCount = val
			}
		}
		if row.LowCount != nil {
			if val, ok := row.LowCount.(int64); ok {
				d.LowCount = val
			}
		}
		if row.OpenIssues != nil {
			if val, ok := row.OpenIssues.(int64); ok {
				d.OpenIssues = val
			}
		}

		if row.IpAddress.Valid {
			d.IPAddress = &row.IpAddress.String
		}
		if row.Description.Valid {
			d.Description = &row.Description.String
		}
		if row.RiskScore.Valid {
			// RiskScore is stored as DECIMAL in DB, returned as string by sqlc
			var score float64
			fmt.Sscanf(row.RiskScore.String, "%f", &score)
			d.RiskScore = &score
		}
		if row.IsActive.Valid {
			d.IsActive = row.IsActive.Bool
		}
		if row.LastScanned.Valid {
			d.LastScanned = &row.LastScanned.Time
		}
		if row.CreatedAt.Valid {
			d.CreatedAt = row.CreatedAt.Time
		}
		if row.UpdatedAt.Valid {
			d.UpdatedAt = row.UpdatedAt.Time
		}
		if row.SlaCompliance != nil {
			if val, ok := row.SlaCompliance.(float64); ok {
				d.SLACompliance = val
			}
		}

		domains = append(domains, d)
	}

	return domains, nil
}

// CountDomains counts total active domains
func (r *DomainsRepository) CountDomains(ctx context.Context) (int64, error) {
	count, err := r.queries.CountDomains(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count domains: %w", err)
	}
	return count, nil
}

// GetDomainByID retrieves a domain by ID
func (r *DomainsRepository) GetDomainByID(ctx context.Context, id string) (*domain.Domain, error) {
	domainID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid domain ID: %w", err)
	}

	d, err := r.queries.GetDomainByID(ctx, domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	return r.toDomain(&d), nil
}

// CreateDomain creates a new domain
func (r *DomainsRepository) CreateDomain(ctx context.Context, params *domain.CreateDomainParams) (*domain.Domain, error) {
	domainID, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid domain ID: %w", err)
	}

	var ipAddress, description sql.NullString
	if params.IPAddress != nil {
		ipAddress = sql.NullString{String: *params.IPAddress, Valid: true}
	}
	if params.Description != nil {
		description = sql.NullString{String: *params.Description, Valid: true}
	}

	var riskScore sql.NullString
	if params.RiskScore != nil {
		riskScore = sql.NullString{String: fmt.Sprintf("%.1f", *params.RiskScore), Valid: true}
	}

	var lastScanned sql.NullTime
	if params.LastScanned != nil {
		lastScanned = sql.NullTime{Time: *params.LastScanned, Valid: true}
	}

	sqlcParams := sqlc.CreateDomainParams{
		ID:          domainID,
		DomainName:  params.DomainName,
		IpAddress:   ipAddress,
		Description: description,
		RiskScore:   riskScore,
		IsActive:    sql.NullBool{Bool: params.IsActive, Valid: true},
		LastScanned: lastScanned,
		CreatedAt:   sql.NullTime{Time: params.CreatedAt, Valid: true},
		UpdatedAt:   sql.NullTime{Time: params.UpdatedAt, Valid: true},
	}

	d, err := r.queries.CreateDomain(ctx, sqlcParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}

	return r.toDomain(&d), nil
}

// UpdateDomain updates a domain
func (r *DomainsRepository) UpdateDomain(ctx context.Context, params *domain.UpdateDomainParams) (*domain.Domain, error) {
	domainID, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid domain ID: %w", err)
	}

	// Get existing domain
	existing, err := r.queries.GetDomainByID(ctx, domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing domain: %w", err)
	}

	// Use existing values if not provided in update
	domainName := existing.DomainName
	if params.DomainName != nil {
		domainName = *params.DomainName
	}

	ipAddress := existing.IpAddress
	if params.IPAddress != nil {
		ipAddress = sql.NullString{String: *params.IPAddress, Valid: true}
	}

	description := existing.Description
	if params.Description != nil {
		description = sql.NullString{String: *params.Description, Valid: true}
	}

	riskScore := existing.RiskScore
	if params.RiskScore != nil {
		riskScore = sql.NullString{String: fmt.Sprintf("%.1f", *params.RiskScore), Valid: true}
	}

	lastScanned := existing.LastScanned
	if params.LastScanned != nil {
		lastScanned = sql.NullTime{Time: *params.LastScanned, Valid: true}
	}

	sqlcParams := sqlc.UpdateDomainParams{
		ID:          domainID,
		DomainName:  domainName,
		IpAddress:   ipAddress,
		Description: description,
		RiskScore:   riskScore,
		LastScanned: lastScanned,
		UpdatedAt:   sql.NullTime{Time: params.UpdatedAt, Valid: true},
	}

	d, err := r.queries.UpdateDomain(ctx, sqlcParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update domain: %w", err)
	}

	return r.toDomain(&d), nil
}

// DeleteDomain soft deletes a domain
func (r *DomainsRepository) DeleteDomain(ctx context.Context, id string) error {
	domainID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid domain ID: %w", err)
	}

	err = r.queries.DeleteDomain(ctx, sqlc.DeleteDomainParams{
		ID:        domainID,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}

	return nil
}

// GetSecurityMetrics retrieves security metrics
func (r *DomainsRepository) GetSecurityMetrics(ctx context.Context) ([]domain.SecurityMetric, error) {
	rows, err := r.queries.GetSecurityMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get security metrics: %w", err)
	}

	metrics := make([]domain.SecurityMetric, 0, len(rows))
	for _, row := range rows {
		avgValue := 0.0
		if row.AvgValue != nil {
			switch v := row.AvgValue.(type) {
			case float64:
				avgValue = v
			case string:
				fmt.Sscanf(v, "%f", &avgValue)
			case []byte:
				fmt.Sscanf(string(v), "%f", &avgValue)
			}
		}

		metrics = append(metrics, domain.SecurityMetric{
			MetricName: row.MetricName,
			AvgValue:   avgValue,
		})
	}

	return metrics, nil
}

// CreateSecurityMetric creates a new security metric
func (r *DomainsRepository) CreateSecurityMetric(ctx context.Context, params *domain.CreateSecurityMetricParams) (*domain.SecurityMetric, error) {
	metricID, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %w", err)
	}

	domainID, err := uuid.Parse(params.DomainID)
	if err != nil {
		return nil, fmt.Errorf("invalid domain ID: %w", err)
	}

	_, err = r.queries.CreateSecurityMetric(ctx, sqlc.CreateSecurityMetricParams{
		ID:          metricID,
		DomainID:    uuid.NullUUID{UUID: domainID, Valid: true},
		MetricName:  params.MetricName,
		MetricValue: sql.NullString{String: fmt.Sprintf("%.2f", params.MetricValue), Valid: true},
		MeasuredAt:  sql.NullTime{Time: params.MeasuredAt, Valid: true},
		CreatedAt:   sql.NullTime{Time: params.CreatedAt, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create security metric: %w", err)
	}

	return &domain.SecurityMetric{
		MetricName: params.MetricName,
		AvgValue:   params.MetricValue,
	}, nil
}

// GetSLABreachAnalysis retrieves SLA breach analysis
func (r *DomainsRepository) GetSLABreachAnalysis(ctx context.Context) ([]domain.SLABreachDomain, error) {
	rows, err := r.queries.GetSLABreachAnalysis(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get SLA breach analysis: %w", err)
	}

	domains := make([]domain.SLABreachDomain, 0, len(rows))
	for _, row := range rows {
		compliance := 0.0
		if row.SlaCompliancePercent != nil {
			switch v := row.SlaCompliancePercent.(type) {
			case float64:
				compliance = v
			case string:
				fmt.Sscanf(v, "%f", &compliance)
			case []byte:
				fmt.Sscanf(string(v), "%f", &compliance)
			}
		}

		domains = append(domains, domain.SLABreachDomain{
			DomainName:           row.DomainName,
			SLACompliancePercent: compliance,
		})
	}

	return domains, nil
}

// toDomain converts sqlc.Domain to domain.Domain
func (r *DomainsRepository) toDomain(d *sqlc.Domain) *domain.Domain {
	domainModel := &domain.Domain{
		ID:         d.ID,
		DomainName: d.DomainName,
	}

	if d.IpAddress.Valid {
		domainModel.IPAddress = &d.IpAddress.String
	}
	if d.RiskScore.Valid {
		// RiskScore is stored as DECIMAL in DB, returned as string by sqlc
		var score float64
		fmt.Sscanf(d.RiskScore.String, "%f", &score)
		domainModel.RiskScore = &score
	}
	if d.LastScanned.Valid {
		domainModel.LastScannedAt = &d.LastScanned.Time
	}
	if d.CreatedAt.Valid {
		domainModel.CreatedAt = d.CreatedAt.Time
	}
	if d.UpdatedAt.Valid {
		domainModel.UpdatedAt = d.UpdatedAt.Time
	}

	return domainModel
}
