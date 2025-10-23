package devconsole

import (
	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/systems/actor"
)

// System is the ECS component that integrates the console
// with the engine’s Update and Draw lifecycle.
type System struct {
	state *ConsoleState
}

// NewSystem initializes a new developer console bound to
// the shared actor registry.
func NewSystem(reg *actor.Registry) *System {
	cs := NewConsoleState(reg)
	return &System{state: cs}
}

// Layer ensures the console renders last in the overlay stack.
func (s *System) Layer() ecs.DrawLayer { return ecs.LayerConsole }

// Update handles input and command execution each frame.
func (s *System) Update(w *ecs.World) {
	if s.state == nil {
		return
	}
	s.state.UpdateInput(w)
}

// Draw renders the console overlay if active.
func (s *System) Draw(w *ecs.World, screen *platform.Image) {
	if s.state == nil {
		return
	}
	s.state.Render(w, screen)
}

