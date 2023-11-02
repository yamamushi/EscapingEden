package game

import (
	"errors"
	"github.com/yamamushi/EscapingEden/edentypes"
	"log"
)

func (gm *GameManager) HandleDigRequest(itemID string, charID string, deltaX, deltaY int) error {
	item := edentypes.Item{}
	err := gm.DB.One("Items", "ID", itemID, &item)
	if err != nil {
		//gm.Log.Println("Failed to get item", err.Error())
		return errors.New("failed to get item")
	}

	character, err := gm.GetCharacter(charID)
	if err != nil {
		log.Println("Failed to get character", err.Error())
		return errors.New("character not found")
	}
	_, tile, _, _, _ := gm.GetTileFromCharacter(charID, character.Position.X+deltaX, character.Position.Y+deltaY, 0)
	if tile == nil {
		log.Println("Tile not found")
		return errors.New("tile not found")
	}
	if tile.TileType == "floor" {
		log.Println("Tile is already floor")
		return errors.New("tile is already floor")
	}
	tile.TileType = "floor"
	gm.FixWallAlignment(charID, character.Position.X+deltaX, character.Position.Y+deltaY, 0, 3)

	return nil
}
