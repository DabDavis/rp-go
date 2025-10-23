package data

import (
	"sync"

	"rp-go/engine/events"
	"rp-go/engine/ecs"
)

// Subscriber is a function called when a data file reloads.
type Subscriber func(event events.DataReloaded)

// DataSubscriber manages subscriptions to DataReloaded events.
// It safely dispatches reload notifications to dependent systems.
type DataSubscriber struct {
	mu          sync.RWMutex
	subscribers map[string][]Subscriber // keyed by data type ("render_config", "actor_db", etc.)
}

// NewDataSubscriber creates a new event subscriber manager.
func NewDataSubscriber() *DataSubscriber {
	return &DataSubscriber{
		subscribers: make(map[string][]Subscriber),
	}
}

// Register binds a reload handler for a given data type (e.g. "actor_db").
func (d *DataSubscriber) Register(dataType string, fn Subscriber) {
	if fn == nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.subscribers[dataType] = append(d.subscribers[dataType], fn)
}

// Notify dispatches reload events to all registered handlers for the type.
func (d *DataSubscriber) Notify(e events.DataReloaded) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if subs, ok := d.subscribers[e.Type]; ok {
		for _, fn := range subs {
			go fn(e) // dispatch asynchronously to avoid blocking Update loop
		}
	}
}

// BindToWorld connects the subscriber to the ECS world event bus.
func (d *DataSubscriber) BindToWorld(world *ecs.World) {
	if world == nil || world.EventBus == nil {
		return
	}
	if bus, ok := world.EventBus.(*events.TypedBus); ok {
		events.Subscribe(bus, func(e events.DataReloaded) {
			d.Notify(e)
		})
	}
}

// Example integration:
//
// func (s *AISystem) OnReload(e events.DataReloaded) {
//     if e.Type == "actor_db" {
//         s.ReloadBehaviorTemplates()
//     }
// }
//
// func init() {
//     dataSubscriber.Register("actor_db", aiSystem.OnReload)
// }

