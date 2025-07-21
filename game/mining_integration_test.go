package game

import (
	"testing"
)

// TestMiningActionQueuing tests that mining actions are properly queued
func TestMiningActionQueuing(t *testing.T) {
	// Create test action manager
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"
	actionType := "mine"
	actionData := map[string]interface{}{
		"itemID": "iron_ore",
		"toolID": "pickaxe",
		"deltaX": 1,
		"deltaY": 0,
	}

	// Queue mining action
	err := actionManager.QueueAction(characterID, actionType, actionData)
	if err != nil {
		t.Fatalf("Failed to queue mining action: %v", err)
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

	if len(status.QueuedActionTypes) != 1 || status.QueuedActionTypes[0] != "mine" {
		t.Errorf("Expected queued action type 'mine', got %v", status.QueuedActionTypes)
	}
}

// TestMiningActionValidation tests mining action validation
func TestMiningActionValidation(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Test invalid action type
	err := actionManager.QueueAction(characterID, "invalid_mine", nil)
	if err == nil {
		t.Error("Expected error for invalid action type")
	}

	// Test valid mining action
	actionData := map[string]interface{}{
		"itemID": "stone",
		"toolID": "pickaxe",
		"deltaX": 0,
		"deltaY": 1,
	}

	err = actionManager.QueueAction(characterID, "mine", actionData)
	if err != nil {
		t.Errorf("Expected valid mining action to succeed, got error: %v", err)
	}
}

// TestMiningActionProcessing tests that mining actions are processed correctly
func TestMiningActionProcessing(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "iron_ore",
		"toolID": "pickaxe",
		"deltaX": -1,
		"deltaY": 0,
	}

	// Queue mining action
	err := actionManager.QueueAction(characterID, "mine", actionData)
	if err != nil {
		t.Fatalf("Failed to queue mining action: %v", err)
	}

	// Start the action
	queue := actionManager.GetCharacterQueue(characterID)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil {
		t.Fatal("Expected action to start")
	}

	// Simulate ticks until action completes (20 ticks for mining)
	for tick := uint64(101); tick <= 120; tick++ {
		actionManager.OnTick(tick)
	}

	// Verify action was executed
	if len(mockGameManager.mineRequests) == 0 {
		t.Error("Expected HandleMineRequest to be called")
	}

	// Verify correct parameters were passed
	if len(mockGameManager.mineRequests) > 0 {
		request := mockGameManager.mineRequests[0]
		if request.CharacterID != characterID {
			t.Errorf("Expected character ID %s, got %s", characterID, request.CharacterID)
		}
		if request.ItemID != "iron_ore" {
			t.Errorf("Expected item ID 'iron_ore', got %s", request.ItemID)
		}
		if request.ToolID != "pickaxe" {
			t.Errorf("Expected tool ID 'pickaxe', got %s", request.ToolID)
		}
		if request.DeltaX != -1 || request.DeltaY != 0 {
			t.Errorf("Expected delta (-1, 0), got (%d, %d)", request.DeltaX, request.DeltaY)
		}
	}
}

// TestMiningActionNonInterruptible tests that mining actions cannot be interrupted
func TestMiningActionNonInterruptible(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "stone",
		"toolID": "",
		"deltaX": 0,
		"deltaY": -1,
	}

	// Queue and start mining action
	err := actionManager.QueueAction(characterID, "mine", actionData)
	if err != nil {
		t.Fatalf("Failed to queue mining action: %v", err)
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

	// Try to interrupt the action (should fail)
	err = actionManager.InterruptCurrentAction(characterID)
	if err == nil {
		t.Error("Expected error when trying to interrupt non-interruptible mining action")
	}

	// Verify action is still running
	status = queue.GetQueueStatus()
	if !status.HasCurrentAction {
		t.Error("Expected current action to still be running after failed interruption")
	}
}

// TestMiningActionFailure tests handling of mining action failures
func TestMiningActionFailure(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	// Configure mock to return error
	mockGameManager.shouldFail = true

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "iron_ore",
		"toolID": "pickaxe",
		"deltaX": 1,
		"deltaY": 1,
	}

	// Queue and process mining action
	err := actionManager.QueueAction(characterID, "mine", actionData)
	if err != nil {
		t.Fatalf("Failed to queue mining action: %v", err)
	}

	queue := actionManager.GetCharacterQueue(characterID)
	queue.StartNextAction(100)

	// Process until completion
	for tick := uint64(101); tick <= 120; tick++ {
		actionManager.OnTick(tick)
	}

	// Verify action was attempted but failed
	if len(mockGameManager.mineRequests) > 0 {
		t.Error("Expected no mine requests to be recorded when mining fails")
	}

	// Check metrics for failed action
	metrics := actionManager.GetMetrics()
	if metrics.ActionsFailed == 0 {
		t.Error("Expected failed action count to be incremented")
	}
}

// TestMiningWithOtherActions tests mining actions with other concurrent actions
func TestMiningWithOtherActions(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Queue multiple actions including mining
	mineData := map[string]interface{}{
		"itemID": "stone",
		"toolID": "pickaxe",
		"deltaX": 0,
		"deltaY": 1,
	}

	moveData := map[string]interface{}{
		"direction": "north",
		"deltaX":    0,
		"deltaY":    -1,
	}

	// Queue mining action first
	err := actionManager.QueueAction(characterID, "mine", mineData)
	if err != nil {
		t.Fatalf("Failed to queue mining action: %v", err)
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

	// Start first action (mining)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil || startedAction.Type != "mine" {
		t.Error("Expected mining action to start first")
	}

	// Process mining action to completion (20 ticks)
	for tick := uint64(101); tick <= 120; tick++ {
		actionManager.OnTick(tick)
	}

	// Verify next action (movement) starts automatically
	status = queue.GetQueueStatus()
	if !status.HasCurrentAction || status.CurrentActionType != "move" {
		t.Error("Expected movement action to start after mining completes")
	}
}

// TestMiningActionTiming tests that mining actions take the correct number of ticks
func TestMiningActionTiming(t *testing.T) {
	actionManager, registry, mockGameManager := createTestActionManager()

	// Verify mining action is configured for 20 ticks
	actionDef, err := registry.GetActionDefinition("mine")
	if err != nil {
		t.Fatalf("Failed to get mining action definition: %v", err)
	}

	if actionDef.TickCost != 20 {
		t.Errorf("Expected mining action to cost 20 ticks, got %d", actionDef.TickCost)
	}

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "iron_ore",
		"toolID": "pickaxe",
		"deltaX": 1,
		"deltaY": 0,
	}

	// Queue and start mining action
	err = actionManager.QueueAction(characterID, "mine", actionData)
	if err != nil {
		t.Fatalf("Failed to queue mining action: %v", err)
	}

	queue := actionManager.GetCharacterQueue(characterID)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil {
		t.Fatal("Expected action to start")
	}

	// Verify action is not completed before 20 ticks
	actionManager.OnTick(119) // 19 ticks elapsed
	if len(mockGameManager.mineRequests) > 0 {
		t.Error("Mining action should not complete before 20 ticks")
	}

	// Verify action completes after 20 ticks
	actionManager.OnTick(120) // 20 ticks elapsed
	if len(mockGameManager.mineRequests) == 0 {
		t.Error("Mining action should complete after 20 ticks")
	}
}

// TestMiningToolValidation tests mining tool validation
func TestMiningToolValidation(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	testCases := []struct {
		name       string
		actionData map[string]interface{}
		shouldFail bool
	}{
		{
			name: "valid_mining_data_with_tool",
			actionData: map[string]interface{}{
				"itemID": "iron_ore",
				"toolID": "pickaxe",
				"deltaX": 1,
				"deltaY": 0,
			},
			shouldFail: false,
		},
		{
			name: "valid_mining_data_without_tool",
			actionData: map[string]interface{}{
				"itemID": "stone",
				"toolID": "",
				"deltaX": 1,
				"deltaY": 0,
			},
			shouldFail: false,
		},
		{
			name: "missing_itemID",
			actionData: map[string]interface{}{
				"toolID": "pickaxe",
				"deltaX": 1,
				"deltaY": 0,
			},
			shouldFail: false, // ActionManager doesn't validate data structure, only action type
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := actionManager.QueueAction(characterID, "mine", tc.actionData)
			if tc.shouldFail && err == nil {
				t.Error("Expected action queuing to fail")
			}
			if !tc.shouldFail && err != nil {
				t.Errorf("Expected action queuing to succeed, got error: %v", err)
			}
		})
	}
}

// TestMiningQueueLimit tests that mining actions respect queue limits
func TestMiningQueueLimit(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "stone",
		"toolID": "pickaxe",
		"deltaX": 0,
		"deltaY": 1,
	}

	// Fill the queue to its limit (3 actions)
	for i := 0; i < 3; i++ {
		err := actionManager.QueueAction(characterID, "mine", actionData)
		if err != nil {
			t.Fatalf("Failed to queue mining action %d: %v", i+1, err)
		}
	}

	// Try to queue one more action (should fail)
	err := actionManager.QueueAction(characterID, "mine", actionData)
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

// TestMiningActionDataExtraction tests that mining action data is properly extracted
func TestMiningActionDataExtraction(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	characterID := "test-character"

	testCases := []struct {
		name         string
		actionData   map[string]interface{}
		expectedItem string
		expectedTool string
		expectedDX   int
		expectedDY   int
	}{
		{
			name: "complete_data",
			actionData: map[string]interface{}{
				"itemID": "iron_ore",
				"toolID": "pickaxe",
				"deltaX": 2,
				"deltaY": -1,
			},
			expectedItem: "iron_ore",
			expectedTool: "pickaxe",
			expectedDX:   2,
			expectedDY:   -1,
		},
		{
			name: "no_tool",
			actionData: map[string]interface{}{
				"itemID": "stone",
				"deltaX": -2,
				"deltaY": 3,
			},
			expectedItem: "stone",
			expectedTool: "",
			expectedDX:   -2,
			expectedDY:   3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear previous requests
			mockGameManager.mineRequests = nil

			// Queue and process action
			err := actionManager.QueueAction(characterID, "mine", tc.actionData)
			if err != nil {
				t.Fatalf("Failed to queue mining action: %v", err)
			}

			queue := actionManager.GetCharacterQueue(characterID)
			queue.StartNextAction(100)

			// Process to completion
			for tick := uint64(101); tick <= 120; tick++ {
				actionManager.OnTick(tick)
			}

			// Verify data extraction
			if len(mockGameManager.mineRequests) == 0 {
				t.Fatal("Expected mining request to be made")
			}

			request := mockGameManager.mineRequests[0]
			if request.ItemID != tc.expectedItem {
				t.Errorf("Expected item ID %s, got %s", tc.expectedItem, request.ItemID)
			}
			if request.ToolID != tc.expectedTool {
				t.Errorf("Expected tool ID %s, got %s", tc.expectedTool, request.ToolID)
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

// TestMiningResourceGeneration tests mining resource generation logic
func TestMiningResourceGeneration(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	testCases := []struct {
		name           string
		requestedItem  string
		expectedResult bool
	}{
		{
			name:           "request_iron_ore",
			requestedItem:  "iron_ore",
			expectedResult: true,
		},
		{
			name:           "request_stone",
			requestedItem:  "stone",
			expectedResult: true,
		},
		{
			name:           "request_empty",
			requestedItem:  "",
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actionData := map[string]interface{}{
				"itemID": tc.requestedItem,
				"toolID": "pickaxe",
				"deltaX": 1,
				"deltaY": 0,
			}

			err := actionManager.QueueAction(characterID, "mine", actionData)
			if tc.expectedResult && err != nil {
				t.Errorf("Expected mining action to succeed, got error: %v", err)
			}
			if !tc.expectedResult && err == nil {
				t.Error("Expected mining action to fail")
			}
		})
	}
}

// TestMiningWithDifferentTileTypes tests mining different types of tiles
func TestMiningWithDifferentTileTypes(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Test mining different resource types
	resourceTypes := []string{"stone", "iron_ore", "coal", "gold_ore"}

	for _, resourceType := range resourceTypes {
		t.Run("mine_"+resourceType, func(t *testing.T) {
			actionData := map[string]interface{}{
				"itemID": resourceType,
				"toolID": "pickaxe",
				"deltaX": 1,
				"deltaY": 0,
			}

			err := actionManager.QueueAction(characterID, "mine", actionData)
			if err != nil {
				t.Errorf("Failed to queue mining action for %s: %v", resourceType, err)
			}

			// Verify action was queued
			queue := actionManager.GetCharacterQueue(characterID)
			status := queue.GetQueueStatus()
			if status.QueuedCount == 0 {
				t.Errorf("Expected mining action for %s to be queued", resourceType)
			}

			// Clear the queue for next test
			queue.Clear()
		})
	}
}
