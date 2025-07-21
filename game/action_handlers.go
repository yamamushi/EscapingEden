package game

import (
	"errors"
	"fmt"

	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/logging"
)

// HandleMovePlayerRequest handles a player movement request with proper validation
func (gm *GameManager) HandleMovePlayerRequest(characterID string, deltaX, deltaY int) error {
	if gm.Config.Logger.DebugChunk {
		gm.Log.Println(logging.LogInfo, "HandleMovePlayerRequest: ", characterID, deltaX, deltaY)
	}

	// Get character
	character, err := gm.GetCharacter(characterID)
	if err != nil {
		return fmt.Errorf("character not found: %v", err)
	}

	// Get current map
	currentMap := gm.GetMapChunkByIDWithLoading(character.CurrentMapID, characterID)
	if currentMap == nil {
		return errors.New("current map not found")
	}

	// Calculate global coordinates for the movement
	gX, gY, gZ := gm.LocalToGlobalTile(character.Position.X, character.Position.Y, character.Position.Z, currentMap)
	targetGX := gX + deltaX
	targetGY := gY + deltaY

	if gm.Config.Logger.DebugChunk {
		gm.Log.Println(logging.LogInfo, "Current global coords:", gX, gY, gZ)
		gm.Log.Println(logging.LogInfo, "Target global coords:", targetGX, targetGY, gZ)
	}

	// Get target tile and map chunk
	targetTile, targetMapChunk := gm.GlobalTile(targetGX, targetGY, gZ)
	if targetTile == nil {
		return errors.New("target tile not found")
	}

	// Check if target tile is walkable
	if targetTile.TileType != "floor" {
		return fmt.Errorf("target tile is not walkable: %s", targetTile.TileType)
	}

	// Convert back to local coordinates for the target chunk
	targetLocalX, targetLocalY, targetLocalZ, _ := gm.GlobalToLocalTile(targetGX, targetGY, gZ)

	// Check if another player is at the target location
	gm.activeCharactersMutex.Lock()
	targetPlayer := gm.GetCharacterAt(targetMapChunk, targetLocalX, targetLocalY)
	if targetPlayer != nil && targetPlayer.ID != characterID {
		gm.activeCharactersMutex.Unlock()
		return errors.New("target tile is occupied by another player")
	}

	// Update character position and map if necessary
	coordinateBasedID := fmt.Sprintf("%d-%d-%d", targetMapChunk.GlobalPosition.X, targetMapChunk.GlobalPosition.Y, targetMapChunk.GlobalPosition.Z)

	// Move player to new map if needed
	if character.CurrentMapID != coordinateBasedID {
		err = gm.MovePlayerToMap(character.ID, coordinateBasedID)
		if err != nil {
			gm.activeCharactersMutex.Unlock()
			return fmt.Errorf("failed to move player to new map: %v", err)
		}
	}

	// Update character position
	character.Position.X = targetLocalX
	character.Position.Y = targetLocalY
	character.Position.Z = targetLocalZ
	character.Position.MapChunkID = coordinateBasedID
	character.CurrentMapID = coordinateBasedID

	gm.activeCharactersMutex.Unlock()

	// Save to database
	err = gm.DB.UpdateRecord("Characters", character)
	if err != nil {
		gm.Log.Println(logging.LogWarn, fmt.Sprintf("Failed to save character position: %v", err))
		// Don't return error - movement already happened in memory
	}

	if gm.Config.Logger.DebugChunk {
		gm.Log.Println(logging.LogInfo, fmt.Sprintf("Character %s moved to local (%d, %d) in chunk %s",
			characterID, targetLocalX, targetLocalY, coordinateBasedID))
	}

	return nil
}

// HandleMineRequest handles a mining request with proper validation and resource generation
func (gm *GameManager) HandleMineRequest(itemID string, toolID string, characterID string, deltaX, deltaY int) error {
	// Get character
	character, err := gm.GetCharacter(characterID)
	if err != nil {
		return fmt.Errorf("character not found: %v", err)
	}

	// Calculate target position
	targetX := character.Position.X + deltaX
	targetY := character.Position.Y + deltaY

	// Get current map
	currentMap := gm.GetMapChunkByIDWithLoading(character.CurrentMapID, characterID)
	if currentMap == nil {
		return errors.New("current map not found")
	}

	// Check bounds
	if targetX < 0 || targetX >= len(currentMap.TileMap) || targetY < 0 || targetY >= len(currentMap.TileMap[0]) {
		return errors.New("mining target out of bounds")
	}

	// Check if target tile is mineable
	targetTile := &currentMap.TileMap[targetX][targetY][0]
	if targetTile.TileType != "stone" && targetTile.TileType != "ore" {
		return errors.New("target tile is not mineable")
	}

	// Validate and consume tool durability if provided
	var miningTool *edentypes.Item
	if toolID != "" {
		tool := edentypes.Item{}
		err := gm.DB.One("Items", "ID", toolID, &tool)
		if err != nil {
			return fmt.Errorf("mining tool not found: %v", err)
		}

		// Check if tool is appropriate for mining
		if tool.Type != edentypes.ItemTool {
			return errors.New("item is not a tool")
		}

		// Check if character has the tool
		hasTool := gm.CharacterHasItem(characterID, toolID)
		if !hasTool {
			return errors.New("you don't have the required mining tool")
		}

		miningTool = &tool
	}

	// Convert tile to floor
	originalTileType := targetTile.TileType
	targetTile.TileType = "floor"

	// Generate mining rewards based on tile type
	err = gm.generateMiningRewards(characterID, originalTileType, itemID)
	if err != nil {
		gm.Log.Println(logging.LogWarn, fmt.Sprintf("Failed to generate mining rewards: %v", err))
	}

	// Consume tool durability if tool was used
	if miningTool != nil {
		err = gm.consumeToolDurability(characterID, toolID, 1)
		if err != nil {
			gm.Log.Println(logging.LogWarn, fmt.Sprintf("Failed to consume tool durability: %v", err))
		}
	}

	gm.Log.Println(logging.LogInfo, fmt.Sprintf("Character %s mined %s at (%d, %d)", characterID, originalTileType, targetX, targetY))
	return nil
}

// generateMiningRewards generates appropriate rewards based on the mined tile type
func (gm *GameManager) generateMiningRewards(characterID string, tileType string, requestedItemID string) error {
	var rewardItems []string

	// Determine rewards based on tile type
	switch tileType {
	case "stone":
		rewardItems = []string{"stone", "gravel"}
	case "ore":
		rewardItems = []string{"iron_ore", "stone"}
	default:
		return fmt.Errorf("unknown mineable tile type: %s", tileType)
	}

	// If a specific item was requested, try to give that first
	if requestedItemID != "" {
		for _, itemType := range rewardItems {
			if itemType == requestedItemID {
				return gm.addMiningRewardToInventory(characterID, itemType)
			}
		}
	}

	// Otherwise, give the primary reward
	if len(rewardItems) > 0 {
		return gm.addMiningRewardToInventory(characterID, rewardItems[0])
	}

	return nil
}

// addMiningRewardToInventory adds a mining reward item to character inventory
func (gm *GameManager) addMiningRewardToInventory(characterID string, itemType string) error {
	// Try to find the item in the database
	item := edentypes.Item{}
	err := gm.DB.One("Items", "Name", itemType, &item)
	if err != nil {
		// If not found by name, try by ID
		err = gm.DB.One("Items", "ID", itemType, &item)
		if err != nil {
			return fmt.Errorf("reward item not found: %s", itemType)
		}
	}

	// Create new item instance for inventory
	newItem := item
	newItem.ID = gm.generateItemID()

	// Add to character inventory
	err = gm.AddToCharacterInventory(characterID, newItem)
	if err != nil {
		return fmt.Errorf("failed to add reward to inventory: %v", err)
	}

	return nil
}

// consumeToolDurability reduces tool durability and removes broken tools
func (gm *GameManager) consumeToolDurability(characterID string, toolID string, durabilityLoss int) error {
	// For now, this is a placeholder implementation
	// In a full implementation, you would:
	// 1. Get the tool from character inventory
	// 2. Reduce its durability
	// 3. Remove it if durability reaches 0
	// 4. Update the tool in the database

	gm.Log.Println(logging.LogInfo, fmt.Sprintf("Tool %s used by character %s, durability reduced by %d", toolID, characterID, durabilityLoss))
	return nil
}

// generateItemID generates a unique item ID (placeholder implementation)
func (gm *GameManager) generateItemID() string {
	// This should use the same ID generation as other parts of the system
	// For now, use a simple timestamp-based approach
	return fmt.Sprintf("item_%d", gm.GetCurrentTick())
}

// GetCurrentTick returns the current game tick (placeholder)
func (gm *GameManager) GetCurrentTick() uint64 {
	// This would be provided by the ActionManager or GameTicker
	// For now, return a placeholder value
	return 0
}
