package actor

import (
	"fmt"
	"sync"

	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/events"
)

/*───────────────────────────────────────────────*
 | ACTOR SYSTEM                                  |
 *───────────────────────────────────────────────*/

// System maintains and updates the global Actor registry each frame.
// It synchronizes ECS entities with their Actor components and listens
// for data reload events to refresh template-based logic.
type System struct {
	registry    *Registry
	templates   data.ActorDatabase // cached actor definitions
	mu          sync.RWMutex
	initialized bool
}

/*───────────────────────────────────────────────*
 | LIFECYCLE                                     |
 *───────────────────────────────────────────────*/

// NewSystem constructs a fresh Actor system with an empty registry.
func NewSystem() *System {
	return &System{
		registry: NewRegistry(),
	}
}

// Registry exposes the actor registry for other systems to query.
func (s *System) Registry() *Registry {
	if s.registry == nil {
		s.registry = NewRegistry()
	}
	return s.registry
}

// Update rebuilds the actor registry from the ECS world every frame.
// Ensures that only one PlayerInput is active and avoids AI/input conflicts.
func (s *System) Update(w *ecs.World) {
	if s.registry == nil {
		s.registry = NewRegistry()
	}
	s.registry.Reset()

	primaryAssigned := false

	manager := w.EntitiesManager()
	if manager == nil {
		return
	}

	manager.ForEach(func(e *ecs.Entity) {
		actor, ok := e.Get("Actor").(*ecs.Actor)
		if !ok || actor == nil {
			return
		}

		// Index this actor
		s.registry.Add(actor, e)

		// Enforce PlayerInput exclusivity and AI coordination
		if controller, ok := e.Get("PlayerInput").(*ecs.PlayerInput); ok && controller != nil {
			if _, hasAI := e.Get("AIController").(*ecs.AIController); hasAI {
				controller.Enabled = false
				return
			}
			if !primaryAssigned {
				controller.Enabled = true
				primaryAssigned = true
			} else {
				controller.Enabled = false
			}
		}
	})

	if !s.initialized {
		fmt.Printf("[ACTOR] Registry initialized with %d entities\n", len(s.registry.all))
		s.initialized = true
	}
}

/*───────────────────────────────────────────────*
 | DATA RELOAD HOOK                              |
 *───────────────────────────────────────────────*/

// OnDataReload reacts to live data updates (e.g. actors.json changed).
// This is registered automatically via data.System.Subscriber().
func (s *System) OnDataReload(e events.DataReloaded, newDB data.ActorDatabase) {
	if e.Type != "actor_db" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.templates = newDB
	fmt.Printf("[ACTOR] Reloaded actor templates (%d entries)\n", len(newDB.Actors))
	s.RefreshActors()
}

/*───────────────────────────────────────────────*
 | REFRESH ACTORS                                |
 *───────────────────────────────────────────────*/

// RefreshActors re-applies updated template data to live entities.
func (s *System) RefreshActors() {
	if s.registry == nil || len(s.templates.Actors) == 0 {
		return
	}

	templateMap := make(map[string]data.ActorTemplate)
	for _, tpl := range s.templates.Actors {
		templateMap[tpl.Name] = tpl
	}

	for _, e := range s.registry.Entities() {
		act, ok := e.Get("Actor").(*ecs.Actor)
		if !ok || act == nil {
			continue
		}

		// Use actor.ID as template name if no dedicated Template field
		templateName := act.ID
		if tpl, ok := templateMap[templateName]; ok {
			act.Archetype = tpl.Archetype

			// Update Velocity if component exists
			if vel, ok := e.Get("Velocity").(*ecs.Velocity); ok {
				if tpl.Velocity != nil {
					vel.VX = tpl.Velocity.VX
					vel.VY = tpl.Velocity.VY
				}
			}
		}
	}

	fmt.Printf("[ACTOR] Refreshed %d live actors from templates\n", len(s.registry.all))
}

/*───────────────────────────────────────────────*
 | HELPERS                                       |
 *───────────────────────────────────────────────*/

// LoadTemplates allows manual injection of actor DB (used by core init).
func (s *System) LoadTemplates(db data.ActorDatabase) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.templates = db
}

// Templates returns a safe copy of the current actor database.
func (s *System) Templates() data.ActorDatabase {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.templates
}

