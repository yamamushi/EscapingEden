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
	//gm.MapChunks[0].Mutex.Lock()
	//defer gm.MapChunks[0].Mutex.Unlock()
	// Get the character's position
	character, err := gm.GetCharacter(charID)
	if err != nil {
		//gm.Log.Println("Failed to get character", err.Error())
		return errors.New("failed to get character")
	}

	targetMap := gm.GetMapChunkByID(character.CurrentMapID)

	if character.Position.X+deltaX > len(targetMap.TileMap)-1 || character.Position.X+deltaX < 0 || character.Position.Y+deltaY > len(targetMap.TileMap[0])-1 || character.Position.Y+deltaY < 0 {
		dX, dY := 0, 0
		x, y, z := character.Position.X+deltaX, character.Position.Y+deltaY, 0
		if character.Position.X+deltaX > len(targetMap.TileMap)-1 {
			dX = 1
			x = 0
		}
		if character.Position.X+deltaX < 0 {
			dX = -1
			x = len(targetMap.TileMap) - 1
		}
		if character.Position.Y+deltaY > len(targetMap.TileMap[0])-1 {
			dY = 1
			y = 0
		}
		if character.Position.Y+deltaY < 0 {
			dY = -1
			y = len(targetMap.TileMap[0]) - 1
		}

		targetMap = gm.GetMapChunkFrom(targetMap, dX, dY, 0)
		tile := targetMap.TileMap[x][y][z]
		if tile.Passable {
			return errors.New("tile is passable")
		}
		targetMap.TileMap[x][y][z].Passable = true

	} else {
		// Get the tile at the delta position
		tile := targetMap.TileMap[character.Position.X+deltaX][character.Position.Y+deltaY][0]
		// Check if the tile is diggable/passable
		if tile.Passable {
			return errors.New("tile is passable")
		}
		// Change the tile to passable
		targetMap.TileMap[character.Position.X+deltaX][character.Position.Y+deltaY][0].Passable = true
	}

	// Log
	//gm.Log.Println(logging.LogInfo, "Tile at", character.Position.X+deltaX, character.Position.Y+deltaY, "is now passable")
	// return nil
	return nil
}
