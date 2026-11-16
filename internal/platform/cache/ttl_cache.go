package cache

import (
	"sync"
	"time"
)

type cacheItem[V any] struct {
	value     V
	expiresAt time.Time
}

// TTLCache provides a thread-safe, generic in-memory cache with time-to-live expiration.
type TTLCache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]cacheItem[V]
	ttl   time.Duration
}

// NewTTLCache creates an initialized cache instance with a default TTL duration.
func NewTTLCache[K comparable, V any](defaultTTL time.Duration) *TTLCache[K, V] {
	return &TTLCache[K, V]{
		items: make(map[K]cacheItem[V]),
		ttl:   defaultTTL,
	}
}

// Set stores a value with the default TTL.
func (c *TTLCache[K, V]) Set(key K, value V) {
	c.SetWithTTL(key, value, c.ttl)
}

// SetWithTTL stores a value with a specific TTL.
func (c *TTLCache[K, V]) SetWithTTL(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem[V]{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves a value if present and not expired.
func (c *TTLCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		var zero V
		return zero, false
	}

	if time.Now().After(item.expiresAt) {
		// Lazily clean expired entry
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()

		var zero V
		return zero, false
	}

	return item.value, true
}

// Delete removes an item from the cache.
func (c *TTLCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Size returns the count of active items in the cache.
func (c *TTLCache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
