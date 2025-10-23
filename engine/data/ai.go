package data

// ActorAITemplate defines a full AI configuration for an actor.
// Each field corresponds to a possible autonomous behavior, any of which may be omitted.
//
// JSON example:
//
//  {
//    "speed": 2.5,
//    "pursue": { "target": "player", "engage_distance": 200 },
//    "patrol": { "variant": "loop", "waypoints": [{ "x": 64, "y": 128 }] }
//  }
type ActorAITemplate struct {
    Speed   float64         `json:"speed"`
    Follow  *ActorAIFollow  `json:"follow"`
    Pursue  *ActorAIPursue  `json:"pursue"`
    Patrol  *ActorAIPatrol  `json:"patrol"`
    Retreat *ActorAIRetreat `json:"retreat"`
    Travel  *ActorAITravel  `json:"travel"`
}

// --- Follow Behavior --------------------------------------------------------

// ActorAIFollow defines how an actor maintains position near a moving target.
type ActorAIFollow struct {
    Target      string  `json:"target"`        // Target entity (ID, archetype:, or template:)
    OffsetX     float64 `json:"offset_x"`      // Horizontal offset from target
    OffsetY     float64 `json:"offset_y"`      // Vertical offset from target
    MinDistance float64 `json:"min_distance"`  // Stop following if within this distance
    MaxDistance float64 `json:"max_distance"`  // Gradually slow if within this distance
    Speed       float64 `json:"speed"`         // Optional speed override
}

// --- Pursue Behavior --------------------------------------------------------

// ActorAIPursue defines aggressive chase logic toward a target entity.
type ActorAIPursue struct {
    Target         string  `json:"target"`          // Target entity ID or query
    EngageDistance float64 `json:"engage_distance"` // Max range to engage pursuit
    Speed          float64 `json:"speed"`           // Optional speed override
}

// --- Patrol Behavior --------------------------------------------------------

// ActorAIPatrol defines a local area patrol route.
// Variants: "loop", "pingpong", "once", "random"
type ActorAIPatrol struct {
    Variant   string            `json:"variant"`   // Path type
    Waypoints []ActorAIWaypoint `json:"waypoints"` // Ordered waypoint list
    Speed     float64           `json:"speed"`     // Optional speed override
}

// --- Travel Behavior --------------------------------------------------------

// ActorAITravel defines a long-distance path for travel-type motion.
// Shares identical schema to Patrol for convenience.
type ActorAITravel struct {
    Variant   string            `json:"variant"`
    Waypoints []ActorAIWaypoint `json:"waypoints"`
    Speed     float64           `json:"speed"`
}

// --- Retreat Behavior -------------------------------------------------------

// ActorAIRetreat defines how an actor flees from a target until safe.
type ActorAIRetreat struct {
    Target          string  `json:"target"`           // Entity or query to avoid
    TriggerDistance float64 `json:"trigger_distance"` // Distance that triggers fleeing
    SafeDistance    float64 `json:"safe_distance"`    // Distance that counts as "safe"
    Speed           float64 `json:"speed"`            // Optional speed override
}

// --- Shared Data ------------------------------------------------------------

// ActorAIWaypoint defines a single navigation coordinate (tile or pixel space).
type ActorAIWaypoint struct {
    X float64 `json:"x"`
    Y float64 `json:"y"`
}

