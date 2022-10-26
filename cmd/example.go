package main

import (
	"log"
	"time"

	tlcache "github.com/JamesYYang/go-ttl-lru"
)

func main() {

	// lruExample()

	// lruTTLExample()

	ttlExample()
}

func lruExample() {
	cache := tlcache.NewLRUCache(5)

	cache.Add(1, "this is test 1")
	cache.Add(2, "this is test 2")
	cache.Add(3, "this is test 3")
	cache.Add(4, "this is test 4")
	cache.Add(5, "this is test 5")
	cache.Add(6, "this is test 6")

	log.Println("init LRU cache finish")
	log.Printf("cache size: %d", cache.Size())

	if v, ok := cache.Get(2); ok {
		log.Printf("get key (%d) success, value: %s", 2, v)
	}

	if v, ok := cache.Get(3); ok {
		log.Printf("get key (%d) success, value: %s", 3, v)
	}

	cache.Add(7, "this is test 7")
	cache.Add(8, "this is test 8")

	if _, ok := cache.Get(4); !ok {
		log.Printf("get key (%d) failed", 4)
	}
}

func lruTTLExample() {
	cache := tlcache.NewLRUWithTTLCache(5, 5*time.Second)

	cache.Add(1, "this is test 1")
	cache.Add(2, "this is test 2")
	cache.Add(3, "this is test 3")
	cache.Add(4, "this is test 4")
	cache.Add(5, "this is test 5")
	cache.Add(6, "this is test 6")

	log.Println("init LRU TTL cache finish")
	log.Printf("cache size: %d", cache.Size())

	if v, ok := cache.Get(2); ok {
		log.Printf("get key (%d) success, value: %s", 2, v)
	}

	time.Sleep(5 * time.Second)

	if _, ok := cache.Get(2); !ok {
		log.Printf("get key (%d) failed", 2)
	}
}

func ttlExample() {
	cache := tlcache.NewTTLCache(5, 5*time.Second, true)

	cache.Add(1, "this is test 1")
	time.Sleep(3 * time.Millisecond)
	cache.Add(2, "this is test 2")
	time.Sleep(3 * time.Millisecond)
	cache.Add(3, "this is test 3")
	time.Sleep(3 * time.Millisecond)
	cache.Add(4, "this is test 4")
	time.Sleep(3 * time.Millisecond)
	cache.Add(5, "this is test 5")
	time.Sleep(3 * time.Millisecond)
	cache.Add(6, "this is test 6")
	time.Sleep(3 * time.Millisecond)

	log.Println("init TTL cache finish")
	log.Printf("cache size: %d", cache.Size())

	time.Sleep(3 * time.Second)

	if v, ok := cache.Get(2); ok {
		log.Printf("get key (%d) success, value: %s", 2, v)
	}

	time.Sleep(3 * time.Second)

	if v, ok := cache.Get(2); ok {
		log.Printf("get key (%d) success, value: %s", 2, v)
	}

	if _, ok := cache.Get(3); !ok {
		log.Printf("get key (%d) failed", 3)
	}
}
