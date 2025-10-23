package ecs

import (
	"fmt"
	"rp-go/engine/platform"
)

/*───────────────────────────────────────────────*
 | COMPONENT INTERFACE                           |
 *───────────────────────────────────────────────*/

// Component is any data-only object that can be attached to an Entity.
type Component interface {
	Name() string
}

// Tag is a simple zero-data marker component.
type Tag string

func (t Tag) Name() string { return string(t) }

/*───────────────────────────────────────────────*
 | TRANSFORM COMPONENTS                          |
 *───────────────────────────────────────────────*/

type Position struct{ X, Y float64 }
func (p *Position) Name() string { return "Position" }

type Velocity struct{ VX, VY float64 }
func (v *Velocity) Name() string { return "Velocity" }

/*───────────────────────────────────────────────*
 | SPRITE RENDERING                              |
 *───────────────────────────────────────────────*/

type Sprite struct {
	Image          *platform.Image
	Width, Height  int
	Rotation       float64
	FlipHorizontal bool
	PixelPerfect   bool

	cachedImage        *platform.Image
	cachedSourceWidth  int
	cachedSourceHeight int
	cachedTargetWidth  int
	cachedTargetHeight int
	cachedScale        float64
}

func (s *Sprite) Name() string { return "Sprite" }

func (s *Sprite) ensureCache() {
	if s == nil {
		return
	}
	if s.cachedImage != s.Image {
		s.cachedSourceWidth, s.cachedSourceHeight = 0, 0
		s.cachedScale = 0
	}
	targetW, targetH := s.Width, s.Height
	if s.cachedSourceWidth == 0 || s.cachedTargetWidth != targetW || s.cachedTargetHeight != targetH {
		if s.Image != nil {
			b := s.Image.Bounds()
			s.cachedSourceWidth, s.cachedSourceHeight = b.Dx(), b.Dy()
		}
		if s.cachedSourceWidth > 0 {
			s.cachedScale = float64(targetW) / float64(s.cachedSourceWidth)
		}
		s.cachedTargetWidth, s.cachedTargetHeight = targetW, targetH
		s.cachedImage = s.Image
	}
}

func (s *Sprite) NativeSize() (float64, float64) {
	s.ensureCache()
	return float64(s.cachedSourceWidth), float64(s.cachedSourceHeight)
}

func (s *Sprite) PixelScale() float64 {
	s.ensureCache()
	return s.cachedScale
}

func (s *Sprite) DrawSize() (float64, float64) {
	s.ensureCache()
	return float64(s.cachedTargetWidth), float64(s.cachedTargetHeight)
}

/*───────────────────────────────────────────────*
 | CAMERA COMPONENTS                             |
 *───────────────────────────────────────────────*/

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

/*───────────────────────────────────────────────*
 | ACTOR & INPUT                                 |
 *───────────────────────────────────────────────*/

type Actor struct {
	ID         string
	Archetype  string
	Persistent bool
}
func (a *Actor) Name() string { return "Actor" }

type PlayerInput struct{ Enabled bool }
func (p *PlayerInput) Name() string { return "PlayerInput" }

/*───────────────────────────────────────────────*
 | HEALTH COMPONENT                              |
 *───────────────────────────────────────────────*/

type Health struct {
	Current float64
	Max     float64
}
func (h *Health) Name() string { return "Health" }

func (h *Health) Fraction() float64 {
	if h.Max <= 0 {
		return 0
	}
	return h.Current / h.Max
}
func (h *Health) ApplyDamage(amount float64) {
	h.Current -= amount
	if h.Current < 0 {
		h.Current = 0
	}
}
func (h *Health) Heal(amount float64) {
	h.Current += amount
	if h.Current > h.Max {
		h.Current = h.Max
	}
}

/*───────────────────────────────────────────────*
 | SCRIPT STATE (for AI scripts)                 |
 *───────────────────────────────────────────────*/

type ScriptState struct {
	Step   int
	Timer  float64
	Active bool
}
func (s *ScriptState) Name() string { return "AIScriptState" }

/*───────────────────────────────────────────────*
 | ENTITY HELPERS                                |
 *───────────────────────────────────────────────*/

// AddNamed safely adds or replaces a component by explicit key.
func (e *Entity) AddNamed(name string, c Component) {
	if e == nil || c == nil {
		return
	}
	if _, exists := e.Components[name]; exists {
		fmt.Printf("[ECS] Warning: replacing component '%s' on entity %d\n", name, e.ID)
	}
	e.Components[name] = c
}

/*───────────────────────────────────────────────*
 | GENERIC COMPONENT ACCESS                      |
 *───────────────────────────────────────────────*/

// GetTyped is a standalone helper that retrieves a component as a typed value.
// Example:
//   vel := ecs.GetTyped[*ecs.Velocity](entity, "Velocity")
func GetTyped[T Component](e *Entity, name string) T {
	var zero T
	if e == nil {
		return zero
	}
	c := e.Get(name)
	if c == nil {
		return zero
	}
	val, ok := c.(T)
	if !ok {
		return zero
	}
	return val
}

