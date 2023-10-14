package game

import "github.com/yamamushi/EscapingEden/logging"

func (gm *GameManager) MovePlayer(charID string, deltax, deltay int) {
	//gm.Log.Println(logging.LogInfo, "Game Manager MovePlayer: ", deltax, deltay)

	character, err := gm.GetCharacter(charID)
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	if err != nil {
		gm.Log.Println(logging.LogInfo, "error getting character", err.Error())
		return
	}
	if (character.Position.X+deltax) < 0 || (character.Position.X+deltax) > len(gm.MapChunks[0].TileMap)-1 {
		return
	}
	if (character.Position.Y+deltay) < 0 || (character.Position.Y+deltay) > len(gm.MapChunks[0].TileMap[0]) {
		return
	}
	if gm.MapChunks[0].TileMap[character.Position.X+deltax][character.Position.Y+deltay][0].Passable {
		character.Position.X += deltax
		character.Position.Y += deltay
	}

	//gm.Log.Println(logging.LogInfo, "Game Manager Moved Player To: ", character.Position.X, character.Position.Y)
}
