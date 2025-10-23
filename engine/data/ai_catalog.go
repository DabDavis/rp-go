package data

// AIActionCatalog is the top-level JSON schema for ai.json.
// It defines reusable named AI behaviors that can be attached to actors.
type AIActionCatalog struct {
	Actions []AIActionTemplate `json:"actions"`
}

// AIActionTemplate defines one reusable AI behavior.
// It can represent a basic action ("pursue"), a conditional, or a scripted sequence.
type AIActionTemplate struct {
	Name       string                 `json:"name"`       // Unique internal name (used by actors)
	Type       string                 `json:"type"`       // Behavior type (from catalog registry)
	Priority   int                    `json:"priority"`   // Order in execution
	Conditions map[string]any         `json:"conditions"` // Optional conditional filters
	Params     map[string]any         `json:"params"`     // Behavior parameters (target, speed, etc.)
}

