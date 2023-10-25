package game

import (
	"errors"
	"github.com/yamamushi/EscapingEden/edenitems"
)

func (gm *GameManager) HandleDigRequest(itemID string, charID string, deltaX, deltaY int) error {
	item := edenitems.Item{}
	err := gm.DB.One("Items", "ID", itemID, &item)
	if err != nil {
		//gm.Log.Println("Failed to get item", err.Error())
		return errors.New("failed to get item")
	}

	//gm.Log.Println("Digging with", item.Name, "in", deltaX, deltaY, "direction")
	// Lock the map
	gm.MapChunks[0].Mutex.Lock()
	defer gm.MapChunks[0].Mutex.Unlock()
	// Get the character's position
	character, err := gm.GetCharacter(charID)
	if err != nil {
		//gm.Log.Println("Failed to get character", err.Error())
		return errors.New("failed to get character")
	}
	// Get the tile at the delta position
	tile := gm.MapChunks[0].TileMap[character.Position.X+deltaX][character.Position.Y+deltaY][0]
	// Check if the tile is diggable/passable
	if tile.Passable {
		return errors.New("tile is passable")
	}
	// Change the tile to passable
	gm.MapChunks[0].TileMap[character.Position.X+deltaX][character.Position.Y+deltaY][0].Passable = true
	// Log
	//gm.Log.Println(logging.LogInfo, "Tile at", character.Position.X+deltaX, character.Position.Y+deltaY, "is now passable")
	// return nil
	return nil
}
