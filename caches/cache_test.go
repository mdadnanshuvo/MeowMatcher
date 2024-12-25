package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	cache := NewCache(5 * time.Second)
	assert.NotNil(t, cache)
	assert.NotNil(t, cache.data)
	assert.Equal(t, 5*time.Second, cache.ttl)
}

func TestCacheSetAndGet(t *testing.T) {
	cache := NewCache(2 * time.Second)

	// Test setting a value
	cache.Set("key1", "value1")
	val, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value1", val)

	// Test overwriting a value
	cache.Set("key1", "newValue1")
	val, found = cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "newValue1", val)
}

func TestCacheExpiration(t *testing.T) {
	cache := NewCache(1 * time.Second)

	cache.Set("key1", "value1")
	val, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value1", val)

	// Wait for the TTL to expire
	time.Sleep(2 * time.Second)
	val, found = cache.Get("key1")
	assert.False(t, found)
	assert.Nil(t, val)
}

func TestCacheNonExistentKey(t *testing.T) {
	cache := NewCache(2 * time.Second)

	val, found := cache.Get("nonexistent")
	assert.False(t, found)
	assert.Nil(t, val)
}

func TestCacheConcurrentAccess(t *testing.T) {
	cache := NewCache(2 * time.Second)
	var wg sync.WaitGroup

	// Writer goroutine
	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set("key", i)
		}
		wg.Done()
	}()

	// Reader goroutine
	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			_, _ = cache.Get("key")
		}
		wg.Done()
	}()

	wg.Wait()

	// Ensure the cache is not corrupted
	val, found := cache.Get("key")
	assert.True(t, found)
	assert.NotNil(t, val)
}

func TestCacheOverwriteWithNewTTL(t *testing.T) {
	cache := NewCache(1 * time.Second)

	cache.Set("key1", "value1")
	cache.Set("key1", "value2")

	// Ensure the value and TTL are updated
	val, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value2", val)

	// Check expiration behavior
	time.Sleep(2 * time.Second)
	val, found = cache.Get("key1")
	assert.False(t, found)
	assert.Nil(t, val)
}
