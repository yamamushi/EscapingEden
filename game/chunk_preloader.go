package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

// ChunkPreloader handles intelligent chunk preloading
type ChunkPreloader struct {
	gm             *GameManager
	preloadRadius  int
	maxConcurrent  int
	preloadQueue   chan PreloadRequest
	activePreloads map[string]bool
	preloadMutex   sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
}

// PreloadRequest represents a chunk preload request
type PreloadRequest struct {
	CenterX, CenterY, CenterZ int
	Radius                    int
	Priority                  int
}

// NewChunkPreloader creates a new chunk preloader
func NewChunkPreloader(gm *GameManager, preloadRadius, maxConcurrent int) *ChunkPreloader {
	ctx, cancel := context.WithCancel(context.Background())

	cp := &ChunkPreloader{
		gm:             gm,
		preloadRadius:  preloadRadius,
		maxConcurrent:  maxConcurrent,
		preloadQueue:   make(chan PreloadRequest, 100),
		activePreloads: make(map[string]bool),
		ctx:            ctx,
		cancel:         cancel,
	}

	// Start preload workers
	for i := 0; i < maxConcurrent; i++ {
		go cp.preloadWorker(i)
	}

	return cp
}

// RequestPreload requests chunks to be preloaded around a position
func (cp *ChunkPreloader) RequestPreload(centerX, centerY, centerZ, radius, priority int) {
	select {
	case cp.preloadQueue <- PreloadRequest{
		CenterX:  centerX,
		CenterY:  centerY,
		CenterZ:  centerZ,
		Radius:   radius,
		Priority: priority,
	}:
		// Request queued successfully
	default:
		// Queue is full, skip this request
		cp.gm.Log.Println(logging.LogDebug, "Preload queue full, skipping request")
	}
}

// preloadWorker processes preload requests
func (cp *ChunkPreloader) preloadWorker(workerID int) {
	cp.gm.Log.Println(logging.LogDebug, fmt.Sprintf("Preload worker %d started", workerID))

	for {
		select {
		case <-cp.ctx.Done():
			cp.gm.Log.Println(logging.LogDebug, fmt.Sprintf("Preload worker %d stopping", workerID))
			return
		case request := <-cp.preloadQueue:
			cp.processPreloadRequest(request, workerID)
		}
	}
}

// processPreloadRequest processes a single preload request
func (cp *ChunkPreloader) processPreloadRequest(request PreloadRequest, workerID int) {
	startTime := time.Now()
	loaded := 0
	skipped := 0

	for x := request.CenterX - request.Radius; x <= request.CenterX+request.Radius; x++ {
		for y := request.CenterY - request.Radius; y <= request.CenterY+request.Radius; y++ {
			for z := request.CenterZ - request.Radius; z <= request.CenterZ+request.Radius; z++ {
				select {
				case <-cp.ctx.Done():
					return
				default:
					chunkID := fmt.Sprintf("%d-%d-%d", x, y, z)

					// Check if already being preloaded
					cp.preloadMutex.RLock()
					isActive := cp.activePreloads[chunkID]
					cp.preloadMutex.RUnlock()

					if isActive {
						skipped++
						continue
					}

					// Check if already in cache
					if _, ok := cp.gm.ChunkCache[chunkID]; ok {
						skipped++
						continue
					}

					// Mark as being preloaded
					cp.preloadMutex.Lock()
					cp.activePreloads[chunkID] = true
					cp.preloadMutex.Unlock()

					// Preload the chunk
					_, err := cp.gm.LazyLoader.GetChunkAsync(x, y, z)

					// Unmark
					cp.preloadMutex.Lock()
					delete(cp.activePreloads, chunkID)
					cp.preloadMutex.Unlock()

					if err != nil {
						cp.gm.Log.Println(logging.LogDebug, fmt.Sprintf("Worker %d failed to preload chunk %s: %v", workerID, chunkID, err))
					} else {
						loaded++
					}
				}
			}
		}
	}

	duration := time.Since(startTime)
	cp.gm.Log.Println(logging.LogDebug, fmt.Sprintf("Worker %d preloaded %d chunks (%d skipped) in %v",
		workerID, loaded, skipped, duration))
}

// Stop stops the chunk preloader
func (cp *ChunkPreloader) Stop() {
	cp.cancel()
	close(cp.preloadQueue)
}

// GetStats returns preloader statistics
func (cp *ChunkPreloader) GetStats() map[string]interface{} {
	cp.preloadMutex.RLock()
	defer cp.preloadMutex.RUnlock()

	return map[string]interface{}{
		"queue_length":    len(cp.preloadQueue),
		"active_preloads": len(cp.activePreloads),
		"max_concurrent":  cp.maxConcurrent,
		"preload_radius":  cp.preloadRadius,
	}
}
