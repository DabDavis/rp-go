package ecs

import "rp-go/engine/platform"

/*───────────────────────────────────────────────*
 | BASIC COMPONENTS                              |
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

// DrawSize returns the sprite’s current rendered dimensions in pixels.
func (s *Sprite) DrawSize() (float64, float64) {
	s.ensureCache()
	return float64(s.cachedTargetWidth), float64(s.cachedTargetHeight)
}

/*───────────────────────────────────────────────*
 | CAMERA                                        |
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
 | ACTOR + INPUT                                 |
 *───────────────────────────────────────────────*/

type Actor struct {
	ID         string
	Archetype  string
	Persistent bool
}
func (a *Actor) Name() string { return "Actor" }

type PlayerInput struct{ Enabled bool }
func (p *PlayerInput) Name() string { return "PlayerInput" }

