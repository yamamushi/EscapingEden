package game

import (
	"log"

	"github.com/yamamushi/EscapingEden/logging"
)

func (gm *GameManager) GlobalTile(x, y, z int) (*Tile, *MapChunk) {
	chunkSize := gm.Config.WorldGen.ChunkSize
	worldDimensionsString := gm.Config.WorldGen.Dimensions
	worldX, worldY, worldZ := gm.ParseWorldDimensions(worldDimensionsString)

	// Wrap coordinates
	if x < 0 {
		x = (worldX * chunkSize) - 1
	}
	if x > (worldX*chunkSize)-1 {
		x = 0
	}
	if y < 0 {
		y = (worldY * chunkSize) - 1
	}
	if y > (worldY*chunkSize)-1 {
		y = 0
	}
	if z < 0 {
		z = (worldZ * chunkSize) - 1
	}
	if z > (worldZ*chunkSize)-1 {
		z = 0
	}

	// Calculate chunk coordinates
	globalX := x / chunkSize
	globalY := y / chunkSize
	globalZ := z / chunkSize

	// Load chunk lazily
	mapChunk, err := gm.GetMapChunk(globalX, globalY, globalZ)
	if err != nil {
		log.Println("Failed to get map chunk:", err)
		return nil, nil
	}

	// Access tile inside the chunk
	tile := &mapChunk.TileMap[x%chunkSize][y%chunkSize][z%chunkSize]
	return tile, mapChunk
}

func (gm *GameManager) GlobalToLocalTile(x, y, z int) (X, Y, Z int, chunk *MapChunk) {
	chunkSize := gm.Config.WorldGen.ChunkSize
	worldDimensionsString := gm.Config.WorldGen.Dimensions
	worldX, worldY, worldZ := gm.ParseWorldDimensions(worldDimensionsString)
	if x < 0 {
		x = (worldX * chunkSize) - 1
	}
	if x > (worldX*chunkSize)-1 {
		x = 0
	}
	if y < 0 {
		y = (worldY * chunkSize) - 1
	}
	if y > (worldY*chunkSize)-1 {
		y = 0
	}
	if z < 0 {
		z = (worldZ * chunkSize) - 1
	}
	if z > (worldZ*chunkSize)-1 {
		z = 0
	}

	// Deduce the map chunk 0,0,0 from the given coordinates
	globalX := x / chunkSize
	globalY := y / chunkSize
	globalZ := z / chunkSize

	// Get the map chunk at the given coordinates
	mapChunk, err := gm.MapChunkByPos(globalX, globalY, globalZ)
	if err != nil {
		log.Println("Failed to get map chunk", err.Error())
		return 0, 0, 0, nil
	}

	localX, localY, localZ := x%chunkSize, y%chunkSize, z%chunkSize
	if gm.Config.Logger.DebugChunk {
		gm.Log.Println(logging.LogDebug, "GlobalToLocal: global", x, y, z, "-> local", localX, localY, localZ, "chunk", globalX, globalY, globalZ)
	}
	return localX, localY, localZ, mapChunk
}

// Takes a local tile position and returns the global tile position with wrapping
func (gm *GameManager) LocalToGlobalTile(x, y, z int, mapChunk *MapChunk) (X, Y, Z int) {
	chunkSize := gm.Config.WorldGen.ChunkSize
	globalBaseX := mapChunk.GlobalPosition.X * chunkSize
	globalBaseY := mapChunk.GlobalPosition.Y * chunkSize
	globalBaseZ := mapChunk.GlobalPosition.Z * chunkSize

	// Debug logging (only if enabled)
	if gm.Config.Logger.DebugChunk {
		gm.Log.Println(logging.LogDebug, "LocalToGlobal: local", x, y, z, "chunk pos", mapChunk.GlobalPosition.X, mapChunk.GlobalPosition.Y, mapChunk.GlobalPosition.Z, "chunk size", chunkSize)
	}

	worldDimensionsString := gm.Config.WorldGen.Dimensions
	worldX, worldY, worldZ := gm.ParseWorldDimensions(worldDimensionsString)
	if x < 0 {
		x = (worldX * chunkSize) - 1
	}
	if x > (worldX*chunkSize)-1 {
		x = 0
	}
	if y < 0 {
		y = (worldY * chunkSize) - 1
	}
	if y > (worldY*chunkSize)-1 {
		y = 0
	}
	if z < 0 {
		z = (worldZ * chunkSize) - 1
	}
	if z > (worldZ*chunkSize)-1 {
		z = 0
	}

	X = globalBaseX + x
	Y = globalBaseY + y
	Z = globalBaseZ + z

	if gm.Config.Logger.DebugChunk {
		gm.Log.Println(logging.LogDebug, "LocalToGlobal result: global", X, Y, Z)
	}
	return X, Y, Z
}

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

// GetSurroundingTiles returns the tiles surrounding the given position and tileMap which represents the chunk the tile is in
func (gm *GameManager) GetSurroundingTiles(tileMap *MapChunk, x, y, z int) (*Tile, *Tile, *Tile, *Tile, *Tile, *Tile, *Tile, *Tile) {
	n, ne, e, se, s, sw, w, nw := &Tile{}, &Tile{}, &Tile{}, &Tile{}, &Tile{}, &Tile{}, &Tile{}, &Tile{}
	//	x, y, z = gm.LocalToGlobalTile(x, y, z, tileMap)

	// Cover the top row first, gathering a new map chunk if necessary
	n, _ = gm.GlobalTile(x, y-1, z)
	e, _ = gm.GlobalTile(x+1, y, z)
	s, _ = gm.GlobalTile(x, y+1, z)
	w, _ = gm.GlobalTile(x-1, y, z)
	return n, ne, e, se, s, sw, w, nw
}
