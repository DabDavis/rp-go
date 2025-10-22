package render

import (
	"math"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

type System struct{}

// Ensure this system only runs in the world pass
func (s *System) Layer() ecs.DrawLayer { return ecs.LayerWorld }

func (s *System) Update(*ecs.World) {}

// Draw renders all entities with Position + Sprite components using the active Camera.
func (s *System) Draw(w *ecs.World, screen *platform.Image) {
	var cam *ecs.Camera
	for _, e := range w.Entities {
		if c, ok := e.Get("Camera").(*ecs.Camera); ok {
			cam = c
			break
		}
	}
	if cam == nil {
		return
	}

	bounds := screen.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	for _, e := range w.Entities {
		pos, ok1 := e.Get("Position").(*ecs.Position)
		sprite, ok2 := e.Get("Sprite").(*ecs.Sprite)
		if !ok1 || !ok2 || sprite.Image == nil {
			continue
		}

		imgW, imgH := sprite.NativeSize()
		if imgW == 0 || imgH == 0 {
			continue
		}

		entityScale := sprite.PixelScale()
		if entityScale <= 0 {
			continue
		}

		effectiveScale := cam.Scale
		if sprite.PixelPerfect {
			effectiveScale = math.Max(1, math.Round(cam.Scale))
		}

		totalScale := math.Max(0.01, effectiveScale*entityScale)

		op := platform.NewDrawImageOptions()
		op.SetFilter(platform.FilterNearest)

		// Center-origin transform
		op.Translate(-imgW/2, -imgH/2)

		// Flip around center
		if sprite.FlipHorizontal {
			op.Scale(-totalScale, totalScale)
		} else {
			op.Scale(totalScale, totalScale)
		}

		// Rotate around center
		op.Rotate(sprite.Rotation)

		// Translate to world position (centered on entity)
		drawX := (pos.X - cam.X) * effectiveScale
		drawY := (pos.Y - cam.Y) * effectiveScale

		if sprite.PixelPerfect {
			drawX = math.Round(drawX)
			drawY = math.Round(drawY)
		}

		finalX := drawX + halfW
		finalY := drawY + halfH

		if sprite.PixelPerfect {
			finalX = math.Round(finalX)
			finalY = math.Round(finalY)
		}

		op.Translate(finalX, finalY)

		screen.DrawImage(sprite.Image, op)
	}
}
