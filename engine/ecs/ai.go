package ecs

/*───────────────────────────────────────────────*
 | AI CONTROLLER STRUCTURE                       |
 *───────────────────────────────────────────────*/

// DefaultAISpeed defines the base movement speed for AI-driven actors
// when no behavior-specific speed is set.
const DefaultAISpeed = 2.5

// AIController encapsulates all AI-driven behaviors an actor can perform.
// It is attached as an ECS component and managed by the AI system.
type AIController struct {
	Active  bool               // Whether this controller is currently active
	Actions []AIActionInstance // Runtime behavior list (populated from ai.json)

	// Global movement speed for this actor (used if no behavior-specific override is provided).
	Speed float64

	// Optional legacy-style behavior blocks (for direct template use)
	Follow  *AIFollowBehavior
	Pursue  *AIPursueBehavior
	Patrol  *AIPathBehavior
	Retreat *AIRetreatBehavior
	Travel  *AIPathBehavior

	// Runtime state tracking for path-based behaviors.
	PatrolState AIPathState
	TravelState AIPathState
}

// Name returns the ECS component identifier.
func (c *AIController) Name() string { return "AIController" }

/*───────────────────────────────────────────────*
 | AI ACTION INSTANCE (RUNTIME)                  |
 *───────────────────────────────────────────────*/

// AIActionInstance represents a single behavior currently active or queued.
type AIActionInstance struct {
	Name     string
	Type     string
	Priority int
	Params   map[string]any
}

/*───────────────────────────────────────────────*
 | SHARED BEHAVIOR DEFINITIONS                   |
 *───────────────────────────────────────────────*/

// AIFollowBehavior keeps an actor near a target entity.
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

// AIRetreatBehavior makes an actor flee from a target until safe.
type AIRetreatBehavior struct {
	Target          string
	TriggerDistance float64
	SafeDistance    float64
	Speed           float64
}

// AIPathBehavior defines waypoint navigation for patrols/travel.
type AIPathBehavior struct {
	Variant   string       // loop, pingpong, once, random
	Waypoints []AIWaypoint // coordinate list
	Speed     float64
}

// AIWaypoint defines a navigation target.
type AIWaypoint struct {
	X float64
	Y float64
}

// AIPathState tracks runtime waypoint progress.
type AIPathState struct {
	Index     int  // Current waypoint index
	Forward   bool // Direction for pingpong
	Completed bool // True if traversal finished
}

// Reset reinitializes the path traversal state.
func (s *AIPathState) Reset() {
	s.Index = 0
	s.Forward = true
	s.Completed = false
}

