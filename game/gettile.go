package game

import "github.com/yamamushi/EscapingEden/logging"

func (gm *GameManager) GetTileFromCharacter(charID string, x, y, z int) *Tile {
	character, err := gm.GetCharacter(charID)
	if err != nil {
		gm.Log.Println(logging.LogError, "Failed to get character", err.Error())
		return nil
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
		return &targetMap.TileMap[x][y][z]

	} else {
		// Get the tile at the requested position
		return &targetMap.TileMap[x][y][z]
	}

	return nil
}

func (gm *GameManager) GetSurroundingTiles(charID string, x, y, z int) (*Tile, *Tile, *Tile, *Tile, *Tile, *Tile, *Tile, *Tile) {
	character, err := gm.GetCharacter(charID)
	if err != nil {
		gm.Log.Println(logging.LogError, "Failed to get character", err.Error())
		return nil, nil, nil, nil, nil, nil, nil, nil
	}

	n, ne, e, se, s, sw, w, nw := &Tile{}, &Tile{}, &Tile{}, &Tile{}, &Tile{}, &Tile{}, &Tile{}, &Tile{}

	targetMap := gm.GetMapChunkByID(character.CurrentMapID)
	// North
	if y-1 < 0 {
		nMap := gm.GetMapChunkFrom(targetMap, 0, -1, 0)
		n = &nMap.TileMap[x][len(nMap.TileMap[0])-1][z]
	} else {
		n = &targetMap.TileMap[x][y-1][z]
	}

	// North East
	if y-1 < 0 {
		if x+1 > len(targetMap.TileMap)-1 {
			neMap := gm.GetMapChunkFrom(targetMap, 1, -1, 0)
			ne = &neMap.TileMap[0][len(neMap.TileMap[0])-1][z]
		} else {
			neMap := gm.GetMapChunkFrom(targetMap, 1, -1, 0)
			ne = &neMap.TileMap[x+1][len(neMap.TileMap[0])-1][z]
		}
	} else {
		if x+1 > len(targetMap.TileMap)-1 {
			neMap := gm.GetMapChunkFrom(targetMap, 1, 0, 0)
			ne = &neMap.TileMap[0][y-1][z]
		} else {
			ne = &targetMap.TileMap[x+1][y-1][z]
		}
	}

	// East
	if x+1 > len(targetMap.TileMap)-1 {
		eMap := gm.GetMapChunkFrom(targetMap, 1, 0, 0)
		e = &eMap.TileMap[0][y][z]
	} else {
		e = &targetMap.TileMap[x+1][y][z]
	}

	// South East
	if y+1 > len(targetMap.TileMap[0])-1 {
		if x+1 > len(targetMap.TileMap)-1 {
			seMap := gm.GetMapChunkFrom(targetMap, 1, 1, 0)
			se = &seMap.TileMap[0][0][z]
		} else {
			seMap := gm.GetMapChunkFrom(targetMap, 1, 1, 0)
			se = &seMap.TileMap[x+1][0][z]
		}
	} else {
		if x+1 > len(targetMap.TileMap)-1 {
			seMap := gm.GetMapChunkFrom(targetMap, 1, 0, 0)
			se = &seMap.TileMap[0][y+1][z]
		} else {
			se = &targetMap.TileMap[x+1][y+1][z]
		}
	}

	// South
	if y+1 > len(targetMap.TileMap[0])-1 {
		sMap := gm.GetMapChunkFrom(targetMap, 0, 1, 0)
		s = &sMap.TileMap[x][0][z]
	} else {
		s = &targetMap.TileMap[x][y+1][z]
	}

	// South West
	if y+1 > len(targetMap.TileMap[0])-1 {
		if x-1 < 0 {
			swMap := gm.GetMapChunkFrom(targetMap, -1, 1, 0)
			sw = &swMap.TileMap[len(swMap.TileMap)-1][0][z]
		} else {
			swMap := gm.GetMapChunkFrom(targetMap, -1, 1, 0)
			sw = &swMap.TileMap[x-1][0][z]
		}
	} else {
		if x-1 < 0 {
			swMap := gm.GetMapChunkFrom(targetMap, -1, 0, 0)
			sw = &swMap.TileMap[len(swMap.TileMap)-1][y+1][z]
		} else {
			sw = &targetMap.TileMap[x-1][y+1][z]
		}
	}

	// West
	if x-1 < 0 {
		wMap := gm.GetMapChunkFrom(targetMap, -1, 0, 0)
		w = &wMap.TileMap[len(wMap.TileMap)-1][y][z]
	} else {
		w = &targetMap.TileMap[x-1][y][z]
	}

	// North West
	if y-1 < 0 {
		if x-1 < 0 {
			nwMap := gm.GetMapChunkFrom(targetMap, -1, -1, 0)
			nw = &nwMap.TileMap[len(nwMap.TileMap)-1][len(nwMap.TileMap[0])-1][z]
		} else {
			nwMap := gm.GetMapChunkFrom(targetMap, -1, -1, 0)
			nw = &nwMap.TileMap[x-1][len(nwMap.TileMap[0])-1][z]
		}
	} else {
		if x-1 < 0 {
			nwMap := gm.GetMapChunkFrom(targetMap, -1, 0, 0)
			nw = &nwMap.TileMap[len(nwMap.TileMap)-1][y-1][z]
		} else {
			nw = &targetMap.TileMap[x-1][y-1][z]
		}
	}
	return n, ne, e, se, s, sw, w, nw
}
