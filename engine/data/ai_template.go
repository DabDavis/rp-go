package data

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed ai.json
var embeddedAIConfig []byte

// AIConfigDatabase defines the full JSON structure for AI templates.
// Each entry defines behavior presets keyed by template name or archetype.
type AIConfigDatabase struct {
	Templates []AIBehaviorTemplate `json:"templates"`
}

// AIBehaviorTemplate describes one AI configuration preset.
type AIBehaviorTemplate struct {
	Name string            `json:"name"`  // e.g., "enemy_scout"
	AI   ActorAITemplate   `json:"ai"`    // Reuses shared AI schema
	Tags map[string]string `json:"tags"`  // Optional metadata (role, faction, etc.)
}

// LoadAIConfig loads an AI configuration JSON file (with fallback).
func LoadAIConfig(path string) AIConfigDatabase {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("[DATA] Using embedded ai.json (missing %s)\n", path)
		data = embeddedAIConfig
	}

	var cfg AIConfigDatabase
	if err := json.Unmarshal(data, &cfg); err != nil {
		panic(fmt.Errorf("failed to parse AI config %q: %w", path, err))
	}
	return cfg
}

