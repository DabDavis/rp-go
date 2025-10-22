package ecs

import "rp-go/engine/platform"

type System interface {
	Update(world *World)
	Draw(world *World, screen *platform.Image)
}

// Optional: control execution order and render layer
type PrioritizedSystem interface {
	System
	Priority() int
}

// Optional: indicate whether the system should draw in world or overlay layer
type LayeredSystem interface {
	System
	Layer() DrawLayer
}

type DrawLayer int

const (
	LayerBackground DrawLayer = iota // parallax, backgrounds
	LayerWorld // affected by camera transforms
	LayerOverlay                 // drawn directly to final screen (HUD, debug)
)

