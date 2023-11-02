package game

import (
	"errors"
	"github.com/yamamushi/EscapingEden/edentypes"
	"log"
	"strings"
)

func (gm *GameManager) HandleBuildWallRequest(itemID string, toolID string, charID string, deltaX, deltaY int) error {
	item := edentypes.Item{}
	err := gm.DB.One("Items", "ID", itemID, &item)
	if err != nil {
		//gm.Log.Println("Failed to get item", err.Error())
		return errors.New("failed to get item")
	}

	wallTypePrefix := strings.ToLower(item.Name) + "_wall_"

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
	_, tile, _, _, _ := gm.GetTileFromCharacter(charID, character.Position.X+deltaX, character.Position.Y+deltaY, 0)
	if tile == nil {
		//log.Println("Tile not found")
		return errors.New("tile not found")
	}
	if strings.Contains(tile.TileType, "wall") {
		//log.Println("Tile is already wall")
		return errors.New("tile is already wall")
	}
	tile.TileType = wallTypePrefix
	// Now we want to check the surrounding to tiles to see if any are also walls
	gm.FixWallAlignment(charID, character.Position.X+deltaX, character.Position.Y+deltaY, 0, 3)
	return nil
}

func (gm *GameManager) FixWallAlignment(charID string, x, y, z int, recurseCount int) {
	if recurseCount <= 0 {
		return
	}
	recurseCount -= 1
	//gm.Log.Println(logging.LogInfo, "Fixing wall alignment")
	mapChunk, thisTile, x, y, z := gm.GetTileFromCharacter(charID, x, y, z)
	tileVars := strings.Split(thisTile.TileType, "_")
	if len(tileVars) < 2 {
		return
	}
	if strings.Split(thisTile.TileType, "_")[1] != "wall" {
		return
	}
	thisTile.TileType = strings.Split(thisTile.TileType, "_")[0]

	n, _, e, _, s, _, w, _ := gm.GetSurroundingTiles(mapChunk, x, y, z)
	pN, pE, pS, pW := false, false, false, false
	// Now we're going to walk through the surrounding tiles and check if they are walls
	// North
	if strings.Contains(n.TileType, "_wall") {
		pN = true
	}
	// East
	if strings.Contains(e.TileType, "_wall") {
		pE = true
	}
	// South
	if strings.Contains(s.TileType, "_wall") {
		pS = true
	}

	// West
	if strings.Contains(w.TileType, "_wall") {
		pW = true
	}

	if pN && !pE && !pS && !pW {
		thisTile.TileType += "_wall_vertical_north"
		gm.FixWallAlignment(charID, x, y-1, z, recurseCount)
		return
	}
	if !pN && pE && !pS && !pW {
		thisTile.TileType += "_wall_horizontal_east"
		gm.FixWallAlignment(charID, x+1, y, z, recurseCount)
		return
	}
	if !pN && !pE && pS && !pW {
		thisTile.TileType += "_wall_vertical_south"
		gm.FixWallAlignment(charID, x, y+1, z, recurseCount)
		return
	}
	if pN && !pE && pS && !pW {
		thisTile.TileType += "_wall_vertical"
		gm.FixWallAlignment(charID, x, y-1, z, recurseCount) // North
		gm.FixWallAlignment(charID, x, y+1, z, recurseCount) // South
		return
	}
	if !pN && !pE && !pS && pW {
		thisTile.TileType += "_wall_horizontal_west"
		gm.FixWallAlignment(charID, x-1, y, z, recurseCount) // West
		return
	}
	if !pN && pE && !pS && pW {
		thisTile.TileType += "_wall_horizontal"
		gm.FixWallAlignment(charID, x+1, y, z, recurseCount) // East
		gm.FixWallAlignment(charID, x-1, y, z, recurseCount) // West
		return
	}
	if pN && pE && !pS && !pW {
		thisTile.TileType += "_wall_bottom_left"
		gm.FixWallAlignment(charID, x+1, y, z, recurseCount) // East
		gm.FixWallAlignment(charID, x, y-1, z, recurseCount) // North
		return
	}
	if pN && !pE && !pS && pW {
		thisTile.TileType += "_wall_bottom_right"
		gm.FixWallAlignment(charID, x, y-1, z, recurseCount) // North
		gm.FixWallAlignment(charID, x-1, y, z, recurseCount) // West
		return
	}
	if !pN && pE && pS && !pW {
		thisTile.TileType += "_wall_top_left"
		gm.FixWallAlignment(charID, x+1, y, z, recurseCount) // East
		gm.FixWallAlignment(charID, x, y+1, z, recurseCount) // South
		return
	}
	if !pN && !pE && pS && pW {
		thisTile.TileType += "_wall_top_right"
		gm.FixWallAlignment(charID, x, y+1, z, recurseCount) // South
		gm.FixWallAlignment(charID, x-1, y, z, recurseCount) // West
		return
	}
	if pN && pE && pS && !pW {
		thisTile.TileType += "_wall_north_east_south"
		gm.FixWallAlignment(charID, x, y-1, z, recurseCount) // North
		gm.FixWallAlignment(charID, x+1, y, z, recurseCount) // East
		gm.FixWallAlignment(charID, x, y+1, z, recurseCount) // South
		return
	}
	if pN && !pE && pS && pW {
		thisTile.TileType += "_wall_north_west_south"
		gm.FixWallAlignment(charID, x, y-1, z, recurseCount) // North
		gm.FixWallAlignment(charID, x-1, y, z, recurseCount) // West
		gm.FixWallAlignment(charID, x, y+1, z, recurseCount) // South
		return
	}
	if !pN && pE && pS && pW {
		thisTile.TileType += "_wall_south_east_west"
		gm.FixWallAlignment(charID, x, y+1, z, recurseCount) // South
		gm.FixWallAlignment(charID, x+1, y, z, recurseCount) // East
		gm.FixWallAlignment(charID, x-1, y, z, recurseCount) // West
		return
	}
	if pN && pE && !pS && pW {
		thisTile.TileType += "_wall_north_east_west"
		gm.FixWallAlignment(charID, x, y-1, z, recurseCount) // North
		gm.FixWallAlignment(charID, x+1, y, z, recurseCount) // East
		gm.FixWallAlignment(charID, x-1, y, z, recurseCount) // West
		return
	}
	if pN && pE && pS && pW {
		thisTile.TileType += "_wall_north_east_south_west"
		gm.FixWallAlignment(charID, x, y-1, z, recurseCount) // North
		gm.FixWallAlignment(charID, x+1, y, z, recurseCount) // East
		gm.FixWallAlignment(charID, x, y+1, z, recurseCount) // South
		gm.FixWallAlignment(charID, x-1, y, z, recurseCount) // West
		return
	}

	// If we get here, then we're surrounded by floor tiles
	thisTile.TileType += "_wall"
}
