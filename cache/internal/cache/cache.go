package cache

import (
	"sync"
)

// CacheService defines the interface for interacting with the in-memory cache
type CacheService interface {
	Get(string) ([]byte, bool)
	Set(string, []byte)
	Delete(string)
}

// Cache is a simple in-memory kv cache
// TODO: implement TTL
// TODO: Refactor to just use a sync.Map
type Cache struct {
	mu   sync.Mutex
	data map[string][]byte
}

// New initializes a new in-memory kv Cache
// TODO: pass in TTL
func NewCache() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

// Get retrieves a key-value pair record matching key if exists
// returns nil, false if not found
// TODO: separate out Get and GetOK
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.data[key]
	return v, ok
}

// Set unconditionally adds key-value pair to cache, overwriting any existing records
// TODO: create Add which only adds k-v pair if doesn't exist
func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

// Delete removes matching key-value record
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}
