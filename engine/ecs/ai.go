package ecs

const DefaultAISpeed = 2.5

// AIController drives non-player actors according to their configured behavior.
type AIController struct {
	Speed   float64
	Follow  *AIFollowBehavior
	Pursue  *AIPursueBehavior
	Patrol  *AIPathBehavior
	Retreat *AIRetreatBehavior
	Travel  *AIPathBehavior

	PatrolState AIPathState
	TravelState AIPathState
}

func (c *AIController) Name() string { return "AIController" }

// SpeedFor returns a usable speed for a behavior, honoring per-behavior overrides.
func (c *AIController) SpeedFor(override float64) float64 {
	if override > 0 {
		return override
	}
	if c != nil && c.Speed > 0 {
		return c.Speed
	}
	return DefaultAISpeed
}

// AIFollowBehavior keeps an actor near a moving target.
type AIFollowBehavior struct {
	Target      string
	OffsetX     float64
	OffsetY     float64
	MinDistance float64
	MaxDistance float64
	Speed       float64
}

// AIPursueBehavior aggressively chases a target.
type AIPursueBehavior struct {
	Target         string
	EngageDistance float64
	Speed          float64
}

// AIRetreatBehavior makes an actor flee from a target until reaching safety.
type AIRetreatBehavior struct {
	Target          string
	TriggerDistance float64
	SafeDistance    float64
	Speed           float64
}

// AIPathBehavior represents a waypoint-driven path (patrol or travel).
type AIPathBehavior struct {
	Variant   string
	Waypoints []AIWaypoint
	Speed     float64
}

// AIWaypoint identifies a coordinate the AI should navigate to.
type AIWaypoint struct {
	X float64
	Y float64
}

// AIPathState tracks runtime progress through a waypoint list.
type AIPathState struct {
	Index     int
	Forward   bool
	Completed bool
}

// Reset initializes the state for a new traversal.
func (s *AIPathState) Reset() {
	s.Index = 0
	s.Forward = true
	s.Completed = false
}
