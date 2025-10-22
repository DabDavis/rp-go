package ecs

import "rp-go/engine/platform"

type System interface {
	Update(world *World)
	Draw(world *World, screen *platform.Image)
}

// Optional extension for ordering systems later
type PrioritizedSystem interface {
	System
	Priority() int
}
