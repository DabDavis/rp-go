package devconsole

import (
	"rp-go/engine/ecs"
	"rp-go/engine/systems/actor"
)

// Config controls layout of the console window.
type Config struct {
	Margin         int
	ViewportWidth  int
	ViewportHeight int
	HeightRatio    float64
	MinWidth       int
	MinHeight      int
}

// normalize fills zero-valued fields with defaults.
func (c *Config) normalize() {
	if c.Margin <= 0 {
		c.Margin = 16
	}
	if c.ViewportWidth <= 0 {
		c.ViewportWidth = 640
	}
	if c.ViewportHeight <= 0 {
		c.ViewportHeight = 360
	}
	if c.HeightRatio <= 0 {
		c.HeightRatio = 0.33
	}
	if c.MinWidth <= 0 {
		c.MinWidth = 360
	}
	if c.MinHeight <= 0 {
		c.MinHeight = 180
	}
}

// System integrates the developer console into the ECS lifecycle.
type System struct {
	state *ConsoleState
	cfg   Config
}

// NewSystem initializes a new developer console bound to the shared actor registry.
func NewSystem(reg *actor.Registry, cfg Config) *System {
	cfg.normalize()
	cs := NewConsoleState(reg)
	return &System{state: cs, cfg: cfg}
}

// Update handles input, window binding, and layout.
func (s *System) Update(w *ecs.World) {
	if s == nil || s.state == nil {
		return
	}
	s.state.ensureWindow(w, s.cfg)
	s.state.UpdateInput(w)
	s.state.syncWindowVisibility()
	s.state.applyLayout(s.cfg)
}
