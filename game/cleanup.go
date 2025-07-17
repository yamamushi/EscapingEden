package game

import (
	"fmt"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

func (gm *GameManager) Cleanup() {
	gm.Log.Println(logging.LogInfo, "Starting server cleanup...")

	// Save all active characters with timeout protection
	gm.Log.Println(logging.LogInfo, "Saving active characters during cleanup...")

	// Create a channel to signal completion of character saving
	charactersDone := make(chan struct{})

	// Start character saving in a goroutine
	go func() {
		defer close(charactersDone)

		gm.activeCharactersMutex.Lock()
		defer gm.activeCharactersMutex.Unlock()

		characterCount := len(gm.ActiveCharacters)
		savedCount := 0
		errorCount := 0

		gm.Log.Println(logging.LogInfo, "Found", characterCount, "active characters to save")

		for _, character := range gm.ActiveCharacters {
			err := gm.DB.UpdateRecord("Characters", character.Record)
			if err != nil {
				gm.Log.Println(logging.LogError, "Failed to save character", character.Name, "during cleanup:", err.Error())
				errorCount++
			} else {
				gm.Log.Println(logging.LogInfo, "Saved character", character.Name, "at position", character.Record.Position.X, character.Record.Position.Y, "in map", character.Record.CurrentMapID, "chunk", character.Record.Position.MapChunkID)
				savedCount++
			}
		}

		gm.Log.Println(logging.LogInfo, "Character save completed:", savedCount, "saved,", errorCount, "errors out of", characterCount, "total")
	}()

	// Wait for character saving with timeout
	characterTimeout := 10 * time.Second
	select {
	case <-charactersDone:
		gm.Log.Println(logging.LogInfo, "All characters saved successfully")
	case <-time.After(characterTimeout):
		gm.Log.Println(logging.LogWarn, "Character save timed out after", characterTimeout, "- forcing shutdown to prevent hanging")
	}

	// Save all map chunks with timeout protection
	gm.Log.Println(logging.LogInfo, "Saving map chunks during cleanup...")

	// Create a channel to signal completion of chunk saving
	chunksDone := make(chan struct{})

	// Start chunk saving in a goroutine
	go func() {
		defer close(chunksDone)

		chunkCount := len(gm.ChunkCache)
		savedCount := 0
		errorCount := 0

		gm.Log.Println(logging.LogInfo, "Found", chunkCount, "chunks to save")

		for id, mapChunk := range gm.ChunkCache {
			filename := fmt.Sprintf("%d-%d-%d.map", mapChunk.GlobalPosition.X, mapChunk.GlobalPosition.Y, mapChunk.GlobalPosition.Z)
			err := gm.SaveMapChunk(*mapChunk, filename, true)
			if err != nil {
				gm.Log.Println(logging.LogError, "Failed to save map chunk", id, ":", err.Error())
				errorCount++
			} else {
				savedCount++
			}
		}

		gm.Log.Println(logging.LogInfo, "Chunk save completed:", savedCount, "saved,", errorCount, "errors out of", chunkCount, "total")
	}()

	// Wait for chunk saving with timeout
	chunkTimeout := 20 * time.Second
	select {
	case <-chunksDone:
		gm.Log.Println(logging.LogInfo, "All chunks saved successfully")
	case <-time.After(chunkTimeout):
		gm.Log.Println(logging.LogWarn, "Chunk save timed out after", chunkTimeout, "- forcing shutdown to prevent hanging")
	}

	gm.Log.Println(logging.LogInfo, "Cleanup completed - server ready for shutdown")
}
