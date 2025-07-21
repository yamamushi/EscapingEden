package game

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// MockGameManager is a mock implementation for testing
type MockGameManager struct {
	moveRequests  []MoveRequest
	buildRequests []BuildRequest
	mineRequests  []MineRequest
	digRequests   []DigRequest
	mutex         sync.Mutex
	shouldFail    bool
}

type MoveRequest struct {
	CharacterID string
	DeltaX      int
	DeltaY      int
}

type BuildRequest struct {
	ItemID      string
	ToolID      string
	CharacterID string
	DeltaX      int
	DeltaY      int
}

type MineRequest struct {
	ItemID      string
	ToolID      string
	CharacterID string
	DeltaX      int
	DeltaY      int
}

type DigRequest struct {
	ItemID      string
	CharacterID string
	DeltaX      int
	DeltaY      int
}

func (mgm *MockGameManager) HandleMovePlayerRequest(characterID string, deltaX, deltaY int) error {
	mgm.mutex.Lock()
	defer mgm.mutex.Unlock()

	if mgm.shouldFail {
		return fmt.Errorf("mock move failure")
	}

	mgm.moveRequests = append(mgm.moveRequests, MoveRequest{
		CharacterID: characterID,
		DeltaX:      deltaX,
		DeltaY:      deltaY,
	})
	return nil
}

func (mgm *MockGameManager) HandleBuildWallRequest(itemID, toolID, characterID string, deltaX, deltaY int) error {
	mgm.mutex.Lock()
	defer mgm.mutex.Unlock()

	if mgm.shouldFail {
		return fmt.Errorf("mock build failure")
	}

	mgm.buildRequests = append(mgm.buildRequests, BuildRequest{
		ItemID:      itemID,
		ToolID:      toolID,
		CharacterID: characterID,
		DeltaX:      deltaX,
		DeltaY:      deltaY,
	})
	return nil
}

func (mgm *MockGameManager) HandleMineRequest(itemID, toolID, characterID string, deltaX, deltaY int) error {
	mgm.mutex.Lock()
	defer mgm.mutex.Unlock()

	if mgm.shouldFail {
		return fmt.Errorf("mock mine failure")
	}

	mgm.mineRequests = append(mgm.mineRequests, MineRequest{
		ItemID:      itemID,
		ToolID:      toolID,
		CharacterID: characterID,
		DeltaX:      deltaX,
		DeltaY:      deltaY,
	})
	return nil
}

func (mgm *MockGameManager) HandleDigRequest(itemID, characterID string, deltaX, deltaY int) error {
	mgm.mutex.Lock()
	defer mgm.mutex.Unlock()

	if mgm.shouldFail {
		return fmt.Errorf("mock dig failure")
	}

	mgm.digRequests = append(mgm.digRequests, DigRequest{
		ItemID:      itemID,
		CharacterID: characterID,
		DeltaX:      deltaX,
		DeltaY:      deltaY,
	})
	return nil
}

func (mgm *MockGameManager) GetMoveRequests() []MoveRequest {
	mgm.mutex.Lock()
	defer mgm.mutex.Unlock()
	result := make([]MoveRequest, len(mgm.moveRequests))
	copy(result, mgm.moveRequests)
	return result
}

func (mgm *MockGameManager) GetBuildRequests() []BuildRequest {
	mgm.mutex.Lock()
	defer mgm.mutex.Unlock()
	result := make([]BuildRequest, len(mgm.buildRequests))
	copy(result, mgm.buildRequests)
	return result
}

func (mgm *MockGameManager) SetShouldFail(shouldFail bool) {
	mgm.mutex.Lock()
	defer mgm.mutex.Unlock()
	mgm.shouldFail = shouldFail
}

// createTestActionManager creates an ActionManager with test dependencies
func createTestActionManager() (*ActionManager, *ActionRegistry, *MockGameManager) {
	mockLogger := &MockLogger{}

	// Create registry with test configuration
	registry := NewActionRegistry(mockLogger)
	registry.TickRate = 200
	registry.MaxQueueSize = 3
	registry.Actions = map[string]*ActionDefinition{
		"move": {
			TickCost:      5,
			Description:   "Moving one tile",
			CanQueue:      true,
			Interruptible: true,
			Category:      "movement",
		},
		"build_wall": {
			TickCost:      15,
			Description:   "Building a wall",
			CanQueue:      true,
			Interruptible: false,
			Category:      "construction",
		},
		"mine": {
			TickCost:      20,
			Description:   "Mining a tile",
			CanQueue:      true,
			Interruptible: false,
			Category:      "resource",
		},
		"dig": {
			TickCost:      10,
			Description:   "Digging a tile",
			CanQueue:      true,
			Interruptible: true,
			Category:      "terrain",
		},
	}

	mockGameManager := &MockGameManager{
		moveRequests:  make([]MoveRequest, 0),
		buildRequests: make([]BuildRequest, 0),
		mineRequests:  make([]MineRequest, 0),
		digRequests:   make([]DigRequest, 0),
	}

	actionManager := NewActionManager(registry, mockGameManager, mockLogger)
	return actionManager, registry, mockGameManager
}

// TestActionManager_QueueAction tests basic action queuing
func TestActionManager_QueueAction(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Test queuing a valid action
	err := actionManager.QueueAction(characterID, "move", map[string]interface{}{
		"direction": "north",
	})
	if err != nil {
		t.Fatalf("Failed to queue valid action: %v", err)
	}

	// Verify action was queued
	queue := actionManager.GetCharacterQueue(characterID)
	if queue == nil {
		t.Fatal("Character queue was not created")
	}

	status := queue.GetQueueStatus()
	if status.QueuedCount != 1 {
		t.Errorf("Expected 1 queued action, got %d", status.QueuedCount)
	}

	// Test queuing invalid action type
	err = actionManager.QueueAction(characterID, "invalid_action", nil)
	if err == nil {
		t.Error("Expected error when queuing invalid action type")
	}

	// Test queuing when disabled
	actionManager.SetEnabled(false)
	err = actionManager.QueueAction(characterID, "move", map[string]interface{}{
		"direction": "south",
	})
	if err == nil {
		t.Error("Expected error when queuing action while disabled")
	}
}

// TestActionManager_QueueLimit tests queue size limits
func TestActionManager_QueueLimit(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Fill queue to capacity (3 actions)
	for i := 0; i < 3; i++ {
		err := actionManager.QueueAction(characterID, "move", map[string]interface{}{
			"direction": "north",
		})
		if err != nil {
			t.Fatalf("Failed to queue action %d: %v", i, err)
		}
	}

	// Try to exceed limit
	err := actionManager.QueueAction(characterID, "move", map[string]interface{}{
		"direction": "south",
	})
	if err == nil {
		t.Error("Expected error when exceeding queue limit")
	}

	// Verify queue size is still at limit
	queue := actionManager.GetCharacterQueue(characterID)
	status := queue.GetQueueStatus()
	if status.QueuedCount != 3 {
		t.Errorf("Expected queue size 3, got %d", status.QueuedCount)
	}
}

// TestActionManager_TickProcessing tests action processing on tick
func TestActionManager_TickProcessing(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	characterID := "test-character"

	// Queue a quick action (5 ticks)
	err := actionManager.QueueAction(characterID, "move", map[string]interface{}{
		"direction": "north",
	})
	if err != nil {
		t.Fatalf("Failed to queue action: %v", err)
	}

	// Process first tick - action should start
	actionManager.OnTick(1)

	queue := actionManager.GetCharacterQueue(characterID)
	if queue.GetCurrentAction() == nil {
		t.Error("Action should have started processing")
	}

	status := queue.GetQueueStatus()
	if status.QueuedCount != 0 {
		t.Error("Action should have been moved from queue to current")
	}

	// Process ticks 2-5 - action should still be running
	for tick := uint64(2); tick <= 5; tick++ {
		actionManager.OnTick(tick)
		if queue.GetCurrentAction() == nil {
			t.Errorf("Action should still be processing at tick %d", tick)
		}
	}

	// Process tick 6 - action should complete
	actionManager.OnTick(6)

	if queue.GetCurrentAction() != nil {
		t.Error("Action should have completed")
	}

	// Verify the action was executed
	moveRequests := mockGameManager.GetMoveRequests()
	if len(moveRequests) != 1 {
		t.Errorf("Expected 1 move request, got %d", len(moveRequests))
	}

	if moveRequests[0].CharacterID != characterID {
		t.Errorf("Expected character ID %s, got %s", characterID, moveRequests[0].CharacterID)
	}

	if moveRequests[0].DeltaX != 0 || moveRequests[0].DeltaY != -1 {
		t.Errorf("Expected north movement (0, -1), got (%d, %d)", moveRequests[0].DeltaX, moveRequests[0].DeltaY)
	}
}

// TestActionManager_ActionInterruption tests action interruption
func TestActionManager_ActionInterruption(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Queue and start an interruptible action
	actionManager.QueueAction(characterID, "move", map[string]interface{}{
		"direction": "north",
	})
	actionManager.OnTick(1) // Start the action

	// Interrupt the action
	err := actionManager.InterruptCurrentAction(characterID)
	if err != nil {
		t.Errorf("Failed to interrupt interruptible action: %v", err)
	}

	queue := actionManager.GetCharacterQueue(characterID)
	if queue.GetCurrentAction() != nil {
		t.Error("Interruptible action should have been cancelled")
	}

	// Test trying to interrupt a non-interruptible action
	actionManager.QueueAction(characterID, "build_wall", map[string]interface{}{
		"itemID": "test-item",
		"deltaX": 1,
		"deltaY": 0,
	})
	actionManager.OnTick(2) // Start the action

	err = actionManager.InterruptCurrentAction(characterID)
	if err == nil {
		t.Error("Expected error when trying to interrupt non-interruptible action")
	}

	if queue.GetCurrentAction() == nil {
		t.Error("Non-interruptible action should still be running")
	}

	// Test interrupting when no current action
	queue.Clear()
	err = actionManager.InterruptCurrentAction(characterID)
	if err == nil {
		t.Error("Expected error when trying to interrupt with no current action")
	}
}

// TestActionManager_CharacterLogout tests queue cleanup on logout
func TestActionManager_CharacterLogout(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Queue some actions and start one
	actionManager.QueueAction(characterID, "move", map[string]interface{}{
		"direction": "north",
	})
	actionManager.QueueAction(characterID, "move", map[string]interface{}{
		"direction": "south",
	})
	actionManager.OnTick(1) // Start first action

	// Verify queue exists and has actions
	queue := actionManager.GetCharacterQueue(characterID)
	if queue == nil {
		t.Fatal("Queue should exist before logout")
	}

	status := queue.GetQueueStatus()
	if !status.HasCurrentAction || status.QueuedCount == 0 {
		t.Error("Queue should have current action and queued actions")
	}

	// Simulate logout
	actionManager.OnCharacterLogout(characterID)

	// Verify queue is cleaned up
	queue = actionManager.GetCharacterQueue(characterID)
	if queue != nil {
		t.Error("Queue should be removed after logout")
	}
}

// TestActionManager_ActionExecution tests different action types
func TestActionManager_ActionExecution(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	characterID := "test-character"

	testCases := []struct {
		actionType string
		data       map[string]interface{}
		duration   int
	}{
		{
			actionType: "move",
			data: map[string]interface{}{
				"direction": "east",
			},
			duration: 5,
		},
		{
			actionType: "build_wall",
			data: map[string]interface{}{
				"itemID": "wall-item",
				"toolID": "hammer",
				"deltaX": 1,
				"deltaY": 0,
			},
			duration: 15,
		},
		{
			actionType: "mine",
			data: map[string]interface{}{
				"itemID": "ore-item",
				"toolID": "pickaxe",
				"deltaX": 0,
				"deltaY": 1,
			},
			duration: 20,
		},
		{
			actionType: "dig",
			data: map[string]interface{}{
				"itemID": "dirt-item",
				"deltaX": -1,
				"deltaY": 0,
			},
			duration: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.actionType, func(t *testing.T) {
			// Queue the action
			err := actionManager.QueueAction(characterID, tc.actionType, tc.data)
			if err != nil {
				t.Fatalf("Failed to queue %s action: %v", tc.actionType, err)
			}

			// Start the action
			actionManager.OnTick(1)

			// Complete the action
			actionManager.OnTick(uint64(tc.duration + 1))

			// Verify the action was executed based on type
			switch tc.actionType {
			case "move":
				moveRequests := mockGameManager.GetMoveRequests()
				if len(moveRequests) == 0 {
					t.Error("Expected move request to be executed")
				}
			case "build_wall":
				buildRequests := mockGameManager.GetBuildRequests()
				if len(buildRequests) == 0 {
					t.Error("Expected build request to be executed")
				}
			}

			// Clear for next test
			actionManager.OnCharacterLogout(characterID)
		})
	}
}

// TestActionManager_ActionFailure tests handling of action execution failures
func TestActionManager_ActionFailure(t *testing.T) {
	actionManager, _, mockGameManager := createTestActionManager()

	characterID := "test-character"

	// Set mock to fail
	mockGameManager.SetShouldFail(true)

	// Queue an action
	err := actionManager.QueueAction(characterID, "move", map[string]interface{}{
		"direction": "north",
	})
	if err != nil {
		t.Fatalf("Failed to queue action: %v", err)
	}

	// Process action to completion
	actionManager.OnTick(1) // Start
	actionManager.OnTick(6) // Complete

	// Verify metrics show failure
	metrics := actionManager.GetMetrics()
	if metrics.ActionsFailed == 0 {
		t.Error("Expected failed action to be recorded in metrics")
	}
}

// TestActionManager_ConcurrentAccess tests thread safety
func TestActionManager_ConcurrentAccess(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	var wg sync.WaitGroup
	numGoroutines := 20
	actionsPerGoroutine := 5

	// Concurrently queue actions for multiple characters
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			characterID := fmt.Sprintf("character-%d", goroutineID)

			for j := 0; j < actionsPerGoroutine; j++ {
				err := actionManager.QueueAction(characterID, "move", map[string]interface{}{
					"direction": "north",
				})
				if err != nil {
					// Some may fail due to queue limits, which is expected
					continue
				}

				// Read operations
				_ = actionManager.GetCharacterQueue(characterID)
				_ = actionManager.GetAllQueueStatuses()
				_ = actionManager.GetMetrics()
			}
		}(i)
	}

	// Concurrently process ticks
	wg.Add(1)
	go func() {
		defer wg.Done()
		for tick := uint64(1); tick <= 50; tick++ {
			actionManager.OnTick(tick)
			time.Sleep(1 * time.Millisecond) // Small delay
		}
	}()

	wg.Wait()

	// Verify system integrity
	metrics := actionManager.GetMetrics()
	if metrics.CurrentTick != 50 {
		t.Errorf("Expected current tick 50, got %d", metrics.CurrentTick)
	}
}

// TestActionManager_Metrics tests metrics collection
func TestActionManager_Metrics(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	characterID := "test-character"

	// Initial metrics
	metrics := actionManager.GetMetrics()
	if metrics.ActionsCompleted != 0 || metrics.ActionsFailed != 0 {
		t.Error("Expected zero metrics initially")
	}

	// Queue and complete an action
	actionManager.QueueAction(characterID, "move", map[string]interface{}{
		"direction": "north",
	})

	actionManager.OnTick(1) // Start
	actionManager.OnTick(6) // Complete

	// Check updated metrics
	metrics = actionManager.GetMetrics()
	if metrics.ActionsCompleted != 1 {
		t.Errorf("Expected 1 completed action, got %d", metrics.ActionsCompleted)
	}

	if metrics.ActiveQueues == 0 {
		t.Error("Expected at least 1 active queue")
	}
}

// TestActionManager_EnableDisable tests enabling/disabling functionality
func TestActionManager_EnableDisable(t *testing.T) {
	actionManager, _, _ := createTestActionManager()

	// Test initial state
	if !actionManager.IsEnabled() {
		t.Error("ActionManager should be enabled by default")
	}

	// Test disabling
	actionManager.SetEnabled(false)
	if actionManager.IsEnabled() {
		t.Error("ActionManager should be disabled after SetEnabled(false)")
	}

	// Test that disabled manager doesn't process ticks
	characterID := "test-character"
	actionManager.QueueAction(characterID, "move", map[string]interface{}{
		"direction": "north",
	}) // This should fail

	// Enable and test normal operation
	actionManager.SetEnabled(true)
	if !actionManager.IsEnabled() {
		t.Error("ActionManager should be enabled after SetEnabled(true)")
	}

	err := actionManager.QueueAction(characterID, "move", map[string]interface{}{
		"direction": "north",
	})
	if err != nil {
		t.Errorf("Should be able to queue action when enabled: %v", err)
	}
}

// BenchmarkActionManager_QueueAction benchmarks action queuing performance
func BenchmarkActionManager_QueueAction(b *testing.B) {
	actionManager, _, _ := createTestActionManager()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		characterID := fmt.Sprintf("character-%d", i%100) // Cycle through 100 characters
		actionManager.QueueAction(characterID, "move", map[string]interface{}{
			"direction": "north",
		})
	}
}

// BenchmarkActionManager_OnTick benchmarks tick processing performance
func BenchmarkActionManager_OnTick(b *testing.B) {
	actionManager, _, _ := createTestActionManager()

	// Pre-populate with actions
	for i := 0; i < 100; i++ {
		characterID := fmt.Sprintf("character-%d", i)
		for j := 0; j < 2; j++ {
			actionManager.QueueAction(characterID, "move", map[string]interface{}{
				"direction": "north",
			})
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		actionManager.OnTick(uint64(i + 1))
	}
}
