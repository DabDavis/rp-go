package ecs

import "rp-go/engine/platform"

// Position represents a 2D coordinate in world space.
type Position struct{ X, Y float64 }

func (p *Position) Name() string { return "Position" }

// Velocity represents a rate of change in position.
type Velocity struct{ VX, VY float64 }

func (v *Velocity) Name() string { return "Velocity" }

// Sprite is the renderable visual attached to an entity.
// It stores a texture reference and visual transform data.
type Sprite struct {
	Image    *platform.Image // pointer to GPU texture
	Width    int
	Height   int
	Rotation float64
}

func (s *Sprite) Name() string { return "Sprite" }

// Camera defines the view transform for rendering the world.
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

// CameraTarget marks an entity as being tracked by the active camera.
type CameraTarget struct{}

func (c *CameraTarget) Name() string { return "CameraTarget" }

// Actor defines metadata for players, NPCs, hostiles, and ships.
type Actor struct {
	ID         string // logical ID (e.g., "player", "npc_guard")
	Archetype  string // "player", "npc", "enemy", "ship"
	Persistent bool   // whether to keep across scene transitions
}

func (a *Actor) Name() string { return "Actor" }

func (a *Actor) Name() string { return "Actor" }
