package game

import (
	"fmt"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

// WorldStats provides statistics about the world and chunk loading
type WorldStats struct {
	TotalChunks     int           `json:"total_chunks"`
	ValidChunks     int           `json:"valid_chunks"`
	CachedChunks    int           `json:"cached_chunks"`
	LoadingChunks   int           `json:"loading_chunks"`
	GenerationTime  time.Duration `json:"generation_time_ms"`
	ValidationTime  time.Duration `json:"validation_time_ms"`
	LastValidation  time.Time     `json:"last_validation"`
	WorldSize       string        `json:"world_size"`
	ChunkSize       int           `json:"chunk_size"`
	RegistryVersion int           `json:"registry_version"`
}

// GetWorldStats returns current world statistics
func (gm *GameManager) GetWorldStats() WorldStats {
	worldX, worldY, worldZ := gm.ParseWorldDimensions(gm.Config.WorldGen.Dimensions)
	totalChunks := worldX * worldY * worldZ

	var validChunks int
	if gm.ChunkRegistry != nil {
		validChunks = gm.ChunkRegistry.GetValidChunkCount(worldX, worldY, worldZ)
	}

	var cacheStats map[string]interface{}
	if gm.LazyLoader != nil {
		cacheStats = gm.LazyLoader.GetCacheStats()
	}

	stats := WorldStats{
		TotalChunks:    totalChunks,
		ValidChunks:    validChunks,
		CachedChunks:   len(gm.ChunkCache),
		WorldSize:      gm.Config.WorldGen.Dimensions,
		ChunkSize:      gm.Config.WorldGen.ChunkSize,
		LastValidation: time.Now(),
	}

	if cacheStats != nil {
		if loadingChunks, ok := cacheStats["loading_chunks"].(int); ok {
			stats.LoadingChunks = loadingChunks
		}
	}

	if gm.ChunkRegistry != nil {
		stats.RegistryVersion = gm.ChunkRegistry.Version
	}

	return stats
}

// GetWorldHealth returns a health status for the world
func (gm *GameManager) GetWorldHealth() string {
	stats := gm.GetWorldStats()

	if !gm.WorldValidated {
		return "validating"
	}

	if stats.ValidChunks == 0 {
		return "empty"
	}

	if stats.ValidChunks < stats.TotalChunks {
		return "incomplete"
	}

	return "healthy"
}

// PrintWorldStats prints world statistics to the log
func (gm *GameManager) PrintWorldStats() {
	stats := gm.GetWorldStats()
	health := gm.GetWorldHealth()

	gm.Log.Println(logging.LogInfo, fmt.Sprintf("World Stats - Health: %s", health))
	gm.Log.Println(logging.LogInfo, fmt.Sprintf("  Chunks: %d/%d valid (%d cached, %d loading)",
		stats.ValidChunks, stats.TotalChunks, stats.CachedChunks, stats.LoadingChunks))
	gm.Log.Println(logging.LogInfo, fmt.Sprintf("  World Size: %s, Chunk Size: %d",
		stats.WorldSize, stats.ChunkSize))

	if gm.LazyLoader != nil {
		cacheStats := gm.LazyLoader.GetCacheStats()
		if usage, ok := cacheStats["cache_usage_pct"].(float64); ok {
			gm.Log.Println(logging.LogInfo, fmt.Sprintf("  Cache Usage: %.1f%%", usage))
		}
	}
}
