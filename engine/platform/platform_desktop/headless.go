//go:build !headless

package platform_desktop

import "github.com/hajimehoshi/ebiten/v2"

type Game interface {
	Update() error
	Draw(screen *Image)
	Layout(outsideWidth, outsideHeight int) (int, int)
}

type gameAdapter struct {
	game Game
}

func (g *gameAdapter) Update() error              { return g.game.Update() }
func (g *gameAdapter) Draw(screen *ebiten.Image)  { g.game.Draw(newImageFromNative(screen)) }
func (g *gameAdapter) Layout(w, h int) (int, int) { return g.game.Layout(w, h) }

func RunGame(game Game) error {
	return ebiten.RunGame(&gameAdapter{game: game})
}

func RunHeadless(game Game, frames, width, height int) error {
	if frames <= 0 {
		return nil
	}
	offscreen := NewImage(width, height)
	for i := 0; i < frames; i++ {
		if err := game.Update(); err != nil {
			return err
		}
		offscreen.Clear()
		game.Draw(offscreen)
	}
	return nil
}
