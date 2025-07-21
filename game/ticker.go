package game

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

// TickSubscriber interface for components that need tick updates
type TickSubscriber interface {
	OnTick(tickNumber uint64)
	GetSubscriberID() string
}

// GameTicker manages the global game tick system
type GameTicker struct {
	tickRate    time.Duration
	currentTick uint64
	running     bool
	subscribers []TickSubscriber
	mutex       sync.RWMutex
	stopChan    chan struct{}
	ticker      *time.Ticker
	log         logging.LoggerType
}

// NewGameTicker creates a new GameTicker with the specified tick rate
func NewGameTicker(tickRate time.Duration, log logging.LoggerType) *GameTicker {
	return &GameTicker{
		tickRate:    tickRate,
		currentTick: 0,
		running:     false,
		subscribers: make([]TickSubscriber, 0),
		stopChan:    make(chan struct{}),
		log:         log,
	}
}

// Start begins the tick processing loop
func (gt *GameTicker) Start() error {
	gt.mutex.Lock()
	defer gt.mutex.Unlock()

	if gt.running {
		return errors.New("ticker is already running")
	}

	gt.running = true
	gt.ticker = time.NewTicker(gt.tickRate)
	gt.stopChan = make(chan struct{})

	gt.log.Println(logging.LogInfo, fmt.Sprintf("GameTicker starting with %v interval", gt.tickRate))

	// Start the tick processing goroutine
	go gt.tickLoop()

	return nil
}

// Stop gracefully shuts down the tick processing
func (gt *GameTicker) Stop() error {
	gt.mutex.Lock()
	defer gt.mutex.Unlock()

	if !gt.running {
		return errors.New("ticker is not running")
	}

	gt.log.Println(logging.LogInfo, "GameTicker stopping...")

	gt.running = false

	if gt.ticker != nil {
		gt.ticker.Stop()
	}

	// Signal the tick loop to stop
	close(gt.stopChan)

	gt.log.Println(logging.LogInfo, "GameTicker stopped successfully")
	return nil
}

// Subscribe adds a new subscriber to receive tick notifications
func (gt *GameTicker) Subscribe(subscriber TickSubscriber) error {
	if subscriber == nil {
		return errors.New("subscriber cannot be nil")
	}

	gt.mutex.Lock()
	defer gt.mutex.Unlock()

	// Check if subscriber already exists
	subscriberID := subscriber.GetSubscriberID()
	for _, existing := range gt.subscribers {
		if existing.GetSubscriberID() == subscriberID {
			return fmt.Errorf("subscriber with ID '%s' already exists", subscriberID)
		}
	}

	gt.subscribers = append(gt.subscribers, subscriber)
	gt.log.Println(logging.LogInfo, fmt.Sprintf("Subscriber '%s' added to GameTicker", subscriberID))

	return nil
}

// Unsubscribe removes a subscriber from tick notifications
func (gt *GameTicker) Unsubscribe(subscriberID string) error {
	gt.mutex.Lock()
	defer gt.mutex.Unlock()

	for i, subscriber := range gt.subscribers {
		if subscriber.GetSubscriberID() == subscriberID {
			// Remove subscriber by slicing
			gt.subscribers = append(gt.subscribers[:i], gt.subscribers[i+1:]...)
			gt.log.Println(logging.LogInfo, fmt.Sprintf("Subscriber '%s' removed from GameTicker", subscriberID))
			return nil
		}
	}

	return fmt.Errorf("subscriber with ID '%s' not found", subscriberID)
}

// GetCurrentTick returns the current tick number
func (gt *GameTicker) GetCurrentTick() uint64 {
	gt.mutex.RLock()
	defer gt.mutex.RUnlock()
	return gt.currentTick
}

// GetSubscriberCount returns the number of active subscribers
func (gt *GameTicker) GetSubscriberCount() int {
	gt.mutex.RLock()
	defer gt.mutex.RUnlock()
	return len(gt.subscribers)
}

// IsRunning returns whether the ticker is currently running
func (gt *GameTicker) IsRunning() bool {
	gt.mutex.RLock()
	defer gt.mutex.RUnlock()
	return gt.running
}

// tickLoop is the main tick processing loop (runs in its own goroutine)
func (gt *GameTicker) tickLoop() {
	defer func() {
		if r := recover(); r != nil {
			gt.log.Println(logging.LogError, fmt.Sprintf("GameTicker panic recovered: %v", r))
			// Attempt to restart the ticker
			go gt.handleTickerRestart()
		}
	}()

	for {
		select {
		case <-gt.ticker.C:
			gt.processTick()
		case <-gt.stopChan:
			gt.log.Println(logging.LogInfo, "GameTicker tick loop stopping")
			return
		}
	}
}

// processTick handles a single tick and notifies all subscribers
func (gt *GameTicker) processTick() {
	startTime := time.Now()

	// Increment tick counter
	gt.mutex.Lock()
	gt.currentTick++
	currentTick := gt.currentTick

	// Create a copy of subscribers to avoid holding the lock during notifications
	subscribersCopy := make([]TickSubscriber, len(gt.subscribers))
	copy(subscribersCopy, gt.subscribers)
	gt.mutex.Unlock()

	// Notify all subscribers
	for _, subscriber := range subscribersCopy {
		func() {
			defer func() {
				if r := recover(); r != nil {
					gt.log.Println(logging.LogError, fmt.Sprintf("Subscriber '%s' panic during tick %d: %v",
						subscriber.GetSubscriberID(), currentTick, r))
				}
			}()

			subscriber.OnTick(currentTick)
		}()
	}

	// Log performance warnings if tick processing takes too long
	processingTime := time.Since(startTime)
	maxProcessingTime := gt.tickRate / 4 // 25% of tick rate

	if processingTime > maxProcessingTime {
		gt.log.Println(logging.LogWarn, fmt.Sprintf("Tick %d processing took %v (threshold: %v)",
			currentTick, processingTime, maxProcessingTime))
	}
}

// handleTickerRestart attempts to restart the ticker after a panic
func (gt *GameTicker) handleTickerRestart() {
	gt.log.Println(logging.LogWarn, "Attempting to restart GameTicker after panic")

	// Wait a bit before restarting
	time.Sleep(time.Second)

	gt.mutex.Lock()
	if gt.running {
		gt.mutex.Unlock()

		// Create a new ticker and restart the loop
		gt.ticker = time.NewTicker(gt.tickRate)
		go gt.tickLoop()

		gt.log.Println(logging.LogInfo, "GameTicker restarted successfully")
	} else {
		gt.mutex.Unlock()
		gt.log.Println(logging.LogInfo, "GameTicker restart cancelled - ticker was stopped")
	}
}

// GetTickRate returns the current tick rate
func (gt *GameTicker) GetTickRate() time.Duration {
	gt.mutex.RLock()
	defer gt.mutex.RUnlock()
	return gt.tickRate
}

// SetTickRate updates the tick rate (requires restart to take effect)
func (gt *GameTicker) SetTickRate(newRate time.Duration) error {
	if newRate <= 0 {
		return errors.New("tick rate must be positive")
	}

	gt.mutex.Lock()
	defer gt.mutex.Unlock()

	gt.tickRate = newRate
	gt.log.Println(logging.LogInfo, fmt.Sprintf("GameTicker tick rate updated to %v", newRate))

	// If running, restart with new rate
	if gt.running {
		gt.ticker.Stop()
		gt.ticker = time.NewTicker(gt.tickRate)
	}

	return nil
}
