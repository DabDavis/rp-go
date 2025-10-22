package gfx

import (
	"encoding/json"
	"os"
)

type RenderConfig struct {
	TileSize int     `json:"tile_size"`
	Scale    float64 `json:"scale"`
}

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

