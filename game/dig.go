package game

import (
	"errors"
	"fmt"
	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/logging"
	"strings"
)

func (gm *GameManager) HandleDigRequest(itemID string, charID string, deltaX, deltaY int) error {
	// Get character
	character, err := gm.GetCharacter(charID)
	if err != nil {
		return errors.New("character not found")
	}

	// Validate digging tool
	item := edentypes.Item{}
	err = gm.DB.One("Items", "ID", itemID, &item)
	if err != nil {
		return errors.New("digging tool not found")
	}

	// Check if character has the tool in inventory
	hasTool := gm.CharacterHasItem(charID, itemID)
	if !hasTool {
		return errors.New("you don't have the required digging tool")
	}

	// Check if tool is appropriate for digging
	if item.Type != edentypes.ItemTool {
		return errors.New("item is not a tool")
	}

	// Get target tile information
	mapChunk, tile, x, y, z := gm.GetTileFromCharacter(charID, character.Position.X+deltaX, character.Position.Y+deltaY, 0)
	if tile == nil {
		return errors.New("target location not found")
	}

	// Check if target tile can be dug
	if tile.TileType == "floor" {
		return errors.New("tile is already floor")
	}

	// Check if tile is diggable (not certain wall types that shouldn't be dug)
	if !gm.isTileDiggable(tile.TileType) {
		return errors.New("this tile cannot be dug")
	}

	// Check if another player is standing at the target location
	gm.activeCharactersMutex.Lock()
	targetPlayer := gm.GetCharacterAt(mapChunk, character.Position.X+deltaX, character.Position.Y+deltaY)
	gm.activeCharactersMutex.Unlock()

	if targetPlayer != nil && targetPlayer.ID != charID {
		return errors.New("cannot dig on occupied tile")
	}

	// Convert tile to floor
	originalTileType := tile.TileType
	tile.TileType = "floor"

	// Fix wall alignment for surrounding tiles
	globalX, globalY, globalZ := gm.LocalToGlobalTile(x, y, z, mapChunk)
	gm.FixWallAlignment(globalX, globalY, globalZ, 1, false)

	// Consume tool durability
	err = gm.consumeToolDurability(charID, itemID, 1)
	if err != nil {
		gm.Log.Println(logging.LogWarn, fmt.Sprintf("Failed to consume tool durability: %v", err))
	}

	gm.Log.Println(logging.LogInfo, fmt.Sprintf("Character %s dug %s at (%d, %d)", charID, originalTileType, x, y))
	return nil
}

// isTileDiggable checks if a tile type can be dug
func (gm *GameManager) isTileDiggable(tileType string) bool {
	// Define which tile types can be dug
	diggableTiles := map[string]bool{
		"stone":  true,
		"dirt":   true,
		"sand":   true,
		"gravel": true,
		"clay":   true,
		"ore":    true,
		"coal":   true,
		"rock":   true,
	}

	// Check for wall types that can be dug
	if strings.Contains(tileType, "_wall") {
		// Most walls can be dug, but some special walls might not be
		specialWalls := map[string]bool{
			"bedrock_wall":    false,
			"reinforced_wall": false,
			"magic_wall":      false,
		}

		if !specialWalls[tileType] {
			return true // Most walls can be dug
		}
		return false
	}

	return diggableTiles[tileType]
}
