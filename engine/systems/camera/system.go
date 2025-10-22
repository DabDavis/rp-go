package camera

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
)

// Config controls runtime camera zoom limits and responsiveness.
type Config struct {
	MinScale float64
	MaxScale float64
	ZoomStep float64
	ZoomLerp float64
}

func (c Config) normalized() Config {
	if c.MinScale <= 0 {
		c.MinScale = 0.5
	}
	if c.MaxScale <= 0 {
		c.MaxScale = 3
	}
	if c.MaxScale < c.MinScale {
		c.MaxScale = c.MinScale
	}
	if c.ZoomStep <= 0 {
		c.ZoomStep = 0.1
	}
	if c.ZoomLerp < 0 {
		c.ZoomLerp = 0
	}
	return c
}

type System struct {
	cfg        Config
	subscribed bool
}

func NewSystem(cfg Config) *System {
	return &System{cfg: cfg.normalized()}
}

func (s *System) Update(w *ecs.World) {
	var cam *ecs.Camera
	var target *ecs.Position
	var targetSprite *ecs.Sprite

	for _, e := range w.Entities {
		if c, ok := e.Get("Camera").(*ecs.Camera); ok {
			cam = c
		}
		if e.Has("CameraTarget") {
			if pos, ok := e.Get("Position").(*ecs.Position); ok {
				target = pos
			}
			if spr, ok := e.Get("Sprite").(*ecs.Sprite); ok {
				targetSprite = spr
			}
		}
	}

	if cam == nil || target == nil {
		return
	}

	// One-time wiring for zoom events.
	if !s.subscribed {
		if bus, ok := w.EventBus.(*events.TypedBus); ok && bus != nil {
			events.Subscribe(bus, func(ev events.CameraZoomEvent) {
				cam.TargetScale = clamp(ev.NewScale, cam.MinScale, cam.MaxScale)
			})
			s.subscribed = true
		}
	}

	// Ensure the camera carries sane zoom limits and defaults.
	if cam.MinScale <= 0 {
		cam.MinScale = s.cfg.MinScale
	}
	if cam.MaxScale <= 0 {
		cam.MaxScale = s.cfg.MaxScale
	}
	if cam.MaxScale < cam.MinScale {
		cam.MaxScale = cam.MinScale
	}
	if cam.DefaultScale <= 0 {
		cam.DefaultScale = clamp(cam.Scale, cam.MinScale, cam.MaxScale)
	}
	if cam.TargetScale <= 0 {
		cam.TargetScale = clamp(cam.Scale, cam.MinScale, cam.MaxScale)
	} else {
		cam.TargetScale = clamp(cam.TargetScale, cam.MinScale, cam.MaxScale)
	}

	// Handle zoom input (keyboard + mouse wheel).
	zoomDelta := 0.0
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) || inpututil.IsKeyJustPressed(ebiten.KeyKPSubtract) {
		zoomDelta -= s.cfg.ZoomStep
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) || inpututil.IsKeyJustPressed(ebiten.KeyKPAdd) {
		zoomDelta += s.cfg.ZoomStep
	}
	if inpututil.IsKeyJustPressed(ebiten.Key0) || inpututil.IsKeyJustPressed(ebiten.KeyKP0) {
		cam.TargetScale = clamp(cam.DefaultScale, cam.MinScale, cam.MaxScale)
	}
	if wheel := ebiten.Wheel(); wheel.Y != 0 {
		zoomDelta += wheel.Y * s.cfg.ZoomStep
	}
	if zoomDelta != 0 {
		cam.TargetScale = clamp(cam.TargetScale+zoomDelta, cam.MinScale, cam.MaxScale)
	}

	// Smooth follow
	cam.X += (target.X - cam.X) * 0.1
	cam.Y += (target.Y - cam.Y) * 0.1

	// Smooth zoom
	if s.cfg.ZoomLerp <= 0 {
		cam.Scale = cam.TargetScale
	} else {
		cam.Scale += (cam.TargetScale - cam.Scale) * math.Min(1, s.cfg.ZoomLerp)
		if math.Abs(cam.Scale-cam.TargetScale) < 1e-4 {
			cam.Scale = cam.TargetScale
		}
	}

	// Manual rotation (Q / E)
	const rotSpeed = 0.03
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		cam.Rotation -= rotSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		cam.Rotation += rotSpeed
	}

	// Sync player sprite rotation to camera
	if targetSprite != nil {
		targetSprite.Rotation = cam.Rotation
	}

	// Normalize rotation
	if cam.Rotation > math.Pi*2 {
		cam.Rotation -= math.Pi * 2
	} else if cam.Rotation < 0 {
		cam.Rotation += math.Pi * 2
	}
}

func (s *System) Draw(*ecs.World, *ebiten.Image) {}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
