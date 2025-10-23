package debug

import (
	"rp-go/engine/ecs"
	"rp-go/engine/events"

	"rp-go/engine/systems/aicomposer"
)

/*───────────────────────────────────────────────*
 | SYSTEM STRUCTURE                              |
 *───────────────────────────────────────────────*/

// System controls all debug overlay windows (stats, entity lists, systems, etc.)
// and now includes an AI Composer inspector window.
type System struct {
	cfg Config
	bus *events.TypedBus

	enabled bool

	stats     *StatsWindow
	entities  *EntitiesWindow
	systems   *SystemWindow
	toolbar   *ToolbarWindow
	composer  *aicomposer.DebugWindow // ✅ new AI Composer debug window
	composerS *aicomposer.System      // reference for live updates
}

/*───────────────────────────────────────────────*
 | INITIALIZATION                                |
 *───────────────────────────────────────────────*/

// NewSystem initializes the debug overlay controller.
func NewSystem(cfg Config) *System {
	cfg.normalize()
	return &System{
		cfg:     cfg,
		enabled: true,
	}
}

/*───────────────────────────────────────────────*
 | UPDATE LOOP                                   |
 *───────────────────────────────────────────────*/

// Update runs every frame to manage debug visibility and window content.
func (s *System) Update(world *ecs.World) {
	if world == nil {
		return
	}

	// Bind event bus if not already done
	if s.bus == nil {
		s.bus, _ = world.EventBus.(*events.TypedBus)
		if s.bus != nil {
			s.initSubscriptions()
		}
	}

	// Skip content updates if disabled
	if !s.enabled {
		return
	}

	// Lazy creation of all windows
	s.ensureWindows(world)

	// Refresh content per frame
	if s.stats != nil {
		s.stats.Update(world)
	}
	if s.entities != nil {
		s.entities.Update(world)
	}
	if s.systems != nil {
		s.systems.Update(world)
	}
	if s.composer != nil && s.composerS != nil {
		s.composer.Update(world, s.composerS)
	}
}

/*───────────────────────────────────────────────*
 | WINDOW MANAGEMENT                             |
 *───────────────────────────────────────────────*/

func (s *System) ensureWindows(world *ecs.World) {
	if s.stats == nil {
		s.stats = NewStatsWindow(s.cfg)
	}
	if s.entities == nil {
		s.entities = NewEntitiesWindow(s.cfg)
	}
	if s.systems == nil {
		s.systems = NewSystemWindow(s.cfg)
	}
	if s.toolbar == nil {
		s.toolbar = NewToolbarWindow(s.cfg, s.bus)
	}
	if s.composer == nil {
		s.composer = aicomposer.NewDebugWindow()
	}

	s.stats.Ensure(world)
	s.entities.Ensure(world)
	s.systems.Ensure(world)
	s.toolbar.Ensure(world)
	s.composer.Ensure(world)
}

// AttachComposer binds the active AIComposerSystem for live debugging.
func (s *System) AttachComposer(sys *aicomposer.System) {
	s.composerS = sys
}

/*───────────────────────────────────────────────*
 | EVENT SUBSCRIPTIONS                           |
 *───────────────────────────────────────────────*/

func (s *System) initSubscriptions() {
	if s.bus == nil {
		return
	}

	// Global toggle of debug overlay visibility
	events.Subscribe(s.bus, func(e events.DebugToggleEvent) {
		s.enabled = e.Enabled
	})

	// Hide relevant windows when closed manually
	events.Subscribe(s.bus, func(e events.WindowClosedEvent) {
		switch e.ID {
		case "debug.stats", "debug.entities", "debug.systems", "debug.aicomposer":
			s.enabled = false
		}
	})

	// Optional: allow specific keybinding toggles for AI Composer
	events.Subscribe(s.bus, func(e events.DebugKeyToggleEvent) {
		if e.Key == "F8" { // Example: F8 toggles AI composer window
			if s.composer != nil {
				s.composer.Hide()
			}
		}
	})
}

/*───────────────────────────────────────────────*
 | TOGGLING & VISIBILITY                         |
 *───────────────────────────────────────────────*/

func (s *System) Toggle(world *ecs.World) {
	s.enabled = !s.enabled

	if s.bus != nil {
		events.Queue(s.bus, events.DebugToggleEvent{Enabled: s.enabled})
	}

	if !s.enabled {
		s.hideAll(world)
	}
}

func (s *System) hideAll(world *ecs.World) {
	if s.stats != nil {
		s.stats.Hide(world)
	}
	if s.entities != nil {
		s.entities.Hide(world)
	}
	if s.systems != nil {
		s.systems.Hide(world)
	}
	if s.composer != nil {
		s.composer.Hide()
	}
}

