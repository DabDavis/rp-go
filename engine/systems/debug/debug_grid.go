package debug

import (
	"image/color"
	"math"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"
)

// DrawDebugGrid renders a faint grid that follows the camera.
func DrawDebugGrid(w *ecs.World, screen *platform.Image) {
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
	sw := float64(bounds.Dx())
	sh := float64(bounds.Dy())
	tileSize := 32.0
	gridColor := color.RGBA{90, 90, 110, 100}

	left := cam.X - (sw / 2 / cam.Scale)
	right := cam.X + (sw / 2 / cam.Scale)
	top := cam.Y - (sh / 2 / cam.Scale)
	bottom := cam.Y + (sh / 2 / cam.Scale)

	startX := math.Floor(left/tileSize) * tileSize
	startY := math.Floor(top/tileSize) * tileSize

	line := platform.NewImage(1, 1)
	line.Fill(gridColor)

	for x := startX; x <= right; x += tileSize {
		screenX := (x-cam.X)*cam.Scale + sw/2
		if screenX < 0 || screenX > sw {
			continue
		}
		op := platform.NewDrawImageOptions()
		op.Scale(1, sh)
		op.Translate(screenX, 0)
		screen.DrawImage(line, op)
	}

	for y := startY; y <= bottom; y += tileSize {
		screenY := (y-cam.Y)*cam.Scale + sh/2
		if screenY < 0 || screenY > sh {
			continue
		}
		op := platform.NewDrawImageOptions()
		op.Scale(sw, 1)
		op.Translate(0, screenY)
		screen.DrawImage(line, op)
	}

	// Subscribe to debug toggle event (show/hide grid)
	if bus, ok := w.EventBus.(*events.TypedBus); ok && bus != nil {
		events.Subscribe(bus, func(e events.DebugToggleEvent) {
			if e.Enabled {
				gridColor.A = 120
			} else {
				gridColor.A = 0
			}
		})
	}
}
