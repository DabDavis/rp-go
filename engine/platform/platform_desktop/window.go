//go:build !headless

package platform_desktop

import "github.com/hajimehoshi/ebiten/v2"

func Wheel() (float64, float64) { return ebiten.Wheel() }
func ActualFPS() float64        { return ebiten.ActualFPS() }

func SetWindowSize(w, h int)  { ebiten.SetWindowSize(w, h) }
func SetWindowTitle(t string) { ebiten.SetWindowTitle(t) }
