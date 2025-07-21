package game

import (
	"testing"
)

// TestMovementActionQueuing tests that movement actions are properly queued
func TestMovementActionQueuing(t *testing.T) {
	// Create test action manager
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"
	actionType := "move"
	actionData := map[string]interface{}{
		"direction": "north",
		"deltaX":    0,
		"deltaY":    -1,
	}

	// Queue movement action
	err := actionManager.QueueAction(characterID, actionType, actionData)
	if err != nil {
		t.Fatalf("Failed to queue movement action: %v", err)
	}

	// Verify action was queued
	queue := actionManager.GetCharacterQueue(characterID)
	if queue == nil {
		t.Fatal("Expected character queue to exist")
	}

	status := queue.GetQueueStatus()
	if status.QueuedCount != 1 {
		t.Errorf("Expected 1 queued action, got %d", status.QueuedCount)
	}

	if len(status.QueuedActionTypes) != 1 || status.QueuedActionTypes[0] != "move" {
		t.Errorf("Expected queued action type 'move', got %v", status.QueuedActionTypes)
	}
}

// TestMovementActionValidation tests movement action validation
func TestMovementActionValidation(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Test invalid action type
	err := actionManager.QueueAction(characterID, "invalid_action", nil)
	if err == nil {
		t.Error("Expected error for invalid action type")
	}

	// Test valid movement action
	actionData := map[string]interface{}{
		"direction": "south",
		"deltaX":    0,
		"deltaY":    1,
	}

	err = actionManager.QueueAction(characterID, "move", actionData)
	if err != nil {
		t.Errorf("Expected valid movement action to succeed, got error: %v", err)
	}
}

// TestMovementActionProcessing tests that movement actions are processed correctly
func TestMovementActionProcessing(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"direction": "east",
		"deltaX":    1,
		"deltaY":    0,
	}

	// Queue movement action
	err := actionManager.QueueAction(characterID, "move", actionData)
	if err != nil {
		t.Fatalf("Failed to queue movement action: %v", err)
	}

	// Start the action
	queue := actionManager.GetCharacterQueue(characterID)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil {
		t.Fatal("Expected action to start")
	}

	// Simulate ticks until action completes (5 ticks for movement)
	for tick := uint64(101); tick <= 105; tick++ {
		actionManager.OnTick(tick)
	}

	// Verify action was executed
	if len(mockGameManager.moveRequests) == 0 {
		t.Error("Expected HandleMovePlayerRequest to be called")
	}

	// Verify correct parameters were passed
	if len(mockGameManager.moveRequests) > 0 {
		request := mockGameManager.moveRequests[0]
		if request.CharacterID != characterID {
			t.Errorf("Expected character ID %s, got %s", characterID, request.CharacterID)
		}
		if request.DeltaX != 1 || request.DeltaY != 0 {
			t.Errorf("Expected delta (1, 0), got (%d, %d)", request.DeltaX, request.DeltaY)
		}
	}
}

// TestMovementActionInterruption tests movement action interruption
func TestMovementActionInterruption(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"direction": "west",
		"deltaX":    -1,
		"deltaY":    0,
	}

	// Queue and start movement action
	err := actionManager.QueueAction(characterID, "move", actionData)
	if err != nil {
		t.Fatalf("Failed to queue movement action: %v", err)
	}

	queue := actionManager.GetCharacterQueue(characterID)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil {
		t.Fatal("Expected action to start")
	}

	// Verify action is running
	status := queue.GetQueueStatus()
	if !status.HasCurrentAction {
		t.Error("Expected current action to be running")
	}

	// Interrupt the action
	err = actionManager.InterruptCurrentAction(characterID)
	if err != nil {
		t.Errorf("Failed to interrupt movement action: %v", err)
	}

	// Verify action was interrupted
	status = queue.GetQueueStatus()
	if status.HasCurrentAction {
		t.Error("Expected current action to be interrupted")
	}
}

// TestMovementActionFailure tests handling of movement action failures
func TestMovementActionFailure(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	// Configure mock to return error
	mockGameManager.shouldFail = true

	characterID := "test-character"
	actionData := map[string]interface{}{
		"direction": "north",
		"deltaX":    0,
		"deltaY":    -1,
	}

	// Queue and process movement action
	err := actionManager.QueueAction(characterID, "move", actionData)
	if err != nil {
		t.Fatalf("Failed to queue movement action: %v", err)
	}

	queue := actionManager.GetCharacterQueue(characterID)
	queue.StartNextAction(100)

	// Process until completion
	for tick := uint64(101); tick <= 105; tick++ {
		actionManager.OnTick(tick)
	}

	// Verify action was attempted but failed
	if len(mockGameManager.moveRequests) > 0 {
		t.Error("Expected no move requests to be recorded when movement fails")
	}

	// Check metrics for failed action
	metrics := actionManager.GetMetrics()
	if metrics.ActionsFailed == 0 {
		t.Error("Expected failed action count to be incremented")
	}
}

// TestMovementWithOtherActions tests movement actions with other concurrent actions
func TestMovementWithOtherActions(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Queue multiple actions including movement
	moveData := map[string]interface{}{
		"direction": "north",
		"deltaX":    0,
		"deltaY":    -1,
	}

	digData := map[string]interface{}{
		"itemID": "test-tool",
		"deltaX": 1,
		"deltaY": 0,
	}

	// Queue movement action
	err := actionManager.QueueAction(characterID, "move", moveData)
	if err != nil {
		t.Fatalf("Failed to queue movement action: %v", err)
	}

	// Queue dig action
	err = actionManager.QueueAction(characterID, "dig", digData)
	if err != nil {
		t.Fatalf("Failed to queue dig action: %v", err)
	}

	// Verify both actions are queued
	queue := actionManager.GetCharacterQueue(characterID)
	status := queue.GetQueueStatus()
	if status.QueuedCount != 2 {
		t.Errorf("Expected 2 queued actions, got %d", status.QueuedCount)
	}

	// Start first action (movement)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil || startedAction.Type != "move" {
		t.Error("Expected movement action to start first")
	}

	// Process movement action to completion
	for tick := uint64(101); tick <= 105; tick++ {
		actionManager.OnTick(tick)
	}

	// Verify next action (dig) starts automatically
	status = queue.GetQueueStatus()
	if !status.HasCurrentAction || status.CurrentActionType != "dig" {
		t.Error("Expected dig action to start after movement completes")
	}
}

// TestMovementDirectionConversion tests direction string to delta conversion
func TestMovementDirectionConversion(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	testCases := []struct {
		direction string
		expectedX int
		expectedY int
	}{
		{"north", 0, -1},
		{"south", 0, 1},
		{"east", 1, 0},
		{"west", -1, 0},
		{"northeast", 1, -1},
		{"northwest", -1, -1},
		{"southeast", 1, 1},
		{"southwest", -1, 1},
	}

	for _, tc := range testCases {
		deltaX, deltaY := actionManager.getMovementDelta(tc.direction)
		if deltaX != tc.expectedX || deltaY != tc.expectedY {
			t.Errorf("Direction %s: expected (%d, %d), got (%d, %d)",
				tc.direction, tc.expectedX, tc.expectedY, deltaX, deltaY)
		}
	}
}

// TestMovementActionTiming tests that movement actions take the correct number of ticks
func TestMovementActionTiming(t *testing.T) {
	actionManager, registry, mockGameManager := createTestActionManager()

	// Verify movement action is configured for 5 ticks
	actionDef, err := registry.GetActionDefinition("move")
	if err != nil {
		t.Fatalf("Failed to get movement action definition: %v", err)
	}

	if actionDef.TickCost != 5 {
		t.Errorf("Expected movement action to cost 5 ticks, got %d", actionDef.TickCost)
	}

	characterID := "test-character"
	actionData := map[string]interface{}{
		"direction": "north",
		"deltaX":    0,
		"deltaY":    -1,
	}

	// Queue and start movement action
	err = actionManager.QueueAction(characterID, "move", actionData)
	if err != nil {
		t.Fatalf("Failed to queue movement action: %v", err)
	}

	queue := actionManager.GetCharacterQueue(characterID)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil {
		t.Fatal("Expected action to start")
	}

	// Verify action is not completed before 5 ticks
	actionManager.OnTick(104) // 4 ticks elapsed
	if len(mockGameManager.moveRequests) > 0 {
		t.Error("Movement action should not complete before 5 ticks")
	}

	// Verify action completes after 5 ticks
	actionManager.OnTick(105) // 5 ticks elapsed
	if len(mockGameManager.moveRequests) == 0 {
		t.Error("Movement action should complete after 5 ticks")
	}
}

// TestMovementQueueLimit tests that movement actions respect queue limits
func TestMovementQueueLimit(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"direction": "north",
		"deltaX":    0,
		"deltaY":    -1,
	}

	// Fill the queue to its limit (3 actions)
	for i := 0; i < 3; i++ {
		err := actionManager.QueueAction(characterID, "move", actionData)
		if err != nil {
			t.Fatalf("Failed to queue movement action %d: %v", i+1, err)
		}
	}

	// Try to queue one more action (should fail)
	err := actionManager.QueueAction(characterID, "move", actionData)
	if err == nil {
		t.Error("Expected error when queuing action beyond limit")
	}

	// Verify queue is at capacity
	queue := actionManager.GetCharacterQueue(characterID)
	status := queue.GetQueueStatus()
	if status.QueuedCount != 3 {
		t.Errorf("Expected queue to be at capacity (3), got %d", status.QueuedCount)
	}
}
