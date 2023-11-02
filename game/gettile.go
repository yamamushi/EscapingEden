package game

import (
	"github.com/yamamushi/EscapingEden/logging"
	"log"
)

func (gm *GameManager) GetTileFromCharacter(charID string, x, y, z int) (mapChunk *MapChunk, tile *Tile, X, Y, Z int) {
	character, err := gm.GetCharacter(charID)
	if err != nil {
		gm.Log.Println(logging.LogError, "Failed to get character", err.Error())
		return nil, nil, 0, 0, 0
	}

	targetMap := gm.GetMapChunkByID(character.CurrentMapID)
	if x > len(targetMap.TileMap)-1 || x < 0 || y > len(targetMap.TileMap[0])-1 || y < 0 {
		dX, dY := 0, 0
		if x > len(targetMap.TileMap)-1 {
			dX = 1
			x = 0
		}
		if x < 0 {
			dX = -1
			x = len(targetMap.TileMap) - 1
		}
		if y > len(targetMap.TileMap[0])-1 {
			dY = 1
			y = 0
		}
		if y < 0 {
			dY = -1
			y = len(targetMap.TileMap[0]) - 1
		}

		targetMap = gm.GetMapChunkFrom(targetMap, dX, dY, 0)
		return targetMap, &targetMap.TileMap[x][y][z], x, y, z

	} else {
		// Get the tile at the requested position
		return targetMap, &targetMap.TileMap[x][y][z], x, y, z
	}
}

func (gm *GameManager) GetSurroundingTiles(tileMap *MapChunk, x, y, z int) (*Tile, *Tile, *Tile, *Tile, *Tile, *Tile, *Tile, *Tile) {
	log.Println("Checking surrounding tiles from, ", x, y, z)

	n, ne, e, se, s, sw, w, nw := &Tile{}, &Tile{}, &Tile{}, &Tile{}, &Tile{}, &Tile{}, &Tile{}, &Tile{}

	// Cover the top row first, gathering a new map chunk if necessary
	n = gm.GetNorthTile(tileMap, x, y, z)
	e = gm.GetEastTile(tileMap, x, y, z)
	s = gm.GetSouthTile(tileMap, x, y, z)
	w = gm.GetWestTile(tileMap, x, y, z)
	return n, ne, e, se, s, sw, w, nw
}

// Returns the tile north of the given position
func (gm *GameManager) GetNorthTile(tileMap *MapChunk, x, y, z int) *Tile {
	if y == 0 {
		// Get the map chunk north of the current one
		mapChunk := gm.GetMapChunkFrom(tileMap, 0, -1, 0)
		return &mapChunk.TileMap[x][len(mapChunk.TileMap[0])-1][z]
	} else {
		return &tileMap.TileMap[x][y-1][z]
	}
}

// Returns the tile east of the given position
func (gm *GameManager) GetEastTile(tileMap *MapChunk, x, y, z int) *Tile {
	if x == len(tileMap.TileMap)-1 {
		// Get the map chunk east of the current one
		mapChunk := gm.GetMapChunkFrom(tileMap, 1, 0, 0)
		log.Println("East tile is", mapChunk.TileMap[0][y][z])
		return &mapChunk.TileMap[0][y][z]
	} else {
		return &tileMap.TileMap[x+1][y][z]
	}
}

// Returns the tile south of the given position
func (gm *GameManager) GetSouthTile(tileMap *MapChunk, x, y, z int) *Tile {
	if y == len(tileMap.TileMap[0])-1 {
		// Get the map chunk south of the current one
		mapChunk := gm.GetMapChunkFrom(tileMap, 0, 1, 0)
		return &mapChunk.TileMap[x][0][z]
	} else {
		return &tileMap.TileMap[x][y+1][z]
	}
}

// Returns the tile west of the given position
func (gm *GameManager) GetWestTile(tileMap *MapChunk, x, y, z int) *Tile {
	if x == 0 {
		// Get the map chunk west of the current one
		mapChunk := gm.GetMapChunkFrom(tileMap, -1, 0, 0)
		return &mapChunk.TileMap[len(mapChunk.TileMap)-1][y][z]
	} else {
		return &tileMap.TileMap[x-1][y][z]
	}
}
