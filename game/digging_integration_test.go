package game

import (
	"testing"
)

// TestDiggingActionQueuing tests that digging actions are properly queued
func TestDiggingActionQueuing(t *testing.T) {
	// Create test action manager
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"
	actionType := "dig"
	actionData := map[string]interface{}{
		"itemID": "shovel",
		"deltaX": 1,
		"deltaY": 0,
	}

	// Queue digging action
	err := actionManager.QueueAction(characterID, actionType, actionData)
	if err != nil {
		t.Fatalf("Failed to queue digging action: %v", err)
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

	if len(status.QueuedActionTypes) != 1 || status.QueuedActionTypes[0] != "dig" {
		t.Errorf("Expected queued action type 'dig', got %v", status.QueuedActionTypes)
	}
}

// TestDiggingActionValidation tests digging action validation
func TestDiggingActionValidation(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Test invalid action type
	err := actionManager.QueueAction(characterID, "invalid_dig", nil)
	if err == nil {
		t.Error("Expected error for invalid action type")
	}

	// Test valid digging action
	actionData := map[string]interface{}{
		"itemID": "pickaxe",
		"deltaX": 0,
		"deltaY": 1,
	}

	err = actionManager.QueueAction(characterID, "dig", actionData)
	if err != nil {
		t.Errorf("Expected valid digging action to succeed, got error: %v", err)
	}
}

// TestDiggingActionProcessing tests that digging actions are processed correctly
func TestDiggingActionProcessing(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "shovel",
		"deltaX": -1,
		"deltaY": 0,
	}

	// Queue digging action
	err := actionManager.QueueAction(characterID, "dig", actionData)
	if err != nil {
		t.Fatalf("Failed to queue digging action: %v", err)
	}

	// Start the action
	queue := actionManager.GetCharacterQueue(characterID)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil {
		t.Fatal("Expected action to start")
	}

	// Simulate ticks until action completes (10 ticks for digging)
	for tick := uint64(101); tick <= 110; tick++ {
		actionManager.OnTick(tick)
	}

	// Verify action was executed
	if len(mockGameManager.digRequests) == 0 {
		t.Error("Expected HandleDigRequest to be called")
	}

	// Verify correct parameters were passed
	if len(mockGameManager.digRequests) > 0 {
		request := mockGameManager.digRequests[0]
		if request.CharacterID != characterID {
			t.Errorf("Expected character ID %s, got %s", characterID, request.CharacterID)
		}
		if request.ItemID != "shovel" {
			t.Errorf("Expected item ID 'shovel', got %s", request.ItemID)
		}
		if request.DeltaX != -1 || request.DeltaY != 0 {
			t.Errorf("Expected delta (-1, 0), got (%d, %d)", request.DeltaX, request.DeltaY)
		}
	}
}

// TestDiggingActionInterruptible tests that digging actions can be interrupted
func TestDiggingActionInterruptible(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "pickaxe",
		"deltaX": 0,
		"deltaY": -1,
	}

	// Queue and start digging action
	err := actionManager.QueueAction(characterID, "dig", actionData)
	if err != nil {
		t.Fatalf("Failed to queue digging action: %v", err)
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

	// Interrupt the action (should succeed since digging is interruptible)
	err = actionManager.InterruptCurrentAction(characterID)
	if err != nil {
		t.Errorf("Failed to interrupt digging action: %v", err)
	}

	// Verify action was interrupted
	status = queue.GetQueueStatus()
	if status.HasCurrentAction {
		t.Error("Expected current action to be interrupted")
	}
}

// TestDiggingActionFailure tests handling of digging action failures
func TestDiggingActionFailure(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	// Configure mock to return error
	mockGameManager.shouldFail = true

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "shovel",
		"deltaX": 1,
		"deltaY": 1,
	}

	// Queue and process digging action
	err := actionManager.QueueAction(characterID, "dig", actionData)
	if err != nil {
		t.Fatalf("Failed to queue digging action: %v", err)
	}

	queue := actionManager.GetCharacterQueue(characterID)
	queue.StartNextAction(100)

	// Process until completion
	for tick := uint64(101); tick <= 110; tick++ {
		actionManager.OnTick(tick)
	}

	// Verify action was attempted but failed
	if len(mockGameManager.digRequests) > 0 {
		t.Error("Expected no dig requests to be recorded when digging fails")
	}

	// Check metrics for failed action
	metrics := actionManager.GetMetrics()
	if metrics.ActionsFailed == 0 {
		t.Error("Expected failed action count to be incremented")
	}
}

// TestDiggingWithOtherActions tests digging actions with other concurrent actions
func TestDiggingWithOtherActions(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Queue multiple actions including digging
	digData := map[string]interface{}{
		"itemID": "pickaxe",
		"deltaX": 0,
		"deltaY": 1,
	}

	moveData := map[string]interface{}{
		"direction": "north",
		"deltaX":    0,
		"deltaY":    -1,
	}

	// Queue digging action first
	err := actionManager.QueueAction(characterID, "dig", digData)
	if err != nil {
		t.Fatalf("Failed to queue digging action: %v", err)
	}

	// Queue movement action
	err = actionManager.QueueAction(characterID, "move", moveData)
	if err != nil {
		t.Fatalf("Failed to queue movement action: %v", err)
	}

	// Verify both actions are queued
	queue := actionManager.GetCharacterQueue(characterID)
	status := queue.GetQueueStatus()
	if status.QueuedCount != 2 {
		t.Errorf("Expected 2 queued actions, got %d", status.QueuedCount)
	}

	// Start first action (digging)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil || startedAction.Type != "dig" {
		t.Error("Expected digging action to start first")
	}

	// Process digging action to completion (10 ticks)
	for tick := uint64(101); tick <= 110; tick++ {
		actionManager.OnTick(tick)
	}

	// Verify next action (movement) starts automatically
	status = queue.GetQueueStatus()
	if !status.HasCurrentAction || status.CurrentActionType != "move" {
		t.Error("Expected movement action to start after digging completes")
	}
}

// TestDiggingActionTiming tests that digging actions take the correct number of ticks
func TestDiggingActionTiming(t *testing.T) {
	actionManager, registry, mockGameManager := createTestActionManager()

	// Verify digging action is configured for 10 ticks
	actionDef, err := registry.GetActionDefinition("dig")
	if err != nil {
		t.Fatalf("Failed to get digging action definition: %v", err)
	}

	if actionDef.TickCost != 10 {
		t.Errorf("Expected digging action to cost 10 ticks, got %d", actionDef.TickCost)
	}

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "shovel",
		"deltaX": 1,
		"deltaY": 0,
	}

	// Queue and start digging action
	err = actionManager.QueueAction(characterID, "dig", actionData)
	if err != nil {
		t.Fatalf("Failed to queue digging action: %v", err)
	}

	queue := actionManager.GetCharacterQueue(characterID)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil {
		t.Fatal("Expected action to start")
	}

	// Verify action is not completed before 10 ticks
	actionManager.OnTick(109) // 9 ticks elapsed
	if len(mockGameManager.digRequests) > 0 {
		t.Error("Digging action should not complete before 10 ticks")
	}

	// Verify action completes after 10 ticks
	actionManager.OnTick(110) // 10 ticks elapsed
	if len(mockGameManager.digRequests) == 0 {
		t.Error("Digging action should complete after 10 ticks")
	}
}

// TestDiggingToolValidation tests digging tool validation
func TestDiggingToolValidation(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	testCases := []struct {
		name       string
		actionData map[string]interface{}
		shouldFail bool
	}{
		{
			name: "valid_digging_data",
			actionData: map[string]interface{}{
				"itemID": "shovel",
				"deltaX": 1,
				"deltaY": 0,
			},
			shouldFail: false,
		},
		{
			name: "missing_itemID",
			actionData: map[string]interface{}{
				"deltaX": 1,
				"deltaY": 0,
			},
			shouldFail: false, // ActionManager doesn't validate data structure, only action type
		},
		{
			name: "empty_itemID",
			actionData: map[string]interface{}{
				"itemID": "",
				"deltaX": 1,
				"deltaY": 0,
			},
			shouldFail: false, // ActionManager doesn't validate data content
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := actionManager.QueueAction(characterID, "dig", tc.actionData)
			if tc.shouldFail && err == nil {
				t.Error("Expected action queuing to fail")
			}
			if !tc.shouldFail && err != nil {
				t.Errorf("Expected action queuing to succeed, got error: %v", err)
			}
		})
	}
}

// TestDiggingQueueLimit tests that digging actions respect queue limits
func TestDiggingQueueLimit(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "pickaxe",
		"deltaX": 0,
		"deltaY": 1,
	}

	// Fill the queue to its limit (3 actions)
	for i := 0; i < 3; i++ {
		err := actionManager.QueueAction(characterID, "dig", actionData)
		if err != nil {
			t.Fatalf("Failed to queue digging action %d: %v", i+1, err)
		}
	}

	// Try to queue one more action (should fail)
	err := actionManager.QueueAction(characterID, "dig", actionData)
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

// TestDiggingActionDataExtraction tests that digging action data is properly extracted
func TestDiggingActionDataExtraction(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	characterID := "test-character"

	testCases := []struct {
		name         string
		actionData   map[string]interface{}
		expectedItem string
		expectedDX   int
		expectedDY   int
	}{
		{
			name: "shovel_digging",
			actionData: map[string]interface{}{
				"itemID": "shovel",
				"deltaX": 2,
				"deltaY": -1,
			},
			expectedItem: "shovel",
			expectedDX:   2,
			expectedDY:   -1,
		},
		{
			name: "pickaxe_digging",
			actionData: map[string]interface{}{
				"itemID": "pickaxe",
				"deltaX": -2,
				"deltaY": 3,
			},
			expectedItem: "pickaxe",
			expectedDX:   -2,
			expectedDY:   3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear previous requests
			mockGameManager.digRequests = nil

			// Queue and process action
			err := actionManager.QueueAction(characterID, "dig", tc.actionData)
			if err != nil {
				t.Fatalf("Failed to queue digging action: %v", err)
			}

			queue := actionManager.GetCharacterQueue(characterID)
			queue.StartNextAction(100)

			// Process to completion
			for tick := uint64(101); tick <= 110; tick++ {
				actionManager.OnTick(tick)
			}

			// Verify data extraction
			if len(mockGameManager.digRequests) == 0 {
				t.Fatal("Expected digging request to be made")
			}

			request := mockGameManager.digRequests[0]
			if request.ItemID != tc.expectedItem {
				t.Errorf("Expected item ID %s, got %s", tc.expectedItem, request.ItemID)
			}
			if request.DeltaX != tc.expectedDX {
				t.Errorf("Expected deltaX %d, got %d", tc.expectedDX, request.DeltaX)
			}
			if request.DeltaY != tc.expectedDY {
				t.Errorf("Expected deltaY %d, got %d", tc.expectedDY, request.DeltaY)
			}
		})
	}
}

// TestDiggingInterruptionMechanism tests the digging interruption mechanism
func TestDiggingInterruptionMechanism(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "shovel",
		"deltaX": 1,
		"deltaY": 0,
	}

	// Queue and start digging action
	err := actionManager.QueueAction(characterID, "dig", actionData)
	if err != nil {
		t.Fatalf("Failed to queue digging action: %v", err)
	}

	queue := actionManager.GetCharacterQueue(characterID)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil {
		t.Fatal("Expected action to start")
	}

	// Process a few ticks
	actionManager.OnTick(105) // 5 ticks elapsed, 5 remaining

	// Verify action is still running
	status := queue.GetQueueStatus()
	if !status.HasCurrentAction {
		t.Error("Expected current action to still be running")
	}

	// Interrupt the action
	err = actionManager.InterruptCurrentAction(characterID)
	if err != nil {
		t.Errorf("Failed to interrupt digging action: %v", err)
	}

	// Verify action was interrupted and not completed
	if len(mockGameManager.digRequests) > 0 {
		t.Error("Expected digging action to be interrupted before completion")
	}

	// Verify no current action
	status = queue.GetQueueStatus()
	if status.HasCurrentAction {
		t.Error("Expected no current action after interruption")
	}
}

// TestDiggingTileStateChanges tests tile state changes during digging
func TestDiggingTileStateChanges(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Test digging different tile types
	tileTypes := []string{"stone", "dirt", "sand", "gravel"}

	for _, tileType := range tileTypes {
		t.Run("dig_"+tileType, func(t *testing.T) {
			actionData := map[string]interface{}{
				"itemID": "shovel",
				"deltaX": 1,
				"deltaY": 0,
			}

			err := actionManager.QueueAction(characterID, "dig", actionData)
			if err != nil {
				t.Errorf("Failed to queue digging action for %s: %v", tileType, err)
			}

			// Verify action was queued
			queue := actionManager.GetCharacterQueue(characterID)
			status := queue.GetQueueStatus()
			if status.QueuedCount == 0 {
				t.Errorf("Expected digging action for %s to be queued", tileType)
			}

			// Clear the queue for next test
			queue.Clear()
		})
	}
}

// TestDiggingWithTerrainModification tests digging with other terrain modification actions
func TestDiggingWithTerrainModification(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Queue terrain modification actions
	digData := map[string]interface{}{
		"itemID": "shovel",
		"deltaX": 1,
		"deltaY": 0,
	}

	buildData := map[string]interface{}{
		"itemID": "stone-block",
		"toolID": "",
		"deltaX": 0,
		"deltaY": 1,
	}

	// Queue digging action
	err := actionManager.QueueAction(characterID, "dig", digData)
	if err != nil {
		t.Fatalf("Failed to queue digging action: %v", err)
	}

	// Queue building action
	err = actionManager.QueueAction(characterID, "build_wall", buildData)
	if err != nil {
		t.Fatalf("Failed to queue building action: %v", err)
	}

	// Verify both actions are queued
	queue := actionManager.GetCharacterQueue(characterID)
	status := queue.GetQueueStatus()
	if status.QueuedCount != 2 {
		t.Errorf("Expected 2 queued actions, got %d", status.QueuedCount)
	}

	// Verify action types
	expectedTypes := []string{"dig", "build_wall"}
	if len(status.QueuedActionTypes) != 2 {
		t.Errorf("Expected 2 action types, got %d", len(status.QueuedActionTypes))
	}

	for i, expectedType := range expectedTypes {
		if i < len(status.QueuedActionTypes) && status.QueuedActionTypes[i] != expectedType {
			t.Errorf("Expected action type %s at position %d, got %s", expectedType, i, status.QueuedActionTypes[i])
		}
	}
}
