package game

import (
	"fmt"

	"github.com/yamamushi/EscapingEden/logging"
)

func (gm *GameManager) GetMapChunkByID(id string) *MapChunk {
	return gm.GetMapChunkByIDWithLoading(id, "")
}

func (gm *GameManager) GetMapChunkByIDWithLoading(id string, characterID string) *MapChunk {
	// First check if chunk is already in cache by UUID
	for _, chunk := range gm.ChunkCache {
		if chunk.ID == id {
			return chunk
		}
	}

	// If the ID looks like coordinates (e.g., "0-0-0"), check cache by coordinate key
	var globalX, globalY, globalZ int
	n, err := fmt.Sscanf(id, "%d-%d-%d", &globalX, &globalY, &globalZ)
	if err == nil && n == 3 {
		// This is a coordinate-based ID, check cache directly
		if chunk, ok := gm.ChunkCache[id]; ok {
			return chunk
		}

		// Chunk needs to be loaded from disk - send loading message if we have a character ID
		if characterID != "" {
			gm.SendLoadingMessage(characterID, "Loading world data, please wait...")
		}

		// Try to load from disk using coordinates
		if gm.LazyLoader != nil {
			chunk, err := gm.LazyLoader.GetChunkAsync(globalX, globalY, globalZ)
			if err != nil {
				gm.Log.Println(logging.LogWarn, "Failed to load chunk", id, "from disk:", err.Error())
				if characterID != "" {
					gm.SendLoadingMessage(characterID, "Failed to load world data. Please try again.")
				}
				return nil
			}
			gm.Log.Println(logging.LogInfo, "Successfully loaded chunk", id, "from disk")
			if characterID != "" {
				gm.SendLoadingMessage(characterID, "World loaded successfully!")
			}
			return chunk
		}

		// Fallback: try to load directly from disk
		filename := id + ".map"
		chunk, err := gm.LoadMapChunk(filename)
		if err != nil {
			gm.Log.Println(logging.LogWarn, "Failed to load chunk", id, "from disk:", err.Error())
			if characterID != "" {
				gm.SendLoadingMessage(characterID, "Failed to load world data. Please try again.")
			}
			return nil
		}

		// Add to cache and return
		gm.ChunkCache[id] = chunk
		gm.Log.Println(logging.LogInfo, "Successfully loaded chunk", id, "from disk (fallback)")
		if characterID != "" {
			gm.SendLoadingMessage(characterID, "World loaded successfully!")
		}
		return chunk
	}

	// This is a UUID-based ID, but chunk is not in cache
	// We need to search all chunks on disk to find the one with this UUID
	// This is expensive, so we should log it
	gm.Log.Println(logging.LogWarn, "Chunk with UUID", id, "not found in cache - this may indicate a data consistency issue")
	if characterID != "" {
		gm.SendLoadingMessage(characterID, "World data inconsistency detected. Moving to safe location...")
	}

	// For now, return nil to trigger the fallback system
	// TODO: In the future, we could implement a UUID->coordinates mapping
	return nil
}
