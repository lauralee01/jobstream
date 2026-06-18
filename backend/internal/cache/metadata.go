package cache

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// MetadataCache provides thread-safe caching for categories and platforms
// with separate expiration times and duplicate-request suppression.
type MetadataCache struct {
	mu sync.RWMutex

	categories          []string
	categoriesExpiresAt time.Time

	platforms          []string
	platformsExpiresAt time.Time

	ttl time.Duration

	group singleflight.Group
}

// NewMetadataCache creates a new cache with the specified TTL.
func NewMetadataCache(ttl time.Duration) *MetadataCache {
	return &MetadataCache{
		ttl: ttl,
	}
}

func (c *MetadataCache) GetCategories(
	ctx context.Context,
	fetch func(context.Context) ([]string, error),
) ([]string, error) {
	c.mu.RLock()
	if time.Now().Before(c.categoriesExpiresAt) && c.categories != nil {
		categories := c.categories
		c.mu.RUnlock()
		return categories, nil
	}
	c.mu.RUnlock()

	result, err, _ := c.group.Do("categories", func() (interface{}, error) {
		c.mu.RLock()
		if time.Now().Before(c.categoriesExpiresAt) && c.categories != nil {
			categories := c.categories
			c.mu.RUnlock()
			return categories, nil
		}
		c.mu.RUnlock()

		categories, err := fetch(ctx)
		if err != nil {
			return nil, err
		}

		c.mu.Lock()
		c.categories = categories
		c.categoriesExpiresAt = time.Now().Add(c.ttl)
		c.mu.Unlock()

		return categories, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]string), nil
}

func (c *MetadataCache) GetPlatforms(
	ctx context.Context,
	fetch func(context.Context) ([]string, error),
) ([]string, error) {
	c.mu.RLock()
	if time.Now().Before(c.platformsExpiresAt) && c.platforms != nil {
		platforms := c.platforms
		c.mu.RUnlock()
		return platforms, nil
	}
	c.mu.RUnlock()

	result, err, _ := c.group.Do("platforms", func() (interface{}, error) {
		c.mu.RLock()
		if time.Now().Before(c.platformsExpiresAt) && c.platforms != nil {
			platforms := c.platforms
			c.mu.RUnlock()
			return platforms, nil
		}
		c.mu.RUnlock()

		platforms, err := fetch(ctx)
		if err != nil {
			return nil, err
		}

		c.mu.Lock()
		c.platforms = platforms
		c.platformsExpiresAt = time.Now().Add(c.ttl)
		c.mu.Unlock()

		return platforms, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]string), nil
}

func (c *MetadataCache) Invalidate() {
	c.mu.Lock()
	c.categoriesExpiresAt = time.Now()
	c.platformsExpiresAt = time.Now()
	c.mu.Unlock()
}

func (c *MetadataCache) InvalidateCategories() {
	c.mu.Lock()
	c.categoriesExpiresAt = time.Now()
	c.mu.Unlock()
}

func (c *MetadataCache) InvalidatePlatforms() {
	c.mu.Lock()
	c.platformsExpiresAt = time.Now()
	c.mu.Unlock()
}

func (c *MetadataCache) InvalidateAfter(duration time.Duration) {
	time.AfterFunc(duration, c.Invalidate)
}
