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
	Image          *platform.Image // pointer to GPU texture
	Width          int
	Height         int
	Rotation       float64
	FlipHorizontal bool

	PixelPerfect bool

	cachedImage        *platform.Image
	cachedSourceWidth  int
	cachedSourceHeight int
	cachedTargetWidth  int
	cachedTargetHeight int
	cachedScale        float64
}

func (s *Sprite) Name() string { return "Sprite" }

// ensureCache synchronizes cached sprite metrics with the current image
// and desired output size. It avoids expensive Bounds() calls on every draw
// by recomputing only when the source image or target dimensions change.
func (s *Sprite) ensureCache() {
	if s == nil {
		return
	}

	if s.cachedImage != s.Image {
		s.cachedSourceWidth = 0
		s.cachedSourceHeight = 0
		s.cachedScale = 0
	}

	targetW := s.Width
	targetH := s.Height

	if s.cachedSourceWidth == 0 || s.cachedSourceHeight == 0 ||
		s.cachedTargetWidth != targetW || s.cachedTargetHeight != targetH {
		if s.Image != nil {
			bounds := s.Image.Bounds()
			s.cachedSourceWidth = bounds.Dx()
			s.cachedSourceHeight = bounds.Dy()
		} else {
			s.cachedSourceWidth = 0
			s.cachedSourceHeight = 0
		}

		if s.cachedSourceWidth > 0 {
			s.cachedScale = float64(targetW) / float64(s.cachedSourceWidth)
		} else {
			s.cachedScale = 0
		}

		s.cachedTargetWidth = targetW
		s.cachedTargetHeight = targetH
		s.cachedImage = s.Image
	}
}

// NativeSize returns the intrinsic image dimensions for the sprite.
func (s *Sprite) NativeSize() (float64, float64) {
	if s == nil {
		return 0, 0
	}
	s.ensureCache()
	return float64(s.cachedSourceWidth), float64(s.cachedSourceHeight)
}

// PixelScale returns the ratio between the desired sprite width and the
// intrinsic image width. This value is cached so we only recompute when the
// sprite dimensions actually change.
func (s *Sprite) PixelScale() float64 {
	if s == nil {
		return 0
	}
	s.ensureCache()
	return s.cachedScale
}

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
