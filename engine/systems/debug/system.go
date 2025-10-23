package debug

import (
	"fmt"
	"image/color"
	"sort"

	"golang.org/x/image/font/basicfont"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/ui/window"
)

// Config defines layout properties for the debug overlay windows.
type Config struct {
	Margin         int
	ViewportWidth  int
	ViewportHeight int
	StatsWidth     int
	EntitiesWidth  int
	MaxEntities    int
}

// normalize fills zero-valued fields with sensible defaults.
func (c *Config) normalize() {
	if c.Margin <= 0 {
		c.Margin = 16
	}
	if c.ViewportWidth <= 0 {
		c.ViewportWidth = 640
	}
	if c.ViewportHeight <= 0 {
		c.ViewportHeight = 360
	}
	if c.StatsWidth <= 0 {
		c.StatsWidth = 260
	}
	if c.EntitiesWidth <= 0 {
		c.EntitiesWidth = 360
	}
	if c.MaxEntities <= 0 {
		c.MaxEntities = 10
	}
}

// System maintains the debug overlay windows and populates their content.
type System struct {
	cfg Config

	frameCounter int

	statsEntity    *ecs.Entity
	statsComponent *window.Component
	statsContent   *statsWindowContent

	entitiesEntity    *ecs.Entity
	entitiesComponent *window.Component
	entitiesContent   *entitiesWindowContent
}

// NewSystem constructs a debug system with the provided configuration.
func NewSystem(cfg Config) *System {
	cfg.normalize()
	return &System{
		cfg: cfg,
		statsContent: &statsWindowContent{
			lineHeight:     16,
			baselineOffset: 12,
		},
		entitiesContent: &entitiesWindowContent{
			lineHeight:     16,
			baselineOffset: 12,
		},
	}
}

// Update refreshes window content each frame.
func (s *System) Update(world *ecs.World) {
	if world == nil {
		return
	}

	s.frameCounter++
	s.ensureWindows(world)
	s.refreshStats(world)
	s.refreshEntities(world)
	s.layoutWindows()
}

func (s *System) ensureWindows(world *ecs.World) {
	if s.statsComponent == nil || s.statsEntity == nil || !s.statsEntity.Has("Window") {
		entity := world.NewEntity()
		component := window.NewComponent("debug.stats", "Debug Stats", window.Bounds{
			X:      s.cfg.Margin,
			Y:      s.cfg.Margin,
			Width:  s.cfg.StatsWidth,
			Height: 0,
		}, s.statsContent)
		component.Layer = ecs.LayerDebug
		component.Order = 10
		component.Padding = 12
		component.TitleBarHeight = 26
		component.Background = color.RGBA{12, 16, 24, 220}
		component.Border = color.RGBA{90, 130, 200, 200}
		component.TitleBar = color.RGBA{30, 50, 90, 230}
		component.TitleColor = color.RGBA{230, 240, 255, 255}

		entity.Add(component)
		s.statsEntity = entity
		s.statsComponent = component
	}

	if s.entitiesComponent == nil || s.entitiesEntity == nil || !s.entitiesEntity.Has("Window") {
		entity := world.NewEntity()
		component := window.NewComponent("debug.entities", "Entity Diagnostics", window.Bounds{
			X:      s.cfg.ViewportWidth - s.cfg.EntitiesWidth - s.cfg.Margin,
			Y:      s.cfg.Margin,
			Width:  s.cfg.EntitiesWidth,
			Height: 0,
		}, s.entitiesContent)
		component.Layer = ecs.LayerDebug
		component.Order = 20
		component.Padding = 12
		component.TitleBarHeight = 26
		component.Background = color.RGBA{12, 16, 24, 220}
		component.Border = color.RGBA{120, 160, 235, 200}
		component.TitleBar = color.RGBA{36, 58, 110, 230}
		component.TitleColor = color.RGBA{235, 245, 255, 255}

		entity.Add(component)
		s.entitiesEntity = entity
		s.entitiesComponent = component
	}
}

func (s *System) refreshStats(world *ecs.World) {
	if s.statsComponent == nil || s.statsContent == nil {
		return
	}
	s.statsContent.Refresh(world)
	padding := s.statsComponent.Padding
	if padding < 0 {
		padding = 0
	}
	title := s.statsComponent.TitleBarHeight
	if title < 0 {
		title = 0
	}
	lines := len(s.statsContent.lines)
	total := title + padding*2 + lines*s.statsContent.lineHeight
	if total < title+padding*2+s.statsContent.lineHeight {
		total = title + padding*2 + s.statsContent.lineHeight
	}
	s.statsComponent.Bounds.Width = s.cfg.StatsWidth
	s.statsComponent.Bounds.Height = total
}

func (s *System) refreshEntities(world *ecs.World) {
	if s.entitiesComponent == nil || s.entitiesContent == nil {
		return
	}
	s.entitiesContent.Refresh(world, s.frameCounter, s.cfg.MaxEntities)
	padding := s.entitiesComponent.Padding
	if padding < 0 {
		padding = 0
	}
	title := s.entitiesComponent.TitleBarHeight
	if title < 0 {
		title = 0
	}
	lines := len(s.entitiesContent.lines)
	total := title + padding*2 + lines*s.entitiesContent.lineHeight
	if total < title+padding*2+s.entitiesContent.lineHeight {
		total = title + padding*2 + s.entitiesContent.lineHeight
	}
	s.entitiesComponent.Bounds.Width = s.cfg.EntitiesWidth
	s.entitiesComponent.Bounds.Height = total
}

func (s *System) layoutWindows() {
	if s.statsComponent != nil {
		s.statsComponent.Bounds.X = s.cfg.Margin
		s.statsComponent.Bounds.Y = s.cfg.Margin
	}
	if s.entitiesComponent != nil {
		x := s.cfg.ViewportWidth - s.cfg.EntitiesWidth - s.cfg.Margin
		if x < s.cfg.Margin {
			x = s.cfg.Margin
		}
		s.entitiesComponent.Bounds.X = x
		s.entitiesComponent.Bounds.Y = s.cfg.Margin
	}
}

type statsWindowContent struct {
	lines          []string
	lineHeight     int
	baselineOffset int
}

func (c *statsWindowContent) Refresh(world *ecs.World) {
	lines := make([]string, 0, 8)
	fps := platform.ActualFPS()
	lines = append(lines, fmt.Sprintf("FPS: %.0f", fps))

	manager := world.EntitiesManager()
	if manager != nil {
		lines = append(lines, fmt.Sprintf("Entities: %d", manager.Count()))
	}

	var cam *ecs.Camera
	if manager != nil {
		if _, comp := manager.FirstComponent("Camera"); comp != nil {
			cam, _ = comp.(*ecs.Camera)
		}
	}
	if cam != nil {
		targetScale := cam.TargetScale
		if targetScale <= 0 {
			targetScale = cam.Scale
		}
		minScale := cam.MinScale
		if minScale <= 0 {
			minScale = cam.Scale
		}
		maxScale := cam.MaxScale
		if maxScale <= 0 {
			maxScale = cam.Scale
		}
		lines = append(lines,
			fmt.Sprintf("Camera: (%.1f, %.1f)", cam.X, cam.Y),
			fmt.Sprintf("Scale: %.2f → %.2f", cam.Scale, targetScale),
			fmt.Sprintf("Bounds: %.2f – %.2f", minScale, maxScale),
		)
	}

	var playerPos *ecs.Position
	if manager != nil {
		manager.ForEach(func(entity *ecs.Entity) {
			if playerPos != nil {
				return
			}
			if entity.Has("CameraTarget") {
				if pos, ok := entity.Get("Position").(*ecs.Position); ok {
					playerPos = pos
				}
			}
		})
	}
	if playerPos != nil {
		lines = append(lines, fmt.Sprintf("Player: (%.1f, %.1f)", playerPos.X, playerPos.Y))
	}

	c.lines = lines
}

func (c *statsWindowContent) Draw(_ *ecs.World, canvas *platform.Image, bounds window.Bounds) {
	if canvas == nil || len(c.lines) == 0 {
		return
	}
	baseline := bounds.Y + c.baselineOffset
	textX := bounds.X
	if textX < 0 {
		textX = 0
	}
	for _, line := range c.lines {
		platform.DrawText(canvas, line, basicfont.Face7x13, textX, baseline, color.White)
		baseline += c.lineHeight
	}
}

type entitiesWindowContent struct {
	lines           []string
	lineHeight      int
	baselineOffset  int
	lastLoggedFrame int
}

func (c *entitiesWindowContent) Refresh(world *ecs.World, frame int, maxLines int) {
	lines := make([]string, 0, maxLines)
	manager := world.EntitiesManager()
	if manager != nil {
		collected := make([]string, 0, manager.Count())
		manager.ForEach(func(e *ecs.Entity) {
			pos, _ := e.Get("Position").(*ecs.Position)
			sprite, _ := e.Get("Sprite").(*ecs.Sprite)
			if pos == nil {
				return
			}
			spriteSize := "?"
			if sprite != nil && sprite.Image != nil {
				w, h := sprite.NativeSize()
				scale := sprite.PixelScale()
				spriteSize = fmt.Sprintf("%dx%d x%.2f", w, h, scale)
			}
			collected = append(collected, fmt.Sprintf("%d @ (%.0f, %.0f) %s", e.ID, pos.X, pos.Y, spriteSize))
		})
		sort.Strings(collected)
		if len(collected) > maxLines {
			collected = collected[:maxLines]
		}
		lines = append(lines, collected...)
		if frame-c.lastLoggedFrame >= 15 {
			for _, line := range collected {
				fmt.Println("[ENTITY DEBUG]", line)
			}
			c.lastLoggedFrame = frame
		}
	}
	if len(lines) == 0 {
		lines = append(lines, "No entities with Position component")
	}
	c.lines = lines
}

func (c *entitiesWindowContent) Draw(_ *ecs.World, canvas *platform.Image, bounds window.Bounds) {
	if canvas == nil || len(c.lines) == 0 {
		return
	}
	baseline := bounds.Y + c.baselineOffset
	textX := bounds.X
	if textX < 0 {
		textX = 0
	}
	for _, line := range c.lines {
		platform.DrawText(canvas, line, basicfont.Face7x13, textX, baseline, color.RGBA{200, 220, 255, 255})
		baseline += c.lineHeight
	}
}
