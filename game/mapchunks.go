package game

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/yamamushi/EscapingEden/logging"
	"os"
)

func (gm *GameManager) SaveMapChunk(data MapChunk, filename string) error {

	filename = "./assets/world/" + filename
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
	filename = "./assets/world/" + filename
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
				// Initialize each point as a floor tile
				tiles[i][j][k] = Tile{
					TileType: "floor",
				}
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

func (gm *GameManager) GetTileInfo(tile *Tile) TileInfo {
	return gm.TileTypes[tile.TileType]
}

func (gm *GameManager) LoadTileTypes() {
	gm.TileTypes = make(map[string]TileInfo)
	gm.LoadTileTypesFromAssets("./assets/tiles")
}

func (gm *GameManager) LoadTileTypesFromAssets(directoryPath string) {
	// Check that tiles directory exists in assets/ and is not empty
	// If it is, panic and exit
	// If it isn't, load all the tiles into the map

	_, err := os.Stat(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			gm.Log.Println(logging.LogError, "Directory does not exist\n", directoryPath)
		} else {
			gm.Log.Println(logging.LogError, "Error checking directory: %v\n", err)
		}
		os.Exit(1)
	}

	files, err := os.ReadDir(directoryPath)
	if err != nil {
		gm.Log.Println(logging.LogError, "Error reading directory: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		gm.Log.Println(logging.LogError, "Directory is empty\n", directoryPath)
		os.Exit(1)
	}

	// Load all the tiles into the map, they are TileInfo json files
	// We will use the TileInfo.TileType as the key
	// We will use the TileInfo as the value

	for _, file := range files {
		// Load the file
		// Decode the file into a TileInfo struct
		// Add the TileInfo to the map
		if file.IsDir() {
			gm.LoadTileTypesFromAssets(directoryPath + "/" + file.Name())
		} else {
			content, err := os.ReadFile(directoryPath + "/" + file.Name())
			if err != nil {
				gm.Log.Println(logging.LogError, "Error reading file while loading tile types: \n", err)
				os.Exit(1)
			}

			var tileInfo []TileInfo
			err = json.Unmarshal(content, &tileInfo)
			if err != nil {
				gm.Log.Println(logging.LogError, "Error unmarshalling tile info: \n", err)
				os.Exit(1)
			}

			for _, tileType := range tileInfo {
				gm.TileTypes[tileType.TileType] = tileType
			}

		}
	}
}
