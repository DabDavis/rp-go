package camera

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"rp-go/engine/ecs"
)

type System struct{}

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

	// Smooth follow
	cam.X += (target.X - cam.X) * 0.1
	cam.Y += (target.Y - cam.Y) * 0.1

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

