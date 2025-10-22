package ecs

import "github.com/hajimehoshi/ebiten/v2"

type Position struct{ X, Y float64 }

func (p *Position) Name() string { return "Position" }

type Velocity struct{ VX, VY float64 }

func (v *Velocity) Name() string { return "Velocity" }

type Sprite struct {
	Image    *ebiten.Image
	Width    int
	Height   int
	Rotation float64
}

func (s *Sprite) Name() string { return "Sprite" }

type Camera struct {
	X, Y         float64
	Scale        float64
	Rotation     float64
	Target       *Entity
	TargetScale  float64
	MinScale     float64
	MaxScale     float64
	DefaultScale float64
}

func (c *Camera) Name() string { return "Camera" }

type CameraTarget struct{}

func (c *CameraTarget) Name() string { return "CameraTarget" }

// Actor defines metadata for players, NPCs, hostiles, and ships.
type Actor struct {
	ID         string // logical ID (e.g., "player", "npc_guard")
	Archetype  string // "player", "npc", "enemy", "ship"
	Persistent bool   // whether to keep across scene transitions
}

func (a *Actor) Name() string { return "Actor" }
