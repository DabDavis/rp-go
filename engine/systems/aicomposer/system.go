package aicomposer

import (
	"fmt"
	"sync"

	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/systems/ai"
)

/*───────────────────────────────────────────────*
 | SYSTEM STRUCTURE                              |
 *───────────────────────────────────────────────*/

// System automatically binds AIControllers to entities with Actor.AIRefs.
// It also listens for AI catalog reloads to clear its cache.
type System struct {
	data     *data.System // global data (actors, AI catalog, etc.)
	ai       *ai.System   // AI logic and behavior execution
	mu       sync.RWMutex
	processed map[ecs.EntityID]bool // cache of entities already composed
	reloadFlag bool                  // set to true when ai.json is reloaded
}

/*───────────────────────────────────────────────*
 | CONSTRUCTOR                                   |
 *───────────────────────────────────────────────*/

// NewSystem initializes an AIComposer instance.
func NewSystem(dataSys *data.System, aiSys *ai.System) *System {
	return &System{
		data:      dataSys,
		ai:        aiSys,
		processed: make(map[ecs.EntityID]bool),
	}
}

/*───────────────────────────────────────────────*
 | ECS UPDATE LOOP                               |
 *───────────────────────────────────────────────*/

// Update scans entities and attaches AIControllers if missing.
// Uses caching to skip already processed entities.
func (s *System) Update(w *ecs.World) {
	if w == nil || s.data == nil || s.ai == nil {
		return
	}

	manager := w.EntitiesManager()
	if manager == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.reloadFlag {
		// AI catalog changed — rebind all entities
		fmt.Println("[AICOMPOSER] Reload detected — clearing cache for rebind")
		s.processed = make(map[ecs.EntityID]bool)
		s.reloadFlag = false
	}

	manager.ForEach(func(e *ecs.Entity) {
		id := e.ID

		// Skip if already processed
		if s.processed[id] {
			return
		}

		actor, _ := e.Get("Actor").(*ecs.Actor)
		if actor == nil || len(actor.AIRefs) == 0 {
			s.processed[id] = true // mark as known
			return
		}

		// Skip if AIController already exists
		if _, exists := e.Get("AIController").(*ai.AIController); exists {
			s.processed[id] = true
			return
		}

		ctrl := s.ai.BuildControllerFromRefs(actor.AIRefs)
		if ctrl == nil {
			return
		}

		e.AddNamed("AIController", ctrl)
		s.processed[id] = true
		fmt.Printf("[AICOMPOSER] Bound %d AI actions to %q (entity %d)\n", len(ctrl.Actions), actor.ID, e.ID)
	})
}

/*───────────────────────────────────────────────*
 | HOT RELOAD SUPPORT                            |
 *───────────────────────────────────────────────*/

// OnDataReload refreshes AI references when ai.json changes.
// This ensures that all controllers will be rebuilt next frame.
func (s *System) OnDataReload(e events.DataReloaded) {
	if e.Type != "ai_catalog" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil || s.ai == nil {
		return
	}

	fmt.Println("[AICOMPOSER] AI catalog updated — reinitializing controller catalog")

	// Rebuild AI lookup
	if s.data.AICatalog.Actions != nil {
		s.ai.OnDataReload(e, s.data.AICatalog)
	}

	// Mark for next-frame rebuild
	s.reloadFlag = true
}

/*───────────────────────────────────────────────*
 | DRAW (NO-OP)                                  |
 *───────────────────────────────────────────────*/

func (s *System) Draw(_ *ecs.World, _ *ai.System) {}

