package cache

import (
	"context"
	"sync"
	"time"
)

// MetadataCache provides thread-safe caching for categories and platforms
// with automatic expiration (TTL).
type MetadataCache struct {
	mu         sync.RWMutex
	categories []string
	platforms  []string
	expiresAt  time.Time
	ttl        time.Duration
}

// NewMetadataCache creates a new cache with the specified TTL
func NewMetadataCache(ttl time.Duration) *MetadataCache {
	return &MetadataCache{
		ttl: ttl,
	}
}

// GetCategories returns cached categories if available and not expired,
// otherwise calls fetch() to get fresh data and updates the cache
func (c *MetadataCache) GetCategories(
	ctx context.Context,
	fetch func(context.Context) ([]string, error),
) ([]string, error) {
	c.mu.RLock()
	if time.Now().Before(c.expiresAt) && c.categories != nil {
		defer c.mu.RUnlock()
		return c.categories, nil
	}
	c.mu.RUnlock()

	// Cache miss or expired — fetch fresh data
	categories, err := fetch(ctx)
	if err != nil {
		return nil, err
	}

	// Update cache
	c.mu.Lock()
	c.categories = categories
	c.expiresAt = time.Now().Add(c.ttl)
	c.mu.Unlock()

	return categories, nil
}

// GetPlatforms returns cached platforms if available and not expired,
// otherwise calls fetch() to get fresh data and updates the cache
func (c *MetadataCache) GetPlatforms(
	ctx context.Context,
	fetch func(context.Context) ([]string, error),
) ([]string, error) {
	c.mu.RLock()
	if time.Now().Before(c.expiresAt) && c.platforms != nil {
		defer c.mu.RUnlock()
		return c.platforms, nil
	}
	c.mu.RUnlock()

	// Cache miss or expired — fetch fresh data
	platforms, err := fetch(ctx)
	if err != nil {
		return nil, err
	}

	// Update cache
	c.mu.Lock()
	c.platforms = platforms
	c.expiresAt = time.Now().Add(c.ttl)
	c.mu.Unlock()

	return platforms, nil
}

// Invalidate forces a cache refresh on next request
func (c *MetadataCache) Invalidate() {
	c.mu.Lock()
	c.expiresAt = time.Now()
	c.mu.Unlock()
}

// InvalidateAfter schedules automatic invalidation after duration
// Useful for periodic refreshes
func (c *MetadataCache) InvalidateAfter(duration time.Duration) {
	time.AfterFunc(duration, c.Invalidate)
}