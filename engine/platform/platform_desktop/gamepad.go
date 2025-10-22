package platform_desktop

import "github.com/hajimehoshi/ebiten/v2"

type GamepadID = ebiten.GamepadID
type StandardGamepadButton = ebiten.StandardGamepadButton
type StandardGamepadAxis = ebiten.StandardGamepadAxis

func IsStandardGamepadLayoutAvailable(id GamepadID) bool {
	return ebiten.IsStandardGamepadLayoutAvailable(id)
}

const (
	StandardGamepadAxisLeftStickHorizontal = ebiten.StandardGamepadAxisLeftStickHorizontal
	StandardGamepadAxisLeftStickVertical   = ebiten.StandardGamepadAxisLeftStickVertical
)

func GamepadIDs() []GamepadID                            { return ebiten.GamepadIDs() }
func StandardGamepadAxisValue(id GamepadID, axis StandardGamepadAxis) float64 {
	return ebiten.StandardGamepadAxisValue(id, axis)
}
func IsStandardGamepadButtonPressed(id GamepadID, button StandardGamepadButton) bool {
	return ebiten.IsStandardGamepadButtonPressed(id, button)
}

func IsGamepadLeft(id GamepadID) bool  { return ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) < -0.5 }
func IsGamepadRight(id GamepadID) bool { return ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) > 0.5 }
func IsGamepadUp(id GamepadID) bool    { return ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical) < -0.5 }
func IsGamepadDown(id GamepadID) bool  { return ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical) > 0.5 }

