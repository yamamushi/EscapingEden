package ui

import (
	"sync"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/terminals"
	"github.com/yamamushi/EscapingEden/ui/renderer"
)

// OptimizedConsole provides efficient terminal rendering with minimal data transmission
type OptimizedConsole struct {
	// Basic properties
	Width  int
	Height int
	Log    logging.LoggerType

	// Terminal interface
	Terminal terminals.TerminalType

	// Rendering system
	deltaRenderer *renderer.DeltaRenderer
	windowManager *renderer.WindowManager

	// State management
	mutex            sync.RWMutex
	initialized      bool
	resizeActive     bool
	forceFullRefresh bool
	lastRenderTime   time.Time

	// Performance tracking
	renderStats RenderingStats

	// Connection management
	connectionID string
	sendChannel  chan []byte
}

// RenderingStats tracks performance metrics
type RenderingStats struct {
	TotalRenders      int64
	BytesSent         int64
	AverageRenderTime time.Duration
	LastRenderSize    int
	EfficiencyRatio   float64
	WindowUpdates     int64
}

// NewOptimizedConsole creates a new optimized console
func NewOptimizedConsole(width, height int, terminal terminals.TerminalType, log logging.LoggerType, connectionID string) *OptimizedConsole {
	console := &OptimizedConsole{
		Width:        width,
		Height:       height,
		Terminal:     terminal,
		Log:          log,
		connectionID: connectionID,
		sendChannel:  make(chan []byte, 100), // Buffered channel for async sending
	}

	console.deltaRenderer = renderer.NewDeltaRenderer(width, height, terminal)
	console.windowManager = renderer.NewWindowManager(console.deltaRenderer)

	return console
}

// Initialize sets up the console for first use
func (c *OptimizedConsole) Initialize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.initialized {
		return nil
	}

	// Send initial terminal setup
	c.sendRaw([]byte(c.Terminal.ClearTerminal()))

	c.initialized = true
	c.forceFullRefresh = true

	return nil
}

// Resize handles terminal resize efficiently
func (c *OptimizedConsole) Resize(newWidth, newHeight int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.Width == newWidth && c.Height == newHeight {
		return
	}

	c.Width = newWidth
	c.Height = newHeight
	c.resizeActive = true

	// Resize the renderer
	c.deltaRenderer.Resize(newWidth, newHeight)

	// Force full refresh after resize
	c.forceFullRefresh = true

	c.Log.Println(logging.LogInfo, "Console resized to", newWidth, "x", newHeight)
}

// CreateWindow creates a new optimized window
func (c *OptimizedConsole) CreateWindow(id string, x, y, width, height int) *renderer.OptimizedWindow {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	window := renderer.NewOptimizedWindow(id, x, y, width, height)
	c.windowManager.AddWindow(window)

	return window
}

// GetWindow retrieves a window by ID
func (c *OptimizedConsole) GetWindow(id string) *renderer.OptimizedWindow {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.windowManager.GetWindow(id)
}

// BringWindowToFront moves a window to the front
func (c *OptimizedConsole) BringWindowToFront(id string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.windowManager.BringToFront(id)
}

// SetText efficiently sets text at a specific position
func (c *OptimizedConsole) SetText(x, y int, text string, style uint32) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.deltaRenderer.SetText(x, y, text, style)
}

// SetCell efficiently sets a single character
func (c *OptimizedConsole) SetCell(x, y int, char rune, style uint32) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.deltaRenderer.SetCell(x, y, char, style)
}

// ClearRegion efficiently clears a rectangular area
func (c *OptimizedConsole) ClearRegion(x, y, width, height int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.deltaRenderer.ClearRegion(x, y, width, height)
}

// Render generates the minimal output needed to update the terminal
func (c *OptimizedConsole) Render() []byte {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	startTime := time.Now()

	// Handle full refresh if needed
	if c.forceFullRefresh {
		c.sendRaw([]byte(c.Terminal.ClearTerminal()))
		c.forceFullRefresh = false
		c.resizeActive = false
	}

	// Render all windows and get delta output
	output := c.windowManager.RenderAll()

	// Update statistics
	renderTime := time.Since(startTime)
	c.updateStats(len(output), renderTime)

	c.lastRenderTime = time.Now()

	return output
}

// RenderAsync renders and sends output asynchronously
func (c *OptimizedConsole) RenderAsync() {
	output := c.Render()
	if len(output) > 0 {
		select {
		case c.sendChannel <- output:
			// Sent successfully
		default:
			// Channel full, log warning
			c.Log.Println(logging.LogWarn, "Send channel full for connection", c.connectionID)
		}
	}
}

// GetSendChannel returns the channel for async output
func (c *OptimizedConsole) GetSendChannel() <-chan []byte {
	return c.sendChannel
}

// sendRaw sends data immediately (for initialization, etc.)
func (c *OptimizedConsole) sendRaw(data []byte) {
	select {
	case c.sendChannel <- data:
		// Sent successfully
	default:
		// Channel full, this is critical data so we log an error
		c.Log.Println(logging.LogError, "Failed to send critical data for connection", c.connectionID)
	}
}

// updateStats updates rendering performance statistics
func (c *OptimizedConsole) updateStats(outputSize int, renderTime time.Duration) {
	c.renderStats.TotalRenders++
	c.renderStats.BytesSent += int64(outputSize)
	c.renderStats.LastRenderSize = outputSize

	// Update average render time
	if c.renderStats.TotalRenders == 1 {
		c.renderStats.AverageRenderTime = renderTime
	} else {
		// Exponential moving average
		alpha := 0.1
		c.renderStats.AverageRenderTime = time.Duration(
			float64(c.renderStats.AverageRenderTime)*(1-alpha) +
				float64(renderTime)*alpha,
		)
	}

	// Update efficiency ratio from renderer
	rendererStats := c.deltaRenderer.GetStats()
	c.renderStats.EfficiencyRatio = rendererStats.EfficiencyRatio
}

// GetStats returns comprehensive performance statistics
func (c *OptimizedConsole) GetStats() ConsoleStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	windowStats := c.windowManager.GetStats()

	return ConsoleStats{
		ConnectionID:   c.connectionID,
		Width:          c.Width,
		Height:         c.Height,
		Initialized:    c.initialized,
		LastRenderTime: c.lastRenderTime,
		RenderingStats: c.renderStats,
		WindowStats:    windowStats,
	}
}

// ConsoleStats provides comprehensive console statistics
type ConsoleStats struct {
	ConnectionID   string
	Width, Height  int
	Initialized    bool
	LastRenderTime time.Time
	RenderingStats RenderingStats
	WindowStats    renderer.WindowManagerStats
}

// IsValidSize checks if the console size meets minimum requirements
func (c *OptimizedConsole) IsValidSize() bool {
	return c.Width >= MINWIDTH && c.Height >= MINHEIGHT
}

// ForceRefresh forces a complete screen refresh
func (c *OptimizedConsole) ForceRefresh() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.forceFullRefresh = true
}

// Close cleanly shuts down the console
func (c *OptimizedConsole) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Send cursor restore and cleanup
	c.sendRaw([]byte(c.Terminal.Reset()))

	// Close send channel
	close(c.sendChannel)
}

// GetEfficiencyReport returns a detailed efficiency report
func (c *OptimizedConsole) GetEfficiencyReport() EfficiencyReport {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	totalCells := c.Width * c.Height
	avgBytesPerRender := float64(0)
	if c.renderStats.TotalRenders > 0 {
		avgBytesPerRender = float64(c.renderStats.BytesSent) / float64(c.renderStats.TotalRenders)
	}

	// Theoretical maximum bytes per full screen update
	maxBytesPerFullScreen := float64(totalCells * 20) // Rough estimate with escape codes

	return EfficiencyReport{
		TotalCells:            totalCells,
		AverageBytesPerRender: avgBytesPerRender,
		MaxBytesPerFullScreen: maxBytesPerFullScreen,
		EfficiencyRatio:       c.renderStats.EfficiencyRatio,
		DataReduction:         1.0 - (avgBytesPerRender / maxBytesPerFullScreen),
		RenderFrequency:       float64(c.renderStats.TotalRenders) / time.Since(c.lastRenderTime).Seconds(),
	}
}

// EfficiencyReport provides detailed efficiency metrics
type EfficiencyReport struct {
	TotalCells            int
	AverageBytesPerRender float64
	MaxBytesPerFullScreen float64
	EfficiencyRatio       float64
	DataReduction         float64 // Percentage of data saved vs full screen updates
	RenderFrequency       float64 // Renders per second
}

// Use existing constants from console.go
