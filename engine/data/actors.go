package data

// ActorDatabase represents the full JSON dataset of all actor templates.
type ActorDatabase struct {
	Actors []ActorTemplate `json:"actors"`
}

// ActorTemplate defines how an actor is instantiated in the world.
// It includes appearance, physics, and AI behavior.
type ActorTemplate struct {
	Name       string              `json:"name"`        // Template name, e.g. "drone"
	Archetype  string              `json:"archetype"`   // Type group for filtering or spawning
	Persistent bool                `json:"persistent"`  // If true, survives scene transitions
	Sprite     ActorSpriteTemplate `json:"sprite"`      // Visual properties
	Velocity   *ActorVelocity      `json:"velocity"`    // Optional starting motion vector
	AI         *ActorAITemplate    `json:"ai"`          // Optional AI behavior preset
}

// ActorSpriteTemplate defines the sprite used by an actor, including image and transform data.
type ActorSpriteTemplate struct {
	Image          string  `json:"image"`           // Path to sprite asset
	Width          int     `json:"width"`           // Optional explicit width
	Height         int     `json:"height"`          // Optional explicit height
	Rotation       float64 `json:"rotation"`        // Initial rotation (radians)
	FlipHorizontal bool    `json:"flip_horizontal"` // Mirror horizontally
	PixelPerfect   bool    `json:"pixel_perfect"`   // Disable smoothing/filtering
}

// ActorVelocity defines an initial velocity vector for a spawned actor.
type ActorVelocity struct {
	VX float64 `json:"vx"`
	VY float64 `json:"vy"`
}

// ActorVelocityPreset is an alias for ActorVelocity for backward compatibility
// with legacy tests and serialized data.
type ActorVelocityPreset = ActorVelocity

