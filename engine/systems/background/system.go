package background

import (
	"image/color"
	"math"
	"math/rand"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

// A simple star definition
type star struct {
	X, Y       float64
	Brightness uint8
}

// System draws a parallax starfield in the background layer.
type System struct {
	stars []star
}

// This system renders in the background layer.
func (s *System) Layer() ecs.DrawLayer { return ecs.LayerBackground }

func (s *System) Update(*ecs.World) {}

// Lazy initialize stars once
func (s *System) ensureStars(screen *platform.Image) {
	if len(s.stars) > 0 {
		return
	}
	bounds := screen.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	numStars := (width * height) / 2000 // density factor

	s.stars = make([]star, numStars)
	for i := range s.stars {
		s.stars[i] = star{
			X:          rand.Float64() * float64(width),
			Y:          rand.Float64() * float64(height),
			Brightness: uint8(155 + rand.Intn(100)), // 155–255
		}
	}
}

// Draw renders the starfield slightly offset by camera position for parallax.
func (s *System) Draw(w *ecs.World, screen *platform.Image) {
	s.ensureStars(screen)

	var cam *ecs.Camera
	if manager := w.EntitiesManager(); manager != nil {
		_, comp := manager.FirstComponent("Camera")
		cam, _ = comp.(*ecs.Camera)
	}

	bounds := screen.Bounds()
	width, height := float64(bounds.Dx()), float64(bounds.Dy())

	// Clear the screen to deep space (almost black)
	screen.Fill(color.RGBA{5, 5, 10, 255})

	// Parallax offset (camera position * small factor)
	offsetX, offsetY := 0.0, 0.0
	if cam != nil {
		offsetX = cam.X * 0.05
		offsetY = cam.Y * 0.05
	}

	// Draw stars
	for _, star := range s.stars {
		x := math.Mod(star.X-offsetX, width)
		y := math.Mod(star.Y-offsetY, height)

		if x < 0 {
			x += width
		}
		if y < 0 {
			y += height
		}

		screen.FillRect(int(x), int(y), 1, 1, color.RGBA{
			R: star.Brightness,
			G: star.Brightness,
			B: star.Brightness,
			A: 255,
		})
	}
}
