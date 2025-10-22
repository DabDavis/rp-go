package render

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"rp-go/engine/ecs"
)

type System struct{}

func (s *System) Update(*ecs.World) {}

func (s *System) Draw(w *ecs.World, screen *ebiten.Image) {
	var cam *ecs.Camera
	for _, e := range w.Entities {
		if c, ok := e.Get("Camera").(*ecs.Camera); ok {
			cam = c
			break
		}
	}
	if cam == nil {
		fmt.Println("[RENDER] ⚠️ No camera found")
		return
	}

	sw := float64(screen.Bounds().Dx())
	sh := float64(screen.Bounds().Dy())
	halfW := sw / 2
	halfH := sh / 2

	drawn := 0
	for _, e := range w.Entities {
		pos, ok1 := e.Get("Position").(*ecs.Position)
		sprite, ok2 := e.Get("Sprite").(*ecs.Sprite)
		if !ok1 || !ok2 || sprite.Image == nil {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterNearest

		entityScale := float64(sprite.Width) / float64(sprite.Image.Bounds().Dx())
		totalScale := math.Max(0.01, cam.Scale*entityScale)
		op.GeoM.Scale(totalScale, totalScale)
		op.GeoM.Rotate(sprite.Rotation + cam.Rotation)

		drawX := (pos.X-cam.X)*cam.Scale + halfW - (float64(sprite.Width)/2)*cam.Scale
		drawY := (pos.Y-cam.Y)*cam.Scale + halfH - (float64(sprite.Height)/2)*cam.Scale
		op.GeoM.Translate(drawX, drawY)

		screen.DrawImage(sprite.Image, op)
		drawn++
	}

	if drawn == 0 {
		fmt.Println("[RENDER] ⚠️ No sprites drawn this frame")
	} else {
		fmt.Printf("[RENDER] ✅ Drew %d entities\n", drawn)
	}
}

