package ecs

// DefaultAISpeed defines the base movement speed for AI-driven actors
// when no behavior-specific speed is set.
const DefaultAISpeed = 2.5

// --- AI Controller ----------------------------------------------------------

// AIController encapsulates all AI-driven behaviors an actor can perform.
// Each behavior type (Follow, Pursue, Patrol, Retreat, Travel) can be
// independently configured, and the system chooses which to apply each frame.
type AIController struct {
    // Global movement speed for this actor (used if no behavior-specific override is provided).
    Speed float64

    // Optional behaviors — only active if non-nil.
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

// SpeedFor returns the effective movement speed for a given behavior,
// prioritizing per-behavior overrides, then the controller’s Speed,
// then falling back to DefaultAISpeed.
func (c *AIController) SpeedFor(override float64) float64 {
    if override > 0 {
        return override
    }
    if c != nil && c.Speed > 0 {
        return c.Speed
    }
    return DefaultAISpeed
}

// --- Behavior: Follow -------------------------------------------------------

// AIFollowBehavior keeps an actor near a target entity, maintaining an
// offset and preferred distance range.
type AIFollowBehavior struct {
    // Target entity identifier (can use archetype:, template:, or direct ID).
    Target string

    // Desired positional offset relative to the target.
    OffsetX, OffsetY float64

    // Distance constraints.
    MinDistance float64 // Minimum allowed distance before stopping.
    MaxDistance float64 // Maximum allowed distance before slowing down.

    // Optional movement speed override.
    Speed float64
}

// --- Behavior: Pursue -------------------------------------------------------

// AIPursueBehavior aggressively chases a target until it enters engagement range.
type AIPursueBehavior struct {
    Target         string  // Target entity to pursue.
    EngageDistance float64 // Maximum distance to begin pursuit.
    Speed          float64 // Optional speed override.
}

// --- Behavior: Retreat ------------------------------------------------------

// AIRetreatBehavior makes an actor flee from a target until it reaches safety.
type AIRetreatBehavior struct {
    Target          string  // Target entity to avoid.
    TriggerDistance float64 // Distance at which retreat begins.
    SafeDistance    float64 // Distance considered "safe" to stop.
    Speed           float64 // Optional speed override.
}

// --- Behavior: Path (Patrol / Travel) --------------------------------------

// AIPathBehavior represents a waypoint-driven route used for patrols or travel.
// Supports variants: "loop", "pingpong", "once", and "random".
type AIPathBehavior struct {
    Variant   string        // Path variant (loop, pingpong, once, random).
    Waypoints []AIWaypoint  // Ordered list of waypoints.
    Speed     float64       // Optional movement speed override.
}

// AIWaypoint identifies a coordinate in the world that the AI should navigate to.
type AIWaypoint struct {
    X float64
    Y float64
}

// --- Path State -------------------------------------------------------------

// AIPathState tracks runtime progress through a list of waypoints.
// It is updated internally by the AI system during traversal.
type AIPathState struct {
    Index     int  // Current waypoint index.
    Forward   bool // Direction of traversal (for pingpong paths).
    Completed bool // True if traversal is finished for "once" variant.
}

// Reset reinitializes the path state to its starting position and direction.
func (s *AIPathState) Reset() {
    s.Index = 0
    s.Forward = true
    s.Completed = false
}

