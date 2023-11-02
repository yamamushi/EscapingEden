package game

import (
	"errors"
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/logging"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func (gm *GameManager) LoadWorld() {
	//

	gm.Log.Println(logging.LogInfo, "Loading World...")
	//
	worldDimensionsString := gm.Config.World.Dimensions
	worldX, worldY, worldZ := gm.ParseWorldDimensions(worldDimensionsString)
	gm.Log.Println(logging.LogInfo, "Checking for existing map...")
	err := gm.ValidateWorldFiles(worldX, worldY, worldZ)
	if err != nil {
		if err.Error() == "empty" {
			gm.Log.Println(logging.LogInfo, "Generating new world...")
			gm.Log.Println(logging.LogInfo, "Creating Map Chunks...")

			for i := 0; i < worldX; i++ {
				for j := 0; j < worldY; j++ {
					for k := 0; k < worldZ; k++ {
						mapChunkID := uuid.New().String()
						mapChunk := gm.CreateMapChunk(255, 255, 255, i, j, k, mapChunkID)

						err := gm.SaveMapChunk(mapChunk, strconv.Itoa(mapChunk.GlobalPosition.X)+"-"+strconv.Itoa(mapChunk.GlobalPosition.Y)+"-"+strconv.Itoa(mapChunk.GlobalPosition.Z)+".map")
						if err != nil {
							gm.Log.Println(logging.LogError, "Failed to save map chunk:", err.Error())
						}
						gm.MapChunks = append(gm.MapChunks, mapChunk)
					}
				}
			}
		} else {
			panic(err)
		}
	} else {
		gm.Log.Println(logging.LogInfo, "Loading existing world...")
		for i := 0; i < worldX; i++ {
			for j := 0; j < worldY; j++ {
				for k := 0; k < worldZ; k++ {
					mapChunkFilename := strconv.Itoa(i) + "-" + strconv.Itoa(j) + "-" + strconv.Itoa(k) + ".map"
					mapChunk, err := gm.LoadMapChunk(mapChunkFilename)
					if err != nil {
						panic(err)
					}
					gm.MapChunks = append(gm.MapChunks, *mapChunk)
				}
			}
		}
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
	if len(loaded.TileMap[0][0]) == 0 {
		panic("Map Chunk at [0][0] is empty!")
	}
	// append to the map chunks
	for _, mapChunk := range gm.MapChunks {
		for i := 0; i < 10; i++ {
			// generate random x1, y1, x2, y2 between 0 and 100 using go rand
			x1 := rand.Intn(200) // 0 - 100
			y1 := rand.Intn(200) // 0 - 100
			x2 := rand.Intn(200) // 0 - 100
			y2 := rand.Intn(200) // 0 - 100
			gm.DrawRect(&mapChunk, x1, y1, x2, y2)

		}
		gm.DrawRect(&mapChunk, 0, 0, 200, 200)
	}

	/*
		gm.Log.Println(logging.LogInfo, "Checking Map Chunk Sizes...")
		if len(mapChunk.TileMap) != len(loadedChunk.TileMap) {
			panic("Map Chunk sizes do not match!")
		}
	*/

	gm.Log.Println(logging.LogInfo, "Loaded World!")
}

func (gm *GameManager) ValidateWorldFiles(x, y, z int) error {
	worldDir := "assets/world/"
	// check if world directory exists
	// if not, create it
	_, err := os.Stat(worldDir)
	if err != nil {
		if os.IsNotExist(err) {
			// create directory
			err := os.Mkdir(worldDir, 0755)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	var missingFiles []string
	// Check for all map chunks in range to validate they exist
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			for k := 0; k < z; k++ {
				mapChunkFilename := strconv.Itoa(i) + "-" + strconv.Itoa(j) + "-" + strconv.Itoa(k) + ".map"
				_, err := os.Stat(worldDir + mapChunkFilename)
				if err != nil {
					if os.IsNotExist(err) {
						missingFiles = append(missingFiles, mapChunkFilename)
					} else {
						return err
					}
				}
			}
		}
	}

	if len(missingFiles) > 0 && len(missingFiles) != (x)*(y)*(z) {
		return errors.New("world files are missing, will not regenerate missing chunks. Aborting")
	}
	if len(missingFiles) == (x)*(y)*(z) {
		gm.Log.Println(logging.LogInfo, "No world files found!")
		return errors.New("empty")
	}
	return nil
}

func (gm *GameManager) ParseWorldDimensions(dimensions string) (int, int, int) {
	stringCollection := strings.Split(dimensions, ",")
	if len(stringCollection) != 3 {
		panic("World dimensions are invalid, please check config!")
	}

	x, err := strconv.Atoi(stringCollection[0])
	if err != nil {
		panic(err)
	}
	if x < 1 {
		panic("World dimensions are invalid, expecting x larger than 0!")
	}
	y, err := strconv.Atoi(stringCollection[1])
	if err != nil {
		panic(err)
	}
	if y < 1 {
		panic("World dimensions are invalid, expecting y larger than 0!")
	}
	z, err := strconv.Atoi(stringCollection[2])
	if err != nil {
		panic(err)
	}
	if z < 1 {
		panic("World dimensions are invalid, expecting z larger than 0!")
	}
	return x, y, z
}
