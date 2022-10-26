package tlcache

import (
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
		ee.value = value
		c.changeTTL(key)
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

func (c *ttlCache) changeTTL(key Key) {
	if ee, ok := c.cache[key]; ok {
		delete(c.expiration, ee.ttl.UnixNano())
		ee.ttl = time.Now().Add(c.expiry)
		exp := ee.ttl.UnixNano()
		c.expiration[exp] = key
	}
}

func (c *ttlCache) size() int {
	if c.cache == nil {
		return 0
	}
	return len(c.cache)
}

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
			c.changeTTL(key)
		}
		return ele.value, hit
	}
	return
}

func (c *ttlCache) remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

func (c *ttlCache) purgeToCapacity() {
	expKeys := make([]int64, 0, len(c.expiration))
	for k := range c.expiration {
		expKeys = append(expKeys, k)
	}
	sort.Slice(expKeys, func(i, j int) bool { return expKeys[i] < expKeys[j] })
	for _, k := range expKeys {
		c.remove(c.expiration[k])

		if len(c.cache) <= c.maxEntries {
			return
		}
	}
}

func (c *ttlCache) removeElement(e *entry) {
	delete(c.expiration, e.ttl.UnixNano())
	delete(c.cache, e.key)
}

func (c *ttlCache) clear() {
	c.cache = nil
	c.expiration = nil
}
