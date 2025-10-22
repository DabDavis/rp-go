//go:build !headless

package platform_desktop

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Key = ebiten.Key

const (
	KeyArrowLeft  Key = ebiten.KeyArrowLeft
	KeyArrowRight Key = ebiten.KeyArrowRight
	KeyArrowUp    Key = ebiten.KeyArrowUp
	KeyArrowDown  Key = ebiten.KeyArrowDown
	KeyA          Key = ebiten.KeyA
	KeyD          Key = ebiten.KeyD
	KeyW          Key = ebiten.KeyW
	KeyS          Key = ebiten.KeyS
	KeyQ          Key = ebiten.KeyQ
	KeyE          Key = ebiten.KeyE
	KeyMinus      Key = ebiten.KeyMinus
	KeyEqual      Key = ebiten.KeyEqual
	Key0          Key = ebiten.Key0
	KeyKP0        Key = ebiten.KeyKP0
	KeyKPAdd      Key = ebiten.KeyKPAdd
	KeyKPSubtract Key = ebiten.KeyKPSubtract
)

func IsKeyPressed(k Key) bool     { return ebiten.IsKeyPressed(k) }
func IsKeyJustPressed(k Key) bool { return inpututil.IsKeyJustPressed(k) }
