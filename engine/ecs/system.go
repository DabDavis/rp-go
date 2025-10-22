package ecs

import "rp-go/engine/platform"

type System interface {
	Update(world *World)
	Draw(world *World, screen *platform.Image)
}

// OverlaySystem can render UI elements after the world has been
// composited to the final window surface.
type OverlaySystem interface {
	DrawOverlay(world *World, screen *platform.Image)
}

// Optional extension for ordering systems later
type PrioritizedSystem interface {
	System
	Priority() int
}
