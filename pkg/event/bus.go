package event

import (
	"context"
	"log/slog"
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

func Subscribe[T any](eb *Bus, handler func(context.Context, T)) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	var zero T
	eventType := reflect.TypeOf(zero)
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

func Publish[T any](ctx context.Context, eb *Bus, event T) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eventType := reflect.TypeOf(event)
	for _, h := range eb.subscribers[eventType] {
		handler, _ := h.(func(context.Context, T))
		go func(handler func(context.Context, T)) {
			defer func() {
				if p := recover(); p != nil {
					slog.Error("Event handler panic", "event", event, "panic", p)
				}
			}()

			if ctx.Err() != nil {
				return
			}

			handler(ctx, event)
		}(handler)
	}
}
