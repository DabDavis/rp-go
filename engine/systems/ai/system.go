package ai

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"
)

type System struct {
	mu       sync.RWMutex
	rng      *rand.Rand
	catalog  *AIActionCatalogLookup
	lastLoad time.Time
}

// NewSystem constructs an AI system and initializes behavior catalog.
func NewSystem(cat data.AIActionCatalog) *System {
	sys := &System{
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
		catalog: NewCatalogLookup(cat),
	}
	RegisterDefaultBehaviors(sys)
	return sys
}

// OnDataReload updates AI catalog on file reload.
func (s *System) OnDataReload(e events.DataReloaded, cat data.AIActionCatalog) {
	if e.Type != "ai_catalog" {
		return
	}
	s.mu.Lock()
	s.catalog = NewCatalogLookup(cat)
	s.lastLoad = time.Now()
	s.mu.Unlock()
	fmt.Printf("[AI] Reloaded %d actions from ai.json\n", len(cat.Actions))
}

// BuildControllerFromRefs instantiates a controller from action names.
func (s *System) BuildControllerFromRefs(refs []string) *AIController {
	if len(refs) == 0 {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	ctrl := &AIController{Active: true}
	for _, name := range refs {
		if tpl, ok := s.catalog.Get(name); ok {
			ctrl.Actions = append(ctrl.Actions, AIActionInstance{
				Name:     tpl.Name,
				Type:     tpl.Type,
				Priority: tpl.Priority,
				Params:   tpl.Params,
			})
		}
	}
	sort.SliceStable(ctrl.Actions, func(i, j int) bool {
		return ctrl.Actions[i].Priority < ctrl.Actions[j].Priority
	})
	return ctrl
}

func (s *System) Update(w *ecs.World) {
	if w == nil {
		return
	}
	s.ensureRNG()

	manager := w.EntitiesManager()
	if manager == nil {
		return
	}

	manager.ForEach(func(e *ecs.Entity) {
		ctrl, ok := e.Get("AIController").(*AIController)
		if !ok || ctrl == nil || !ctrl.Active {
			return
		}

		pos, _ := e.Get("Position").(*ecs.Position)
		vel, _ := e.Get("Velocity").(*ecs.Velocity)
		if pos == nil || vel == nil {
			return
		}

		vel.VX, vel.VY = 0, 0
		for _, act := range ctrl.Actions {
			if s.executeAction(w, e, pos, vel, act) {
				break
			}
		}
	})
}

func (s *System) Draw(*ecs.World, *platform.Image) {}

