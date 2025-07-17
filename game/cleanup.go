package game

import (
	"fmt"
	"github.com/yamamushi/EscapingEden/logging"
)

func (gm *GameManager) Cleanup() {
	for id, mapChunk := range gm.ChunkCache {
		filename := fmt.Sprintf("%d-%d-%d.map", mapChunk.GlobalPosition.X, mapChunk.GlobalPosition.Y, mapChunk.GlobalPosition.Z)
		err := gm.SaveMapChunk(*mapChunk, filename, true)
		if err != nil {
			gm.Log.Println(logging.LogError, "Failed to save map chunk", id, ":", err.Error())
		}
	}
}
