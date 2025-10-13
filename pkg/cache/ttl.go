package cache

import (
	"sync"
	"time"
)

type TTLCache[T any] struct {
	items map[string]*item[T]
	mu    *sync.RWMutex
	ttl   time.Duration
	stop  chan struct{}
}

type item[T any] struct {
	value     T
	expiresAt time.Time
}

func NewTTLCache[T any](ttl, interval time.Duration) *TTLCache[T] {
	cache := &TTLCache[T]{
		items: make(map[string]*item[T]),
		mu:    &sync.RWMutex{},
		ttl:   ttl,
		stop:  make(chan struct{}),
	}

	go cache.cleanup(interval)
	return cache
}

func (c *TTLCache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &item[T]{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *TTLCache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok || time.Now().After(item.expiresAt) {
		var zero T
		return zero, false
	}

	return item.value, true
}

func (c *TTLCache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *TTLCache[T]) List() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	items := make([]T, 0, len(c.items))
	for _, item := range c.items {
		items = append(items, item.value)
	}
	return items
}

func (c *TTLCache[T]) cleanup(interval time.Duration) {
	ticker := time.Tick(interval)

	for {
		select {
		case <-ticker:
			c.mu.Lock()
			for k, v := range c.items {
				if time.Now().After(v.expiresAt) {
					delete(c.items, k)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			return
		}
	}
}

func (c *TTLCache[T]) Close() {
	close(c.stop)
}
