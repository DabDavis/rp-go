package events

import (
	"fmt"
	"sync"
)

// TypedBus provides type-safe pub/sub communication between systems.
type TypedBus struct {
	mu          sync.RWMutex
	subscribers map[string][]func(any)
	queue       []any
}

// NewBus creates a new typed event bus.
func NewBus() *TypedBus {
	return &TypedBus{
		subscribers: make(map[string][]func(any)),
		queue:       []any{},
	}
}

// Subscribe registers a handler for an event type T.
func Subscribe[T any](b *TypedBus, handler func(T)) {
	var t T
	key := typeKey(t)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[key] = append(b.subscribers[key], func(e any) {
		if ev, ok := e.(T); ok {
			handler(ev)
		}
	})
}

// Publish immediately dispatches an event to all handlers.
func Publish[T any](b *TypedBus, event T) {
	key := typeKey(event)
	b.mu.RLock()
	handlers := b.subscribers[key]
	b.mu.RUnlock()
	for _, h := range handlers {
		h(event)
	}
}

// Queue schedules an event to run next frame (safe during updates).
func Queue[T any](b *TypedBus, event T) {
	b.mu.Lock()
	b.queue = append(b.queue, event)
	b.mu.Unlock()
}

// Flush processes all queued events in FIFO order.
func (b *TypedBus) Flush() {
	b.mu.Lock()
	q := b.queue
	b.queue = []any{}
	b.mu.Unlock()

	for _, e := range q {
		key := typeKey(e)
		b.mu.RLock()
		handlers := b.subscribers[key]
		b.mu.RUnlock()
		for _, h := range handlers {
			h(e)
		}
	}
}

// typeKey returns a unique key string for a given type.
func typeKey(v any) string {
	return fmt.Sprintf("%T", v)
}

