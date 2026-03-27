package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
)

// CachedDomainsUseCase wraps DomainsUseCase with caching
type CachedDomainsUseCase struct {
	usecase domain.DomainsUseCase
	cache   *ristretto.Cache
}

// NewCachedDomainsUseCase creates a new CachedDomainsUseCase
func NewCachedDomainsUseCase(usecase domain.DomainsUseCase, cache *ristretto.Cache) *CachedDomainsUseCase {
	return &CachedDomainsUseCase{
		usecase: usecase,
		cache:   cache,
	}
}

// GetDomainsStats retrieves domains stats with caching
func (uc *CachedDomainsUseCase) GetDomainsStats(ctx context.Context) (*dto.DomainsStatsResponse, error) {
	cacheKey := "domains:stats"

	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKey); found {
		if stats, ok := cached.(*dto.DomainsStatsResponse); ok {
			return stats, nil
		}
	}

	// Get from use case
	stats, err := uc.usecase.GetDomainsStats(ctx)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	uc.cache.SetWithTTL(cacheKey, stats, 1, 5*time.Minute)

	return stats, nil
}

// ListDomains retrieves domains list with caching
func (uc *CachedDomainsUseCase) ListDomains(ctx context.Context, req *dto.ListDomainsRequest) (*dto.ListDomainsResponse, error) {
	cacheKey := fmt.Sprintf("domains:list:%d:%d", req.Page, req.PerPage)

	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKey); found {
		if domains, ok := cached.(*dto.ListDomainsResponse); ok {
			return domains, nil
		}
	}

	// Get from use case
	domains, err := uc.usecase.ListDomains(ctx, req)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	uc.cache.SetWithTTL(cacheKey, domains, 1, 5*time.Minute)

	return domains, nil
}

// GetDomainByID retrieves a domain by ID with caching
func (uc *CachedDomainsUseCase) GetDomainByID(ctx context.Context, id string) (*dto.DomainResponse, error) {
	cacheKey := fmt.Sprintf("domains:id:%s", id)

	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKey); found {
		if domain, ok := cached.(*dto.DomainResponse); ok {
			return domain, nil
		}
	}

	// Get from use case
	domain, err := uc.usecase.GetDomainByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache for 10 minutes
	uc.cache.SetWithTTL(cacheKey, domain, 1, 10*time.Minute)

	return domain, nil
}

// CreateDomain creates a domain and invalidates cache
func (uc *CachedDomainsUseCase) CreateDomain(ctx context.Context, req *dto.CreateDomainRequest) (*dto.DomainResponse, error) {
	domain, err := uc.usecase.CreateDomain(ctx, req)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	uc.invalidateCache()

	return domain, nil
}

// UpdateDomain updates a domain and invalidates cache
func (uc *CachedDomainsUseCase) UpdateDomain(ctx context.Context, id string, req *dto.UpdateDomainRequest) (*dto.DomainResponse, error) {
	domain, err := uc.usecase.UpdateDomain(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	uc.invalidateCache()
	uc.cache.Del(fmt.Sprintf("domains:id:%s", id))

	return domain, nil
}

// DeleteDomain deletes a domain and invalidates cache
func (uc *CachedDomainsUseCase) DeleteDomain(ctx context.Context, id string) error {
	err := uc.usecase.DeleteDomain(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache
	uc.invalidateCache()
	uc.cache.Del(fmt.Sprintf("domains:id:%s", id))

	return nil
}

// GetSecurityMetrics retrieves security metrics with caching
func (uc *CachedDomainsUseCase) GetSecurityMetrics(ctx context.Context) (*dto.SecurityMetricsResponse, error) {
	cacheKey := "domains:security_metrics"

	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKey); found {
		if metrics, ok := cached.(*dto.SecurityMetricsResponse); ok {
			return metrics, nil
		}
	}

	// Get from use case
	metrics, err := uc.usecase.GetSecurityMetrics(ctx)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	uc.cache.SetWithTTL(cacheKey, metrics, 1, 5*time.Minute)

	return metrics, nil
}

// CreateSecurityMetric creates a security metric and invalidates cache
func (uc *CachedDomainsUseCase) CreateSecurityMetric(ctx context.Context, req *dto.CreateSecurityMetricRequest) error {
	err := uc.usecase.CreateSecurityMetric(ctx, req)
	if err != nil {
		return err
	}

	// Invalidate security metrics cache
	uc.cache.Del("domains:security_metrics")

	return nil
}

// GetSLABreachAnalysis retrieves SLA breach analysis with caching
func (uc *CachedDomainsUseCase) GetSLABreachAnalysis(ctx context.Context) (*dto.SLABreachAnalysisResponse, error) {
	cacheKey := "domains:sla_breach"

	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKey); found {
		if analysis, ok := cached.(*dto.SLABreachAnalysisResponse); ok {
			return analysis, nil
		}
	}

	// Get from use case
	analysis, err := uc.usecase.GetSLABreachAnalysis(ctx)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	uc.cache.SetWithTTL(cacheKey, analysis, 1, 5*time.Minute)

	return analysis, nil
}

// invalidateCache invalidates all domains-related cache entries
func (uc *CachedDomainsUseCase) invalidateCache() {
	// Invalidate stats and lists
	uc.cache.Del("domains:stats")
	uc.cache.Del("domains:sla_breach")
	
	// Note: We can't easily invalidate all list cache entries without tracking them
	// In production, consider using a cache key prefix pattern or cache tags
}

// Helper function to serialize/deserialize cache data (if needed)
func serializeToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func deserializeFromJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

