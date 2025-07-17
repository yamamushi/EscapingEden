package game

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

// ChunkMonitor provides HTTP endpoints for monitoring chunk loading performance
type ChunkMonitor struct {
	gm *GameManager
}

// NewChunkMonitor creates a new chunk monitor
func NewChunkMonitor(gm *GameManager) *ChunkMonitor {
	return &ChunkMonitor{gm: gm}
}

// RegisterEndpoints registers HTTP endpoints for chunk monitoring
func (cm *ChunkMonitor) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/chunks/stats", cm.handleChunkStats)
	mux.HandleFunc("/chunks/health", cm.handleChunkHealth)
	mux.HandleFunc("/chunks/cache", cm.handleCacheStats)
	mux.HandleFunc("/chunks/preload", cm.handlePreloadRequest)
	mux.HandleFunc("/chunks/registry", cm.handleRegistryInfo)
}

// handleChunkStats returns detailed chunk statistics
func (cm *ChunkMonitor) handleChunkStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats := cm.gm.GetWorldStats()

	response := map[string]interface{}{
		"world_stats": stats,
		"timestamp":   time.Now(),
	}

	if cm.gm.LazyLoader != nil {
		response["cache_stats"] = cm.gm.LazyLoader.GetCacheStats()
	}

	json.NewEncoder(w).Encode(response)
}

// handleChunkHealth returns chunk system health
func (cm *ChunkMonitor) handleChunkHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	health := cm.gm.GetWorldHealth()
	stats := cm.gm.GetWorldStats()

	response := map[string]interface{}{
		"status":      health,
		"world_ready": cm.gm.WorldValidated,
		"chunk_ratio": float64(stats.ValidChunks) / float64(stats.TotalChunks),
		"cache_usage": float64(stats.CachedChunks) / 100.0, // Assuming max cache of 100
		"timestamp":   time.Now(),
	}

	// Set HTTP status based on health
	switch health {
	case "healthy":
		w.WriteHeader(http.StatusOK)
	case "incomplete", "validating":
		w.WriteHeader(http.StatusAccepted)
	case "empty":
		w.WriteHeader(http.StatusServiceUnavailable)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(response)
}

// handleCacheStats returns detailed cache statistics
func (cm *ChunkMonitor) handleCacheStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"cached_chunks": len(cm.gm.ChunkCache),
		"chunk_ids":     make([]string, 0, len(cm.gm.ChunkCache)),
		"timestamp":     time.Now(),
	}

	// Add chunk IDs
	chunkIDs := response["chunk_ids"].([]string)
	for chunkID := range cm.gm.ChunkCache {
		chunkIDs = append(chunkIDs, chunkID)
	}
	response["chunk_ids"] = chunkIDs

	if cm.gm.LazyLoader != nil {
		cacheStats := cm.gm.LazyLoader.GetCacheStats()
		response["loader_stats"] = cacheStats
	}

	json.NewEncoder(w).Encode(response)
}

// handlePreloadRequest handles chunk preload requests
func (cm *ChunkMonitor) handlePreloadRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse parameters
	centerX, _ := strconv.Atoi(r.URL.Query().Get("x"))
	centerY, _ := strconv.Atoi(r.URL.Query().Get("y"))
	centerZ, _ := strconv.Atoi(r.URL.Query().Get("z"))
	radius, _ := strconv.Atoi(r.URL.Query().Get("radius"))

	if radius == 0 {
		radius = 2 // Default radius
	}

	// Trigger preload
	if cm.gm.LazyLoader != nil {
		cm.gm.LazyLoader.PreloadChunksAsync(centerX, centerY, centerZ, radius)

		response := map[string]interface{}{
			"status":    "preload_started",
			"center":    fmt.Sprintf("%d,%d,%d", centerX, centerY, centerZ),
			"radius":    radius,
			"timestamp": time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Lazy loader not available", http.StatusServiceUnavailable)
	}
}

// handleRegistryInfo returns chunk registry information
func (cm *ChunkMonitor) handleRegistryInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if cm.gm.ChunkRegistry == nil {
		http.Error(w, "Chunk registry not available", http.StatusServiceUnavailable)
		return
	}

	worldX, worldY, worldZ := cm.gm.ParseWorldDimensions(cm.gm.Config.WorldGen.Dimensions)
	missingChunks := cm.gm.ChunkRegistry.GetMissingChunks(worldX, worldY, worldZ)

	response := map[string]interface{}{
		"total_chunks":      cm.gm.ChunkRegistry.GetChunkCount(),
		"missing_chunks":    len(missingChunks),
		"missing_chunk_ids": missingChunks,
		"registry_version":  cm.gm.ChunkRegistry.Version,
		"world_dir":         cm.gm.ChunkRegistry.WorldDir,
		"timestamp":         time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

// StartPerformanceMonitoring starts background performance monitoring
func (cm *ChunkMonitor) StartPerformanceMonitoring() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			cm.gm.PrintWorldStats()
		}
	}()

	cm.gm.Log.Println(logging.LogInfo, "Chunk performance monitoring started")
}
