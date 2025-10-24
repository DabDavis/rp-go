package aicomposer

import (
	"fmt"
	"sync"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/systems/ai"
	dataSys "rp-go/engine/systems/data"
)

/*───────────────────────────────────────────────*
 | SYSTEM STRUCTURE                              |
 *───────────────────────────────────────────────*/

// System automatically binds AIControllers to entities with Actor.AIRefs.
// It listens for data reloads and clears its internal cache when needed.
type System struct {
	data       *dataSys.System
	ai         *ai.System
	mu         sync.RWMutex
	processed  map[ecs.EntityID]bool // cache of already composed entities
	reloadFlag bool                  // true when ai.json is reloaded
}

/*───────────────────────────────────────────────*
 | CONSTRUCTOR                                   |
 *───────────────────────────────────────────────*/

// NewSystem links the AIComposer to the data + AI systems.
func NewSystem(data *dataSys.System, ai *ai.System) *System {
	if data == nil || ai == nil {
		fmt.Println("[AICOMPOSER] Warning: constructed with nil dependencies")
	}
	return &System{
		data:      data,
		ai:        ai,
		processed: make(map[ecs.EntityID]bool),
	}
}

/*───────────────────────────────────────────────*
 | ECS UPDATE LOOP                               |
 *───────────────────────────────────────────────*/

// Update scans the ECS world for new entities with Actor.AIRefs and
// attaches AIController components generated from the AI catalog.
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

	// Rebuild on data reload
	if s.reloadFlag {
		fmt.Println("[AICOMPOSER] Data reload detected — resetting AI bindings")
		s.processed = make(map[ecs.EntityID]bool)
		s.reloadFlag = false
	}

	manager.ForEach(func(e *ecs.Entity) {
		id := e.ID
		if s.processed[id] {
			return
		}

		actor, _ := e.Get("Actor").(*ecs.Actor)
		if actor == nil || len(actor.AIRefs) == 0 {
			s.processed[id] = true
			return
		}

		// Skip if already has a controller
		if _, ok := e.Get("AIController").(*ecs.AIController); ok {
			s.processed[id] = true
			return
		}

		// Build controller
		ctrl := s.ai.BuildControllerFromRefs(actor.AIRefs)
		if ctrl == nil {
			return
		}

		e.AddNamed("AIController", ctrl)
		s.processed[id] = true
		fmt.Printf("[AICOMPOSER] Bound %d AI actions to %q (entity %d)\n",
			len(ctrl.Actions), actor.ID, e.ID)
	})
}

/*───────────────────────────────────────────────*
 | DATA RELOAD HOOK                              |
 *───────────────────────────────────────────────*/

// OnDataReload ensures AIComposer responds to AI catalog changes.
func (s *System) OnDataReload(e events.DataReloaded) {
	if e.Type != "ai_catalog" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("[AICOMPOSER] AI catalog updated — rebuilding controller catalog")

	if s.data != nil && s.ai != nil {
		s.ai.OnDataReload(e, s.data.AICatalog)
	}
	s.reloadFlag = true
}

/*───────────────────────────────────────────────*
 | DRAW (NO-OP)                                  |
 *───────────────────────────────────────────────*/

func (s *System) Draw(_ *ecs.World, _ *ecs.World) {}

