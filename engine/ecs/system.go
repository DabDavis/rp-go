package ecs

import "rp-go/engine/platform"

type System interface {
	Update(world *World)
}

type DrawableSystem interface {
	System
	Draw(world *World, screen *platform.Image)
}

// Optional: control execution order and render layer
type PrioritizedSystem interface {
	System
	Priority() int
}

// Optional: allow systems to specify which render layer they draw in.
type LayeredSystem interface {
	DrawableSystem
	Layer() DrawLayer
}

type DrawLayer int

const (
	LayerBackground DrawLayer = iota // parallax, starfields, distant scenery
	LayerWorld                       // primary world entities affected by the camera
	LayerForeground                  // world-space effects rendered after entities
	LayerHUD                         // overlay UI drawn in screen space
	LayerDebug                       // overlay diagnostics drawn in screen space
)
