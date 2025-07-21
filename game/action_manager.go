package game

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

// GameActionHandler defines the interface for executing game actions
type GameActionHandler interface {
	HandleMovePlayerRequest(characterID string, deltaX, deltaY int) error
	HandleBuildWallRequest(itemID, toolID, characterID string, deltaX, deltaY int) error
	HandleMineRequest(itemID, toolID, characterID string, deltaX, deltaY int) error
	HandleDigRequest(itemID, characterID string, deltaX, deltaY int) error
}

// ActionManager is the central coordinator for the tick-based action system
type ActionManager struct {
	registry     *ActionRegistry
	queueManager *ActionQueueManager
	gameHandler  GameActionHandler
	ticker       *GameTicker
	tickCounter  uint64
	mutex        sync.RWMutex
	log          logging.LoggerType

	// Message sending
	sendChannel chan messages.ConnectionManagerMessage

	// Metrics
	actionsProcessed uint64
	actionsCompleted uint64
	actionsFailed    uint64

	// Configuration
	enabled bool
}

// NewActionManager creates a new ActionManager
func NewActionManager(registry *ActionRegistry, gameHandler GameActionHandler, log logging.LoggerType) *ActionManager {
	return &ActionManager{
		registry:     registry,
		queueManager: NewActionQueueManager(log),
		gameHandler:  gameHandler,
		log:          log,
		enabled:      true,
	}
}

// SetSendChannel sets the channel for sending messages to the connection manager
func (am *ActionManager) SetSendChannel(sendChannel chan messages.ConnectionManagerMessage) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	am.sendChannel = sendChannel
}

// SetGameTicker sets the game ticker for this action manager
func (am *ActionManager) SetGameTicker(ticker *GameTicker) error {
	if ticker == nil {
		return errors.New("ticker cannot be nil")
	}

	am.mutex.Lock()
	defer am.mutex.Unlock()

	am.ticker = ticker
	return ticker.Subscribe(am)
}

// GetSubscriberID returns the subscriber ID for the GameTicker
func (am *ActionManager) GetSubscriberID() string {
	return "ActionManager"
}

// OnTick processes all queued actions on each tick
func (am *ActionManager) OnTick(tickNumber uint64) {
	if !am.enabled {
		return
	}

	startTime := time.Now()

	am.mutex.Lock()
	am.tickCounter = tickNumber
	am.mutex.Unlock()

	// Process all character queues
	am.processAllQueues(tickNumber)

	// Log performance if processing takes too long
	processingTime := time.Since(startTime)
	if processingTime > 50*time.Millisecond {
		am.log.Println(logging.LogWarn, fmt.Sprintf("Tick %d action processing took %v", tickNumber, processingTime))
	}
}

// QueueAction queues a new action for a character
func (am *ActionManager) QueueAction(characterID string, actionType string, data interface{}) error {
	if !am.enabled {
		return errors.New("action manager is disabled")
	}

	// Validate action type
	err := am.registry.ValidateAction(actionType)
	if err != nil {
		return fmt.Errorf("action validation failed: %v", err)
	}

	// Get action definition
	actionDef, err := am.registry.GetActionDefinition(actionType)
	if err != nil {
		return fmt.Errorf("failed to get action definition: %v", err)
	}

	// Get or create character queue
	queue := am.queueManager.GetOrCreateQueue(characterID, am.registry.GetMaxQueueSize())

	// Check if queue can accept more actions
	if !queue.CanQueue() {
		return errors.New("action queue is full")
	}

	// Create queued action
	action := &QueuedAction{
		ID:          edenutil.GenerateID(),
		CharacterID: characterID,
		Type:        actionType,
		Duration:    actionDef.TickCost,
		Data:        data,
		CreatedAt:   time.Now(),
		Callback:    am.createActionCallback(characterID, actionType),
	}

	// Add to queue
	err = queue.AddAction(action)
	if err != nil {
		return fmt.Errorf("failed to queue action: %v", err)
	}

	am.log.Println(logging.LogInfo, fmt.Sprintf("Action '%s' queued for character %s", actionType, characterID))

	// Send action queue update
	am.sendActionQueueUpdate(characterID)

	return nil
}

// InterruptCurrentAction interrupts the current action for a character if it's interruptible
func (am *ActionManager) InterruptCurrentAction(characterID string) error {
	queue := am.queueManager.GetQueue(characterID)
	if queue == nil {
		return errors.New("no action queue found for character")
	}

	currentAction := queue.GetCurrentAction()
	if currentAction == nil {
		return errors.New("no current action to interrupt")
	}

	// Check if action is interruptible
	actionDef, err := am.registry.GetActionDefinition(currentAction.Type)
	if err != nil {
		return fmt.Errorf("failed to get action definition: %v", err)
	}

	if !actionDef.Interruptible {
		return fmt.Errorf("action '%s' cannot be interrupted", currentAction.Type)
	}

	// Interrupt the action
	err = queue.InterruptCurrentAction()
	if err != nil {
		return fmt.Errorf("failed to interrupt action: %v", err)
	}

	am.log.Println(logging.LogInfo, fmt.Sprintf("Action '%s' interrupted for character %s", currentAction.Type, characterID))

	// Send action queue update
	am.sendActionQueueUpdate(characterID)

	return nil
}

// OnCharacterLogout cleans up a character's action queue when they log out
func (am *ActionManager) OnCharacterLogout(characterID string) {
	am.queueManager.RemoveQueue(characterID)
	am.log.Println(logging.LogInfo, fmt.Sprintf("Cleaned up action queue for logged out character %s", characterID))
}

// GetCharacterQueue returns the action queue for a character
func (am *ActionManager) GetCharacterQueue(characterID string) *CharacterActionQueue {
	return am.queueManager.GetQueue(characterID)
}

// GetAllQueueStatuses returns status for all character queues
func (am *ActionManager) GetAllQueueStatuses() []ActionQueueStatus {
	return am.queueManager.GetAllQueueStatuses()
}

// processAllQueues processes actions for all character queues
func (am *ActionManager) processAllQueues(tickNumber uint64) {
	statuses := am.queueManager.GetAllQueueStatuses()

	for _, status := range statuses {
		am.processCharacterQueue(status.CharacterID, tickNumber)
	}
}

// processCharacterQueue processes actions for a single character
func (am *ActionManager) processCharacterQueue(characterID string, tickNumber uint64) {
	queue := am.queueManager.GetQueue(characterID)
	if queue == nil {
		return
	}

	// Check if current action is completed
	currentAction := queue.GetCurrentAction()
	if currentAction != nil && currentAction.IsCompleted(tickNumber) {
		am.completeAction(queue, currentAction, tickNumber)
	}

	// Start next action if no current action
	if queue.GetCurrentAction() == nil {
		nextAction := queue.StartNextAction(tickNumber)
		if nextAction != nil {
			am.log.Println(logging.LogInfo, fmt.Sprintf("Started action '%s' for character %s", nextAction.Type, characterID))
			am.sendActionStarted(characterID, nextAction)
			am.sendActionQueueUpdate(characterID)
		}
	}
}

// completeAction handles the completion of an action
func (am *ActionManager) completeAction(queue *CharacterActionQueue, action *QueuedAction, tickNumber uint64) {
	// Mark action as completed in queue
	completedAction := queue.CompleteCurrentAction(tickNumber)
	if completedAction == nil {
		am.log.Println(logging.LogWarn, "Failed to complete action in queue")
		return
	}

	// Execute the action through GameManager
	success := am.executeAction(completedAction)

	// Update metrics
	am.mutex.Lock()
	am.actionsCompleted++
	if !success {
		am.actionsFailed++
	}
	am.mutex.Unlock()

	// Call the action callback if it exists
	if completedAction.Callback != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					am.log.Println(logging.LogError, fmt.Sprintf("Action callback panic: %v", r))
				}
			}()
			completedAction.Callback(success, nil)
		}()
	}

	// Send completion notification and queue update
	am.sendActionCompleted(completedAction.CharacterID, completedAction, success)
	am.sendActionQueueUpdate(completedAction.CharacterID)

	if success {
		am.log.Println(logging.LogInfo, fmt.Sprintf("Action '%s' completed successfully for character %s",
			completedAction.Type, completedAction.CharacterID))
	} else {
		am.log.Println(logging.LogWarn, fmt.Sprintf("Action '%s' failed for character %s",
			completedAction.Type, completedAction.CharacterID))
	}
}

// executeAction executes a completed action through the appropriate game system
func (am *ActionManager) executeAction(action *QueuedAction) bool {
	switch action.Type {
	case "move":
		return am.executeMovementAction(action)
	case "build_wall":
		return am.executeBuildingAction(action)
	case "mine":
		return am.executeMiningAction(action)
	case "dig":
		return am.executeDiggingAction(action)
	default:
		am.log.Println(logging.LogError, fmt.Sprintf("Unknown action type: %s", action.Type))
		return false
	}
}

// executeMovementAction executes a movement action
func (am *ActionManager) executeMovementAction(action *QueuedAction) bool {
	// Extract movement data
	movementData, ok := action.Data.(map[string]interface{})
	if !ok {
		am.log.Println(logging.LogError, "Invalid movement action data")
		return false
	}

	direction, ok := movementData["direction"].(string)
	if !ok {
		am.log.Println(logging.LogError, "Missing direction in movement action")
		return false
	}

	// Calculate delta based on direction
	deltaX, deltaY := am.getMovementDelta(direction)

	// Execute movement through GameManager
	err := am.gameHandler.HandleMovePlayerRequest(action.CharacterID, deltaX, deltaY)
	if err != nil {
		am.log.Println(logging.LogWarn, fmt.Sprintf("Movement failed for character %s: %v", action.CharacterID, err))
		return false
	}

	return true
}

// executeBuildingAction executes a building action
func (am *ActionManager) executeBuildingAction(action *QueuedAction) bool {
	// Extract building data
	buildData, ok := action.Data.(map[string]interface{})
	if !ok {
		am.log.Println(logging.LogError, "Invalid building action data")
		return false
	}

	itemID, ok := buildData["itemID"].(string)
	if !ok {
		am.log.Println(logging.LogError, "Missing itemID in building action")
		return false
	}

	toolID, _ := buildData["toolID"].(string) // Optional
	deltaX, _ := buildData["deltaX"].(int)
	deltaY, _ := buildData["deltaY"].(int)

	// Execute building through GameManager
	err := am.gameHandler.HandleBuildWallRequest(itemID, toolID, action.CharacterID, deltaX, deltaY)
	if err != nil {
		am.log.Println(logging.LogWarn, fmt.Sprintf("Building failed for character %s: %v", action.CharacterID, err))
		return false
	}

	return true
}

// executeMiningAction executes a mining action
func (am *ActionManager) executeMiningAction(action *QueuedAction) bool {
	// Extract mining data
	miningData, ok := action.Data.(map[string]interface{})
	if !ok {
		am.log.Println(logging.LogError, "Invalid mining action data")
		return false
	}

	itemID, ok := miningData["itemID"].(string)
	if !ok {
		am.log.Println(logging.LogError, "Missing itemID in mining action")
		return false
	}

	toolID, _ := miningData["toolID"].(string) // Optional
	deltaX, _ := miningData["deltaX"].(int)
	deltaY, _ := miningData["deltaY"].(int)

	// Execute mining through GameManager
	err := am.gameHandler.HandleMineRequest(itemID, toolID, action.CharacterID, deltaX, deltaY)
	if err != nil {
		am.log.Println(logging.LogWarn, fmt.Sprintf("Mining failed for character %s: %v", action.CharacterID, err))
		return false
	}

	return true
}

// executeDiggingAction executes a digging action
func (am *ActionManager) executeDiggingAction(action *QueuedAction) bool {
	// Extract digging data
	diggingData, ok := action.Data.(map[string]interface{})
	if !ok {
		am.log.Println(logging.LogError, "Invalid digging action data")
		return false
	}

	itemID, ok := diggingData["itemID"].(string)
	if !ok {
		am.log.Println(logging.LogError, "Missing itemID in digging action")
		return false
	}

	deltaX, _ := diggingData["deltaX"].(int)
	deltaY, _ := diggingData["deltaY"].(int)

	// Execute digging through GameManager
	err := am.gameHandler.HandleDigRequest(itemID, action.CharacterID, deltaX, deltaY)
	if err != nil {
		am.log.Println(logging.LogWarn, fmt.Sprintf("Digging failed for character %s: %v", action.CharacterID, err))
		return false
	}

	return true
}

// getMovementDelta converts direction string to delta coordinates
func (am *ActionManager) getMovementDelta(direction string) (int, int) {
	switch direction {
	case "north", "k":
		return 0, -1
	case "south", "j":
		return 0, 1
	case "east", "l":
		return 1, 0
	case "west", "h":
		return -1, 0
	case "northeast", "u":
		return 1, -1
	case "northwest", "y":
		return -1, -1
	case "southeast", "n":
		return 1, 1
	case "southwest", "b":
		return -1, 1
	default:
		am.log.Println(logging.LogWarn, fmt.Sprintf("Unknown movement direction: %s", direction))
		return 0, 0
	}
}

// createActionCallback creates a callback function for an action
func (am *ActionManager) createActionCallback(characterID string, actionType string) func(bool, interface{}) {
	return func(success bool, result interface{}) {
		// This callback can be used to send notifications to the UI
		// or perform other post-action processing
		if success {
			am.log.Println(logging.LogInfo, fmt.Sprintf("Action callback: %s completed for %s", actionType, characterID))
		} else {
			am.log.Println(logging.LogWarn, fmt.Sprintf("Action callback: %s failed for %s", actionType, characterID))
		}
	}
}

// GetMetrics returns current action processing metrics
func (am *ActionManager) GetMetrics() ActionManagerMetrics {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	return ActionManagerMetrics{
		CurrentTick:      am.tickCounter,
		ActionsProcessed: am.actionsProcessed,
		ActionsCompleted: am.actionsCompleted,
		ActionsFailed:    am.actionsFailed,
		ActiveQueues:     am.queueManager.GetQueueCount(),
		Enabled:          am.enabled,
	}
}

// sendActionQueueUpdate sends an action queue update message to the UI
func (am *ActionManager) sendActionQueueUpdate(characterID string) {
	if am.sendChannel == nil {
		return // No send channel configured
	}

	queue := am.queueManager.GetQueue(characterID)
	if queue == nil {
		return // No queue for character
	}

	// Get current action info
	var currentAction *messages.GameActionInfo
	if current := queue.GetCurrentAction(); current != nil {
		actionDef, err := am.registry.GetActionDefinition(current.Type)
		if err == nil {
			currentAction = &messages.GameActionInfo{
				Type:          current.Type,
				Description:   actionDef.Description,
				Progress:      current.GetProgress(am.tickCounter),
				TicksLeft:     current.GetTicksLeft(am.tickCounter),
				TotalTicks:    current.Duration,
				Interruptible: actionDef.Interruptible,
			}
		}
	}

	// Get queued actions info
	queuedActions := make([]*messages.GameActionInfo, 0)
	for _, action := range queue.GetQueuedActions() {
		actionDef, err := am.registry.GetActionDefinition(action.Type)
		if err == nil {
			queuedActions = append(queuedActions, &messages.GameActionInfo{
				Type:          action.Type,
				Description:   actionDef.Description,
				Progress:      0.0, // Queued actions haven't started
				TicksLeft:     action.Duration,
				TotalTicks:    action.Duration,
				Interruptible: actionDef.Interruptible,
			})
		}
	}

	// Create and send message
	updateData := messages.GameActionQueueUpdate{
		CharacterID:   characterID,
		CurrentAction: currentAction,
		QueuedActions: queuedActions,
		MaxQueue:      am.registry.GetMaxQueueSize(),
	}

	message := messages.ConnectionManagerMessage{
		Type:               messages.ConnectManager_Message_GameCommandResponse,
		RecipientConsoleID: am.getConnectionIDForCharacter(characterID),
		Data: messages.GameMessage{
			Type: messages.GM_ActionQueueUpdate,
			Data: messages.GameMessageData{
				CharacterID: characterID,
				Data:        updateData,
			},
		},
	}

	select {
	case am.sendChannel <- message:
		// Message sent successfully
	default:
		// Channel is full, log warning
		am.log.Println(logging.LogWarn, "Failed to send action queue update - channel full")
	}
}

// sendActionStarted sends an action started notification
func (am *ActionManager) sendActionStarted(characterID string, action *QueuedAction) {
	if am.sendChannel == nil {
		return
	}

	actionDef, err := am.registry.GetActionDefinition(action.Type)
	if err != nil {
		return
	}

	startedData := messages.GameActionStarted{
		CharacterID: characterID,
		ActionType:  action.Type,
		Description: actionDef.Description,
		Duration:    action.Duration,
	}

	message := messages.ConnectionManagerMessage{
		Type:               messages.ConnectManager_Message_GameCommandResponse,
		RecipientConsoleID: am.getConnectionIDForCharacter(characterID),
		Data: messages.GameMessage{
			Type: messages.GM_ActionStarted,
			Data: messages.GameMessageData{
				CharacterID: characterID,
				Data:        startedData,
			},
		},
	}

	select {
	case am.sendChannel <- message:
		// Message sent successfully
	default:
		am.log.Println(logging.LogWarn, "Failed to send action started notification - channel full")
	}
}

// sendActionCompleted sends an action completed notification
func (am *ActionManager) sendActionCompleted(characterID string, action *QueuedAction, success bool) {
	if am.sendChannel == nil {
		return
	}

	completedData := messages.GameActionCompleted{
		CharacterID: characterID,
		ActionType:  action.Type,
		Success:     success,
		Result:      nil, // Could include action results in the future
	}

	message := messages.ConnectionManagerMessage{
		Type:               messages.ConnectManager_Message_GameCommandResponse,
		RecipientConsoleID: am.getConnectionIDForCharacter(characterID),
		Data: messages.GameMessage{
			Type: messages.GM_ActionCompleted,
			Data: messages.GameMessageData{
				CharacterID: characterID,
				Data:        completedData,
			},
		},
	}

	select {
	case am.sendChannel <- message:
		// Message sent successfully
	default:
		am.log.Println(logging.LogWarn, "Failed to send action completed notification - channel full")
	}
}

// getConnectionIDForCharacter gets the connection ID for a character
// This is a placeholder - in a real implementation, this would need to be provided
// by the GameManager or another component that tracks character-to-connection mapping
func (am *ActionManager) getConnectionIDForCharacter(characterID string) string {
	// TODO: Implement proper character-to-connection mapping
	// For now, return the character ID as a placeholder
	return characterID
}

// SetEnabled enables or disables the action manager
func (am *ActionManager) SetEnabled(enabled bool) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	am.enabled = enabled
	am.log.Println(logging.LogInfo, fmt.Sprintf("ActionManager enabled: %t", enabled))
}

// IsEnabled returns whether the action manager is enabled
func (am *ActionManager) IsEnabled() bool {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	return am.enabled
}

// CleanupInactiveQueues removes inactive character queues
func (am *ActionManager) CleanupInactiveQueues(maxInactivity time.Duration) int {
	return am.queueManager.CleanupInactiveQueues(maxInactivity)
}

// ActionManagerMetrics holds metrics about action processing
type ActionManagerMetrics struct {
	CurrentTick      uint64
	ActionsProcessed uint64
	ActionsCompleted uint64
	ActionsFailed    uint64
	ActiveQueues     int
	Enabled          bool
}

// String returns a string representation of the metrics
func (amm ActionManagerMetrics) String() string {
	successRate := float64(0)
	if amm.ActionsProcessed > 0 {
		successRate = float64(amm.ActionsCompleted) / float64(amm.ActionsProcessed) * 100
	}

	return fmt.Sprintf("ActionManager: tick=%d, processed=%d, completed=%d, failed=%d, success=%.1f%%, queues=%d, enabled=%t",
		amm.CurrentTick, amm.ActionsProcessed, amm.ActionsCompleted, amm.ActionsFailed, successRate, amm.ActiveQueues, amm.Enabled)
}
