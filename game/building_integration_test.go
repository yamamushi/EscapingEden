package game

import (
	"testing"
)

// TestBuildingActionQueuing tests that building actions are properly queued
func TestBuildingActionQueuing(t *testing.T) {
	// Create test action manager
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"
	actionType := "build_wall"
	actionData := map[string]interface{}{
		"itemID": "stone-block",
		"toolID": "hammer",
		"deltaX": 1,
		"deltaY": 0,
	}

	// Queue building action
	err := actionManager.QueueAction(characterID, actionType, actionData)
	if err != nil {
		t.Fatalf("Failed to queue building action: %v", err)
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

	if len(status.QueuedActionTypes) != 1 || status.QueuedActionTypes[0] != "build_wall" {
		t.Errorf("Expected queued action type 'build_wall', got %v", status.QueuedActionTypes)
	}
}

// TestBuildingActionValidation tests building action validation
func TestBuildingActionValidation(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Test invalid action type
	err := actionManager.QueueAction(characterID, "invalid_build", nil)
	if err == nil {
		t.Error("Expected error for invalid action type")
	}

	// Test valid building action
	actionData := map[string]interface{}{
		"itemID": "wood-plank",
		"toolID": "",
		"deltaX": 0,
		"deltaY": 1,
	}

	err = actionManager.QueueAction(characterID, "build_wall", actionData)
	if err != nil {
		t.Errorf("Expected valid building action to succeed, got error: %v", err)
	}
}

// TestBuildingActionProcessing tests that building actions are processed correctly
func TestBuildingActionProcessing(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "brick",
		"toolID": "trowel",
		"deltaX": -1,
		"deltaY": 0,
	}

	// Queue building action
	err := actionManager.QueueAction(characterID, "build_wall", actionData)
	if err != nil {
		t.Fatalf("Failed to queue building action: %v", err)
	}

	// Start the action
	queue := actionManager.GetCharacterQueue(characterID)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil {
		t.Fatal("Expected action to start")
	}

	// Simulate ticks until action completes (15 ticks for building)
	for tick := uint64(101); tick <= 115; tick++ {
		actionManager.OnTick(tick)
	}

	// Verify action was executed
	if len(mockGameManager.buildRequests) == 0 {
		t.Error("Expected HandleBuildWallRequest to be called")
	}

	// Verify correct parameters were passed
	if len(mockGameManager.buildRequests) > 0 {
		request := mockGameManager.buildRequests[0]
		if request.CharacterID != characterID {
			t.Errorf("Expected character ID %s, got %s", characterID, request.CharacterID)
		}
		if request.ItemID != "brick" {
			t.Errorf("Expected item ID 'brick', got %s", request.ItemID)
		}
		if request.ToolID != "trowel" {
			t.Errorf("Expected tool ID 'trowel', got %s", request.ToolID)
		}
		if request.DeltaX != -1 || request.DeltaY != 0 {
			t.Errorf("Expected delta (-1, 0), got (%d, %d)", request.DeltaX, request.DeltaY)
		}
	}
}

// TestBuildingActionNonInterruptible tests that building actions cannot be interrupted
func TestBuildingActionNonInterruptible(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "stone-block",
		"toolID": "",
		"deltaX": 0,
		"deltaY": -1,
	}

	// Queue and start building action
	err := actionManager.QueueAction(characterID, "build_wall", actionData)
	if err != nil {
		t.Fatalf("Failed to queue building action: %v", err)
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
		t.Error("Expected error when trying to interrupt non-interruptible building action")
	}

	// Verify action is still running
	status = queue.GetQueueStatus()
	if !status.HasCurrentAction {
		t.Error("Expected current action to still be running after failed interruption")
	}
}

// TestBuildingActionFailure tests handling of building action failures
func TestBuildingActionFailure(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	// Configure mock to return error
	mockGameManager.shouldFail = true

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "wood-plank",
		"toolID": "",
		"deltaX": 1,
		"deltaY": 1,
	}

	// Queue and process building action
	err := actionManager.QueueAction(characterID, "build_wall", actionData)
	if err != nil {
		t.Fatalf("Failed to queue building action: %v", err)
	}

	queue := actionManager.GetCharacterQueue(characterID)
	queue.StartNextAction(100)

	// Process until completion
	for tick := uint64(101); tick <= 115; tick++ {
		actionManager.OnTick(tick)
	}

	// Verify action was attempted but failed
	if len(mockGameManager.buildRequests) > 0 {
		t.Error("Expected no build requests to be recorded when building fails")
	}

	// Check metrics for failed action
	metrics := actionManager.GetMetrics()
	if metrics.ActionsFailed == 0 {
		t.Error("Expected failed action count to be incremented")
	}
}

// TestBuildingWithOtherActions tests building actions with other concurrent actions
func TestBuildingWithOtherActions(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Queue multiple actions including building
	buildData := map[string]interface{}{
		"itemID": "stone-block",
		"toolID": "hammer",
		"deltaX": 0,
		"deltaY": 1,
	}

	moveData := map[string]interface{}{
		"direction": "north",
		"deltaX":    0,
		"deltaY":    -1,
	}

	// Queue building action first
	err := actionManager.QueueAction(characterID, "build_wall", buildData)
	if err != nil {
		t.Fatalf("Failed to queue building action: %v", err)
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

	// Start first action (building)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil || startedAction.Type != "build_wall" {
		t.Error("Expected building action to start first")
	}

	// Process building action to completion (15 ticks)
	for tick := uint64(101); tick <= 115; tick++ {
		actionManager.OnTick(tick)
	}

	// Verify next action (movement) starts automatically
	status = queue.GetQueueStatus()
	if !status.HasCurrentAction || status.CurrentActionType != "move" {
		t.Error("Expected movement action to start after building completes")
	}
}

// TestBuildingActionTiming tests that building actions take the correct number of ticks
func TestBuildingActionTiming(t *testing.T) {
	actionManager, registry, mockGameManager := createTestActionManager()

	// Verify building action is configured for 15 ticks
	actionDef, err := registry.GetActionDefinition("build_wall")
	if err != nil {
		t.Fatalf("Failed to get building action definition: %v", err)
	}

	if actionDef.TickCost != 15 {
		t.Errorf("Expected building action to cost 15 ticks, got %d", actionDef.TickCost)
	}

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "brick",
		"toolID": "",
		"deltaX": 1,
		"deltaY": 0,
	}

	// Queue and start building action
	err = actionManager.QueueAction(characterID, "build_wall", actionData)
	if err != nil {
		t.Fatalf("Failed to queue building action: %v", err)
	}

	queue := actionManager.GetCharacterQueue(characterID)
	startedAction := queue.StartNextAction(100)
	if startedAction == nil {
		t.Fatal("Expected action to start")
	}

	// Verify action is not completed before 15 ticks
	actionManager.OnTick(114) // 14 ticks elapsed
	if len(mockGameManager.buildRequests) > 0 {
		t.Error("Building action should not complete before 15 ticks")
	}

	// Verify action completes after 15 ticks
	actionManager.OnTick(115) // 15 ticks elapsed
	if len(mockGameManager.buildRequests) == 0 {
		t.Error("Building action should complete after 15 ticks")
	}
}

// TestBuildingMaterialValidation tests building material validation
func TestBuildingMaterialValidation(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	testCases := []struct {
		name       string
		actionData map[string]interface{}
		shouldFail bool
	}{
		{
			name: "valid_building_data",
			actionData: map[string]interface{}{
				"itemID": "stone-block",
				"toolID": "hammer",
				"deltaX": 1,
				"deltaY": 0,
			},
			shouldFail: false,
		},
		{
			name: "missing_itemID",
			actionData: map[string]interface{}{
				"toolID": "hammer",
				"deltaX": 1,
				"deltaY": 0,
			},
			shouldFail: false, // ActionManager doesn't validate data structure, only action type
		},
		{
			name: "empty_itemID",
			actionData: map[string]interface{}{
				"itemID": "",
				"toolID": "hammer",
				"deltaX": 1,
				"deltaY": 0,
			},
			shouldFail: false, // ActionManager doesn't validate data content
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := actionManager.QueueAction(characterID, "build_wall", tc.actionData)
			if tc.shouldFail && err == nil {
				t.Error("Expected action queuing to fail")
			}
			if !tc.shouldFail && err != nil {
				t.Errorf("Expected action queuing to succeed, got error: %v", err)
			}
		})
	}
}

// TestBuildingQueueLimit tests that building actions respect queue limits
func TestBuildingQueueLimit(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"
	actionData := map[string]interface{}{
		"itemID": "wood-plank",
		"toolID": "",
		"deltaX": 0,
		"deltaY": 1,
	}

	// Fill the queue to its limit (3 actions)
	for i := 0; i < 3; i++ {
		err := actionManager.QueueAction(characterID, "build_wall", actionData)
		if err != nil {
			t.Fatalf("Failed to queue building action %d: %v", i+1, err)
		}
	}

	// Try to queue one more action (should fail)
	err := actionManager.QueueAction(characterID, "build_wall", actionData)
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

// TestBuildingActionDataExtraction tests that building action data is properly extracted
func TestBuildingActionDataExtraction(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	characterID := "test-character"

	testCases := []struct {
		name         string
		actionData   map[string]interface{}
		expectedID   string
		expectedTool string
		expectedDX   int
		expectedDY   int
	}{
		{
			name: "complete_data",
			actionData: map[string]interface{}{
				"itemID": "stone-block",
				"toolID": "hammer",
				"deltaX": 2,
				"deltaY": -1,
			},
			expectedID:   "stone-block",
			expectedTool: "hammer",
			expectedDX:   2,
			expectedDY:   -1,
		},
		{
			name: "no_tool",
			actionData: map[string]interface{}{
				"itemID": "wood-plank",
				"deltaX": -2,
				"deltaY": 3,
			},
			expectedID:   "wood-plank",
			expectedTool: "",
			expectedDX:   -2,
			expectedDY:   3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear previous requests
			mockGameManager.buildRequests = nil

			// Queue and process action
			err := actionManager.QueueAction(characterID, "build_wall", tc.actionData)
			if err != nil {
				t.Fatalf("Failed to queue building action: %v", err)
			}

			queue := actionManager.GetCharacterQueue(characterID)
			queue.StartNextAction(100)

			// Process to completion
			for tick := uint64(101); tick <= 115; tick++ {
				actionManager.OnTick(tick)
			}

			// Verify data extraction
			if len(mockGameManager.buildRequests) == 0 {
				t.Fatal("Expected building request to be made")
			}

			request := mockGameManager.buildRequests[0]
			if request.ItemID != tc.expectedID {
				t.Errorf("Expected item ID %s, got %s", tc.expectedID, request.ItemID)
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
