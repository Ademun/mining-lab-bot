package event

import (
	"reflect"
	"sync"
)

type Bus struct {
	subscribers map[reflect.Type][]interface{}
	mu          sync.RWMutex
}

func NewEventBus() *Bus {
	return &Bus{
		subscribers: make(map[reflect.Type][]interface{}),
		mu:          sync.RWMutex{},
	}
}

func Subscribe[T any](eb *Bus, handler func(T)) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	var zero T
	eventType := reflect.TypeOf(zero)
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

func Publish[T any](eb *Bus, event T) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eventType := reflect.TypeOf(event)
	for _, h := range eb.subscribers[eventType] {
		handler := h.(func(T))
		handler(event)
	}
}
