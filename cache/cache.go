package cache

import (
	"sync"
	"time"
)

type Cache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
	ttl  time.Duration
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewCache initializes a new cache with a specified TTL
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		data: make(map[string]cacheItem),
		ttl:  ttl,
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists || time.Now().After(item.expiration) {
		return nil, false
	}
	return item.value, true
}

// Set stores an item in the cache
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
}
