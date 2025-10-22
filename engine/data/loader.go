package data

import (
	"encoding/json"
	"os"
)

type RenderConfig struct {
	Window struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"window"`

	Viewport struct {
		Width  int     `json:"width"`
		Height int     `json:"height"`
		Scale  float64 `json:"scale"`
	} `json:"viewport"`

	Player struct {
		SpriteWidth  int     `json:"sprite_width"`
		SpriteHeight int     `json:"sprite_height"`
		Scale        float64 `json:"scale"`
	} `json:"player"`

	Terrain struct {
		TileSize int     `json:"tile_size"`
		Scale    float64 `json:"scale"`
	} `json:"terrain"`
}

// LoadRenderConfig reads and parses the JSON config file.
func LoadRenderConfig(path string) RenderConfig {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var cfg RenderConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	return cfg
}

