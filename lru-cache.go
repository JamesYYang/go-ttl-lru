package tlcache

import (
	"container/list"
	"time"
)

type (
	lruCache struct {
		maxEntries int
		ll         *list.List
		cache      map[Key]*list.Element
		expiry     time.Duration
	}
)

func newLRU(maxEntries int, expiry time.Duration) *lruCache {
	return &lruCache{
		maxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[Key]*list.Element),
		expiry:     expiry,
	}
}

func (c *lruCache) add(key Key, value interface{}) {

	if c.cache == nil {
		c.cache = make(map[Key]*list.Element)
		c.ll = list.New()
	}
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).ttl = time.Now().Add(c.expiry)
		ee.Value.(*entry).value = value
		return
	}
	ele := c.ll.PushFront(&entry{key, time.Now().Add(c.expiry), value})
	c.cache[key] = ele
	if c.maxEntries != 0 && c.ll.Len() > c.maxEntries {
		c.removeOldest()
	}
}

func (c *lruCache) size() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Get looks up a key's value from the cache.
func (c *lruCache) get(key Key) (value interface{}, ok bool) {

	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		if time.Now().After(ele.Value.(*entry).ttl) {
			c.remove(key)
			return
		}
		return ele.Value.(*entry).value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *lruCache) remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *lruCache) removeOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *lruCache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
}

// Clear purges all stored items from the cache.
func (c *lruCache) clear() {
	c.ll = nil
	c.cache = nil
}
