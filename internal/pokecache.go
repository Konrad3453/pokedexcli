package pokecache

import (
	"sync"
	"time"
)


type CacheEntry struct {
    createdAt time.Time
	val []byte
}


type Cache struct {
	mu sync.Mutex
	store map[string]CacheEntry
}


func NewCache() *Cache {
	return &Cache{
		store: make(map[string]CacheEntry),
	}
}
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = CacheEntry{
		createdAt: time.Now(),
		val: val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.store[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() { // remove expired entries
	for {
		time.Sleep(time.Minute)
		c.reap()
	}
}
