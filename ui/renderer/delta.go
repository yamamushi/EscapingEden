package renderer

import (
	"strings"

	"github.com/yamamushi/EscapingEden/terminals"
)

// DeltaRenderer efficiently tracks and renders only changes to the terminal
type DeltaRenderer struct {
	width, height int
	terminal      terminals.TerminalType

	// Current state - only store what's actually displayed
	cells [][]Cell

	// Change tracking - much more efficient than full comparison
	dirtyRegions []Region
	lastCursor   Position

	// Output optimization
	outputBuffer strings.Builder
}

// Cell represents a single terminal cell with minimal data
type Cell struct {
	Char  rune
	Style uint32 // Packed color/style info instead of escape strings
}

// Region represents a rectangular area that needs updating
type Region struct {
	X, Y, Width, Height int
}

// Position represents a cursor position
type Position struct {
	X, Y int
}

// NewDeltaRenderer creates an optimized renderer
func NewDeltaRenderer(width, height int, terminal terminals.TerminalType) *DeltaRenderer {
	dr := &DeltaRenderer{
		width:    width,
		height:   height,
		terminal: terminal,
		cells:    make([][]Cell, height),
	}

	// Initialize cell grid
	for y := 0; y < height; y++ {
		dr.cells[y] = make([]Cell, width)
	}

	return dr
}

// SetCell efficiently updates a single cell and marks region as dirty
func (dr *DeltaRenderer) SetCell(x, y int, char rune, style uint32) {
	if x < 0 || x >= dr.width || y < 0 || y >= dr.height {
		return
	}

	// Only update if actually changed
	if dr.cells[y][x].Char != char || dr.cells[y][x].Style != style {
		dr.cells[y][x].Char = char
		dr.cells[y][x].Style = style
		dr.markDirty(x, y, 1, 1)
	}
}

// SetText efficiently updates a string of text
func (dr *DeltaRenderer) SetText(x, y int, text string, style uint32) {
	if y < 0 || y >= dr.height {
		return
	}

	changed := false
	for i, char := range text {
		if x+i >= dr.width {
			break
		}

		if dr.cells[y][x+i].Char != char || dr.cells[y][x+i].Style != style {
			dr.cells[y][x+i].Char = char
			dr.cells[y][x+i].Style = style
			changed = true
		}
	}

	if changed {
		dr.markDirty(x, y, len(text), 1)
	}
}

// ClearRegion efficiently clears a rectangular area
func (dr *DeltaRenderer) ClearRegion(x, y, width, height int) {
	for dy := 0; dy < height; dy++ {
		if y+dy >= dr.height {
			break
		}
		for dx := 0; dx < width; dx++ {
			if x+dx >= dr.width {
				break
			}
			dr.cells[y+dy][x+dx] = Cell{Char: ' ', Style: 0}
		}
	}
	dr.markDirty(x, y, width, height)
}

// markDirty adds a region to the dirty list, merging with adjacent regions
func (dr *DeltaRenderer) markDirty(x, y, width, height int) {
	newRegion := Region{X: x, Y: y, Width: width, Height: height}

	// Try to merge with existing dirty regions
	merged := false
	for i := range dr.dirtyRegions {
		if dr.canMerge(dr.dirtyRegions[i], newRegion) {
			dr.dirtyRegions[i] = dr.mergeRegions(dr.dirtyRegions[i], newRegion)
			merged = true
			break
		}
	}

	if !merged {
		dr.dirtyRegions = append(dr.dirtyRegions, newRegion)
	}
}

// canMerge checks if two regions can be efficiently merged
func (dr *DeltaRenderer) canMerge(a, b Region) bool {
	// Merge if regions overlap or are adjacent
	return !(a.X+a.Width < b.X || b.X+b.Width < a.X ||
		a.Y+a.Height < b.Y || b.Y+b.Height < a.Y)
}

// mergeRegions combines two regions into one
func (dr *DeltaRenderer) mergeRegions(a, b Region) Region {
	minX := min(a.X, b.X)
	minY := min(a.Y, b.Y)
	maxX := max(a.X+a.Width, b.X+b.Width)
	maxY := max(a.Y+a.Height, b.Y+b.Height)

	return Region{
		X:      minX,
		Y:      minY,
		Width:  maxX - minX,
		Height: maxY - minY,
	}
}

// Render generates the minimal terminal output for all changes
func (dr *DeltaRenderer) Render() []byte {
	if len(dr.dirtyRegions) == 0 {
		return nil // No changes
	}

	dr.outputBuffer.Reset()

	// Sort regions by position for optimal cursor movement
	dr.sortRegions()

	currentPos := dr.lastCursor

	for _, region := range dr.dirtyRegions {
		dr.renderRegion(region, &currentPos)
	}

	// Clear dirty regions after rendering
	dr.dirtyRegions = dr.dirtyRegions[:0]
	dr.lastCursor = currentPos

	return []byte(dr.outputBuffer.String())
}

// renderRegion efficiently renders a single dirty region
func (dr *DeltaRenderer) renderRegion(region Region, currentPos *Position) {
	for y := region.Y; y < region.Y+region.Height && y < dr.height; y++ {
		// Find continuous runs of characters with same style
		x := region.X
		for x < region.X+region.Width && x < dr.width {
			startX := x
			currentStyle := dr.cells[y][x].Style

			// Find end of run with same style
			for x < region.X+region.Width && x < dr.width &&
				dr.cells[y][x].Style == currentStyle {
				x++
			}

			// Render this run efficiently
			dr.renderRun(startX, y, x-startX, currentStyle, currentPos)
		}
	}
}

// renderRun renders a continuous run of characters with the same style
func (dr *DeltaRenderer) renderRun(x, y, length int, style uint32, currentPos *Position) {
	// Move cursor only if necessary
	if currentPos.X != x || currentPos.Y != y {
		dr.outputBuffer.WriteString(dr.terminal.MoveCursor(x, y))
		currentPos.X = x
		currentPos.Y = y
	}

	// Apply style if changed
	if style != 0 {
		dr.outputBuffer.WriteString(dr.styleToEscape(style))
	}

	// Write characters
	for i := 0; i < length; i++ {
		dr.outputBuffer.WriteRune(dr.cells[y][x+i].Char)
	}

	// Reset style if needed
	if style != 0 {
		dr.outputBuffer.WriteString(dr.terminal.Reset())
	}

	currentPos.X += length
}

// sortRegions sorts dirty regions for optimal rendering order
func (dr *DeltaRenderer) sortRegions() {
	// Simple sort by Y then X for optimal cursor movement
	for i := 0; i < len(dr.dirtyRegions)-1; i++ {
		for j := i + 1; j < len(dr.dirtyRegions); j++ {
			a, b := dr.dirtyRegions[i], dr.dirtyRegions[j]
			if a.Y > b.Y || (a.Y == b.Y && a.X > b.X) {
				dr.dirtyRegions[i], dr.dirtyRegions[j] = b, a
			}
		}
	}
}

// styleToEscape converts packed style to escape sequence
func (dr *DeltaRenderer) styleToEscape(style uint32) string {
	// Extract color and style info from packed uint32
	// This is much more efficient than storing full escape strings

	if style == 0 {
		return ""
	}

	// Example: extract foreground/background colors and attributes
	fg := (style >> 0) & 0xFF
	bg := (style >> 8) & 0xFF
	attrs := (style >> 16) & 0xFFFF

	var escape strings.Builder

	if fg != 0 {
		escape.WriteString("\033[38;5;")
		escape.WriteString(string(rune('0' + fg)))
		escape.WriteString("m")
	}

	if bg != 0 {
		escape.WriteString("\033[48;5;")
		escape.WriteString(string(rune('0' + bg)))
		escape.WriteString("m")
	}

	if attrs&1 != 0 { // Bold
		escape.WriteString("\033[1m")
	}

	return escape.String()
}

// Resize efficiently handles terminal resize
func (dr *DeltaRenderer) Resize(newWidth, newHeight int) {
	if newWidth == dr.width && newHeight == dr.height {
		return
	}

	// Create new cell grid
	newCells := make([][]Cell, newHeight)
	for y := 0; y < newHeight; y++ {
		newCells[y] = make([]Cell, newWidth)

		// Copy existing data if within bounds
		if y < dr.height {
			copyWidth := min(dr.width, newWidth)
			copy(newCells[y][:copyWidth], dr.cells[y][:copyWidth])
		}
	}

	dr.cells = newCells
	dr.width = newWidth
	dr.height = newHeight

	// Mark entire screen as dirty after resize
	dr.dirtyRegions = []Region{{X: 0, Y: 0, Width: newWidth, Height: newHeight}}
}

// GetStats returns rendering statistics for monitoring
func (dr *DeltaRenderer) GetStats() RendererStats {
	totalDirtyCells := 0
	for _, region := range dr.dirtyRegions {
		totalDirtyCells += region.Width * region.Height
	}

	return RendererStats{
		DirtyRegions:    len(dr.dirtyRegions),
		DirtyCells:      totalDirtyCells,
		TotalCells:      dr.width * dr.height,
		EfficiencyRatio: float64(totalDirtyCells) / float64(dr.width*dr.height),
	}
}

// RendererStats provides performance metrics
type RendererStats struct {
	DirtyRegions    int
	DirtyCells      int
	TotalCells      int
	EfficiencyRatio float64 // Lower is better (less data sent)
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
