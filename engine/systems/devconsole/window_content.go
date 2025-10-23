package devconsole

import (
	"image/color"

	"golang.org/x/image/font/basicfont"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/ui/window"
)

const (
	maxLogEntries     = 12
	maxHistoryStored  = 32
	consoleLineHeight = 16
)

type consoleWindowContent struct {
	state             *ConsoleState
	baselineOffset    int
	promptBaselinePad int
}

func newConsoleWindowContent(state *ConsoleState) *consoleWindowContent {
	return &consoleWindowContent{
		state:             state,
		baselineOffset:    12,
		promptBaselinePad: 6,
	}
}

func (c *consoleWindowContent) Draw(_ *ecs.World, canvas *platform.Image, bounds window.Bounds) {
	if c == nil || c.state == nil || canvas == nil {
		return
	}
	if !c.state.Open {
		return
	}

	textX := bounds.X
	if textX < 0 {
		textX = 0
	}

	availableHeight := bounds.Height
	if availableHeight <= consoleLineHeight {
		prompt := c.state.promptString()
		baseline := bounds.Y + availableHeight - c.promptBaselinePad
		if baseline < bounds.Y+c.baselineOffset {
			baseline = bounds.Y + c.baselineOffset
		}
		platform.DrawText(canvas, prompt, basicfont.Face7x13, textX, baseline, color.White)
		return
	}

	// Reserve space for the prompt line at the bottom.
	logHeight := availableHeight - consoleLineHeight
	maxLines := logHeight / consoleLineHeight
	if maxLines < 1 {
		maxLines = 1
	}

	lines := c.state.LogMessages
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}

	baseline := bounds.Y + c.baselineOffset
	for _, line := range lines {
		platform.DrawText(canvas, line, basicfont.Face7x13, textX, baseline, color.White)
		baseline += consoleLineHeight
	}

	prompt := c.state.promptString()
	promptBaseline := bounds.Y + availableHeight - c.promptBaselinePad
	if promptBaseline < baseline {
		promptBaseline = baseline
	}
	platform.DrawText(canvas, prompt, basicfont.Face7x13, textX, promptBaseline, color.White)
}
