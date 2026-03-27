package cache

import (
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
)

// Cache wraps ristretto cache
type Cache struct {
	client *ristretto.Cache
}

// NewCache creates a new cache instance
func NewCache() (*Cache, error) {
	client, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	return &Cache{client: client}, nil
}

// Set stores a value in cache with TTL
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) bool {
	return c.client.SetWithTTL(key, value, 1, ttl)
}

// Get retrieves a value from cache
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.client.Get(key)
}

// Delete removes a value from cache
func (c *Cache) Delete(key string) {
	c.client.Del(key)
}

// Clear clears all cache entries
func (c *Cache) Clear() {
	c.client.Clear()
}

// Close closes the cache
func (c *Cache) Close() {
	c.client.Close()
}

// GetClient returns the underlying ristretto cache client
func (c *Cache) GetClient() *ristretto.Cache {
	return c.client
}

// Cache key constants
const (
	CacheKeyUsersList   = "users:list:%d:%d" // page:perPage
	CacheKeyUserStats   = "users:stats"
	CacheKeyUserByID    = "users:id:%s"
	CacheKeyUserByEmail = "users:email:%s"

	CacheKeyDomainsList       = "domains:list:%d:%d" // page:perPage
	CacheKeyDomainsStats      = "domains:stats"
	CacheKeyDomainByID        = "domains:id:%s"
	CacheKeySecurityMetrics   = "domains:security_metrics"
	CacheKeySLABreachAnalysis = "domains:sla_breach"

	NotificationsTotalKey = "notifications:total"
	NotificationsListKey  = "notifications:list"
	AlertsListKey         = "alerts:list"
)

// Cache TTL constants
const (
	CacheTTLUsersList = 5 * time.Minute
	CacheTTLUserStats = 5 * time.Minute
	CacheTTLUser      = 10 * time.Minute

	CacheTTLDomainsList       = 5 * time.Minute
	CacheTTLDomainsStats      = 5 * time.Minute
	CacheTTLDomain            = 10 * time.Minute
	CacheTTLSecurityMetrics   = 5 * time.Minute
	CacheTTLSLABreachAnalysis = 5 * time.Minute

	NotificationsTotalTTL = 5 * time.Minute
	NotificationsListTTL  = 5 * time.Minute
	AlertsListTTL         = 5 * time.Minute
)
