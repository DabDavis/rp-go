package data

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed actors.json
var embeddedActors []byte

// ActorDatabase holds template definitions for spawning actors at runtime.
type ActorDatabase struct {
	Actors []ActorTemplate `json:"actors"`
}

// ActorTemplate defines how to construct an actor and its default components.
type ActorTemplate struct {
	Name       string               `json:"name"`
	Archetype  string               `json:"archetype"`
	Persistent bool                 `json:"persistent"`
	Sprite     ActorSpriteTemplate  `json:"sprite"`
	Velocity   *ActorVelocityPreset `json:"velocity"`
}

// ActorSpriteTemplate describes the sprite component attached to a spawned actor.
type ActorSpriteTemplate struct {
	Image          string  `json:"image"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	PixelPerfect   bool    `json:"pixel_perfect"`
	Rotation       float64 `json:"rotation"`
	FlipHorizontal bool    `json:"flip_horizontal"`
}

// ActorVelocityPreset defines the default velocity component to attach on spawn.
type ActorVelocityPreset struct {
	VX float64 `json:"vx"`
	VY float64 `json:"vy"`
}

// LoadActorDatabase reads actor templates from disk, falling back to the embedded copy.
func LoadActorDatabase(path string) ActorDatabase {
	data, err := os.ReadFile(path)
	if err != nil {
		data = embeddedActors
	}

	var db ActorDatabase
	if err := json.Unmarshal(data, &db); err != nil {
		panic(fmt.Errorf("failed to parse actor database: %w", err))
	}
	return db
}

// TemplateByName retrieves an actor template by its unique name.
func (db ActorDatabase) TemplateByName(name string) (ActorTemplate, bool) {
	for _, tpl := range db.Actors {
		if tpl.Name == name {
			return tpl, true
		}
	}
	return ActorTemplate{}, false
}
