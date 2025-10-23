package data

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed ai.json
var embeddedAI []byte

// LoadAICatalog loads and parses ai.json from disk, or falls back to the embedded version.
func LoadAICatalog(path string) AIActionCatalog {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("[DATA] Using embedded ai.json (missing %s)\n", path)
		data = embeddedAI
	}
	var catalog AIActionCatalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		panic(fmt.Errorf("failed to parse ai.json: %w", err))
	}
	fmt.Printf("[DATA] Loaded %d AI actions from %s\n", len(catalog.Actions), path)
	return catalog
}

