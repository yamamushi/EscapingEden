package edenutil

import (
	"context"
	"sync"
)

// ContextManager provides utilities for managing context-based shutdowns
type ContextManager struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewContextManager creates a new context manager
func NewContextManager() *ContextManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ContextManager{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Context returns the context
func (cm *ContextManager) Context() context.Context {
	return cm.ctx
}

// Cancel cancels the context
func (cm *ContextManager) Cancel() {
	cm.cancel()
}

// Add increments the WaitGroup counter
func (cm *ContextManager) Add(delta int) {
	cm.wg.Add(delta)
}

// Done decrements the WaitGroup counter
func (cm *ContextManager) Done() {
	cm.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero
func (cm *ContextManager) Wait() {
	cm.wg.Wait()
}

// GoWithContext runs a function in a goroutine with context awareness
func (cm *ContextManager) GoWithContext(fn func(context.Context)) {
	cm.wg.Add(1)
	go func() {
		defer cm.wg.Done()
		fn(cm.ctx)
	}()
}
