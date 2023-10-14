package game

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"errors"
	"github.com/yamamushi/EscapingEden/logging"
	"math/rand"
	"os"
)

func (gm *GameManager) LoadWorld() {
	//

	gm.Log.Println(logging.LogInfo, "Loading World...")
	//
	gm.Log.Println(logging.LogInfo, "Testing Map Chunk Creation...")

	/*
		mapChunk := gm.CreateMapChunk(255, 255, 255, 0, 0, 0, "Test")

		gm.Log.Println(logging.LogInfo, "Testing Map Chunk Saving...")
		err := gm.SaveMapChunk(mapChunk, "test.map")
		if err != nil {
			panic(err)
		}

	*/

	gm.Log.Println(logging.LogInfo, "Testing Map Chunk Loading...")
	loaded, err := gm.LoadMapChunk("test.map")
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

	/*
		gm.Log.Println(logging.LogInfo, "Checking Map Chunk Sizes...")
		if len(mapChunk.TileMap) != len(loadedChunk.TileMap) {
			panic("Map Chunk sizes do not match!")
		}
	*/

	gm.Log.Println(logging.LogInfo, "Loaded World!")

}

type MapChunk struct {
	ID string
	// The chunk's position in the world as a string [x,y,z]
	GlobalPosition struct {
		X int
		Y int
		Z int
	}
	// The chunk's data
	TileMap [][][]Tile
}

type Tile struct {
	Passable     bool
	BlocksVision bool
}

func (gm *GameManager) SaveMapChunk(data MapChunk, filename string) error {

	var _, err = os.Stat(filename)
	// create file if not exists
	if os.IsNotExist(err) {
		var buf bytes.Buffer

		// Compress using gzip
		zw := gzip.NewWriter(&buf)
		encoder := gob.NewEncoder(zw)

		if err := encoder.Encode(data); err != nil {
			return err
		}

		if err := zw.Close(); err != nil {
			return err
		}

		return os.WriteFile(filename, buf.Bytes(), 0644)
	}
	return errors.New("file " + filename + " already exists")
}

func (gm *GameManager) LoadMapChunk(filename string) (*MapChunk, error) {
	compressedFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	uncompressedFile, err := gzip.NewReader(bytes.NewReader(compressedFile))
	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(uncompressedFile)
	var chunk MapChunk
	if err := decoder.Decode(&chunk); err != nil {
		return nil, err
	}

	return &chunk, nil
}

func (gm *GameManager) CreateMapChunk(x, y, z int, gX, gY, gZ int, ID string) MapChunk { // Length of z is 3.

	// Create the 3D slice.
	tiles := make([][][]Tile, x)
	for i := range tiles {
		tiles[i] = make([][]Tile, y)
		for j := range tiles[i] {
			tiles[i][j] = make([]Tile, z)
			for k := range tiles[i][j] {
				// Initialize each point if needed.
				tiles[i][j][k] = Tile{Passable: true, BlocksVision: true} // Example initialization.
			}
		}
	}

	return MapChunk{
		ID: ID,
		GlobalPosition: struct {
			X int
			Y int
			Z int
		}{X: gX, Y: gY, Z: gZ},
		TileMap: tiles,
	}

}

func (gm *GameManager) DrawRect(x1, y1, x2, y2 int) {

	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	// Draw a rectangle with the given coordinates to gm.MapChunks[0].TileMap
	// We want to get the first layer of the map, z = 0 (the ground layer)
	for x := 0; x < len(gm.MapChunks[0].TileMap); x++ {
		for y := 0; y < len(gm.MapChunks[0].TileMap[0]); y++ {
			// If x or y is the rectangle border, draw a wall.
			if (x == x1 || x == x2) && (y >= y1 && y <= y2) { // Left or right edges
				gm.MapChunks[0].TileMap[x][y][0] = Tile{Passable: false, BlocksVision: true}
			} else if (y == y1 || y == y2) && (x >= x1 && x <= x2) { // Top or bottom edges
				gm.MapChunks[0].TileMap[x][y][0] = Tile{Passable: false, BlocksVision: true}
			} /* else if x > x1 && x < x2 && y > y1 && y < y2 { // Inside the rectangle
				fmt.Print("^")
			}*/
		}
	}
}
