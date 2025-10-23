//go:build !headless

package platform_desktop

import "github.com/hajimehoshi/ebiten/v2"

// Wheel returns the current mouse wheel scroll deltas.
func Wheel() (float64, float64) { return ebiten.Wheel() }

// ActualFPS returns the current rendering FPS as reported by Ebiten.
func ActualFPS() float64 { return ebiten.ActualFPS() }

// SetWindowSize changes the window dimensions.
func SetWindowSize(w, h int) { ebiten.SetWindowSize(w, h) }

// SetWindowTitle changes the window title bar text.
func SetWindowTitle(t string) { ebiten.SetWindowTitle(t) }

// MousePosition returns the current mouse position in screen coordinates.
func MousePosition() (int, int) { return ebiten.CursorPosition() }

// MouseButton type alias for Ebiten's button constants.
type MouseButton = ebiten.MouseButton

const (
    MouseButtonLeft   = ebiten.MouseButtonLeft
    MouseButtonRight  = ebiten.MouseButtonRight
    MouseButtonMiddle = ebiten.MouseButtonMiddle
)

// IsMouseButtonPressed reports whether a mouse button is currently pressed.
func IsMouseButtonPressed(b MouseButton) bool {
    return ebiten.IsMouseButtonPressed(b)
}

