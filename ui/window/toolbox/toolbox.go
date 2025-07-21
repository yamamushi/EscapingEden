package toolbox

import (
	"fmt"
	"strings"

	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/terminals"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/window"
	"sync"
	"time"
)

// ActionInfo represents action information for UI display
type ActionInfo struct {
	Type          string
	Description   string
	Progress      float64 // 0.0 to 1.0
	TicksLeft     int
	TotalTicks    int
	Interruptible bool
}

// ActionQueueDisplay holds the current action queue state for display
type ActionQueueDisplay struct {
	CurrentAction   *ActionInfo
	QueuedActions   []*ActionInfo
	MaxQueue        int
	QueueFull       bool
	LastUpdateTime  time.Time
	AnimationFrame  int
	ConflictMessage string
	ConflictExpiry  time.Time
	mutex           sync.RWMutex
}

// ToolboxWindow is a window that contains a toolbox for misc use
type ToolboxWindow struct {
	window.Window
	twMutex     sync.Mutex
	actionQueue *ActionQueueDisplay
}

// NewToolboxWindow creates a new toolbox window
func NewToolboxWindow(x, y, w, h, consoleWidth, consoleHeight int,
	input, output chan messages.WindowMessage, log logging.LoggerType, term terminals.TerminalType) *ToolboxWindow {
	lw := &ToolboxWindow{}
	lw.Log = log
	lw.Terminal = term
	lw.ID = config.WindowToolBox
	// if x or y are less than 1 set them to 1
	if x < 1 {
		x = 1
	}
	if y < 1 {
		y = 1
	}
	lw.X = x
	lw.Y = y

	// if w or h are less than 1 set them to 1
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	lw.Width = w
	lw.Height = h
	lw.ConsoleWidth = consoleWidth
	lw.ConsoleHeight = consoleHeight
	lw.Bordered = true
	lw.ConsoleReceive = input
	lw.ConsoleSend = output

	// Initialize action queue display
	lw.actionQueue = &ActionQueueDisplay{
		MaxQueue: 3, // Default max queue size
	}

	return lw
}

// HandleInput handles input for the toolbox window
func (tw *ToolboxWindow) HandleInput(input types.Input) {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	if tw.GetActive() {
		tw.Log.Println(logging.LogInfo, "Toolbox Handling input")
	}

	/*
		if len(input.Data) > 0 {
			tw.Log.Println(input.Data)
		}
	*/
}

// UpdateContents updates the contents of the toolbox window
func (tw *ToolboxWindow) UpdateContents() {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	// Time information (existing)
	serverTime := time.Now().Format("15:04:05")
	edenTime := edenutil.EdenTime.CurrentTimeString(edenutil.EdenTime{})
	edenMonth := edenutil.EdenTime.EdenMonth(edenutil.EdenTime{})
	edenDay := edenutil.EdenTime.EdenDay(edenutil.EdenTime{})
	edenYear := edenutil.EdenTime.Year(edenutil.EdenTime{})

	content := fmt.Sprintf("Current Server Time: %s\n  Current Eden Time: %s\n  Eden Month: %s\n  Eden Day: %s\n  Eden Year: %d\n\n",
		serverTime, edenTime, edenMonth.String(), edenDay.String(), edenYear)

	// Action queue information (new)
	content += tw.formatActionQueue()

	tw.SetContents(content)
}

// UpdateActionQueue updates the action queue display with new information
func (tw *ToolboxWindow) UpdateActionQueue(current *ActionInfo, queued []*ActionInfo) {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	tw.actionQueue.mutex.Lock()
	tw.actionQueue.CurrentAction = current
	tw.actionQueue.QueuedActions = queued
	tw.actionQueue.QueueFull = len(queued) >= tw.actionQueue.MaxQueue
	tw.actionQueue.LastUpdateTime = time.Now()
	tw.actionQueue.mutex.Unlock()

	// Trigger content update
	tw.UpdateContents()
}

// ShowQueueConflict displays a temporary conflict message
func (tw *ToolboxWindow) ShowQueueConflict(message string, duration time.Duration) {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	tw.actionQueue.mutex.Lock()
	tw.actionQueue.ConflictMessage = message
	tw.actionQueue.ConflictExpiry = time.Now().Add(duration)
	tw.actionQueue.mutex.Unlock()

	// Trigger immediate content update
	tw.UpdateContents()
}

// ClearQueueConflict clears any active conflict message
func (tw *ToolboxWindow) ClearQueueConflict() {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	tw.actionQueue.mutex.Lock()
	tw.actionQueue.ConflictMessage = ""
	tw.actionQueue.ConflictExpiry = time.Time{}
	tw.actionQueue.mutex.Unlock()

	tw.UpdateContents()
}

// IsQueueFull returns whether the action queue is at capacity
func (tw *ToolboxWindow) IsQueueFull() bool {
	tw.actionQueue.mutex.RLock()
	defer tw.actionQueue.mutex.RUnlock()
	return tw.actionQueue.QueueFull
}

// GetQueueUtilization returns the current queue utilization as a percentage
func (tw *ToolboxWindow) GetQueueUtilization() float64 {
	tw.actionQueue.mutex.RLock()
	defer tw.actionQueue.mutex.RUnlock()

	if tw.actionQueue.MaxQueue == 0 {
		return 0.0
	}

	return float64(len(tw.actionQueue.QueuedActions)) / float64(tw.actionQueue.MaxQueue)
}

// formatActionQueue formats the action queue for display with minimal, clean output
func (tw *ToolboxWindow) formatActionQueue() string {
	tw.actionQueue.mutex.RLock()
	defer tw.actionQueue.mutex.RUnlock()

	var content strings.Builder

	// Simple header
	content.WriteString("No active action - use game commands:\n")
	content.WriteString("    Movement: h,j,k,l (vi keys)\n")
	content.WriteString("    Digging: Ctrl+D\n")
	content.WriteString("    Building: Ctrl+B\n\n")

	// Simple queue status - only essential information
	content.WriteString("QUEUE STATUS:\n")

	queuedCount := len(tw.actionQueue.QueuedActions)
	freeSlots := tw.actionQueue.MaxQueue - queuedCount
	utilization := float64(queuedCount) / float64(tw.actionQueue.MaxQueue) * 100

	content.WriteString(fmt.Sprintf("  Used: %d/%d slots (%d free)\n",
		queuedCount, tw.actionQueue.MaxQueue, freeSlots))
	content.WriteString(fmt.Sprintf("  Utilization: %.0f%%\n", utilization))

	if tw.actionQueue.QueueFull {
		content.WriteString("  Status: Queue full\n")
	} else {
		content.WriteString("  Status: Room for more actions\n")
	}

	return content.String()
}

// getQueueStatusIndicator returns a simple text indicator for queue status
func (tw *ToolboxWindow) getQueueStatusIndicator() string {
	queuedCount := len(tw.actionQueue.QueuedActions)

	if tw.actionQueue.QueueFull {
		return "FULL"
	}

	// Check for conflict messages
	if tw.actionQueue.ConflictMessage != "" && time.Now().Before(tw.actionQueue.ConflictExpiry) {
		return "ERROR"
	}

	hasCurrentAction := tw.actionQueue.CurrentAction != nil

	if hasCurrentAction {
		switch queuedCount {
		case 0:
			return "WORKING"
		case 1:
			return "ACTIVE"
		case 2:
			return "BUSY"
		default:
			return "LOADED"
		}
	}

	// Static indicators for idle state
	switch queuedCount {
	case 0:
		return "IDLE"
	case 1:
		return "QUEUED"
	case 2:
		return "READY"
	default:
		return "WAITING"
	}
}

// formatTimeRemaining formats tick count into human-readable time
func (tw *ToolboxWindow) formatTimeRemaining(ticks int) string {
	if ticks <= 0 {
		return "0s"
	}

	// Assuming 200ms per tick (5 ticks per second)
	seconds := float64(ticks) / 5.0

	if seconds < 1.0 {
		return fmt.Sprintf("%.1fs", seconds)
	} else if seconds < 60.0 {
		return fmt.Sprintf("%.0fs", seconds)
	} else {
		minutes := int(seconds / 60)
		remainingSeconds := int(seconds) % 60
		return fmt.Sprintf("%dm%ds", minutes, remainingSeconds)
	}
}

// getHelpText returns simple help text without Unicode characters
func (tw *ToolboxWindow) getHelpText() string {
	// This function is no longer used in the simplified UI
	// All help information has been moved to documentation
	return ""
}

// createProgressBar creates a simple text-based progress bar
func (tw *ToolboxWindow) createProgressBar(progress float64, width int) string {
	filled := int(progress * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	// Create simple progress bar with ASCII characters only
	var bar strings.Builder

	for i := 0; i < width; i++ {
		if i < filled {
			bar.WriteString("#")
		} else {
			bar.WriteString("-")
		}
	}

	return bar.String()
}

// createAnimatedProgressBar creates a simple progress bar without complex animations
func (tw *ToolboxWindow) createAnimatedProgressBar(progress float64, width int, isActive bool) string {
	return tw.createProgressBar(progress, width)
}

// SetMaxQueueSize sets the maximum queue size for display
func (tw *ToolboxWindow) SetMaxQueueSize(maxSize int) {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	tw.actionQueue.mutex.Lock()
	tw.actionQueue.MaxQueue = maxSize
	tw.actionQueue.mutex.Unlock()
}

// GetActionQueueStatus returns the current action queue status
func (tw *ToolboxWindow) GetActionQueueStatus() (bool, int, int) {
	tw.actionQueue.mutex.RLock()
	defer tw.actionQueue.mutex.RUnlock()

	hasCurrentAction := tw.actionQueue.CurrentAction != nil
	queuedCount := len(tw.actionQueue.QueuedActions)
	maxQueue := tw.actionQueue.MaxQueue

	return hasCurrentAction, queuedCount, maxQueue
}

// ShowActionFeedback displays temporary feedback for user actions
func (tw *ToolboxWindow) ShowActionFeedback(actionType string, success bool, message string) {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	var feedbackMsg string
	var duration time.Duration

	if success {
		switch actionType {
		case "queue":
			feedbackMsg = fmt.Sprintf("OK: %s", message)
			duration = 2 * time.Second
		case "interrupt":
			feedbackMsg = fmt.Sprintf("STOPPED: %s", message)
			duration = 1500 * time.Millisecond
		case "complete":
			feedbackMsg = fmt.Sprintf("DONE: %s", message)
			duration = 2 * time.Second
		default:
			feedbackMsg = fmt.Sprintf("INFO: %s", message)
			duration = 2 * time.Second
		}
	} else {
		feedbackMsg = fmt.Sprintf("ERROR: %s", message)
		duration = 3 * time.Second
	}

	tw.ShowQueueConflict(feedbackMsg, duration)
}

// UpdateWithAnimation triggers a content update with animation refresh
func (tw *ToolboxWindow) UpdateWithAnimation() {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	// Clear expired conflict messages
	tw.actionQueue.mutex.Lock()
	if tw.actionQueue.ConflictMessage != "" && time.Now().After(tw.actionQueue.ConflictExpiry) {
		tw.actionQueue.ConflictMessage = ""
		tw.actionQueue.ConflictExpiry = time.Time{}
	}
	tw.actionQueue.mutex.Unlock()

	// Update contents to refresh animations
	tw.UpdateContents()
}

// GetDetailedQueueInfo returns detailed information about the current queue state
func (tw *ToolboxWindow) GetDetailedQueueInfo() map[string]interface{} {
	tw.actionQueue.mutex.RLock()
	defer tw.actionQueue.mutex.RUnlock()

	info := make(map[string]interface{})

	// Current action details
	if tw.actionQueue.CurrentAction != nil {
		info["current_action"] = map[string]interface{}{
			"type":          tw.actionQueue.CurrentAction.Type,
			"description":   tw.actionQueue.CurrentAction.Description,
			"progress":      tw.actionQueue.CurrentAction.Progress,
			"ticks_left":    tw.actionQueue.CurrentAction.TicksLeft,
			"total_ticks":   tw.actionQueue.CurrentAction.TotalTicks,
			"interruptible": tw.actionQueue.CurrentAction.Interruptible,
		}
	}

	// Queue details
	queueInfo := make([]map[string]interface{}, len(tw.actionQueue.QueuedActions))
	for i, action := range tw.actionQueue.QueuedActions {
		queueInfo[i] = map[string]interface{}{
			"type":          action.Type,
			"description":   action.Description,
			"total_ticks":   action.TotalTicks,
			"interruptible": action.Interruptible,
		}
	}
	info["queued_actions"] = queueInfo

	// Queue statistics
	info["queue_stats"] = map[string]interface{}{
		"max_queue":     tw.actionQueue.MaxQueue,
		"current_count": len(tw.actionQueue.QueuedActions),
		"is_full":       tw.actionQueue.QueueFull,
		"utilization":   tw.GetQueueUtilization(),
	}

	// Status information
	info["status"] = map[string]interface{}{
		"has_conflict":    tw.actionQueue.ConflictMessage != "",
		"conflict_msg":    tw.actionQueue.ConflictMessage,
		"last_update":     tw.actionQueue.LastUpdateTime,
		"animation_frame": tw.actionQueue.AnimationFrame,
	}

	return info
}

// SetQueueTheme allows customization of the queue display theme
func (tw *ToolboxWindow) SetQueueTheme(theme map[string]string) {
	// This could be extended to allow theme customization
	// For now, it's a placeholder for future enhancement
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	// Theme could include custom icons, colors, etc.
	// Implementation would depend on terminal capabilities
}

// GetQueuePerformanceMetrics returns performance metrics for the action queue
func (tw *ToolboxWindow) GetQueuePerformanceMetrics() map[string]float64 {
	tw.actionQueue.mutex.RLock()
	defer tw.actionQueue.mutex.RUnlock()

	metrics := make(map[string]float64)

	// Calculate average queue utilization over time
	metrics["utilization"] = tw.GetQueueUtilization()

	// Calculate estimated completion time for all queued actions
	totalTicks := 0
	if tw.actionQueue.CurrentAction != nil {
		totalTicks += tw.actionQueue.CurrentAction.TicksLeft
	}
	for _, action := range tw.actionQueue.QueuedActions {
		totalTicks += action.TotalTicks
	}
	metrics["estimated_completion_seconds"] = float64(totalTicks) / 5.0 // 5 ticks per second

	// Queue efficiency (how full the queue stays)
	metrics["queue_efficiency"] = float64(len(tw.actionQueue.QueuedActions)) / float64(tw.actionQueue.MaxQueue)

	return metrics
}

// StartAnimationTimer starts a timer for smooth animation updates with adaptive refresh rate
func (tw *ToolboxWindow) StartAnimationTimer() {
	go func() {
		// Use adaptive refresh rate based on activity
		baseTicker := time.NewTicker(200 * time.Millisecond) // 5 FPS base rate
		fastTicker := time.NewTicker(100 * time.Millisecond) // 10 FPS for active animations
		defer baseTicker.Stop()
		defer fastTicker.Stop()

		var currentTicker *time.Ticker = baseTicker
		var lastActivity time.Time

		for {
			select {
			case <-currentTicker.C:
				// Check current activity level
				tw.actionQueue.mutex.RLock()
				hasActiveAction := tw.actionQueue.CurrentAction != nil
				hasConflict := tw.actionQueue.ConflictMessage != "" && time.Now().Before(tw.actionQueue.ConflictExpiry)
				queueFull := tw.actionQueue.QueueFull
				tw.actionQueue.mutex.RUnlock()

				// Determine if we need high-frequency updates
				needsFastUpdate := hasActiveAction || hasConflict || queueFull

				if needsFastUpdate {
					lastActivity = time.Now()
					if currentTicker == baseTicker {
						// Switch to fast updates
						currentTicker = fastTicker
					}
					tw.UpdateWithAnimation()
				} else if time.Since(lastActivity) > 2*time.Second {
					// Switch back to slow updates after period of inactivity
					if currentTicker == fastTicker {
						currentTicker = baseTicker
					}
					// Still update occasionally to clear expired messages
					tw.UpdateWithAnimation()
				}
			}
		}
	}()
}

// HandleActionQueueMessage processes action queue related messages
func (tw *ToolboxWindow) HandleActionQueueMessage(msgType string, data interface{}) {
	switch msgType {
	case "action_started":
		tw.ShowActionFeedback("queue", true, "Action started")
	case "action_completed":
		tw.ShowActionFeedback("complete", true, "Action completed")
	case "action_failed":
		tw.ShowActionFeedback("queue", false, "Action failed")
	case "action_interrupted":
		tw.ShowActionFeedback("interrupt", true, "Action cancelled")
	case "queue_full":
		tw.ShowActionFeedback("queue", false, "Queue is full - wait or cancel actions")
	case "queue_updated":
		// Just refresh the display
		tw.UpdateContents()
	}
}

// GetActionQueueTooltip returns a tooltip for the specified queue position
func (tw *ToolboxWindow) GetActionQueueTooltip(position int) string {
	tw.actionQueue.mutex.RLock()
	defer tw.actionQueue.mutex.RUnlock()

	if position == 0 && tw.actionQueue.CurrentAction != nil {
		action := tw.actionQueue.CurrentAction
		return fmt.Sprintf("Current: %s\nProgress: %.1f%%\nTime left: %s\nCan cancel: %t",
			action.Description,
			action.Progress*100,
			tw.formatTimeRemaining(action.TicksLeft),
			action.Interruptible)
	}

	queueIndex := position - 1
	if queueIndex >= 0 && queueIndex < len(tw.actionQueue.QueuedActions) {
		action := tw.actionQueue.QueuedActions[queueIndex]
		return fmt.Sprintf("Queued #%d: %s\nEstimated time: %s\nCan cancel: %t",
			position,
			action.Description,
			tw.formatTimeRemaining(action.TotalTicks),
			action.Interruptible)
	}

	return fmt.Sprintf("Queue slot #%d: Empty\nUse game commands to add actions", position)
}

// ShowProgressMilestone displays simple feedback when actions reach certain progress milestones
func (tw *ToolboxWindow) ShowProgressMilestone(progress float64, actionType string) {
	// Milestone messages are disabled in the simplified UI
	// Progress information is available in logs if needed
}

// ShowQueueOptimizationTip displays helpful tips based on current queue state
func (tw *ToolboxWindow) ShowQueueOptimizationTip() {
	// Optimization tips are disabled in the simplified UI
	// Tips and strategies are available in the game documentation
}

// ShowActionChainSuggestion suggests logical next actions based on current queue
func (tw *ToolboxWindow) ShowActionChainSuggestion(lastCompletedAction string) {
	// Action chain suggestions are disabled in the simplified UI
	// Strategy guides are available in the game documentation
}

// GetSmartQueueRecommendations returns intelligent recommendations for queue optimization
func (tw *ToolboxWindow) GetSmartQueueRecommendations() []string {
	tw.actionQueue.mutex.RLock()
	defer tw.actionQueue.mutex.RUnlock()

	var recommendations []string

	queuedCount := len(tw.actionQueue.QueuedActions)
	hasCurrentAction := tw.actionQueue.CurrentAction != nil

	// Analyze current queue composition
	if queuedCount == 0 && !hasCurrentAction {
		recommendations = append(recommendations, "Start with movement to position yourself strategically")
		recommendations = append(recommendations, "Queue 2-3 actions at once for better efficiency")
	}

	if queuedCount > 0 {
		// Analyze action types in queue
		interruptibleCount := 0
		longActionCount := 0

		for _, action := range tw.actionQueue.QueuedActions {
			if action.Interruptible {
				interruptibleCount++
			}
			if action.TotalTicks > 15 { // Actions longer than 3 seconds
				longActionCount++
			}
		}

		// Provide specific recommendations
		if interruptibleCount == 0 {
			recommendations = append(recommendations, "Consider adding quick actions (movement/digging) for flexibility")
		}

		if longActionCount == queuedCount && queuedCount > 1 {
			recommendations = append(recommendations, "Mix in some quick actions between long ones")
		}

		if queuedCount < tw.actionQueue.MaxQueue-1 {
			recommendations = append(recommendations, "You have queue space - add more actions for efficiency")
		}
	}

	// Performance-based recommendations
	metrics := tw.GetQueuePerformanceMetrics()
	if efficiency := metrics["queue_efficiency"]; efficiency < 0.5 {
		recommendations = append(recommendations, "Low queue utilization - try to keep 2-3 actions queued")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Excellent queue management! Keep up the good work!")
	}

	return recommendations
}

// ShowContextualHelp displays help text relevant to the current situation
func (tw *ToolboxWindow) ShowContextualHelp() {
	// Contextual help is disabled in the simplified UI
	// Help information is available in the game documentation
}

// UpdateProgressWithSync updates progress with synchronized timing for smoother animations
func (tw *ToolboxWindow) UpdateProgressWithSync(currentAction *ActionInfo, queued []*ActionInfo, tickNumber uint64) {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	tw.actionQueue.mutex.Lock()

	// Store previous progress for milestone detection
	var previousProgress float64
	if tw.actionQueue.CurrentAction != nil {
		previousProgress = tw.actionQueue.CurrentAction.Progress
	}

	// Update action queue data
	tw.actionQueue.CurrentAction = currentAction
	tw.actionQueue.QueuedActions = queued
	tw.actionQueue.QueueFull = len(queued) >= tw.actionQueue.MaxQueue
	tw.actionQueue.LastUpdateTime = time.Now()

	tw.actionQueue.mutex.Unlock()

	// Check for progress milestones
	if currentAction != nil && previousProgress > 0 {
		// Check if we crossed a milestone boundary
		milestones := []float64{0.25, 0.50, 0.75, 0.95}
		for _, milestone := range milestones {
			if previousProgress < milestone && currentAction.Progress >= milestone {
				tw.ShowProgressMilestone(milestone, currentAction.Type)
				break
			}
		}
	}

	// Trigger content update with animation
	tw.UpdateContents()
}

// GetAdvancedQueueMetrics returns detailed metrics for power users
func (tw *ToolboxWindow) GetAdvancedQueueMetrics() map[string]interface{} {
	tw.actionQueue.mutex.RLock()
	defer tw.actionQueue.mutex.RUnlock()

	metrics := make(map[string]interface{})

	// Basic metrics
	queuedCount := len(tw.actionQueue.QueuedActions)
	hasCurrentAction := tw.actionQueue.CurrentAction != nil

	// Calculate action type distribution
	actionTypes := make(map[string]int)
	interruptibleCount := 0
	totalEstimatedTime := 0

	if hasCurrentAction {
		actionTypes[tw.actionQueue.CurrentAction.Type]++
		if tw.actionQueue.CurrentAction.Interruptible {
			interruptibleCount++
		}
		totalEstimatedTime += tw.actionQueue.CurrentAction.TicksLeft
	}

	for _, action := range tw.actionQueue.QueuedActions {
		actionTypes[action.Type]++
		if action.Interruptible {
			interruptibleCount++
		}
		totalEstimatedTime += action.TotalTicks
	}

	// Compile advanced metrics
	metrics["queue_composition"] = actionTypes
	metrics["interruptible_ratio"] = float64(interruptibleCount) / float64(queuedCount+1)
	metrics["estimated_total_time_seconds"] = float64(totalEstimatedTime) / 5.0
	metrics["queue_diversity"] = len(actionTypes) // Number of different action types
	metrics["average_action_duration"] = float64(totalEstimatedTime) / float64(queuedCount+1)

	// Performance indicators
	metrics["efficiency_score"] = tw.calculateEfficiencyScore()
	metrics["flexibility_score"] = float64(interruptibleCount) / float64(queuedCount+1)

	// Timing analysis
	if hasCurrentAction {
		metrics["current_progress_rate"] = tw.actionQueue.CurrentAction.Progress / float64(tw.actionQueue.CurrentAction.TotalTicks-tw.actionQueue.CurrentAction.TicksLeft)
	}

	return metrics
}

// calculateEfficiencyScore calculates a score from 0-100 based on queue utilization and composition
func (tw *ToolboxWindow) calculateEfficiencyScore() float64 {
	queuedCount := len(tw.actionQueue.QueuedActions)
	hasCurrentAction := tw.actionQueue.CurrentAction != nil

	// Base score from utilization
	utilizationScore := float64(queuedCount) / float64(tw.actionQueue.MaxQueue) * 50

	// Bonus for having an active action
	if hasCurrentAction {
		utilizationScore += 25
	}

	// Bonus for action diversity (having different types of actions)
	actionTypes := make(map[string]bool)
	if hasCurrentAction {
		actionTypes[tw.actionQueue.CurrentAction.Type] = true
	}
	for _, action := range tw.actionQueue.QueuedActions {
		actionTypes[action.Type] = true
	}

	diversityBonus := float64(len(actionTypes)) * 5 // Up to 25 points for 5+ different action types
	if diversityBonus > 25 {
		diversityBonus = 25
	}

	totalScore := utilizationScore + diversityBonus
	if totalScore > 100 {
		totalScore = 100
	}

	return totalScore
}
