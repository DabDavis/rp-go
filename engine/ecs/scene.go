package ecs

import "rp-go/engine/platform"

// Scene defines the required methods for game scenes.
type Scene interface {
	Name() string
	Init(world *World)
	Update(world *World)
	Draw(world *World, screen *platform.Image)
	Unload(world *World)
}
