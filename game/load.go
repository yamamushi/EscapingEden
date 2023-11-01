package game

import (
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/logging"
	"math/rand"
)

func (gm *GameManager) LoadWorld() {
	//

	gm.Log.Println(logging.LogInfo, "Loading World...")
	//
	gm.Log.Println(logging.LogInfo, "Testing Map Chunk Creation...")

	mapChunkID := uuid.New().String()
	mapChunk := gm.CreateMapChunk(255, 255, 255, 0, 0, 0, mapChunkID)

	gm.Log.Println(logging.LogInfo, "Testing Map Chunk Saving...")
	err := gm.SaveMapChunk(mapChunk, "0-0-0.map")
	if err != nil {
		gm.Log.Println(logging.LogError, "Failed to save map chunk:", err.Error())
	}

	gm.Log.Println(logging.LogInfo, "Testing Map Chunk Loading...")
	loaded, err := gm.LoadMapChunk("0-0-0.map")
	if err != nil {
		panic(err) // We'll refactor all of this later, we're just implementing stuff quickly for now.
	}
	if len(loaded.TileMap) == 0 {
		panic("Map Chunk is empty!")
	}
	if len(loaded.TileMap[0]) == 0 {
		panic("Map Chunk at [0] is empty!")
	}
	// append to the map chunks
	gm.MapChunks = append(gm.MapChunks, *loaded)
	for i := 0; i < 10; i++ {
		// generate random x1, y1, x2, y2 between 0 and 100 using go rand
		x1 := rand.Intn(200) // 0 - 100
		y1 := rand.Intn(200) // 0 - 100
		x2 := rand.Intn(200) // 0 - 100
		y2 := rand.Intn(200) // 0 - 100
		gm.DrawRect(x1, y1, x2, y2)

	}
	gm.DrawRect(0, 0, 200, 200)

	/*
		gm.Log.Println(logging.LogInfo, "Checking Map Chunk Sizes...")
		if len(mapChunk.TileMap) != len(loadedChunk.TileMap) {
			panic("Map Chunk sizes do not match!")
		}
	*/

	gm.Log.Println(logging.LogInfo, "Loaded World!")

}
