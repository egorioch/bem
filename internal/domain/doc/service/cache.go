package service

import (
	"sync"
	"time"
)

type Cache struct {
	data map[string]CacheItem
	mu   sync.RWMutex
}

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheItem),
	}
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(duration).Unix()
	c.data[key] = CacheItem{Value: value, Expiration: expiration}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.data[key]
	if !found || item.Expiration < time.Now().Unix() {
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
