package data

// ActorAITemplate defines AI behaviors that can be attached to a spawned actor.
type ActorAITemplate struct {
	Speed   float64         `json:"speed"`
	Follow  *ActorAIFollow  `json:"follow"`
	Pursue  *ActorAIPursue  `json:"pursue"`
	Patrol  *ActorAIPatrol  `json:"patrol"`
	Retreat *ActorAIRetreat `json:"retreat"`
	Travel  *ActorAITravel  `json:"travel"`
}

// ActorAIFollow instructs an actor to maintain a position relative to a target.
type ActorAIFollow struct {
	Target      string  `json:"target"`
	OffsetX     float64 `json:"offset_x"`
	OffsetY     float64 `json:"offset_y"`
	MinDistance float64 `json:"min_distance"`
	MaxDistance float64 `json:"max_distance"`
	Speed       float64 `json:"speed"`
}

// ActorAIPursue causes an actor to directly chase a target.
type ActorAIPursue struct {
	Target         string  `json:"target"`
	EngageDistance float64 `json:"engage_distance"`
	Speed          float64 `json:"speed"`
}

// ActorAIPatrol configures a waypoint patrol pattern.
type ActorAIPatrol struct {
	Variant   string            `json:"variant"`
	Waypoints []ActorAIWaypoint `json:"waypoints"`
	Speed     float64           `json:"speed"`
}

// ActorAITravel configures a long-distance travel route between locations.
type ActorAITravel struct {
	Variant   string            `json:"variant"`
	Waypoints []ActorAIWaypoint `json:"waypoints"`
	Speed     float64           `json:"speed"`
}

// ActorAIRetreat makes an actor flee from a specified target when threatened.
type ActorAIRetreat struct {
	Target          string  `json:"target"`
	TriggerDistance float64 `json:"trigger_distance"`
	SafeDistance    float64 `json:"safe_distance"`
	Speed           float64 `json:"speed"`
}

// ActorAIWaypoint represents a single coordinate in a patrol or travel path.
type ActorAIWaypoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
