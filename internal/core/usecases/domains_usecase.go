package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
)

// DomainsUseCase implements the domain.DomainsUseCase interface
type DomainsUseCase struct {
	repo      domain.DomainsRepository
	validator *validator.Validate
}

// NewDomainsUseCase creates a new DomainsUseCase
func NewDomainsUseCase(repo domain.DomainsRepository) *DomainsUseCase {
	return &DomainsUseCase{
		repo:      repo,
		validator: validator.New(),
	}
}

// GetDomainsStats - UC55-58: Fetch and display domains overview statistics
func (uc *DomainsUseCase) GetDomainsStats(ctx context.Context) (*dto.DomainsStatsResponse, error) {
	stats, err := uc.repo.GetDomainsStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get domains stats: %w", err)
	}

	return &dto.DomainsStatsResponse{
		TotalDomains:         stats.TotalDomains,
		AverageRiskScore:     stats.AverageRiskScore,
		CriticalIssues:       stats.CriticalIssues,
		SLACompliancePercent: stats.SLACompliancePercent,
	}, nil
}

// ListDomains - UC59-60: Display all domains with pagination
func (uc *DomainsUseCase) ListDomains(ctx context.Context, req *dto.ListDomainsRequest) (*dto.ListDomainsResponse, error) {
	// Validate request
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Set defaults
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 5
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Get domains
	domains, err := uc.repo.ListDomains(ctx, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}

	// Get total count
	total, err := uc.repo.CountDomains(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count domains: %w", err)
	}

	// Convert to response
	domainResponses := make([]dto.DomainResponse, 0, len(domains))
	for _, d := range domains {
		domainResponses = append(domainResponses, uc.toDomainResponse(&d))
	}

	// Calculate pagination
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return &dto.ListDomainsResponse{
		Domains: domainResponses,
		Pagination: dto.PaginationInfo{
			CurrentPage: page,
			PerPage:     perPage,
			Total:       total,
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrev:     page > 1,
		},
	}, nil
}

// GetDomainByID retrieves a domain by ID
func (uc *DomainsUseCase) GetDomainByID(ctx context.Context, id string) (*dto.DomainResponse, error) {
	d, err := uc.repo.GetDomainByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	// Get domain stats (we need to call ListDomains with filter, but for now return basic info)
	response := &dto.DomainResponse{
		ID:          d.ID.String(),
		DomainName:  d.DomainName,
		IPAddress:   d.IPAddress,
		RiskScore:   d.RiskScore,
		LastScanned: d.LastScannedAt,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}

	// Calculate risk level
	if d.RiskScore != nil {
		response.RiskLevel = calculateRiskLevel(*d.RiskScore)
	}

	return response, nil
}

// CreateDomain creates a new domain
func (uc *DomainsUseCase) CreateDomain(ctx context.Context, req *dto.CreateDomainRequest) (*dto.DomainResponse, error) {
	// Validate request
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Create domain params
	now := time.Now()
	params := &domain.CreateDomainParams{
		ID:          uuid.Must(uuid.NewV7()).String(),
		DomainName:  req.DomainName,
		IPAddress:   req.IPAddress,
		Description: req.Description,
		RiskScore:   req.RiskScore,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Create domain
	d, err := uc.repo.CreateDomain(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}

	response := &dto.DomainResponse{
		ID:          d.ID.String(),
		DomainName:  d.DomainName,
		IPAddress:   d.IPAddress,
		RiskScore:   d.RiskScore,
		LastScanned: d.LastScannedAt,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}

	// Calculate risk level
	if d.RiskScore != nil {
		response.RiskLevel = calculateRiskLevel(*d.RiskScore)
	}

	return response, nil
}

// UpdateDomain updates a domain
func (uc *DomainsUseCase) UpdateDomain(ctx context.Context, id string, req *dto.UpdateDomainRequest) (*dto.DomainResponse, error) {
	// Validate request
	if err := uc.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Update domain params
	params := &domain.UpdateDomainParams{
		ID:          id,
		DomainName:  req.DomainName,
		IPAddress:   req.IPAddress,
		Description: req.Description,
		RiskScore:   req.RiskScore,
		UpdatedAt:   time.Now(),
	}

	// Update domain
	d, err := uc.repo.UpdateDomain(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update domain: %w", err)
	}

	response := &dto.DomainResponse{
		ID:          d.ID.String(),
		DomainName:  d.DomainName,
		IPAddress:   d.IPAddress,
		RiskScore:   d.RiskScore,
		LastScanned: d.LastScannedAt,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}

	// Calculate risk level
	if d.RiskScore != nil {
		response.RiskLevel = calculateRiskLevel(*d.RiskScore)
	}

	return response, nil
}

// DeleteDomain deletes a domain
func (uc *DomainsUseCase) DeleteDomain(ctx context.Context, id string) error {
	return uc.repo.DeleteDomain(ctx, id)
}

// GetSecurityMetrics - UC65: Render security metrics radar chart
func (uc *DomainsUseCase) GetSecurityMetrics(ctx context.Context) (*dto.SecurityMetricsResponse, error) {
	metrics, err := uc.repo.GetSecurityMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get security metrics: %w", err)
	}

	// Initialize response with default values
	response := &dto.SecurityMetricsResponse{
		Authentication:  0,
		Authorization:   0,
		InputValidation: 0,
		Encryption:      0,
		Configuration:   0,
		NetworkSecurity: 0,
	}

	// Map metrics to response
	for _, m := range metrics {
		switch m.MetricName {
		case "authentication":
			response.Authentication = m.AvgValue
		case "authorization":
			response.Authorization = m.AvgValue
		case "input_validation":
			response.InputValidation = m.AvgValue
		case "encryption":
			response.Encryption = m.AvgValue
		case "configuration":
			response.Configuration = m.AvgValue
		case "network_security":
			response.NetworkSecurity = m.AvgValue
		}
	}

	return response, nil
}

// CreateSecurityMetric creates a new security metric
func (uc *DomainsUseCase) CreateSecurityMetric(ctx context.Context, req *dto.CreateSecurityMetricRequest) error {
	// Validate request
	if err := uc.validator.Struct(req); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Create security metric params
	now := time.Now()
	params := &domain.CreateSecurityMetricParams{
		ID:          uuid.Must(uuid.NewV7()).String(),
		DomainID:    req.DomainID,
		MetricName:  req.MetricName,
		MetricValue: req.MetricValue,
		MeasuredAt:  now,
		CreatedAt:   now,
	}

	// Create security metric
	_, err := uc.repo.CreateSecurityMetric(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create security metric: %w", err)
	}

	return nil
}

// GetSLABreachAnalysis - UC66: Display SLA breach analysis
func (uc *DomainsUseCase) GetSLABreachAnalysis(ctx context.Context) (*dto.SLABreachAnalysisResponse, error) {
	domains, err := uc.repo.GetSLABreachAnalysis(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get SLA breach analysis: %w", err)
	}

	items := make([]dto.SLABreachItem, 0, len(domains))
	for _, d := range domains {
		items = append(items, dto.SLABreachItem{
			DomainName:           d.DomainName,
			SLACompliancePercent: d.SLACompliancePercent,
		})
	}

	return &dto.SLABreachAnalysisResponse{
		Domains: items,
	}, nil
}

// toDomainResponse converts domain.DomainWithStats to dto.DomainResponse
func (uc *DomainsUseCase) toDomainResponse(d *domain.DomainWithStats) dto.DomainResponse {
	response := dto.DomainResponse{
		ID:                   d.ID,
		DomainName:           d.DomainName,
		IPAddress:            d.IPAddress,
		Description:          d.Description,
		RiskScore:            d.RiskScore,
		TotalVulnerabilities: d.TotalVulnerabilities,
		CriticalCount:        d.CriticalCount,
		HighCount:            d.HighCount,
		MediumCount:          d.MediumCount,
		LowCount:             d.LowCount,
		SLACompliance:        d.SLACompliance,
		OpenIssues:           d.OpenIssues,
		LastScanned:          d.LastScanned,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}

	// UC61: Display risk score with color coding
	if d.RiskScore != nil {
		response.RiskLevel = calculateRiskLevel(*d.RiskScore)
	} else {
		response.RiskLevel = "Unknown"
	}

	return response
}

// calculateRiskLevel - UC61: Calculate risk level based on risk score
func calculateRiskLevel(score float64) string {
	if score >= 7 {
		return "High Risk"
	} else if score >= 5 {
		return "Medium Risk"
	} else if score >= 3 {
		return "Low Risk"
	}
	return "Minimal Risk"
}
