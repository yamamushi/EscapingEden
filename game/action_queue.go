package game

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

// QueuedAction represents a single action in the queue
type QueuedAction struct {
	ID          string
	CharacterID string
	Type        string
	StartTick   uint64
	Duration    int // ticks required
	Data        interface{}
	Callback    func(success bool, result interface{})
	CreatedAt   time.Time
	StartedAt   time.Time
	CompletedAt time.Time
}

// GetProgress returns the completion progress (0.0 to 1.0)
func (qa *QueuedAction) GetProgress(currentTick uint64) float64 {
	if qa.StartTick == 0 {
		return 0.0 // Not started yet
	}

	if currentTick < qa.StartTick {
		return 0.0 // Shouldn't happen, but handle gracefully
	}

	ticksElapsed := currentTick - qa.StartTick
	if ticksElapsed >= uint64(qa.Duration) {
		return 1.0 // Completed
	}

	return float64(ticksElapsed) / float64(qa.Duration)
}

// GetTicksLeft returns the number of ticks remaining
func (qa *QueuedAction) GetTicksLeft(currentTick uint64) int {
	if qa.StartTick == 0 {
		return qa.Duration // Not started yet
	}

	if currentTick < qa.StartTick {
		return qa.Duration // Shouldn't happen, but handle gracefully
	}

	ticksElapsed := int(currentTick - qa.StartTick)
	ticksLeft := qa.Duration - ticksElapsed

	if ticksLeft < 0 {
		return 0
	}

	return ticksLeft
}

// IsCompleted returns whether the action is completed
func (qa *QueuedAction) IsCompleted(currentTick uint64) bool {
	return qa.StartTick > 0 && currentTick >= qa.StartTick+uint64(qa.Duration)
}

// CharacterActionQueue manages the action queue for a single character
type CharacterActionQueue struct {
	CharacterID   string
	Actions       []*QueuedAction
	CurrentAction *QueuedAction
	MaxQueueSize  int
	mutex         sync.RWMutex
	lastActivity  time.Time
	log           logging.LoggerType
}

// NewCharacterActionQueue creates a new action queue for a character
func NewCharacterActionQueue(characterID string, maxQueueSize int, log logging.LoggerType) *CharacterActionQueue {
	return &CharacterActionQueue{
		CharacterID:  characterID,
		Actions:      make([]*QueuedAction, 0, maxQueueSize),
		MaxQueueSize: maxQueueSize,
		lastActivity: time.Now(),
		log:          log,
	}
}

// CanQueue returns whether a new action can be queued
func (caq *CharacterActionQueue) CanQueue() bool {
	caq.mutex.RLock()
	defer caq.mutex.RUnlock()
	return len(caq.Actions) < caq.MaxQueueSize
}

// AddAction adds a new action to the queue
func (caq *CharacterActionQueue) AddAction(action *QueuedAction) error {
	if action == nil {
		return errors.New("action cannot be nil")
	}

	caq.mutex.Lock()
	defer caq.mutex.Unlock()

	// Check queue capacity
	if len(caq.Actions) >= caq.MaxQueueSize {
		return fmt.Errorf("action queue is full (max %d actions)", caq.MaxQueueSize)
	}

	// Validate action
	if action.CharacterID != caq.CharacterID {
		return fmt.Errorf("action character ID mismatch: expected %s, got %s",
			caq.CharacterID, action.CharacterID)
	}

	if action.Duration <= 0 {
		return errors.New("action duration must be positive")
	}

	// Add to queue
	caq.Actions = append(caq.Actions, action)
	caq.lastActivity = time.Now()

	caq.log.Println(logging.LogInfo, fmt.Sprintf("Action '%s' queued for character %s (queue: %d/%d)",
		action.Type, caq.CharacterID, len(caq.Actions), caq.MaxQueueSize))

	return nil
}

// GetNextAction removes and returns the next action from the queue
func (caq *CharacterActionQueue) GetNextAction() *QueuedAction {
	caq.mutex.Lock()
	defer caq.mutex.Unlock()

	if len(caq.Actions) == 0 {
		return nil
	}

	// Get first action
	action := caq.Actions[0]

	// Remove from queue
	caq.Actions = caq.Actions[1:]
	caq.lastActivity = time.Now()

	return action
}

// StartNextAction moves the next queued action to current action
func (caq *CharacterActionQueue) StartNextAction(currentTick uint64) *QueuedAction {
	caq.mutex.Lock()
	defer caq.mutex.Unlock()

	// Can't start if there's already a current action
	if caq.CurrentAction != nil {
		return nil
	}

	// Get next action from queue
	if len(caq.Actions) == 0 {
		return nil
	}

	action := caq.Actions[0]
	caq.Actions = caq.Actions[1:]

	// Set as current action
	caq.CurrentAction = action
	caq.CurrentAction.StartTick = currentTick
	caq.CurrentAction.StartedAt = time.Now()
	caq.lastActivity = time.Now()

	caq.log.Println(logging.LogInfo, fmt.Sprintf("Action '%s' started for character %s (tick %d, duration %d)",
		action.Type, caq.CharacterID, currentTick, action.Duration))

	return action
}

// CompleteCurrentAction marks the current action as completed
func (caq *CharacterActionQueue) CompleteCurrentAction(currentTick uint64) *QueuedAction {
	caq.mutex.Lock()
	defer caq.mutex.Unlock()

	if caq.CurrentAction == nil {
		return nil
	}

	action := caq.CurrentAction
	action.CompletedAt = time.Now()
	caq.CurrentAction = nil
	caq.lastActivity = time.Now()

	caq.log.Println(logging.LogInfo, fmt.Sprintf("Action '%s' completed for character %s (tick %d)",
		action.Type, caq.CharacterID, currentTick))

	return action
}

// InterruptCurrentAction cancels the current action if it's interruptible
func (caq *CharacterActionQueue) InterruptCurrentAction() error {
	caq.mutex.Lock()
	defer caq.mutex.Unlock()

	if caq.CurrentAction == nil {
		return errors.New("no current action to interrupt")
	}

	action := caq.CurrentAction
	caq.CurrentAction = nil
	caq.lastActivity = time.Now()

	caq.log.Println(logging.LogInfo, fmt.Sprintf("Action '%s' interrupted for character %s",
		action.Type, caq.CharacterID))

	return nil
}

// Clear removes all actions from the queue and current action
func (caq *CharacterActionQueue) Clear() {
	caq.mutex.Lock()
	defer caq.mutex.Unlock()

	queuedCount := len(caq.Actions)
	hasCurrentAction := caq.CurrentAction != nil

	caq.Actions = caq.Actions[:0] // Clear slice but keep capacity
	caq.CurrentAction = nil
	caq.lastActivity = time.Now()

	caq.log.Println(logging.LogInfo, fmt.Sprintf("Action queue cleared for character %s (%d queued, current: %t)",
		caq.CharacterID, queuedCount, hasCurrentAction))
}

// GetQueueStatus returns the current status of the queue
func (caq *CharacterActionQueue) GetQueueStatus() ActionQueueStatus {
	caq.mutex.RLock()
	defer caq.mutex.RUnlock()

	status := ActionQueueStatus{
		CharacterID:      caq.CharacterID,
		QueuedCount:      len(caq.Actions),
		MaxQueueSize:     caq.MaxQueueSize,
		HasCurrentAction: caq.CurrentAction != nil,
		LastActivity:     caq.lastActivity,
	}

	if caq.CurrentAction != nil {
		status.CurrentActionType = caq.CurrentAction.Type
		status.CurrentActionID = caq.CurrentAction.ID
	}

	// Copy queued action types
	status.QueuedActionTypes = make([]string, len(caq.Actions))
	for i, action := range caq.Actions {
		status.QueuedActionTypes[i] = action.Type
	}

	return status
}

// GetCurrentAction returns a copy of the current action (thread-safe)
func (caq *CharacterActionQueue) GetCurrentAction() *QueuedAction {
	caq.mutex.RLock()
	defer caq.mutex.RUnlock()

	if caq.CurrentAction == nil {
		return nil
	}

	// Return a copy to prevent external modification
	return &QueuedAction{
		ID:          caq.CurrentAction.ID,
		CharacterID: caq.CurrentAction.CharacterID,
		Type:        caq.CurrentAction.Type,
		StartTick:   caq.CurrentAction.StartTick,
		Duration:    caq.CurrentAction.Duration,
		Data:        caq.CurrentAction.Data,
		CreatedAt:   caq.CurrentAction.CreatedAt,
		StartedAt:   caq.CurrentAction.StartedAt,
		CompletedAt: caq.CurrentAction.CompletedAt,
		// Note: Don't copy callback to prevent external access
	}
}

// GetQueuedActions returns copies of all queued actions (thread-safe)
func (caq *CharacterActionQueue) GetQueuedActions() []*QueuedAction {
	caq.mutex.RLock()
	defer caq.mutex.RUnlock()

	if len(caq.Actions) == 0 {
		return nil
	}

	// Return copies to prevent external modification
	copies := make([]*QueuedAction, len(caq.Actions))
	for i, action := range caq.Actions {
		copies[i] = &QueuedAction{
			ID:          action.ID,
			CharacterID: action.CharacterID,
			Type:        action.Type,
			StartTick:   action.StartTick,
			Duration:    action.Duration,
			Data:        action.Data,
			CreatedAt:   action.CreatedAt,
			StartedAt:   action.StartedAt,
			CompletedAt: action.CompletedAt,
			// Note: Don't copy callback to prevent external access
		}
	}

	return copies
}

// IsEmpty returns whether the queue has no actions (including current)
func (caq *CharacterActionQueue) IsEmpty() bool {
	caq.mutex.RLock()
	defer caq.mutex.RUnlock()
	return len(caq.Actions) == 0 && caq.CurrentAction == nil
}

// GetLastActivity returns the time of the last queue activity
func (caq *CharacterActionQueue) GetLastActivity() time.Time {
	caq.mutex.RLock()
	defer caq.mutex.RUnlock()
	return caq.lastActivity
}

// ActionQueueStatus represents the status of a character's action queue
type ActionQueueStatus struct {
	CharacterID       string
	QueuedCount       int
	MaxQueueSize      int
	HasCurrentAction  bool
	CurrentActionType string
	CurrentActionID   string
	QueuedActionTypes []string
	LastActivity      time.Time
}

// String returns a string representation of the queue status
func (aqs ActionQueueStatus) String() string {
	if !aqs.HasCurrentAction && aqs.QueuedCount == 0 {
		return fmt.Sprintf("Character %s: idle", aqs.CharacterID)
	}

	status := fmt.Sprintf("Character %s: ", aqs.CharacterID)

	if aqs.HasCurrentAction {
		status += fmt.Sprintf("current=%s", aqs.CurrentActionType)
	} else {
		status += "current=none"
	}

	status += fmt.Sprintf(", queued=%d/%d", aqs.QueuedCount, aqs.MaxQueueSize)

	if aqs.QueuedCount > 0 {
		status += fmt.Sprintf(" [%v]", aqs.QueuedActionTypes)
	}

	return status
}

// ActionQueueManager manages action queues for all characters
type ActionQueueManager struct {
	queues map[string]*CharacterActionQueue
	mutex  sync.RWMutex
	log    logging.LoggerType
}

// NewActionQueueManager creates a new action queue manager
func NewActionQueueManager(log logging.LoggerType) *ActionQueueManager {
	return &ActionQueueManager{
		queues: make(map[string]*CharacterActionQueue),
		mutex:  sync.RWMutex{},
		log:    log,
	}
}

// GetOrCreateQueue gets or creates a queue for a character
func (aqm *ActionQueueManager) GetOrCreateQueue(characterID string, maxQueueSize int) *CharacterActionQueue {
	aqm.mutex.Lock()
	defer aqm.mutex.Unlock()

	queue, exists := aqm.queues[characterID]
	if !exists {
		queue = NewCharacterActionQueue(characterID, maxQueueSize, aqm.log)
		aqm.queues[characterID] = queue
		aqm.log.Println(logging.LogInfo, fmt.Sprintf("Created action queue for character %s", characterID))
	}

	return queue
}

// GetQueue returns the queue for a character (nil if doesn't exist)
func (aqm *ActionQueueManager) GetQueue(characterID string) *CharacterActionQueue {
	aqm.mutex.RLock()
	defer aqm.mutex.RUnlock()
	return aqm.queues[characterID]
}

// RemoveQueue removes a character's queue (e.g., on logout)
func (aqm *ActionQueueManager) RemoveQueue(characterID string) {
	aqm.mutex.Lock()
	defer aqm.mutex.Unlock()

	if queue, exists := aqm.queues[characterID]; exists {
		queue.Clear() // Clear any remaining actions
		delete(aqm.queues, characterID)
		aqm.log.Println(logging.LogInfo, fmt.Sprintf("Removed action queue for character %s", characterID))
	}
}

// GetAllQueueStatuses returns status for all character queues
func (aqm *ActionQueueManager) GetAllQueueStatuses() []ActionQueueStatus {
	aqm.mutex.RLock()
	defer aqm.mutex.RUnlock()

	statuses := make([]ActionQueueStatus, 0, len(aqm.queues))
	for _, queue := range aqm.queues {
		statuses = append(statuses, queue.GetQueueStatus())
	}

	return statuses
}

// CleanupInactiveQueues removes queues that have been inactive for too long
func (aqm *ActionQueueManager) CleanupInactiveQueues(maxInactivity time.Duration) int {
	aqm.mutex.Lock()
	defer aqm.mutex.Unlock()

	cutoffTime := time.Now().Add(-maxInactivity)
	removedCount := 0

	for characterID, queue := range aqm.queues {
		if queue.GetLastActivity().Before(cutoffTime) && queue.IsEmpty() {
			queue.Clear()
			delete(aqm.queues, characterID)
			removedCount++
			aqm.log.Println(logging.LogInfo, fmt.Sprintf("Cleaned up inactive queue for character %s", characterID))
		}
	}

	if removedCount > 0 {
		aqm.log.Println(logging.LogInfo, fmt.Sprintf("Cleaned up %d inactive action queues", removedCount))
	}

	return removedCount
}

// GetQueueCount returns the total number of active queues
func (aqm *ActionQueueManager) GetQueueCount() int {
	aqm.mutex.RLock()
	defer aqm.mutex.RUnlock()
	return len(aqm.queues)
}
