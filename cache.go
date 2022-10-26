package tlcache

import (
	"sync"
	"time"
)

type (
	Key interface{}

	entry struct {
		key   Key
		ttl   time.Time
		value interface{}
	}

	CacheCore interface {
		add(key Key, value interface{})
		size() int
		get(key Key) (value interface{}, ok bool)
		remove(key Key)
		clear()
	}

	Cache struct {
		cache CacheCore
		sync.Mutex
	}
)

const (
	maxDuration time.Duration = 1<<63 - 1
)

func NewLRUCache(maxEntries int) *Cache {
	return NewLRUWithTTLCache(maxEntries, maxDuration)
}

func NewLRUWithTTLCache(maxEntries int, expiry time.Duration) *Cache {
	c := &Cache{
		cache: newLRU(maxEntries, expiry),
	}
	return c
}

func NewTTLCache(maxEntries int, expiry time.Duration, updateAgeOnGet bool) *Cache {
	c := &Cache{
		cache: newTTLCache(maxEntries, expiry, updateAgeOnGet),
	}
	return c
}

func (c *Cache) Add(key Key, value interface{}) {
	c.Lock()
	defer c.Unlock()
	c.cache.add(key, value)
}

func (c *Cache) Size() int {
	c.Lock()
	defer c.Unlock()
	return c.cache.size()
}

func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	c.Lock()
	defer c.Unlock()
	return c.cache.get(key)
}

func (c *Cache) Remove(key Key) {
	c.Lock()
	defer c.Unlock()
	c.cache.remove(key)
}

func (c *Cache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.cache.clear()
}
