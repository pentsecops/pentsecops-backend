package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/internal/infra/cache"
)

// CachedNotificationsUseCase wraps NotificationsUseCase with caching
type CachedNotificationsUseCase struct {
	usecase domain.NotificationsUseCase
	cache   *ristretto.Cache
}

// NewCachedNotificationsUseCase creates a new CachedNotificationsUseCase
func NewCachedNotificationsUseCase(usecase domain.NotificationsUseCase, cacheInstance *ristretto.Cache) *CachedNotificationsUseCase {
	return &CachedNotificationsUseCase{
		usecase: usecase,
		cache:   cacheInstance,
	}
}

// GetTotalNotificationsSent retrieves the total count of notifications sent with caching
func (uc *CachedNotificationsUseCase) GetTotalNotificationsSent(ctx context.Context, createdBy string) (*dto.TotalNotificationsSentResponse, error) {
	cacheKey := fmt.Sprintf("%s:%s", cache.NotificationsTotalKey, createdBy)

	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKey); found {
		if result, ok := cached.(*dto.TotalNotificationsSentResponse); ok {
			return result, nil
		}
	}

	// Get from use case
	result, err := uc.usecase.GetTotalNotificationsSent(ctx, createdBy)
	if err != nil {
		return nil, err
	}

	// Cache the result
	uc.cache.SetWithTTL(cacheKey, result, 1, cache.NotificationsTotalTTL)

	return result, nil
}

// ListNotifications retrieves a paginated list of notifications with caching
func (uc *CachedNotificationsUseCase) ListNotifications(ctx context.Context, createdBy string, req *dto.ListNotificationsRequest) (*dto.ListNotificationsResponse, error) {
	// Create cache key
	reqJSON, _ := json.Marshal(req)
	cacheKey := fmt.Sprintf("%s:%s:%s", cache.NotificationsListKey, createdBy, string(reqJSON))

	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKey); found {
		if result, ok := cached.(*dto.ListNotificationsResponse); ok {
			return result, nil
		}
	}

	// Get from use case
	result, err := uc.usecase.ListNotifications(ctx, createdBy, req)
	if err != nil {
		return nil, err
	}

	// Cache the result
	uc.cache.SetWithTTL(cacheKey, result, 1, cache.NotificationsListTTL)

	return result, nil
}

// CreateNotification creates a new notification and invalidates cache
func (uc *CachedNotificationsUseCase) CreateNotification(ctx context.Context, createdBy string, req *dto.CreateNotificationRequest) (*dto.NotificationResponse, error) {
	// Create notification
	result, err := uc.usecase.CreateNotification(ctx, createdBy, req)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	uc.invalidateNotificationsCache(createdBy)

	return result, nil
}

// ListImportantAlerts retrieves a paginated list of important alerts with caching
func (uc *CachedNotificationsUseCase) ListImportantAlerts(ctx context.Context, req *dto.ListAlertsRequest) (*dto.ListAlertsResponse, error) {
	// Create cache key
	reqJSON, _ := json.Marshal(req)
	cacheKey := fmt.Sprintf("%s:%s", cache.AlertsListKey, string(reqJSON))

	// Try to get from cache
	if cached, found := uc.cache.Get(cacheKey); found {
		if result, ok := cached.(*dto.ListAlertsResponse); ok {
			return result, nil
		}
	}

	// Get from use case
	result, err := uc.usecase.ListImportantAlerts(ctx, req)
	if err != nil {
		return nil, err
	}

	// Cache the result
	uc.cache.SetWithTTL(cacheKey, result, 1, cache.AlertsListTTL)

	return result, nil
}

// invalidateNotificationsCache invalidates all notifications cache for a user
func (uc *CachedNotificationsUseCase) invalidateNotificationsCache(createdBy string) {
	// Wait for cache to process pending operations
	time.Sleep(10 * time.Millisecond)

	// Delete specific keys
	uc.cache.Del(fmt.Sprintf("%s:%s", cache.NotificationsTotalKey, createdBy))

	// Note: For list cache, we would need to track all possible page combinations
	// For simplicity, we're relying on TTL expiration
	// In production, consider using cache tags or patterns for bulk invalidation
}

