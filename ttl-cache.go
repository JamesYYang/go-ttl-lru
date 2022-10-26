package tlcache

import (
	"log"
	"sort"
	"time"
)

type (
	ttlCache struct {
		maxEntries     int
		updateAgeOnGet bool
		cache          map[Key]*entry
		expiration     map[int64]Key
		expiry         time.Duration
	}
)

func newTTLCache(maxEntries int, expiry time.Duration, updateAgeOnGet bool) *ttlCache {
	return &ttlCache{
		maxEntries:     maxEntries,
		updateAgeOnGet: updateAgeOnGet,
		cache:          make(map[Key]*entry),
		expiration:     make(map[int64]Key),
		expiry:         expiry,
	}
}

func (c *ttlCache) add(key Key, value interface{}) {

	if c.cache == nil {
		c.cache = make(map[Key]*entry)
		c.expiration = make(map[int64]Key)
	}
	if ee, ok := c.cache[key]; ok {
		ee.ttl = time.Now().Add(c.expiry)
		ee.value = value

		exp := ee.ttl.UnixNano()
		c.expiration[exp] = key
		return
	}

	ele := &entry{
		ttl:   time.Now().Add(c.expiry),
		value: value,
		key:   key,
	}
	c.cache[key] = ele
	exp := ele.ttl.UnixNano()
	c.expiration[exp] = key

	if c.maxEntries != 0 && len(c.cache) > c.maxEntries {
		c.purgeToCapacity()
	}
}

func (c *ttlCache) size() int {
	if c.cache == nil {
		return 0
	}
	return len(c.cache)
}

// Get looks up a key's value from the cache.
func (c *ttlCache) get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		if time.Now().After(ele.ttl) {
			c.remove(key)
			return
		}
		if c.updateAgeOnGet {
			oldTTL := ele.ttl
			ele.ttl = time.Now().Add(c.expiry)
			delete(c.expiration, oldTTL.UnixNano())
			c.expiration[ele.ttl.UnixNano()] = key
		}
		return ele.value, hit
	}
	return
}

// Remove removes the provided key from the cache.
func (c *ttlCache) remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *ttlCache) purgeToCapacity() {
	log.Println("do purge")
	expKeys := make([]int64, 0, len(c.expiration))
	for k := range c.expiration {
		expKeys = append(expKeys, k)
	}
	sort.Slice(expKeys, func(i, j int) bool { return expKeys[i] < expKeys[j] })
	log.Println(expKeys)
	for k := range expKeys {
		c.remove(c.expiration[int64(k)])

		if len(c.cache) <= c.maxEntries {
			return
		}
	}
}

func (c *ttlCache) removeElement(e *entry) {
	delete(c.expiration, e.ttl.UnixNano())
	delete(c.cache, e.key)
}

// Clear purges all stored items from the cache.
func (c *ttlCache) clear() {
	c.cache = nil
	c.expiration = nil
}
