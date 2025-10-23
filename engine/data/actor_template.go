package data

// ActorTemplate defines one spawnable actor’s configuration.
type ActorTemplate struct {
	Name       string               `json:"name"`
	Archetype  string               `json:"archetype"`
	Persistent bool                 `json:"persistent"`
	Sprite     ActorSpriteTemplate  `json:"sprite"`
	Velocity   *ActorVelocityPreset `json:"velocity,omitempty"`
	AI         *ActorAITemplate     `json:"ai,omitempty"`
}

// ActorSpriteTemplate defines the sprite for an actor.
type ActorSpriteTemplate struct {
	Image          string  `json:"image"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	Rotation       float64 `json:"rotation"`
	FlipHorizontal bool    `json:"flip_horizontal"`
	PixelPerfect   bool    `json:"pixel_perfect"`
}

// ActorVelocityPreset defines a default velocity vector.
type ActorVelocityPreset struct {
	VX float64 `json:"vx"`
	VY float64 `json:"vy"`
}

