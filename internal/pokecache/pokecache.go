package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	sync.Mutex
	cache map[string]CacheEntry
}

type CacheEntry struct {
	createdAt time.Time
	value     []byte
}

func (ps *Cache) Add(key string, val []byte) {
	ps.Lock()
	defer ps.Unlock()
	ps.cache[key] = CacheEntry{
		createdAt: time.Now(),
		value:     val,
	}
}

func (ps *Cache) Get(key string) ([]byte, bool) {

	if _, ok := ps.cache[key]; ok {
		ps.Lock()
		defer ps.Unlock()
		return ps.cache[key].value, true
	}
	return nil, false
}

func (ps *Cache) reapLoop(interval time.Duration) {

	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			ps.Lock()
			for k, entry := range ps.cache {
				if time.Since(entry.createdAt) >= interval {
					delete(ps.cache, k)
				}
			}
			ps.Unlock()
		}
	}()
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		cache: make(map[string]CacheEntry),
	}
	cache.reapLoop(interval)
	return cache
}
