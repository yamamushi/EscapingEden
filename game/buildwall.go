package game

import (
	"errors"
	"github.com/yamamushi/EscapingEden/edentypes"
	"strings"
)

func (gm *GameManager) HandleBuildWallRequest(itemID string, toolID string, charID string, deltaX, deltaY int) error {
	// Get character
	character, err := gm.GetCharacter(charID)
	if err != nil {
		return errors.New("character not found")
	}

	// Validate building material
	item := edentypes.Item{}
	err = gm.DB.One("Items", "ID", itemID, &item)
	if err != nil {
		return errors.New("building material not found")
	}

	// Check if character has the item in inventory
	hasItem := gm.CharacterHasItem(charID, itemID)
	if !hasItem {
		return errors.New("you don't have the required building material")
	}

	// Validate tool if provided
	if toolID != "" {
		tool := edentypes.Item{}
		err = gm.DB.One("Items", "ID", toolID, &tool)
		if err != nil {
			return errors.New("building tool not found")
		}

		// Check if tool is appropriate for building
		if tool.Type != edentypes.ItemTool {
			return errors.New("item is not a tool")
		}
	}

	// Get target tile information
	mapChunk, tile, x, y, z := gm.GetTileFromCharacter(charID, character.Position.X+deltaX, character.Position.Y+deltaY, 0)
	if tile == nil {
		return errors.New("target location not found")
	}

	// Check if target tile is buildable
	if strings.Contains(tile.TileType, "wall") {
		return errors.New("there is already a wall here")
	}

	// Check if another player is standing at the target location
	gm.activeCharactersMutex.Lock()
	targetPlayer := gm.GetCharacterAt(mapChunk, character.Position.X+deltaX, character.Position.Y+deltaY)
	gm.activeCharactersMutex.Unlock()

	if targetPlayer != nil && targetPlayer.ID != charID {
		return errors.New("cannot build on occupied tile")
	}

	// Consume building material from inventory
	err = gm.RemoveFromCharacterInventory(charID, itemID)
	if err != nil {
		return errors.New("failed to consume building material")
	}

	// Build the wall
	wallTypePrefix := strings.ToLower(item.Name) + "_wall"
	tile.TileType = wallTypePrefix

	// Fix wall alignment with surrounding walls
	globalX, globalY, globalZ := gm.LocalToGlobalTile(x, y, z, mapChunk)
	gm.FixWallAlignment(globalX, globalY, globalZ, 3, true)

	return nil
}

func (gm *GameManager) FixWallAlignment(x, y, z int, recurseCount int, checkSource bool) {
	if recurseCount <= 0 {
		return
	}
	recurseCount -= 1
	//gm.Log.Println(logging.LogInfo, "Fixing wall alignment")
	thisTile, mapChunk := gm.GlobalTile(x, y, z)
	//log.Println("Fixing wall alignment for tile", x, y, z, thisTile.TileType)
	tileVars := strings.Split(thisTile.TileType, "_")
	if len(tileVars) < 2 {
		//log.Println("Tile is not a wall type")
		if checkSource {
			return
		}
	} else {
		if tileVars[1] != "wall" {
			//log.Println("Tile is not a wall type, invalid length")
			if checkSource {
				return
			}
		} else {
			if checkSource {
				thisTile.TileType = tileVars[0]
			}
		}
	}

	n, _, e, _, s, _, w, _ := gm.GetSurroundingTiles(mapChunk, x, y, z)
	pN, pE, pS, pW := false, false, false, false
	// Now we're going to walk through the surrounding tiles and check if they are walls
	// North
	if strings.Contains(n.TileType, "_wall") {
		pN = true
		//log.Println("North is wall")
	}
	// East
	if strings.Contains(e.TileType, "_wall") {
		pE = true
		//log.Println("East is wall")
	}
	// South
	if strings.Contains(s.TileType, "_wall") {
		pS = true
		//log.Println("South is wall")
	}
	// West
	if strings.Contains(w.TileType, "_wall") {
		pW = true
		//log.Println("West is wall")
	}

	if pN && !pE && !pS && !pW {
		if checkSource {
			thisTile.TileType += "_wall_vertical_north"
		}
		gm.FixWallAlignment(x, y-1, z, recurseCount, true)
		return
	}
	if !pN && pE && !pS && !pW {
		if checkSource {
			thisTile.TileType += "_wall_horizontal_east"
		}
		gm.FixWallAlignment(x+1, y, z, recurseCount, true)
		return
	}
	if !pN && !pE && pS && !pW {
		if checkSource {
			thisTile.TileType += "_wall_vertical_south"
		}
		gm.FixWallAlignment(x, y+1, z, recurseCount, true)
		return
	}
	if pN && !pE && pS && !pW {
		if checkSource {
			thisTile.TileType += "_wall_vertical"
		}
		gm.FixWallAlignment(x, y-1, z, recurseCount, true) // North
		gm.FixWallAlignment(x, y+1, z, recurseCount, true) // South
		return
	}
	if !pN && !pE && !pS && pW {
		if checkSource {
			thisTile.TileType += "_wall_horizontal_west"
		}
		gm.FixWallAlignment(x-1, y, z, recurseCount, true) // West
		return
	}
	if !pN && pE && !pS && pW {
		if checkSource {
			thisTile.TileType += "_wall_horizontal"
		}
		gm.FixWallAlignment(x+1, y, z, recurseCount, true) // East
		gm.FixWallAlignment(x-1, y, z, recurseCount, true) // West
		return
	}
	if pN && pE && !pS && !pW {
		if checkSource {
			thisTile.TileType += "_wall_bottom_left"
		}
		gm.FixWallAlignment(x+1, y, z, recurseCount, true) // East
		gm.FixWallAlignment(x, y-1, z, recurseCount, true) // North
		return
	}
	if pN && !pE && !pS && pW {
		if checkSource {
			thisTile.TileType += "_wall_bottom_right"
		}
		gm.FixWallAlignment(x, y-1, z, recurseCount, true) // North
		gm.FixWallAlignment(x-1, y, z, recurseCount, true) // West
		return
	}
	if !pN && pE && pS && !pW {
		if checkSource {
			thisTile.TileType += "_wall_top_left"
		}
		gm.FixWallAlignment(x+1, y, z, recurseCount, true) // East
		gm.FixWallAlignment(x, y+1, z, recurseCount, true) // South
		return
	}
	if !pN && !pE && pS && pW {
		if checkSource {
			thisTile.TileType += "_wall_top_right"
		}
		gm.FixWallAlignment(x, y+1, z, recurseCount, true) // South
		gm.FixWallAlignment(x-1, y, z, recurseCount, true) // West
		return
	}
	if pN && pE && pS && !pW {
		if checkSource {
			thisTile.TileType += "_wall_north_east_south"
		}
		gm.FixWallAlignment(x, y-1, z, recurseCount, true) // North
		gm.FixWallAlignment(x+1, y, z, recurseCount, true) // East
		gm.FixWallAlignment(x, y+1, z, recurseCount, true) // South
		return
	}
	if pN && !pE && pS && pW {
		if checkSource {
			thisTile.TileType += "_wall_north_west_south"
		}
		gm.FixWallAlignment(x, y-1, z, recurseCount, true) // North
		gm.FixWallAlignment(x-1, y, z, recurseCount, true) // West
		gm.FixWallAlignment(x, y+1, z, recurseCount, true) // South
		return
	}
	if !pN && pE && pS && pW {
		if checkSource {
			thisTile.TileType += "_wall_south_east_west"
		}
		gm.FixWallAlignment(x, y+1, z, recurseCount, true) // South
		gm.FixWallAlignment(x+1, y, z, recurseCount, true) // East
		gm.FixWallAlignment(x-1, y, z, recurseCount, true) // West
		return
	}
	if pN && pE && !pS && pW {
		if checkSource {
			thisTile.TileType += "_wall_north_east_west"
		}
		gm.FixWallAlignment(x, y-1, z, recurseCount, true) // North
		gm.FixWallAlignment(x+1, y, z, recurseCount, true) // East
		gm.FixWallAlignment(x-1, y, z, recurseCount, true) // West
		return
	}
	if pN && pE && pS && pW {
		if checkSource {
			thisTile.TileType += "_wall_north_east_south_west"
		}
		gm.FixWallAlignment(x, y-1, z, recurseCount, true) // North
		gm.FixWallAlignment(x+1, y, z, recurseCount, true) // East
		gm.FixWallAlignment(x, y+1, z, recurseCount, true) // South
		gm.FixWallAlignment(x-1, y, z, recurseCount, true) // West
		return
	}

	// If we get here, then we're surrounded by floor tiles
	if checkSource {
		thisTile.TileType += "_wall"
	}
}
