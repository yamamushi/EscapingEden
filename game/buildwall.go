package game

import (
	"errors"
	"github.com/yamamushi/EscapingEden/edentypes"
	"log"
)

func (gm *GameManager) HandleBuildWallRequest(itemID string, toolID string, charID string, deltaX, deltaY int) error {
	item := edentypes.Item{}
	err := gm.DB.One("Items", "ID", itemID, &item)
	if err != nil {
		//gm.Log.Println("Failed to get item", err.Error())
		return errors.New("failed to get item")
	}

	// We're not using the tool right now so nil works fine
	/*tool := edentypes.Item{}
	err = gm.DB.One("Items", "ID", toolID, &tool)
	if err != nil {
		//gm.Log.Println("Failed to get item", err.Error())
		return errors.New("failed to get tool")
	}
	*/

	// Lock the map
	//gm.MapChunks[0].Mutex.Lock()
	//defer gm.MapChunks[0].Mutex.Unlock()
	// Get the character's position
	character, err := gm.GetCharacter(charID)
	if err != nil {
		log.Println("Failed to get character", err.Error())
		return errors.New("character not found")
	}
	tile := gm.GetTileFromCharacter(charID, character.Position.X+deltaX, character.Position.Y+deltaY, 0)
	if tile == nil {
		log.Println("Tile not found")
		return errors.New("tile not found")
	}
	if !tile.Passable {
		log.Println("Tile is not passable")
		return errors.New("tile is not passable")
	}
	tile.Passable = false

	// Log
	//gm.Log.Println(logging.LogInfo, "Tile at", character.Position.X+deltaX, character.Position.Y+deltaY, "is now passable")
	// return nil
	return nil
}
