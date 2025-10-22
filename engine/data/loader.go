package data

import (
	_ "embed"
	"encoding/json"
	"os"
)

//go:embed render_config.json
var embeddedRenderConfig []byte

type RenderConfig struct {
	Window struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"window"`

	Viewport struct {
		Width    int     `json:"width"`
		Height   int     `json:"height"`
		Scale    float64 `json:"scale"`
		MinScale float64 `json:"min_scale"`
		MaxScale float64 `json:"max_scale"`
		ZoomStep float64 `json:"zoom_step"`
		ZoomLerp float64 `json:"zoom_lerp"`
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
		data = embeddedRenderConfig
	}
	var cfg RenderConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	return cfg
}
