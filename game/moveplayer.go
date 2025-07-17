package game

import (
	"github.com/yamamushi/EscapingEden/logging"
	"strings"
)

func (gm *GameManager) MovePlayer(charID string, deltax, deltay int) {
	//gm.Log.Println(logging.LogInfo, "Game Manager MovePlayer: ", deltax, deltay)

	character, err := gm.GetCharacter(charID)
	if err != nil {
		gm.Log.Println(logging.LogInfo, "error getting character", err.Error())
		return
	}

	currentMap := gm.GetMapChunkByID(character.CurrentMapID)
	if currentMap == nil {
		gm.Log.Println(logging.LogInfo, "No map loaded for character")
		return
	}

	gX, gY, gZ := gm.LocalToGlobalTile(character.Position.X, character.Position.Y, 0, currentMap)
	tile, mapChunk := gm.GlobalTile(gX+deltax, gY+deltay, gZ)
	if tile == nil {
		gm.Log.Println(logging.LogInfo, "Tile not found")
		return
	}
	if strings.Contains(tile.TileType, "wall") {
		//gm.Log.Println(logging.LogInfo, "Tile is wall")
		return
	}

	err = gm.MovePlayerToMap(character.ID, mapChunk.ID)
	if err != nil {
		gm.Log.Println(logging.LogError, "Failed to move player to new map", err.Error())
		return
	} /* else {
		gm.Log.Println(logging.LogInfo, "Game Manager Moved Player To Map: ", currentMap.GlobalPosition.X, currentMap.GlobalPosition.Y)

	}*/
	lX, lY, lZ, _ := gm.GlobalToLocalTile(gX+deltax, gY+deltay, gZ)

	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	if currentMap.TileMap[lX][lY][lZ].TileType == "floor" {
		character.Position.X = lX
		character.Position.Y = lY
	}
}
