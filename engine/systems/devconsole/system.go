package devconsole

import (
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"

	"golang.org/x/image/font/basicfont"
	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/systems/actor"
	"rp-go/engine/world"
)

const (
	maxLogEntries    = 12
	maxHistoryStored = 32
)

var consoleOpen atomic.Bool

// IsOpen returns true when the developer console overlay is active.
func IsOpen() bool {
	return consoleOpen.Load()
}

// System implements a developer console overlay that can spawn, remove, and
// move actors at runtime. It also provides a simple command log while
// preventing the standard input system from processing player controls.
type System struct {
	actorRegistry *actor.Registry
	actorCreator  *world.ActorCreator

	open       bool
	justOpened bool

	inputBuffer string
	history     []string
	historyIdx  int

	logMessages []string
	cursorTick  int
}

// NewSystem creates a developer console bound to the shared actor registry.
func NewSystem(registry *actor.Registry) *System {
	consoleOpen.Store(false)
	return &System{
		actorRegistry: registry,
		historyIdx:    -1,
	}
}

// Layer ensures the console draws last in the overlay stack.
func (s *System) Layer() ecs.DrawLayer { return ecs.LayerConsole }

func (s *System) Update(w *ecs.World) {
	if platform.IsKeyJustPressed(platform.KeyF12) {
		s.open = !s.open
		s.justOpened = s.open
		if !s.open {
			s.historyIdx = -1
		}
	}

	consoleOpen.Store(s.open)

	if !s.open {
		return
	}

	s.cursorTick++

	if s.justOpened {
		s.log("Developer console opened. Type 'help' for commands.")
		s.justOpened = false
	}

	// Capture text input for the current frame.
	for _, char := range platform.InputChars() {
		if char == '\r' || char == '\n' {
			continue
		}
		s.inputBuffer += string(char)
	}

	if platform.IsKeyJustPressed(platform.KeyBackspace) && len(s.inputBuffer) > 0 {
		s.inputBuffer = s.inputBuffer[:len(s.inputBuffer)-1]
	}

	if platform.IsKeyJustPressed(platform.KeyEscape) {
		s.open = false
		consoleOpen.Store(false)
		s.historyIdx = -1
		return
	}

	// Command history navigation.
	if platform.IsKeyJustPressed(platform.KeyArrowUp) {
		if len(s.history) > 0 {
			if s.historyIdx == -1 {
				s.historyIdx = len(s.history) - 1
			} else if s.historyIdx > 0 {
				s.historyIdx--
			}
			if s.historyIdx >= 0 && s.historyIdx < len(s.history) {
				s.inputBuffer = s.history[s.historyIdx]
			}
		}
	} else if platform.IsKeyJustPressed(platform.KeyArrowDown) {
		if s.historyIdx != -1 {
			if s.historyIdx < len(s.history)-1 {
				s.historyIdx++
				s.inputBuffer = s.history[s.historyIdx]
			} else {
				s.historyIdx = -1
				s.inputBuffer = ""
			}
		}
	}

	if platform.IsKeyJustPressed(platform.KeyEnter) {
		command := strings.TrimSpace(s.inputBuffer)
		if command != "" {
			s.executeCommand(w, command)
			s.pushHistory(command)
		}
		s.inputBuffer = ""
		s.historyIdx = -1
	}
}

func (s *System) Draw(w *ecs.World, screen *platform.Image) {
	if !s.open || screen == nil {
		return
	}

	bounds := screen.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width == 0 || height == 0 {
		return
	}

	consoleHeight := height / 3
	if consoleHeight < 180 {
		consoleHeight = 180
	}

	overlay := platform.NewImage(width, consoleHeight)
	overlay.FillRect(0, 0, width, consoleHeight, color.RGBA{0, 0, 0, 200})

	y := 24
	visible := s.logMessages
	if len(visible) > maxLogEntries {
		visible = visible[len(visible)-maxLogEntries:]
	}
	for _, line := range visible {
		platform.DrawText(overlay, line, basicfont.Face7x13, 12, y, color.White)
		y += 16
	}

	prompt := "> " + s.inputBuffer
	if (s.cursorTick/20)%2 == 0 {
		prompt += "_"
	}
	platform.DrawText(overlay, prompt, basicfont.Face7x13, 12, consoleHeight-16, color.White)

	op := platform.NewDrawImageOptions()
	op.Translate(0, float64(height-consoleHeight))
	screen.DrawImage(overlay, op)
}

func (s *System) executeCommand(w *ecs.World, input string) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return
	}

	switch strings.ToLower(fields[0]) {
	case "help":
		s.log("Commands: help, spawn <template> [x y], remove <actorID>, move <actorID> <x y>, list")
	case "spawn":
		s.handleSpawn(w, fields)
	case "remove", "rm":
		s.handleRemove(w, fields)
	case "move", "teleport":
		s.handleMove(w, fields)
	case "list":
		s.handleList(w)
	default:
		s.log(fmt.Sprintf("Unknown command: %s", fields[0]))
	}
}

func (s *System) handleSpawn(w *ecs.World, fields []string) {
	if len(fields) < 2 {
		s.log("Usage: spawn <template> [x y]")
		return
	}
	if s.actorCreator == nil {
		db := data.LoadActorDatabase("engine/data/actors.json")
		s.actorCreator = world.NewActorCreator(db)
	}

	x, y := 0.0, 0.0
	if len(fields) >= 3 {
		if parsed, err := strconv.ParseFloat(fields[2], 64); err == nil {
			x = parsed
		} else {
			s.log(fmt.Sprintf("Invalid X coordinate: %q", fields[2]))
			return
		}
	}
	if len(fields) >= 4 {
		if parsed, err := strconv.ParseFloat(fields[3], 64); err == nil {
			y = parsed
		} else {
			s.log(fmt.Sprintf("Invalid Y coordinate: %q", fields[3]))
			return
		}
	}

	entity, err := s.actorCreator.Spawn(w, fields[1], ecs.Position{X: x, Y: y})
	if err != nil {
		s.log(err.Error())
		if templates := s.listTemplates(); len(templates) > 0 {
			s.log("Templates: " + strings.Join(templates, ", "))
		}
		return
	}

	if actorComp, ok := entity.Get("Actor").(*ecs.Actor); ok && actorComp != nil {
		s.log(fmt.Sprintf("Spawned %s at (%.1f, %.1f)", actorComp.ID, x, y))
	} else {
		s.log(fmt.Sprintf("Spawned entity %d", entity.ID))
	}
}

func (s *System) handleRemove(w *ecs.World, fields []string) {
	if len(fields) < 2 {
		s.log("Usage: remove <actorID>")
		return
	}
	target := s.findActorByID(w, fields[1])
	if target == nil {
		s.log(fmt.Sprintf("Actor %q not found", fields[1]))
		return
	}
	w.RemoveEntity(target)
	s.log(fmt.Sprintf("Removed actor %s", fields[1]))
}

func (s *System) handleMove(w *ecs.World, fields []string) {
	if len(fields) < 4 {
		s.log("Usage: move <actorID> <x> <y>")
		return
	}
	x, errX := strconv.ParseFloat(fields[2], 64)
	y, errY := strconv.ParseFloat(fields[3], 64)
	if errX != nil || errY != nil {
		s.log("Invalid coordinates. Expected numbers for x and y.")
		return
	}
	target := s.findActorByID(w, fields[1])
	if target == nil {
		s.log(fmt.Sprintf("Actor %q not found", fields[1]))
		return
	}
	pos, ok := target.Get("Position").(*ecs.Position)
	if !ok {
		pos = &ecs.Position{}
		target.Add(pos)
	}
	pos.X = x
	pos.Y = y
	s.log(fmt.Sprintf("Moved %s to (%.1f, %.1f)", fields[1], x, y))
}

func (s *System) handleList(w *ecs.World) {
	entities := s.collectActors(w)
	if len(entities) == 0 {
		s.log("No actors registered.")
		return
	}
	for _, entry := range entities {
		s.log(entry)
	}
}

func (s *System) listTemplates() []string {
	if s.actorCreator == nil {
		return nil
	}
	names := s.actorCreator.Templates()
	if len(names) == 0 {
		return nil
	}
	sort.Strings(names)
	return names
}

func (s *System) pushHistory(command string) {
	if len(s.history) >= maxHistoryStored {
		s.history = s.history[1:]
	}
	s.history = append(s.history, command)
}

func (s *System) log(message string) {
	if message == "" {
		return
	}
	s.logMessages = append(s.logMessages, message)
	if len(s.logMessages) > maxLogEntries {
		s.logMessages = s.logMessages[len(s.logMessages)-maxLogEntries:]
	}
}

func (s *System) findActorByID(w *ecs.World, id string) *ecs.Entity {
	if s.actorRegistry != nil {
		if entity, ok := s.actorRegistry.FindByID(id); ok {
			return entity
		}
	}
	for _, entity := range w.Entities {
		if actorComp, ok := entity.Get("Actor").(*ecs.Actor); ok && actorComp != nil && actorComp.ID == id {
			return entity
		}
	}
	return nil
}

func (s *System) collectActors(w *ecs.World) []string {
	var entities []*ecs.Entity
	if s.actorRegistry != nil {
		entities = s.actorRegistry.All()
	}
	if len(entities) == 0 {
		entities = make([]*ecs.Entity, 0, len(w.Entities))
		for _, entity := range w.Entities {
			if actorComp, ok := entity.Get("Actor").(*ecs.Actor); ok && actorComp != nil {
				entities = append(entities, entity)
			}
		}
	}

	if len(entities) == 0 {
		return nil
	}

	descriptions := make([]string, 0, len(entities))
	for _, entity := range entities {
		actorComp, _ := entity.Get("Actor").(*ecs.Actor)
		pos, _ := entity.Get("Position").(*ecs.Position)
		if actorComp == nil {
			continue
		}
		if pos != nil {
			descriptions = append(descriptions, fmt.Sprintf("%s (%.1f, %.1f)", actorComp.ID, pos.X, pos.Y))
		} else {
			descriptions = append(descriptions, fmt.Sprintf("%s (no position)", actorComp.ID))
		}
	}

	sort.Strings(descriptions)
	if len(descriptions) > maxLogEntries {
		return descriptions[:maxLogEntries]
	}
	return descriptions
}
