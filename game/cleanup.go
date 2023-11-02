package game

import (
	"github.com/yamamushi/EscapingEden/logging"
	"strconv"
)

func (gm *GameManager) Cleanup() {
	for _, mapChunk := range gm.MapChunks {
		err := gm.SaveMapChunk(mapChunk, strconv.Itoa(mapChunk.GlobalPosition.X)+"-"+strconv.Itoa(mapChunk.GlobalPosition.Y)+"-"+strconv.Itoa(mapChunk.GlobalPosition.Z)+".map")
		if err != nil {
			gm.Log.Println(logging.LogError, "Failed to save map chunk:", err.Error())
		}
	}
}
