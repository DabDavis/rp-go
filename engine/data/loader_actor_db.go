package data

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed actors.json
var embeddedActors []byte

// LoadActorDatabase loads and parses an actor database JSON file.
// Falls back to the embedded version if the external file is missing.
func LoadActorDatabase(path string) ActorDatabase {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("[DATA] Using embedded actor database (missing %s)\n", path)
		data = embeddedActors
	}

	var db ActorDatabase
	if err := json.Unmarshal(data, &db); err != nil {
		panic(fmt.Errorf("failed to parse actor database %q: %w", path, err))
	}

	return db
}

