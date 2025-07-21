package game

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestQueuedAction_Progress tests action progress calculation
func TestQueuedAction_Progress(t *testing.T) {
	action := &QueuedAction{
		ID:        "test-action",
		Type:      "move",
		StartTick: 10,
		Duration:  5,
	}

	// Test before start
	progress := action.GetProgress(5)
	if progress != 0.0 {
		t.Errorf("Expected progress 0.0 before start, got %f", progress)
	}

	// Test at start
	progress = action.GetProgress(10)
	if progress != 0.0 {
		t.Errorf("Expected progress 0.0 at start, got %f", progress)
	}

	// Test halfway through
	progress = action.GetProgress(12)
	expected := 2.0 / 5.0 // 2 ticks elapsed out of 5
	if progress != expected {
		t.Errorf("Expected progress %f halfway through, got %f", expected, progress)
	}

	// Test at completion
	progress = action.GetProgress(15)
	if progress != 1.0 {
		t.Errorf("Expected progress 1.0 at completion, got %f", progress)
	}

	// Test after completion
	progress = action.GetProgress(20)
	if progress != 1.0 {
		t.Errorf("Expected progress 1.0 after completion, got %f", progress)
	}
}

// TestQueuedAction_TicksLeft tests remaining tick calculation
func TestQueuedAction_TicksLeft(t *testing.T) {
	action := &QueuedAction{
		ID:        "test-action",
		Type:      "move",
		StartTick: 10,
		Duration:  5,
	}

	// Test before start
	ticksLeft := action.GetTicksLeft(5)
	if ticksLeft != 5 {
		t.Errorf("Expected 5 ticks left before start, got %d", ticksLeft)
	}

	// Test at start
	ticksLeft = action.GetTicksLeft(10)
	if ticksLeft != 5 {
		t.Errorf("Expected 5 ticks left at start, got %d", ticksLeft)
	}

	// Test halfway through
	ticksLeft = action.GetTicksLeft(12)
	if ticksLeft != 3 {
		t.Errorf("Expected 3 ticks left halfway through, got %d", ticksLeft)
	}

	// Test at completion
	ticksLeft = action.GetTicksLeft(15)
	if ticksLeft != 0 {
		t.Errorf("Expected 0 ticks left at completion, got %d", ticksLeft)
	}

	// Test after completion
	ticksLeft = action.GetTicksLeft(20)
	if ticksLeft != 0 {
		t.Errorf("Expected 0 ticks left after completion, got %d", ticksLeft)
	}
}

// TestQueuedAction_IsCompleted tests completion status
func TestQueuedAction_IsCompleted(t *testing.T) {
	action := &QueuedAction{
		ID:        "test-action",
		Type:      "move",
		StartTick: 10,
		Duration:  5,
	}

	// Test before start
	if action.IsCompleted(5) {
		t.Error("Action should not be completed before start")
	}

	// Test during execution
	if action.IsCompleted(12) {
		t.Error("Action should not be completed during execution")
	}

	// Test at completion
	if !action.IsCompleted(15) {
		t.Error("Action should be completed at end tick")
	}

	// Test after completion
	if !action.IsCompleted(20) {
		t.Error("Action should remain completed after end tick")
	}

	// Test action that hasn't started
	unstarted := &QueuedAction{
		ID:        "unstarted",
		Type:      "move",
		StartTick: 0, // Not started
		Duration:  5,
	}

	if unstarted.IsCompleted(100) {
		t.Error("Unstarted action should not be completed")
	}
}

// TestCharacterActionQueue_AddAction tests adding actions to queue
func TestCharacterActionQueue_AddAction(t *testing.T) {
	mockLogger := &MockLogger{}
	queue := NewCharacterActionQueue("test-character", 3, mockLogger)

	action := &QueuedAction{
		ID:          "action-1",
		CharacterID: "test-character",
		Type:        "move",
		Duration:    5,
		CreatedAt:   time.Now(),
	}

	// Test adding valid action
	err := queue.AddAction(action)
	if err != nil {
		t.Fatalf("Failed to add valid action: %v", err)
	}

	status := queue.GetQueueStatus()
	if status.QueuedCount != 1 {
		t.Errorf("Expected 1 queued action, got %d", status.QueuedCount)
	}

	// Test adding nil action
	err = queue.AddAction(nil)
	if err == nil {
		t.Error("Expected error when adding nil action")
	}

	// Test adding action with wrong character ID
	wrongCharAction := &QueuedAction{
		ID:          "action-2",
		CharacterID: "wrong-character",
		Type:        "move",
		Duration:    5,
	}

	err = queue.AddAction(wrongCharAction)
	if err == nil {
		t.Error("Expected error when adding action with wrong character ID")
	}

	// Test adding action with invalid duration
	invalidDurationAction := &QueuedAction{
		ID:          "action-3",
		CharacterID: "test-character",
		Type:        "move",
		Duration:    0,
	}

	err = queue.AddAction(invalidDurationAction)
	if err == nil {
		t.Error("Expected error when adding action with invalid duration")
	}
}

// TestCharacterActionQueue_QueueLimit tests queue size limits
func TestCharacterActionQueue_QueueLimit(t *testing.T) {
	mockLogger := &MockLogger{}
	queue := NewCharacterActionQueue("test-character", 2, mockLogger) // Small limit

	// Fill queue to capacity
	for i := 0; i < 2; i++ {
		action := &QueuedAction{
			ID:          fmt.Sprintf("action-%d", i),
			CharacterID: "test-character",
			Type:        "move",
			Duration:    5,
			CreatedAt:   time.Now(),
		}

		err := queue.AddAction(action)
		if err != nil {
			t.Fatalf("Failed to add action %d: %v", i, err)
		}
	}

	// Verify queue is full
	if queue.CanQueue() {
		t.Error("Queue should be full")
	}

	// Try to exceed limit
	overflowAction := &QueuedAction{
		ID:          "overflow-action",
		CharacterID: "test-character",
		Type:        "move",
		Duration:    5,
		CreatedAt:   time.Now(),
	}

	err := queue.AddAction(overflowAction)
	if err == nil {
		t.Error("Expected error when exceeding queue limit")
	}

	// Verify queue size is still at limit
	status := queue.GetQueueStatus()
	if status.QueuedCount != 2 {
		t.Errorf("Expected queue size 2, got %d", status.QueuedCount)
	}
}

// TestCharacterActionQueue_StartNextAction tests starting queued actions
func TestCharacterActionQueue_StartNextAction(t *testing.T) {
	mockLogger := &MockLogger{}
	queue := NewCharacterActionQueue("test-character", 3, mockLogger)

	// Add actions to queue
	for i := 0; i < 2; i++ {
		action := &QueuedAction{
			ID:          fmt.Sprintf("action-%d", i),
			CharacterID: "test-character",
			Type:        "move",
			Duration:    5,
			CreatedAt:   time.Now(),
		}
		queue.AddAction(action)
	}

	// Start first action
	currentTick := uint64(100)
	startedAction := queue.StartNextAction(currentTick)

	if startedAction == nil {
		t.Fatal("Expected action to be started")
	}

	if startedAction.ID != "action-0" {
		t.Errorf("Expected first action to be started, got %s", startedAction.ID)
	}

	if startedAction.StartTick != currentTick {
		t.Errorf("Expected start tick %d, got %d", currentTick, startedAction.StartTick)
	}

	// Verify queue size decreased
	status := queue.GetQueueStatus()
	if status.QueuedCount != 1 {
		t.Errorf("Expected 1 queued action after start, got %d", status.QueuedCount)
	}

	if !status.HasCurrentAction {
		t.Error("Expected current action to be set")
	}

	// Try to start another action while one is current
	secondStart := queue.StartNextAction(currentTick + 1)
	if secondStart != nil {
		t.Error("Should not be able to start action while one is current")
	}

	// Try to start from empty queue
	emptyQueue := NewCharacterActionQueue("empty-character", 3, mockLogger)
	emptyStart := emptyQueue.StartNextAction(currentTick)
	if emptyStart != nil {
		t.Error("Should not be able to start action from empty queue")
	}
}

// TestCharacterActionQueue_CompleteCurrentAction tests completing actions
func TestCharacterActionQueue_CompleteCurrentAction(t *testing.T) {
	mockLogger := &MockLogger{}
	queue := NewCharacterActionQueue("test-character", 3, mockLogger)

	// Add and start an action
	action := &QueuedAction{
		ID:          "test-action",
		CharacterID: "test-character",
		Type:        "move",
		Duration:    5,
		CreatedAt:   time.Now(),
	}
	queue.AddAction(action)

	startTick := uint64(100)
	queue.StartNextAction(startTick)

	// Complete the action
	completeTick := startTick + 5
	completedAction := queue.CompleteCurrentAction(completeTick)

	if completedAction == nil {
		t.Fatal("Expected action to be completed")
	}

	if completedAction.ID != "test-action" {
		t.Errorf("Expected completed action ID 'test-action', got %s", completedAction.ID)
	}

	// Verify no current action
	status := queue.GetQueueStatus()
	if status.HasCurrentAction {
		t.Error("Expected no current action after completion")
	}

	// Try to complete when no current action
	secondComplete := queue.CompleteCurrentAction(completeTick + 1)
	if secondComplete != nil {
		t.Error("Should not be able to complete when no current action")
	}
}

// TestCharacterActionQueue_InterruptCurrentAction tests action interruption
func TestCharacterActionQueue_InterruptCurrentAction(t *testing.T) {
	mockLogger := &MockLogger{}
	queue := NewCharacterActionQueue("test-character", 3, mockLogger)

	// Add and start an action
	action := &QueuedAction{
		ID:          "test-action",
		CharacterID: "test-character",
		Type:        "move",
		Duration:    5,
		CreatedAt:   time.Now(),
	}
	queue.AddAction(action)
	queue.StartNextAction(100)

	// Interrupt the action
	err := queue.InterruptCurrentAction()
	if err != nil {
		t.Fatalf("Failed to interrupt action: %v", err)
	}

	// Verify no current action
	status := queue.GetQueueStatus()
	if status.HasCurrentAction {
		t.Error("Expected no current action after interruption")
	}

	// Try to interrupt when no current action
	err = queue.InterruptCurrentAction()
	if err == nil {
		t.Error("Expected error when interrupting with no current action")
	}
}

// TestCharacterActionQueue_Clear tests clearing the queue
func TestCharacterActionQueue_Clear(t *testing.T) {
	mockLogger := &MockLogger{}
	queue := NewCharacterActionQueue("test-character", 3, mockLogger)

	// Add actions and start one
	for i := 0; i < 2; i++ {
		action := &QueuedAction{
			ID:          fmt.Sprintf("action-%d", i),
			CharacterID: "test-character",
			Type:        "move",
			Duration:    5,
			CreatedAt:   time.Now(),
		}
		queue.AddAction(action)
	}

	queue.StartNextAction(100)

	// Verify queue has actions
	status := queue.GetQueueStatus()
	if status.QueuedCount == 0 || !status.HasCurrentAction {
		t.Fatal("Expected queue to have actions before clear")
	}

	// Clear the queue
	queue.Clear()

	// Verify queue is empty
	status = queue.GetQueueStatus()
	if status.QueuedCount != 0 {
		t.Errorf("Expected 0 queued actions after clear, got %d", status.QueuedCount)
	}

	if status.HasCurrentAction {
		t.Error("Expected no current action after clear")
	}

	if !queue.IsEmpty() {
		t.Error("Queue should be empty after clear")
	}
}

// TestCharacterActionQueue_GetCurrentAction tests getting current action copy
func TestCharacterActionQueue_GetCurrentAction(t *testing.T) {
	mockLogger := &MockLogger{}
	queue := NewCharacterActionQueue("test-character", 3, mockLogger)

	// Test with no current action
	current := queue.GetCurrentAction()
	if current != nil {
		t.Error("Expected nil when no current action")
	}

	// Add and start an action
	originalAction := &QueuedAction{
		ID:          "test-action",
		CharacterID: "test-character",
		Type:        "move",
		Duration:    5,
		CreatedAt:   time.Now(),
	}
	queue.AddAction(originalAction)
	queue.StartNextAction(100)

	// Get current action copy
	current = queue.GetCurrentAction()
	if current == nil {
		t.Fatal("Expected current action copy")
	}

	if current.ID != "test-action" {
		t.Errorf("Expected action ID 'test-action', got %s", current.ID)
	}

	// Verify it's a copy (modifying shouldn't affect original)
	current.Duration = 999
	realCurrent := queue.GetCurrentAction()
	if realCurrent.Duration != 5 {
		t.Error("Current action should be returned as a copy")
	}
}

// TestCharacterActionQueue_GetQueuedActions tests getting queued actions copies
func TestCharacterActionQueue_GetQueuedActions(t *testing.T) {
	mockLogger := &MockLogger{}
	queue := NewCharacterActionQueue("test-character", 3, mockLogger)

	// Test with empty queue
	queued := queue.GetQueuedActions()
	if queued != nil {
		t.Error("Expected nil for empty queue")
	}

	// Add actions
	for i := 0; i < 2; i++ {
		action := &QueuedAction{
			ID:          fmt.Sprintf("action-%d", i),
			CharacterID: "test-character",
			Type:        "move",
			Duration:    5 + i,
			CreatedAt:   time.Now(),
		}
		queue.AddAction(action)
	}

	// Get queued actions
	queued = queue.GetQueuedActions()
	if len(queued) != 2 {
		t.Errorf("Expected 2 queued actions, got %d", len(queued))
	}

	// Verify order and content
	if queued[0].ID != "action-0" {
		t.Errorf("Expected first action ID 'action-0', got %s", queued[0].ID)
	}

	if queued[1].Duration != 6 {
		t.Errorf("Expected second action duration 6, got %d", queued[1].Duration)
	}

	// Verify they're copies (modifying shouldn't affect originals)
	queued[0].Duration = 999
	newQueued := queue.GetQueuedActions()
	if newQueued[0].Duration != 5 {
		t.Error("Queued actions should be returned as copies")
	}
}

// TestCharacterActionQueue_ConcurrentAccess tests thread safety
func TestCharacterActionQueue_ConcurrentAccess(t *testing.T) {
	mockLogger := &MockLogger{}
	queue := NewCharacterActionQueue("test-character", 10, mockLogger)

	var wg sync.WaitGroup
	numGoroutines := 20

	// Concurrently add actions
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			action := &QueuedAction{
				ID:          fmt.Sprintf("concurrent-action-%d", id),
				CharacterID: "test-character",
				Type:        "move",
				Duration:    5,
				CreatedAt:   time.Now(),
			}

			// Try to add action (some may fail due to queue limit)
			queue.AddAction(action)

			// Read operations
			_ = queue.CanQueue()
			_ = queue.GetQueueStatus()
			_ = queue.GetCurrentAction()
			_ = queue.GetQueuedActions()
			_ = queue.IsEmpty()
		}(i)
	}

	wg.Wait()

	// Verify queue integrity
	status := queue.GetQueueStatus()
	if status.QueuedCount > 10 {
		t.Errorf("Queue exceeded maximum size: %d", status.QueuedCount)
	}
}

// TestActionQueueManager_GetOrCreateQueue tests queue management
func TestActionQueueManager_GetOrCreateQueue(t *testing.T) {
	mockLogger := &MockLogger{}
	manager := NewActionQueueManager(mockLogger)

	characterID := "test-character"
	maxQueueSize := 5

	// Test creating new queue
	queue1 := manager.GetOrCreateQueue(characterID, maxQueueSize)
	if queue1 == nil {
		t.Fatal("Expected queue to be created")
	}

	if queue1.CharacterID != characterID {
		t.Errorf("Expected character ID %s, got %s", characterID, queue1.CharacterID)
	}

	if queue1.MaxQueueSize != maxQueueSize {
		t.Errorf("Expected max queue size %d, got %d", maxQueueSize, queue1.MaxQueueSize)
	}

	// Test getting existing queue
	queue2 := manager.GetOrCreateQueue(characterID, maxQueueSize)
	if queue2 != queue1 {
		t.Error("Expected to get the same queue instance")
	}

	// Verify manager has the queue
	retrievedQueue := manager.GetQueue(characterID)
	if retrievedQueue != queue1 {
		t.Error("Manager should return the same queue instance")
	}
}

// TestActionQueueManager_RemoveQueue tests queue removal
func TestActionQueueManager_RemoveQueue(t *testing.T) {
	mockLogger := &MockLogger{}
	manager := NewActionQueueManager(mockLogger)

	characterID := "test-character"

	// Create queue and add action
	queue := manager.GetOrCreateQueue(characterID, 3)
	action := &QueuedAction{
		ID:          "test-action",
		CharacterID: characterID,
		Type:        "move",
		Duration:    5,
		CreatedAt:   time.Now(),
	}
	queue.AddAction(action)

	// Verify queue exists and has action
	if manager.GetQueue(characterID) == nil {
		t.Fatal("Queue should exist before removal")
	}

	status := queue.GetQueueStatus()
	if status.QueuedCount != 1 {
		t.Fatal("Queue should have action before removal")
	}

	// Remove queue
	manager.RemoveQueue(characterID)

	// Verify queue is removed
	if manager.GetQueue(characterID) != nil {
		t.Error("Queue should be removed")
	}

	// Test removing non-existent queue (should not panic)
	manager.RemoveQueue("non-existent-character")
}

// TestActionQueueManager_CleanupInactiveQueues tests cleanup functionality
func TestActionQueueManager_CleanupInactiveQueues(t *testing.T) {
	mockLogger := &MockLogger{}
	manager := NewActionQueueManager(mockLogger)

	// Create queues with different activity times
	activeCharacter := "active-character"
	inactiveCharacter := "inactive-character"

	activeQueue := manager.GetOrCreateQueue(activeCharacter, 3)
	inactiveQueue := manager.GetOrCreateQueue(inactiveCharacter, 3)

	// Add action to active queue (recent activity)
	activeAction := &QueuedAction{
		ID:          "active-action",
		CharacterID: activeCharacter,
		Type:        "move",
		Duration:    5,
		CreatedAt:   time.Now(),
	}
	activeQueue.AddAction(activeAction)

	// Make inactive queue old by manipulating lastActivity
	// (In real code, this would happen naturally over time)
	inactiveQueue.lastActivity = time.Now().Add(-2 * time.Hour)

	// Cleanup with 1 hour threshold
	removedCount := manager.CleanupInactiveQueues(1 * time.Hour)

	// Verify inactive queue was removed
	if removedCount != 1 {
		t.Errorf("Expected 1 queue to be removed, got %d", removedCount)
	}

	if manager.GetQueue(inactiveCharacter) != nil {
		t.Error("Inactive queue should have been removed")
	}

	if manager.GetQueue(activeCharacter) == nil {
		t.Error("Active queue should not have been removed")
	}
}

// BenchmarkCharacterActionQueue_AddAction benchmarks action queuing
func BenchmarkCharacterActionQueue_AddAction(b *testing.B) {
	mockLogger := &MockLogger{}
	queue := NewCharacterActionQueue("test-character", 1000, mockLogger)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		action := &QueuedAction{
			ID:          fmt.Sprintf("action-%d", i),
			CharacterID: "test-character",
			Type:        "move",
			Duration:    5,
			CreatedAt:   time.Now(),
		}

		// Clear queue periodically to avoid hitting limit
		if i%500 == 0 {
			queue.Clear()
		}

		queue.AddAction(action)
	}
}

// BenchmarkCharacterActionQueue_GetQueueStatus benchmarks status retrieval
func BenchmarkCharacterActionQueue_GetQueueStatus(b *testing.B) {
	mockLogger := &MockLogger{}
	queue := NewCharacterActionQueue("test-character", 10, mockLogger)

	// Pre-populate queue
	for i := 0; i < 5; i++ {
		action := &QueuedAction{
			ID:          fmt.Sprintf("action-%d", i),
			CharacterID: "test-character",
			Type:        "move",
			Duration:    5,
			CreatedAt:   time.Now(),
		}
		queue.AddAction(action)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = queue.GetQueueStatus()
	}
}
