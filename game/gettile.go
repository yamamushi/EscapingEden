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
