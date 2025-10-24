package menu

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"
)

/*───────────────────────────────────────────────*
 | MENU SCENE                                    |
 *───────────────────────────────────────────────*/

type Scene struct {
	init       bool
	flashTimer int
	showText   bool
}

/*───────────────────────────────────────────────*
 | CORE                                           |
 *───────────────────────────────────────────────*/

func (s *Scene) Name() string { return "menu" }

func (s *Scene) Init(w *ecs.World) {
	if s.init {
		return
	}
	s.init = true
	s.flashTimer = 0
	s.showText = true
	fmt.Println("[SCENE] Init: Main Menu")
}

func (s *Scene) Unload(w *ecs.World) {
	fmt.Println("[SCENE] Unload:", s.Name())
}

/*───────────────────────────────────────────────*
 | FRAME UPDATE                                  |
 *───────────────────────────────────────────────*/

func (s *Scene) Update(w *ecs.World) {
	s.flashTimer++
	if s.flashTimer%40 == 0 {
		s.showText = !s.showText
	}

	// Handle input — press Enter to go to space.Scene
	if inpututil.IsKeyJustPressed(platform.KeyEnter) ||
		inpututil.IsKeyJustPressed(platform.KeySpace) {
		fmt.Println("[MENU] Enter pressed → loading space.Scene")

		if bus, ok := w.EventBus.(*events.TypedBus); ok && bus != nil {
			events.Publish(bus, events.SceneChangeEvent{
				Target: "space",
				Scene:  &space.Scene{},
			})
		}
	}
}

/*───────────────────────────────────────────────*
 | FRAME DRAW                                    |
 *───────────────────────────────────────────────*/

func (s *Scene) Draw(w *ecs.World, screen *platform.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	// background
	screen.Fill(color.RGBA{R: 6, G: 12, B: 18, A: 255})

	title := "R P G   P R O J E C T"
	sub := "Press Enter to Start"

	// center text roughly using offsets
	platform.DrawText(screen, title, platform.DefaultFont(), w/2-160, h/2-40, color.White)
	if s.showText {
		platform.DrawText(screen, sub, platform.DefaultFont(), w/2-100, h/2+20, color.RGBA{200, 200, 220, 255})
	}
}

