package main

import (
	"errors"
	"sync"
)

// Cache is a simple in-memory kv cache
// TODO: implement TTL
type Cache[V any] struct {
	data map[string]V
	mu   sync.Mutex
}

// New initializes a new in-memory kv Cache
// TODO: pass in TTL
func New[V any]() *Cache[V] {
	return &Cache[V]{
		data: make(map[string]V),
	}
}

// Get retrieves a key-value pair record matching key if exists
// returns nil, false if not found
func (c *Cache[V]) Get(key string) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.get(key)
}

func (c *Cache[V]) get(key string) (V, bool) {
	v, ok := c.data[key]
	return v, ok
}

// Set unconditionally adds key-value pair to cache, overwriting any existing records
func (c *Cache[V]) Set(key string, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

// Add only adds key-value pair record, if there is no existing record,
// returns error if existing record
func (c *Cache[V]) Add(key string, value V) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.get(key); !ok {
		return errors.New("cache: record already exists")
	}

	c.data[key] = value
	return nil
}

// Delete removes matching key-value record
func (c *Cache[V]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}
