package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/internal/infra/cache"
)

type CachedProjectsUseCase struct {
	usecase domain.ProjectsUseCase
	cache   *cache.Cache
}

func NewCachedProjectsUseCase(usecase domain.ProjectsUseCase, cache *cache.Cache) *CachedProjectsUseCase {
	return &CachedProjectsUseCase{
		usecase: usecase,
		cache:   cache,
	}
}

// Cache keys
const (
	cacheKeyProjectsList  = "projects:list:%d:%d"      // page:perPage
	cacheKeyProjectStats  = "projects:stats"
	cacheKeyPentesters    = "projects:pentesters"
	cacheTTLProjectsList  = 5 * time.Minute
	cacheTTLProjectStats  = 5 * time.Minute
	cacheTTLPentesters    = 10 * time.Minute
)

// CreateProject creates a new project and invalidates cache
func (uc *CachedProjectsUseCase) CreateProject(ctx context.Context, req *dto.CreateProjectRequest, createdBy string) (*dto.CreateProjectResponse, error) {
	// Create project
	response, err := uc.usecase.CreateProject(ctx, req, createdBy)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	uc.invalidateProjectsCache()

	return response, nil
}

// ListProjects retrieves projects with caching
func (uc *CachedProjectsUseCase) ListProjects(ctx context.Context, page, perPage int) (*dto.ListProjectsResponse, error) {
	// Try to get from cache
	cacheKey := fmt.Sprintf(cacheKeyProjectsList, page, perPage)
	if cached, found := uc.cache.Get(cacheKey); found {
		if data, ok := cached.([]byte); ok {
			var response dto.ListProjectsResponse
			if err := json.Unmarshal(data, &response); err == nil {
				return &response, nil
			}
		}
	}

	// Cache miss - get from use case
	response, err := uc.usecase.ListProjects(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if data, err := json.Marshal(response); err == nil {
		uc.cache.Set(cacheKey, data, cacheTTLProjectsList)
	}

	return response, nil
}

// UpdateProject updates a project and invalidates cache
func (uc *CachedProjectsUseCase) UpdateProject(ctx context.Context, projectID string, req *dto.UpdateProjectRequest) (*dto.UpdateProjectResponse, error) {
	// Update project
	response, err := uc.usecase.UpdateProject(ctx, projectID, req)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	uc.invalidateProjectsCache()

	return response, nil
}

// DeleteProject deletes a project and invalidates cache
func (uc *CachedProjectsUseCase) DeleteProject(ctx context.Context, projectID string) error {
	// Delete project
	if err := uc.usecase.DeleteProject(ctx, projectID); err != nil {
		return err
	}

	// Invalidate cache
	uc.invalidateProjectsCache()

	return nil
}

// GetProjectStats retrieves project statistics with caching
func (uc *CachedProjectsUseCase) GetProjectStats(ctx context.Context) (*dto.ProjectStatsResponse, error) {
	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKeyProjectStats); found {
		if data, ok := cached.([]byte); ok {
			var response dto.ProjectStatsResponse
			if err := json.Unmarshal(data, &response); err == nil {
				return &response, nil
			}
		}
	}

	// Cache miss - get from use case
	response, err := uc.usecase.GetProjectStats(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if data, err := json.Marshal(response); err == nil {
		uc.cache.Set(cacheKeyProjectStats, data, cacheTTLProjectStats)
	}

	return response, nil
}

// GetPentesters retrieves pentesters with caching
func (uc *CachedProjectsUseCase) GetPentesters(ctx context.Context) (*dto.GetPentestersResponse, error) {
	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKeyPentesters); found {
		if data, ok := cached.([]byte); ok {
			var response dto.GetPentestersResponse
			if err := json.Unmarshal(data, &response); err == nil {
				return &response, nil
			}
		}
	}

	// Cache miss - get from use case
	response, err := uc.usecase.GetPentesters(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if data, err := json.Marshal(response); err == nil {
		uc.cache.Set(cacheKeyPentesters, data, cacheTTLPentesters)
	}

	return response, nil
}

// invalidateProjectsCache clears all projects-related cache
func (uc *CachedProjectsUseCase) invalidateProjectsCache() {
	// Clear stats cache
	uc.cache.Delete(cacheKeyProjectStats)

	// Clear list cache for common pagination values
	perPageValues := []int{5, 10, 20, 50}
	for page := 1; page <= 10; page++ {
		for _, perPage := range perPageValues {
			cacheKey := fmt.Sprintf(cacheKeyProjectsList, page, perPage)
			uc.cache.Delete(cacheKey)
		}
	}
}

