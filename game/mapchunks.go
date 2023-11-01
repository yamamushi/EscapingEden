package game

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"errors"
	"github.com/google/uuid"
	"os"
)

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
				tiles[i][j][k] = Tile{
					ID: uuid.New().String(),
					Position: struct {
						X int
						Y int
						Z int
					}{X: i, Y: j, Z: k},
					BuilderID: "0",
					WallType:  "floor",
					Draw: struct {
						Character string
						Color     struct {
							FG string
							BG string
						}
					}{
						Character: ".",
						Color: struct {
							FG string
							BG string
						}{FG: "white", BG: "black"},
					},
					Passable:     true,
					BlocksVision: true,
				} // Example initialization.
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
