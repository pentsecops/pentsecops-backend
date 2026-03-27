package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/pentsecops/backend/internal/core/domain/dto"
)

type CachedVulnerabilitiesUseCase struct {
	base  *VulnerabilitiesUseCase
	cache *ristretto.Cache
}

func NewCachedVulnerabilitiesUseCase(base *VulnerabilitiesUseCase, cache *ristretto.Cache) *CachedVulnerabilitiesUseCase {
	return &CachedVulnerabilitiesUseCase{
		base:  base,
		cache: cache,
	}
}

func (uc *CachedVulnerabilitiesUseCase) CreateVulnerability(ctx context.Context, req *dto.CreateVulnerabilityRequest) (*dto.VulnerabilityResponse, error) {
	resp, err := uc.base.CreateVulnerability(ctx, req)
	if err != nil {
		return nil, err
	}

	// Invalidate list and stats caches
	uc.invalidateListCaches()

	return resp, nil
}

func (uc *CachedVulnerabilitiesUseCase) GetVulnerabilityByID(ctx context.Context, id string) (*dto.VulnerabilityResponse, error) {
	cacheKey := fmt.Sprintf("vulnerability:%s", id)

	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKey); found {
		if resp, ok := cached.(*dto.VulnerabilityResponse); ok {
			return resp, nil
		}
	}

	// Get from base use case
	resp, err := uc.base.GetVulnerabilityByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache for 10 minutes
	uc.cache.SetWithTTL(cacheKey, resp, 1, 10*time.Minute)

	return resp, nil
}

func (uc *CachedVulnerabilitiesUseCase) UpdateVulnerability(ctx context.Context, id string, req *dto.UpdateVulnerabilityRequest) (*dto.VulnerabilityResponse, error) {
	resp, err := uc.base.UpdateVulnerability(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// Invalidate caches
	uc.cache.Del(fmt.Sprintf("vulnerability:%s", id))
	uc.invalidateListCaches()

	return resp, nil
}

func (uc *CachedVulnerabilitiesUseCase) DeleteVulnerability(ctx context.Context, id string) error {
	err := uc.base.DeleteVulnerability(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate caches
	uc.cache.Del(fmt.Sprintf("vulnerability:%s", id))
	uc.invalidateListCaches()

	return nil
}

func (uc *CachedVulnerabilitiesUseCase) ListVulnerabilities(ctx context.Context, req *dto.ListVulnerabilitiesRequest) (*dto.ListVulnerabilitiesResponse, error) {
	// Create cache key based on request parameters
	cacheKeyData, _ := json.Marshal(req)
	cacheKey := fmt.Sprintf("vulnerabilities:list:%s", string(cacheKeyData))

	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKey); found {
		if resp, ok := cached.(*dto.ListVulnerabilitiesResponse); ok {
			return resp, nil
		}
	}

	// Get from base use case
	resp, err := uc.base.ListVulnerabilities(ctx, req)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	uc.cache.SetWithTTL(cacheKey, resp, 1, 5*time.Minute)

	return resp, nil
}

func (uc *CachedVulnerabilitiesUseCase) GetVulnerabilityStats(ctx context.Context) (*dto.VulnerabilityStatsResponse, error) {
	cacheKey := "vulnerabilities:stats"

	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKey); found {
		if resp, ok := cached.(*dto.VulnerabilityStatsResponse); ok {
			return resp, nil
		}
	}

	// Get from base use case
	resp, err := uc.base.GetVulnerabilityStats(ctx)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	uc.cache.SetWithTTL(cacheKey, resp, 1, 5*time.Minute)

	return resp, nil
}

func (uc *CachedVulnerabilitiesUseCase) GetSLACompliance(ctx context.Context) (*dto.SLAComplianceResponse, error) {
	cacheKey := "vulnerabilities:sla"

	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKey); found {
		if resp, ok := cached.(*dto.SLAComplianceResponse); ok {
			return resp, nil
		}
	}

	// Get from base use case
	resp, err := uc.base.GetSLACompliance(ctx)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	uc.cache.SetWithTTL(cacheKey, resp, 1, 5*time.Minute)

	return resp, nil
}

func (uc *CachedVulnerabilitiesUseCase) ExportVulnerabilitiesToCSV(ctx context.Context, req *dto.ListVulnerabilitiesRequest) ([]byte, error) {
	// No caching for export
	return uc.base.ExportVulnerabilitiesToCSV(ctx, req)
}

// Helper function to invalidate list-related caches
func (uc *CachedVulnerabilitiesUseCase) invalidateListCaches() {
	// Invalidate stats cache
	uc.cache.Del("vulnerabilities:stats")
	uc.cache.Del("vulnerabilities:sla")

	// Note: We can't easily invalidate all list caches with different parameters
	// In a production system, you might want to use cache tags or a more sophisticated invalidation strategy
}

