package ecs

import "github.com/hajimehoshi/ebiten/v2"

// Scene defines the required methods for game scenes.
type Scene interface {
	Name() string
	Init(world *World)
	Update(world *World)
	Draw(world *World, screen *ebiten.Image)
	Unload(world *World)
}

