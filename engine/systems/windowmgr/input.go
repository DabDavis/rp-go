package windowmgr

import (
	"rp-go/engine/events"
	"rp-go/engine/platform"
	"rp-go/engine/ui/window"
)

// Internal drag state
var (
	dragging     *window.Component
	dragOffsetX  int
	dragOffsetY  int
	lastDragX    int
	lastDragY    int
)

// UpdateWindowInteractions updates mouse-based interactions for all windows
// and emits relevant window events to the event bus.
func UpdateWindowInteractions(windows []*window.Component, bus *events.TypedBus) {
	mx, my := platform.MousePosition()
	leftDown := platform.IsMouseButtonPressed(platform.MouseButtonLeft)

	// --- Release drag if button is up
	if !leftDown && dragging != nil {
		dragging = nil
		return
	}

	for _, win := range windows {
		if win == nil || win.Locked || !win.Visible {
			continue
		}

		// --- Button hitboxes
		closeBtn := window.Bounds{
			X: win.Bounds.X + win.Bounds.Width - 20,
			Y: win.Bounds.Y + 4,
			Width:  12,
			Height: 12,
		}
		minBtn := window.Bounds{
			X: win.Bounds.X + win.Bounds.Width - 38,
			Y: win.Bounds.Y + 4,
			Width:  12,
			Height: 12,
		}
		titleBar := window.Bounds{
			X: win.Bounds.X,
			Y: win.Bounds.Y,
			Width:  win.Bounds.Width,
			Height: win.TitleBarHeight,
		}

		// --- Close button
		if win.Closable && !leftDown && closeBtn.Contains(mx, my) {
			win.Visible = false
			if bus != nil {
				events.Queue(bus, events.WindowClosedEvent{ID: win.ID})
			}
			continue
		}

		// --- Minimize button
		if win.Movable && !leftDown && minBtn.Contains(mx, my) {
			win.Minimized = !win.Minimized
			if bus != nil {
				if win.Minimized {
					events.Queue(bus, events.WindowMinimizedEvent{ID: win.ID})
				} else {
					events.Queue(bus, events.WindowRestoredEvent{ID: win.ID})
				}
			}
			continue
		}

		// --- Begin drag
		if leftDown && titleBar.Contains(mx, my) && win.Movable {
			if dragging == nil {
				dragging = win
				dragOffsetX = mx - win.Bounds.X
				dragOffsetY = my - win.Bounds.Y
				lastDragX, lastDragY = win.Bounds.X, win.Bounds.Y
			}
		}

		// --- Update drag position
		if dragging == win && leftDown {
			newX := mx - dragOffsetX
			newY := my - dragOffsetY

			// Prevent negative coordinates
			if newX < 0 {
				newX = 0
			}
			if newY < 0 {
				newY = 0
			}

			win.Bounds.X = newX
			win.Bounds.Y = newY

			// Emit move event if changed
			if bus != nil && (newX != lastDragX || newY != lastDragY) {
				events.Queue(bus, events.WindowMovedEvent{
					ID: win.ID,
					X:  newX,
					Y:  newY,
				})
				lastDragX, lastDragY = newX, newY
			}
		}
	}
}

