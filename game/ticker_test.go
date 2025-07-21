package game

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

// TestGameTicker_StartStop tests the basic start and stop functionality
func TestGameTicker_StartStop(t *testing.T) {
	// Create a mock logger for testing
	mockLogger := &MockLogger{}
	ticker := NewGameTicker(200*time.Millisecond, mockLogger)

	// Test starting the ticker
	err := ticker.Start()
	if err != nil {
		t.Fatalf("Failed to start ticker: %v", err)
	}

	if !ticker.IsRunning() {
		t.Error("Ticker should be running after Start()")
	}

	// Test stopping the ticker
	err = ticker.Stop()
	if err != nil {
		t.Fatalf("Failed to stop ticker: %v", err)
	}

	if ticker.IsRunning() {
		t.Error("Ticker should not be running after Stop()")
	}
}

// TestGameTicker_TickAccuracy tests the timing accuracy of ticks
func TestGameTicker_TickAccuracy(t *testing.T) {
	tickRate := 100 * time.Millisecond
	mockLogger := &MockLogger{}
	ticker := NewGameTicker(tickRate, mockLogger)

	// Create a test subscriber to measure tick timing
	testSubscriber := &TestTickSubscriber{
		id:        "test-subscriber",
		tickTimes: make([]time.Time, 0),
	}

	ticker.Subscribe(testSubscriber)
	ticker.Start()

	// Let it run for a short period
	time.Sleep(500 * time.Millisecond)
	ticker.Stop()

	// Analyze tick timing
	if len(testSubscriber.tickTimes) < 4 {
		t.Fatalf("Expected at least 4 ticks, got %d", len(testSubscriber.tickTimes))
	}

	// Check intervals between ticks
	for i := 1; i < len(testSubscriber.tickTimes); i++ {
		interval := testSubscriber.tickTimes[i].Sub(testSubscriber.tickTimes[i-1])

		// Allow 20% variance in timing
		minExpected := time.Duration(float64(tickRate) * 0.8)
		maxExpected := time.Duration(float64(tickRate) * 1.2)

		if interval < minExpected || interval > maxExpected {
			t.Errorf("Tick interval %v is outside expected range [%v, %v]", interval, minExpected, maxExpected)
		}
	}
}

// TestGameTicker_SubscriberManagement tests adding and removing subscribers
func TestGameTicker_SubscriberManagement(t *testing.T) {
	mockLogger := &MockLogger{}
	ticker := NewGameTicker(200*time.Millisecond, mockLogger)

	subscriber1 := &TestTickSubscriber{id: "subscriber-1"}
	subscriber2 := &TestTickSubscriber{id: "subscriber-2"}

	// Test adding subscribers
	err := ticker.Subscribe(subscriber1)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	err = ticker.Subscribe(subscriber2)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	if ticker.GetSubscriberCount() != 2 {
		t.Errorf("Expected 2 subscribers, got %d", ticker.GetSubscriberCount())
	}

	// Test removing subscriber
	err = ticker.Unsubscribe("subscriber-1")
	if err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	if ticker.GetSubscriberCount() != 1 {
		t.Errorf("Expected 1 subscriber after unsubscribe, got %d", ticker.GetSubscriberCount())
	}

	// Test removing non-existent subscriber
	err = ticker.Unsubscribe("non-existent")
	if err == nil {
		t.Error("Expected error when unsubscribing non-existent subscriber")
	}
}

// TestGameTicker_ConcurrentAccess tests thread safety
func TestGameTicker_ConcurrentAccess(t *testing.T) {
	mockLogger := &MockLogger{}
	ticker := NewGameTicker(50*time.Millisecond, mockLogger)

	ticker.Start()
	defer ticker.Stop()

	var wg sync.WaitGroup
	numGoroutines := 10

	// Concurrently add and remove subscribers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			subscriber := &TestTickSubscriber{
				id: fmt.Sprintf("concurrent-subscriber-%d", id),
			}

			// Subscribe
			ticker.Subscribe(subscriber)

			// Wait a bit
			time.Sleep(100 * time.Millisecond)

			// Unsubscribe
			ticker.Unsubscribe(subscriber.GetSubscriberID())
		}(i)
	}

	wg.Wait()

	// All subscribers should be removed
	if ticker.GetSubscriberCount() != 0 {
		t.Errorf("Expected 0 subscribers after concurrent operations, got %d", ticker.GetSubscriberCount())
	}
}

// TestGameTicker_ErrorConditions tests error handling
func TestGameTicker_ErrorConditions(t *testing.T) {
	mockLogger := &MockLogger{}
	ticker := NewGameTicker(200*time.Millisecond, mockLogger)

	// Test starting already running ticker
	ticker.Start()
	err := ticker.Start()
	if err == nil {
		t.Error("Expected error when starting already running ticker")
	}

	// Test stopping already stopped ticker
	ticker.Stop()
	err = ticker.Stop()
	if err == nil {
		t.Error("Expected error when stopping already stopped ticker")
	}

	// Test subscribing nil subscriber
	err = ticker.Subscribe(nil)
	if err == nil {
		t.Error("Expected error when subscribing nil subscriber")
	}
}

// TestTickSubscriber is a test implementation of TickSubscriber
type TestTickSubscriber struct {
	id        string
	tickTimes []time.Time
	mutex     sync.Mutex
}

func (ts *TestTickSubscriber) OnTick(tickNumber uint64) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	ts.tickTimes = append(ts.tickTimes, time.Now())
}

func (ts *TestTickSubscriber) GetSubscriberID() string {
	return ts.id
}

// Benchmark tests for performance
func BenchmarkGameTicker_TickProcessing(b *testing.B) {
	ticker := &GameTicker{
		tickRate:    1 * time.Millisecond, // Very fast for benchmarking
		subscribers: make([]TickSubscriber, 0),
		stopChan:    make(chan struct{}),
	}

	// Add multiple subscribers
	for i := 0; i < 100; i++ {
		subscriber := &TestTickSubscriber{
			id: fmt.Sprintf("bench-subscriber-%d", i),
		}
		ticker.Subscribe(subscriber)
	}

	ticker.Start()
	defer ticker.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Simulate tick processing work
		time.Sleep(1 * time.Microsecond)
	}
}

// TestGameTicker_MemoryLeaks tests for memory leaks in long-running scenarios
func TestGameTicker_MemoryLeaks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory leak test in short mode")
	}

	mockLogger := &MockLogger{}
	ticker := NewGameTicker(10*time.Millisecond, mockLogger)

	ticker.Start()
	defer ticker.Stop()

	// Add and remove subscribers repeatedly
	for i := 0; i < 1000; i++ {
		subscriber := &TestTickSubscriber{
			id: fmt.Sprintf("memory-test-subscriber-%d", i),
		}

		ticker.Subscribe(subscriber)

		if i%10 == 0 {
			// Remove some subscribers periodically
			ticker.Unsubscribe(subscriber.GetSubscriberID())
		}
	}

	// Let it run for a while
	time.Sleep(100 * time.Millisecond)

	// Check that we don't have excessive subscribers
	if ticker.GetSubscriberCount() > 900 {
		t.Errorf("Potential memory leak: too many subscribers remaining: %d", ticker.GetSubscriberCount())
	}
}

// MockLogger is a simple mock implementation of logging.LoggerType for testing
type MockLogger struct {
	messages []string
	mutex    sync.Mutex
}

func (ml *MockLogger) GetTypeID() logging.LoggerTypeID {
	return logging.LoggerTypeID(1) // Mock type ID
}

func (ml *MockLogger) Println(level logging.LogLevel, format string, v ...interface{}) {
	ml.mutex.Lock()
	defer ml.mutex.Unlock()
	message := fmt.Sprintf(format, v...)
	ml.messages = append(ml.messages, fmt.Sprintf("[%v] %s", level, message))
}

func (ml *MockLogger) GetMessages() []string {
	ml.mutex.Lock()
	defer ml.mutex.Unlock()
	result := make([]string, len(ml.messages))
	copy(result, ml.messages)
	return result
}

func (ml *MockLogger) ClearMessages() {
	ml.mutex.Lock()
	defer ml.mutex.Unlock()
	ml.messages = ml.messages[:0]
}
