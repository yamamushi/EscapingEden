package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

// LazyChunkLoader handles on-demand chunk loading
type LazyChunkLoader struct {
	gm                     *GameManager
	loadingMutex           sync.RWMutex
	loadingChunks          map[string]chan *MapChunk // Track chunks currently being loaded
	maxCacheSize           int
	cacheEvictionThreshold int
}

// NewLazyChunkLoader creates a new lazy chunk loader
func NewLazyChunkLoader(gm *GameManager, maxCacheSize int) *LazyChunkLoader {
	return &LazyChunkLoader{
		gm:                     gm,
		loadingChunks:          make(map[string]chan *MapChunk),
		maxCacheSize:           maxCacheSize,
		cacheEvictionThreshold: int(float64(maxCacheSize) * 0.8), // Evict when 80% full
	}
}

// GetChunkAsync loads a chunk asynchronously
func (lcl *LazyChunkLoader) GetChunkAsync(globalX, globalY, globalZ int) (*MapChunk, error) {
	chunkID := fmt.Sprintf("%d-%d-%d", globalX, globalY, globalZ)

	// Check cache first
	if chunk, ok := lcl.gm.ChunkCache[chunkID]; ok {
		return chunk, nil
	}

	// Check if chunk is already being loaded
	lcl.loadingMutex.RLock()
	if loadingChan, isLoading := lcl.loadingChunks[chunkID]; isLoading {
		lcl.loadingMutex.RUnlock()

		// Wait for the chunk to finish loading
		select {
		case chunk := <-loadingChan:
			if chunk == nil {
				return nil, fmt.Errorf("failed to load chunk %s", chunkID)
			}
			return chunk, nil
		case <-time.After(30 * time.Second):
			return nil, fmt.Errorf("timeout waiting for chunk %s to load", chunkID)
		}
	}
	lcl.loadingMutex.RUnlock()

	// Start loading the chunk
	return lcl.loadChunk(chunkID, globalX, globalY, globalZ)
}

// loadChunk loads a single chunk
func (lcl *LazyChunkLoader) loadChunk(chunkID string, globalX, globalY, globalZ int) (*MapChunk, error) {
	// Create loading channel
	loadingChan := make(chan *MapChunk, 1)

	lcl.loadingMutex.Lock()
	lcl.loadingChunks[chunkID] = loadingChan
	lcl.loadingMutex.Unlock()

	// Ensure cleanup
	defer func() {
		lcl.loadingMutex.Lock()
		delete(lcl.loadingChunks, chunkID)
		close(loadingChan)
		lcl.loadingMutex.Unlock()
	}()

	// Load the chunk
	filename := chunkID + ".map"
	chunk, err := lcl.gm.LoadMapChunk(filename)
	if err != nil {
		lcl.gm.Log.Println(logging.LogError, "Failed to load chunk", chunkID, ":", err)
		loadingChan <- nil
		return nil, err
	}

	// Check cache size and evict if necessary
	lcl.evictIfNecessary()

	// Store in cache
	lcl.gm.ChunkCache[chunkID] = chunk
	loadingChan <- chunk

	lcl.gm.Log.Println(logging.LogDebug, "Loaded chunk", chunkID)
	return chunk, nil
}

// evictIfNecessary evicts chunks from cache if it's getting full
func (lcl *LazyChunkLoader) evictIfNecessary() {
	if len(lcl.gm.ChunkCache) < lcl.cacheEvictionThreshold {
		return
	}

	// Simple LRU-like eviction - remove some chunks
	// In a real implementation, you'd track access times
	evictCount := len(lcl.gm.ChunkCache) - lcl.maxCacheSize/2

	for chunkID := range lcl.gm.ChunkCache {
		if evictCount <= 0 {
			break
		}

		delete(lcl.gm.ChunkCache, chunkID)
		evictCount--

		lcl.gm.Log.Println(logging.LogDebug, "Evicted chunk from cache:", chunkID)
	}
}

// PreloadChunksAsync preloads chunks around a position
func (lcl *LazyChunkLoader) PreloadChunksAsync(centerX, centerY, centerZ, radius int) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		for x := centerX - radius; x <= centerX+radius; x++ {
			for y := centerY - radius; y <= centerY+radius; y++ {
				for z := centerZ - radius; z <= centerZ+radius; z++ {
					select {
					case <-ctx.Done():
						return
					default:
						// Skip if already in cache
						chunkID := fmt.Sprintf("%d-%d-%d", x, y, z)
						if _, ok := lcl.gm.ChunkCache[chunkID]; ok {
							continue
						}

						// Load chunk asynchronously
						go func(gx, gy, gz int) {
							_, err := lcl.GetChunkAsync(gx, gy, gz)
							if err != nil {
								lcl.gm.Log.Println(logging.LogDebug, "Failed to preload chunk:", err)
							}
						}(x, y, z)
					}
				}
			}
		}
	}()
}

// GetCacheStats returns cache statistics
func (lcl *LazyChunkLoader) GetCacheStats() map[string]interface{} {
	lcl.loadingMutex.RLock()
	defer lcl.loadingMutex.RUnlock()

	return map[string]interface{}{
		"cached_chunks":   len(lcl.gm.ChunkCache),
		"loading_chunks":  len(lcl.loadingChunks),
		"max_cache_size":  lcl.maxCacheSize,
		"cache_usage_pct": float64(len(lcl.gm.ChunkCache)) / float64(lcl.maxCacheSize) * 100,
	}
}
