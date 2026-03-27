package usecases

import (
	"context"
	"fmt"

	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/internal/infra/cache"
)

// CachedUsersUseCase wraps UsersUseCase with caching
type CachedUsersUseCase struct {
	useCase domain.UsersUseCase
	cache   *cache.Cache
}

// NewCachedUsersUseCase creates a new CachedUsersUseCase
func NewCachedUsersUseCase(useCase domain.UsersUseCase, c *cache.Cache) domain.UsersUseCase {
	return &CachedUsersUseCase{
		useCase: useCase,
		cache:   c,
	}
}

// CreateUser creates a new user and invalidates cache
func (uc *CachedUsersUseCase) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	response, err := uc.useCase.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	// Invalidate cache after creating user
	uc.invalidateUsersCache()

	return response, nil
}

// ListUsers retrieves users with caching
func (uc *CachedUsersUseCase) ListUsers(ctx context.Context, page, perPage int) (*dto.ListUsersResponse, error) {
	// Try to get from cache
	cacheKey := fmt.Sprintf(cache.CacheKeyUsersList, page, perPage)
	if cached, found := uc.cache.Get(cacheKey); found {
		if response, ok := cached.(*dto.ListUsersResponse); ok {
			return response, nil
		}
	}

	// Cache miss - get from use case
	response, err := uc.useCase.ListUsers(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	// Store in cache
	uc.cache.Set(cacheKey, response, cache.CacheTTLUsersList)

	return response, nil
}

// UpdateUser updates a user and invalidates cache
func (uc *CachedUsersUseCase) UpdateUser(ctx context.Context, userID string, req *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	response, err := uc.useCase.UpdateUser(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	// Invalidate cache after updating user
	uc.invalidateUsersCache()
	uc.cache.Delete(fmt.Sprintf(cache.CacheKeyUserByID, userID))

	return response, nil
}

// DeleteUser deletes a user and invalidates cache
func (uc *CachedUsersUseCase) DeleteUser(ctx context.Context, userID string) error {
	err := uc.useCase.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	// Invalidate cache after deleting user
	uc.invalidateUsersCache()
	uc.cache.Delete(fmt.Sprintf(cache.CacheKeyUserByID, userID))

	return nil
}

// RefreshUsers refreshes the user list and clears cache
func (uc *CachedUsersUseCase) RefreshUsers(ctx context.Context, page, perPage int) (*dto.ListUsersResponse, error) {
	// Clear cache before refreshing
	uc.invalidateUsersCache()

	// Get fresh data
	return uc.useCase.RefreshUsers(ctx, page, perPage)
}

// GetUserStats retrieves user statistics with caching
func (uc *CachedUsersUseCase) GetUserStats(ctx context.Context) (*dto.UserStatsResponse, error) {
	// Try to get from cache
	cacheKey := cache.CacheKeyUserStats
	if cached, found := uc.cache.Get(cacheKey); found {
		if stats, ok := cached.(*dto.UserStatsResponse); ok {
			return stats, nil
		}
	}

	// Cache miss - get from use case
	stats, err := uc.useCase.GetUserStats(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache
	uc.cache.Set(cacheKey, stats, cache.CacheTTLUserStats)

	return stats, nil
}

// ExportUsersToCSV exports users to CSV (no caching for exports)
func (uc *CachedUsersUseCase) ExportUsersToCSV(ctx context.Context) ([]byte, error) {
	return uc.useCase.ExportUsersToCSV(ctx)
}

// invalidateUsersCache clears all users-related cache entries
func (uc *CachedUsersUseCase) invalidateUsersCache() {
	// Clear list cache for common pagination values
	for page := 1; page <= 10; page++ {
		for _, perPage := range []int{5, 10, 20, 50} {
			cacheKey := fmt.Sprintf(cache.CacheKeyUsersList, page, perPage)
			uc.cache.Delete(cacheKey)
		}
	}

	// Clear stats cache
	uc.cache.Delete(cache.CacheKeyUserStats)
}

