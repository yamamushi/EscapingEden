package renderer

// OptimizedWindow represents a window that renders efficiently to a DeltaRenderer
type OptimizedWindow struct {
	id      string
	x, y    int
	width   int
	height  int
	visible bool

	// Content management
	content   []string
	scrollPos int
	maxScroll int

	// Style information
	borderStyle uint32
	textStyle   uint32
	titleStyle  uint32

	// Change tracking
	contentChanged  bool
	positionChanged bool
	styleChanged    bool

	// Border settings
	hasBorder bool
	title     string
}

// NewOptimizedWindow creates a new optimized window
func NewOptimizedWindow(id string, x, y, width, height int) *OptimizedWindow {
	return &OptimizedWindow{
		id:      id,
		x:       x,
		y:       y,
		width:   width,
		height:  height,
		visible: true,
		content: make([]string, 0),
	}
}

// SetContent updates the window content efficiently
func (w *OptimizedWindow) SetContent(lines []string) {
	// Only update if content actually changed
	if !w.contentEqual(lines) {
		w.content = make([]string, len(lines))
		copy(w.content, lines)
		w.contentChanged = true
		w.updateScrollLimits()
	}
}

// AppendLine adds a single line efficiently (for chat/logs)
func (w *OptimizedWindow) AppendLine(line string) {
	w.content = append(w.content, line)
	w.contentChanged = true
	w.updateScrollLimits()

	// Auto-scroll to bottom if we were already at bottom
	if w.scrollPos == w.maxScroll {
		w.scrollPos = w.maxScroll + 1
	}
}

// SetPosition moves the window
func (w *OptimizedWindow) SetPosition(x, y int) {
	if w.x != x || w.y != y {
		w.x = x
		w.y = y
		w.positionChanged = true
	}
}

// SetSize resizes the window
func (w *OptimizedWindow) SetSize(width, height int) {
	if w.width != width || w.height != height {
		w.width = width
		w.height = height
		w.positionChanged = true // Size change affects rendering
		w.updateScrollLimits()
	}
}

// SetVisible shows/hides the window
func (w *OptimizedWindow) SetVisible(visible bool) {
	if w.visible != visible {
		w.visible = visible
		w.positionChanged = true // Visibility change affects rendering
	}
}

// SetBorder configures window border
func (w *OptimizedWindow) SetBorder(enabled bool, title string, style uint32) {
	if w.hasBorder != enabled || w.title != title || w.borderStyle != style {
		w.hasBorder = enabled
		w.title = title
		w.borderStyle = style
		w.styleChanged = true
	}
}

// SetTextStyle sets the text style
func (w *OptimizedWindow) SetTextStyle(style uint32) {
	if w.textStyle != style {
		w.textStyle = style
		w.styleChanged = true
	}
}

// Scroll moves the content view
func (w *OptimizedWindow) Scroll(delta int) {
	newPos := w.scrollPos + delta
	if newPos < 0 {
		newPos = 0
	}
	if newPos > w.maxScroll {
		newPos = w.maxScroll
	}

	if w.scrollPos != newPos {
		w.scrollPos = newPos
		w.contentChanged = true
	}
}

// RenderTo efficiently renders the window to a DeltaRenderer
func (w *OptimizedWindow) RenderTo(renderer *DeltaRenderer) {
	if !w.visible {
		return
	}

	// Only render if something changed
	if !w.contentChanged && !w.positionChanged && !w.styleChanged {
		return
	}

	// Clear the window area first if position changed
	if w.positionChanged {
		renderer.ClearRegion(w.x, w.y, w.width, w.height)
	}

	// Render border if enabled
	if w.hasBorder {
		w.renderBorder(renderer)
	}

	// Render content
	w.renderContent(renderer)

	// Reset change flags
	w.contentChanged = false
	w.positionChanged = false
	w.styleChanged = false
}

// renderBorder draws the window border
func (w *OptimizedWindow) renderBorder(renderer *DeltaRenderer) {
	// Top border
	renderer.SetCell(w.x, w.y, '┌', w.borderStyle)
	for x := w.x + 1; x < w.x+w.width-1; x++ {
		renderer.SetCell(x, w.y, '─', w.borderStyle)
	}
	renderer.SetCell(w.x+w.width-1, w.y, '┐', w.borderStyle)

	// Title in top border
	if w.title != "" && len(w.title)+4 < w.width {
		titleX := w.x + (w.width-len(w.title))/2
		renderer.SetText(titleX, w.y, w.title, w.titleStyle)
	}

	// Side borders
	for y := w.y + 1; y < w.y+w.height-1; y++ {
		renderer.SetCell(w.x, y, '│', w.borderStyle)
		renderer.SetCell(w.x+w.width-1, y, '│', w.borderStyle)
	}

	// Bottom border
	renderer.SetCell(w.x, w.y+w.height-1, '└', w.borderStyle)
	for x := w.x + 1; x < w.x+w.width-1; x++ {
		renderer.SetCell(x, w.y+w.height-1, '─', w.borderStyle)
	}
	renderer.SetCell(w.x+w.width-1, w.y+w.height-1, '┘', w.borderStyle)
}

// renderContent draws the window content
func (w *OptimizedWindow) renderContent(renderer *DeltaRenderer) {
	contentX := w.x
	contentY := w.y
	contentWidth := w.width
	contentHeight := w.height

	// Adjust for border
	if w.hasBorder {
		contentX++
		contentY++
		contentWidth -= 2
		contentHeight -= 2
	}

	// Render visible lines
	visibleLines := w.getVisibleLines(contentHeight)
	for i, line := range visibleLines {
		y := contentY + i
		if y >= contentY+contentHeight {
			break
		}

		// Truncate line if too long
		if len(line) > contentWidth {
			line = line[:contentWidth]
		}

		// Render the line
		if len(line) > 0 {
			renderer.SetText(contentX, y, line, w.textStyle)
		}

		// Clear rest of line if needed
		if len(line) < contentWidth {
			for x := contentX + len(line); x < contentX+contentWidth; x++ {
				renderer.SetCell(x, y, ' ', 0)
			}
		}
	}

	// Clear any remaining lines in the content area
	for y := contentY + len(visibleLines); y < contentY+contentHeight; y++ {
		for x := contentX; x < contentX+contentWidth; x++ {
			renderer.SetCell(x, y, ' ', 0)
		}
	}
}

// getVisibleLines returns the lines that should be visible given current scroll
func (w *OptimizedWindow) getVisibleLines(maxLines int) []string {
	if len(w.content) == 0 {
		return []string{}
	}

	startLine := len(w.content) - maxLines - w.scrollPos
	if startLine < 0 {
		startLine = 0
	}

	endLine := startLine + maxLines
	if endLine > len(w.content) {
		endLine = len(w.content)
	}

	return w.content[startLine:endLine]
}

// updateScrollLimits recalculates scroll boundaries
func (w *OptimizedWindow) updateScrollLimits() {
	contentHeight := w.height
	if w.hasBorder {
		contentHeight -= 2
	}

	w.maxScroll = len(w.content) - contentHeight
	if w.maxScroll < 0 {
		w.maxScroll = 0
	}
}

// contentEqual efficiently compares content arrays
func (w *OptimizedWindow) contentEqual(other []string) bool {
	if len(w.content) != len(other) {
		return false
	}

	for i, line := range w.content {
		if line != other[i] {
			return false
		}
	}

	return true
}

// GetID returns the window identifier
func (w *OptimizedWindow) GetID() string {
	return w.id
}

// IsVisible returns whether the window is visible
func (w *OptimizedWindow) IsVisible() bool {
	return w.visible
}

// GetBounds returns the window boundaries
func (w *OptimizedWindow) GetBounds() (x, y, width, height int) {
	return w.x, w.y, w.width, w.height
}

// HasChanges returns whether the window needs re-rendering
func (w *OptimizedWindow) HasChanges() bool {
	return w.contentChanged || w.positionChanged || w.styleChanged
}

// WindowManager efficiently manages multiple windows
type WindowManager struct {
	windows  []*OptimizedWindow
	renderer *DeltaRenderer
	zOrder   []string // Window IDs in z-order (front to back)
}

// NewWindowManager creates a new window manager
func NewWindowManager(renderer *DeltaRenderer) *WindowManager {
	return &WindowManager{
		windows:  make([]*OptimizedWindow, 0),
		renderer: renderer,
		zOrder:   make([]string, 0),
	}
}

// AddWindow adds a window to the manager
func (wm *WindowManager) AddWindow(window *OptimizedWindow) {
	wm.windows = append(wm.windows, window)
	wm.zOrder = append(wm.zOrder, window.GetID())
}

// GetWindow retrieves a window by ID
func (wm *WindowManager) GetWindow(id string) *OptimizedWindow {
	for _, window := range wm.windows {
		if window.GetID() == id {
			return window
		}
	}
	return nil
}

// BringToFront moves a window to the front
func (wm *WindowManager) BringToFront(id string) {
	// Remove from current position
	for i, windowID := range wm.zOrder {
		if windowID == id {
			wm.zOrder = append(wm.zOrder[:i], wm.zOrder[i+1:]...)
			break
		}
	}

	// Add to front
	wm.zOrder = append([]string{id}, wm.zOrder...)
}

// RenderAll efficiently renders all windows that have changes
func (wm *WindowManager) RenderAll() []byte {
	// Render windows in z-order (back to front)
	for i := len(wm.zOrder) - 1; i >= 0; i-- {
		window := wm.GetWindow(wm.zOrder[i])
		if window != nil && window.HasChanges() {
			window.RenderTo(wm.renderer)
		}
	}

	return wm.renderer.Render()
}

// GetStats returns rendering statistics
func (wm *WindowManager) GetStats() WindowManagerStats {
	stats := WindowManagerStats{
		TotalWindows:       len(wm.windows),
		VisibleWindows:     0,
		WindowsWithChanges: 0,
	}

	for _, window := range wm.windows {
		if window.IsVisible() {
			stats.VisibleWindows++
		}
		if window.HasChanges() {
			stats.WindowsWithChanges++
		}
	}

	stats.RendererStats = wm.renderer.GetStats()

	return stats
}

// WindowManagerStats provides performance metrics
type WindowManagerStats struct {
	TotalWindows       int
	VisibleWindows     int
	WindowsWithChanges int
	RendererStats      RendererStats
}
