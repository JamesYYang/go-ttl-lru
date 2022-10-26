# go-ttl-lru
Golang LRU Cache with TTL.

This provide cache in three types:
- LRU
- LRU with TTL
- TTL

## Install

```bash
go get -u github.com/JamesYYang/go-ttl-lru
```

# Example

## Use LRU Cache

```go
import tlcache "github.com/JamesYYang/go-ttl-lru"

cache := tlcache.NewLRUCache(5)
cache.Add(1, "this is test 1")

if v, ok := cache.Get(1); ok {
    log.Printf("get key (%d) success, value: %s", 1, v)
}
```

## Use TTL LRU Cache

```go
import tlcache "github.com/JamesYYang/go-ttl-lru"

cache := tlcache.NewLRUWithTTLCache(5, 5*time.Second)
cache.Add(1, "this is test 1")

if v, ok := cache.Get(1); ok {
    log.Printf("get key (%d) success, value: %s", 1, v)
}
time.Sleep(5 * time.Second)
if v, ok := cache.Get(1); !ok {
    log.Printf("get key (%d) failed", 1)
}
```

## Use TTL Cache

```go
import tlcache "github.com/JamesYYang/go-ttl-lru"

cache := tlcache.NewTTLCache(5, 5*time.Second, true)
cache.Add(1, "this is test 1")

time.Sleep(3 * time.Second)
if v, ok := cache.Get(1); ok {
    log.Printf("get key (%d) success, value: %s", 1, v)
}
time.Sleep(3 * time.Second)
if v, ok := cache.Get(1); ok {
    log.Printf("get key (%d) success, value: %s", 1, v)
}
```