//go:build !headless

package platform

import "rp-go/engine/platform/platform_desktop"

// -----------------------------------------------------------------------------
// Type aliases
// -----------------------------------------------------------------------------

type (
	Image                 = platform_desktop.Image
	DrawImageOptions      = platform_desktop.DrawImageOptions
	Filter                = platform_desktop.Filter
	Key                   = platform_desktop.Key
	GamepadID             = platform_desktop.GamepadID
	StandardGamepadAxis   = platform_desktop.StandardGamepadAxis
	StandardGamepadButton = platform_desktop.StandardGamepadButton
	Game                  = platform_desktop.Game
	MouseButton           = platform_desktop.MouseButton
)

// -----------------------------------------------------------------------------
// Constructors and core platform functions
// -----------------------------------------------------------------------------

var (
	NewImage            = platform_desktop.NewImage
	NewImageFromImage   = platform_desktop.NewImageFromImage
	NewDrawImageOptions = platform_desktop.NewDrawImageOptions

	RunGame     = platform_desktop.RunGame
	RunHeadless = platform_desktop.RunHeadless
)

// -----------------------------------------------------------------------------
// Input: Keyboard + Gamepad
// -----------------------------------------------------------------------------

var (
	IsKeyPressed     = platform_desktop.IsKeyPressed
	IsKeyJustPressed = platform_desktop.IsKeyJustPressed
	InputChars       = platform_desktop.InputChars

	GamepadIDs     = platform_desktop.GamepadIDs
	IsGamepadLeft  = platform_desktop.IsGamepadLeft
	IsGamepadRight = platform_desktop.IsGamepadRight
	IsGamepadUp    = platform_desktop.IsGamepadUp
	IsGamepadDown  = platform_desktop.IsGamepadDown

	IsStandardGamepadLayoutAvailable = platform_desktop.IsStandardGamepadLayoutAvailable
	StandardGamepadAxisValue         = platform_desktop.StandardGamepadAxisValue
)

// -----------------------------------------------------------------------------
// Input: Mouse (NEW)
// -----------------------------------------------------------------------------

func MousePosition() (int, int) {
	return platform_desktop.MousePosition()
}

func IsMouseButtonPressed(b MouseButton) bool {
	return platform_desktop.IsMouseButtonPressed(b)
}

const (
	MouseButtonLeft   = platform_desktop.MouseButtonLeft
	MouseButtonRight  = platform_desktop.MouseButtonRight
	MouseButtonMiddle = platform_desktop.MouseButtonMiddle
)

// -----------------------------------------------------------------------------
// Window + Timing utilities
// -----------------------------------------------------------------------------

var (
	SetWindowSize  = platform_desktop.SetWindowSize
	SetWindowTitle = platform_desktop.SetWindowTitle
	ActualFPS      = platform_desktop.ActualFPS
	Wheel          = platform_desktop.Wheel
)

// -----------------------------------------------------------------------------
// Text rendering
// -----------------------------------------------------------------------------

var DrawText = platform_desktop.DrawText

// -----------------------------------------------------------------------------
// Key constants
// -----------------------------------------------------------------------------

const (
	KeyArrowLeft  = platform_desktop.KeyArrowLeft
	KeyArrowRight = platform_desktop.KeyArrowRight
	KeyArrowUp    = platform_desktop.KeyArrowUp
	KeyArrowDown  = platform_desktop.KeyArrowDown
	KeyA          = platform_desktop.KeyA
	KeyD          = platform_desktop.KeyD
	KeyW          = platform_desktop.KeyW
	KeyS          = platform_desktop.KeyS
	KeyQ          = platform_desktop.KeyQ
	KeyE          = platform_desktop.KeyE

	KeyMinus      = platform_desktop.KeyMinus
	KeyEqual      = platform_desktop.KeyEqual
	Key0          = platform_desktop.Key0
	KeyKP0        = platform_desktop.KeyKP0
	KeyKPAdd      = platform_desktop.KeyKPAdd
	KeyKPSubtract = platform_desktop.KeyKPSubtract
	KeyEnter      = platform_desktop.KeyEnter
	KeyEscape     = platform_desktop.KeyEscape
	KeyBackspace  = platform_desktop.KeyBackspace
	KeyF12        = platform_desktop.KeyF12

	StandardGamepadAxisLeftStickHorizontal = platform_desktop.StandardGamepadAxisLeftStickHorizontal
	StandardGamepadAxisLeftStickVertical   = platform_desktop.StandardGamepadAxisLeftStickVertical

	FilterNearest = platform_desktop.FilterNearest
)

