package render

import (
	"fmt"
	"math"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

type System struct{}

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
		fmt.Println("[RENDER] ⚠ No camera found")
		return
	}

	bounds := screen.Bounds()
	sw := float64(bounds.Dx())
	sh := float64(bounds.Dy())
	halfW := sw / 2
	halfH := sh / 2

	drawn := 0
	for _, e := range w.Entities {
		pos, ok1 := e.Get("Position").(*ecs.Position)
		sprite, ok2 := e.Get("Sprite").(*ecs.Sprite)
		if !ok1 || !ok2 || sprite.Image == nil {
			continue
		}

		op := platform.NewDrawImageOptions()
		op.SetFilter(platform.FilterNearest)

		imgBounds := sprite.Image.Bounds()
		imgW := float64(imgBounds.Dx())
		imgH := float64(imgBounds.Dy())

		entityScale := float64(sprite.Width) / imgW
		totalScale := math.Max(0.01, cam.Scale*entityScale)

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
		drawX := (pos.X - cam.X) * cam.Scale
		drawY := (pos.Y - cam.Y) * cam.Scale
		op.Translate(drawX+halfW, drawY+halfH)

		screen.DrawImage(sprite.Image, op)
		drawn++
	}

	if drawn == 0 {
		fmt.Println("[RENDER] ⚠ No sprites drawn this frame")
	} else {
		fmt.Printf("[RENDER] ✅ Drew %d entities\n", drawn)
	}
}

