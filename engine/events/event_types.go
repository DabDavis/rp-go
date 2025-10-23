package events

// --- Core Game Events -------------------------------------------------------

type EntityMovedEvent struct {
	EntityID int
	X, Y     float64
}

type EntitySpawnedEvent struct {
	EntityID int
}

type CameraZoomEvent struct {
	NewScale float64
}

// Toggles entire debug overlay (F12).
type DebugToggleEvent struct {
	Enabled bool
}

// SceneChangeEvent requests a transition to another scene.
type SceneChangeEvent struct {
	Target string // e.g. "space" or "planet"
	Scene  any    // generic; scene.Manager will type-assert to ecs.Scene
}
// --- UI Window Events -------------------------------------------------------

// WindowClosedEvent is emitted when a window's close button is clicked.
type WindowClosedEvent struct {
	ID string
}

// WindowMinimizedEvent is emitted when a window is minimized via its toolbar button.
type WindowMinimizedEvent struct {
	ID string
}

// WindowRestoredEvent is emitted when a window is restored from minimized state.
type WindowRestoredEvent struct {
	ID string
}

// WindowMovedEvent is emitted when a movable window is dragged to a new position.
type WindowMovedEvent struct {
	ID string
	X  int
	Y  int
}

// WindowToggledEvent toggles a specific window's visibility (e.g. from a toolbar).
type WindowToggledEvent struct {
	ID      string // window ID (e.g. "debug.stats")
	Enabled bool
}
// DataReloaded is emitted when a JSON configuration or data file is reloaded.
type DataReloaded struct {
	Path string // full path of reloaded file
	Type string // logical type: render_config, actor_db, etc.
}

