package game

import (
	"github.com/yamamushi/EscapingEden/logging"
)

func (gm *GameManager) MovePlayer(charID string, deltax, deltay int) {
	//gm.Log.Println(logging.LogInfo, "Game Manager MovePlayer: ", deltax, deltay)

	character, err := gm.GetCharacter(charID)
	if err != nil {
		gm.Log.Println(logging.LogInfo, "error getting character", err.Error())
		return
	}
	// NEEDS MAP TRANSFER LOGIC
	currentMap := gm.GetMapChunkByID(character.CurrentMapID)
	if currentMap == nil {
		gm.Log.Println(logging.LogInfo, "No map loaded for character")
		return
	}

	mapDeltaX, mapDeltaY := 0, 0

	if character.Position.X+deltax < 0 {
		mapDeltaX = -1
	}
	if character.Position.X+deltax >= len(currentMap.TileMap) {
		mapDeltaX = 1
	}
	if character.Position.Y+deltay < 0 {
		mapDeltaY = -1
	}
	if character.Position.Y+deltay >= len(currentMap.TileMap[0]) {
		mapDeltaY = 1
	}
	if mapDeltaX != 0 || mapDeltaY != 0 {
		currentMap = gm.GetMapChunkFrom(currentMap, mapDeltaX, mapDeltaY, 0)
		err := gm.MovePlayerToMap(character.ID, currentMap.ID)
		if err != nil {
			gm.Log.Println(logging.LogError, "Failed to move player to new map", err.Error())
			return
		}
	}

	newPosX := character.Position.X + deltax
	newPosY := character.Position.Y + deltay

	if mapDeltaY == -1 {
		newPosY = len(currentMap.TileMap[0]) - 1
	}
	if mapDeltaY == 1 {
		newPosY = 0
	}
	if mapDeltaX == -1 {
		newPosX = len(currentMap.TileMap) - 1
	}
	if mapDeltaX == 1 {
		newPosX = 0
	}

	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	if currentMap.TileMap[newPosX][newPosY][0].TileType == "floor" {
		character.Position.X = newPosX
		character.Position.Y = newPosY
	}
	gm.Log.Println(logging.LogInfo, "Game Manager Moved Player To: ", character.Position.X, character.Position.Y, currentMap.GlobalPosition.X, currentMap.GlobalPosition.Y)
}
