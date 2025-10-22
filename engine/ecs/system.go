package ecs

import "github.com/hajimehoshi/ebiten/v2"

type System interface {
	Update(world *World)
	Draw(world *World, screen *ebiten.Image)
}

// Optional extension for ordering systems later
type PrioritizedSystem interface {
	System
	Priority() int
}

