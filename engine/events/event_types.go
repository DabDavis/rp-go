package events

// Core game event types (pure, no ECS dependency)

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

type DebugToggleEvent struct {
	Enabled bool
}

// SceneChangeEvent is sent to request a transition to another scene.
type SceneChangeEvent struct {
	Target string // e.g. "space" or "planet"
	Scene  any    // stored as generic; scene.Manager will type-assert to ecs.Scene
}

