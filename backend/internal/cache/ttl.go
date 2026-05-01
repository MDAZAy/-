package cache

import (
	"sync"
	"time"
)

type item struct {
	value     any
	expiresAt time.Time
}

type TTLCache struct {
	mu    sync.RWMutex
	items map[string]item
}

func NewTTLCache() *TTLCache {
	return &TTLCache{items: make(map[string]item)}
}

func (c *TTLCache) Get(key string) (any, bool) {
	c.mu.RLock()
	entry, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(entry.expiresAt) {
		if ok {
			c.Delete(key)
		}
		return nil, false
	}
	return entry.value, true
}

func (c *TTLCache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	c.items[key] = item{value: value, expiresAt: time.Now().Add(ttl)}
	c.mu.Unlock()
}

func (c *TTLCache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *TTLCache) DeletePrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key := range c.items {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(c.items, key)
		}
	}
}
