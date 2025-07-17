package game

import (
	"fmt"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

// CharacterSafetyCheck performs safety checks on character positions
func (gm *GameManager) CharacterSafetyCheck(characterID string) error {
	if gm.FallbackManager == nil {
		return nil // No fallback manager, skip safety check
	}

	character, err := gm.GetCharacter(characterID)
	if err != nil {
		return err
	}

	// Check if character's map exists
	currentMap := gm.GetMapChunkByID(character.CurrentMapID)
	if currentMap == nil {
		gm.Log.Println(logging.LogWarn, fmt.Sprintf("Safety check: Character %s has invalid map ID %s",
			characterID, character.CurrentMapID))
		return gm.FallbackManager.HandleMapLoadFailure(characterID, character.CurrentMapID,
			fmt.Errorf("invalid map ID during safety check"))
	}

	// Check if character's position is valid
	if character.Position.X < 0 || character.Position.X >= len(currentMap.TileMap) ||
		character.Position.Y < 0 || character.Position.Y >= len(currentMap.TileMap[0]) ||
		character.Position.Z < 0 || character.Position.Z >= len(currentMap.TileMap[0][0]) {
		gm.Log.Println(logging.LogWarn, fmt.Sprintf("Safety check: Character %s has invalid position %d,%d,%d",
			characterID, character.Position.X, character.Position.Y, character.Position.Z))
		return gm.FallbackManager.HandleMapLoadFailure(characterID, character.CurrentMapID,
			fmt.Errorf("invalid position during safety check"))
	}

	// Check if character's tile is valid
	tile := &currentMap.TileMap[character.Position.X][character.Position.Y][character.Position.Z]
	if tile.TileType == "" {
		gm.Log.Println(logging.LogWarn, fmt.Sprintf("Safety check: Character %s is on invalid tile type",
			characterID))
		return gm.FallbackManager.HandleMapLoadFailure(characterID, character.CurrentMapID,
			fmt.Errorf("invalid tile type during safety check"))
	}

	// Check if character's tile is walkable
	tileInfo := gm.GetTileInfo(tile)
	if tileInfo.Character != "." && tile.TileType != "floor" {
		gm.Log.Println(logging.LogWarn, fmt.Sprintf("Safety check: Character %s is on unwalkable tile %s",
			characterID, tile.TileType))
		return gm.FallbackManager.HandleMapLoadFailure(characterID, character.CurrentMapID,
			fmt.Errorf("unwalkable tile during safety check"))
	}

	return nil
}

// ValidateAllCharacterPositions checks all active characters for safety
func (gm *GameManager) ValidateAllCharacterPositions() {
	if gm.FallbackManager == nil {
		return // No fallback manager, skip validation
	}

	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()

	for _, character := range gm.ActiveCharacters {
		err := gm.CharacterSafetyCheck(character.Record.ID)
		if err != nil {
			gm.Log.Println(logging.LogWarn, fmt.Sprintf("Failed to validate character %s: %v",
				character.Record.ID, err))
		}
	}
}

// ScheduleCharacterSafetyChecks sets up periodic safety checks
func (gm *GameManager) ScheduleCharacterSafetyChecks() {
	if gm.FallbackManager == nil {
		return // No fallback manager, skip scheduling
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			gm.ValidateAllCharacterPositions()
		}
	}()

	gm.Log.Println(logging.LogInfo, "Scheduled character safety checks")
}

// EnsureSafeLogin ensures a character is in a safe location when logging in
func (gm *GameManager) EnsureSafeLogin(characterID string) error {
	if gm.FallbackManager == nil {
		return nil // No fallback manager, skip safety check
	}

	character, err := gm.GetCharacter(characterID)
	if err != nil {
		return err
	}

	// Check if character needs initialization (first login or invalid data)
	needsInitialization := !character.Initialized ||
		character.CurrentMapID == "" ||
		character.Position.X == 0 && character.Position.Y == 0 && character.Position.Z == 0

	if needsInitialization {
		gm.Log.Println(logging.LogInfo, fmt.Sprintf("Character %s needs initialization, placing at safe spawn",
			characterID))

		// Place at first safe spawn point
		if len(gm.FallbackManager.safeSpawnPoints) > 0 {
			spawn := gm.FallbackManager.safeSpawnPoints[0]
			if gm.FallbackManager.tryMoveToSafeSpawn(character, spawn) {
				gm.Log.Println(logging.LogInfo, fmt.Sprintf("Placed character %s at %s",
					characterID, spawn.Description))

				// Mark character as initialized
				character.Initialized = true

				// Update the character in the database
				err = gm.DB.UpdateRecord("Characters", character)
				if err != nil {
					gm.Log.Println(logging.LogWarn, fmt.Sprintf("Failed to update character %s in database: %v",
						characterID, err))
				}

				return nil
			}
		}

		// If safe spawn fails, use emergency chunk
		gm.Log.Println(logging.LogWarn, fmt.Sprintf("Safe spawn failed for character %s, using emergency chunk",
			characterID))
		return gm.FallbackManager.moveToEmergencyChunk(character)
	}

	// For existing characters, perform safety check
	// But first, let's verify their current map exists
	currentMap := gm.GetMapChunkByID(character.CurrentMapID)
	if currentMap == nil {
		gm.Log.Println(logging.LogWarn, fmt.Sprintf("Character %s has invalid map ID %s, resetting to safe spawn",
			characterID, character.CurrentMapID))

		// Reset to safe spawn
		if len(gm.FallbackManager.safeSpawnPoints) > 0 {
			spawn := gm.FallbackManager.safeSpawnPoints[0]
			if gm.FallbackManager.tryMoveToSafeSpawn(character, spawn) {
				gm.Log.Println(logging.LogInfo, fmt.Sprintf("Reset character %s to %s due to invalid map",
					characterID, spawn.Description))

				// Update the character in the database
				err = gm.DB.UpdateRecord("Characters", character)
				if err != nil {
					gm.Log.Println(logging.LogWarn, fmt.Sprintf("Failed to update character %s in database: %v",
						characterID, err))
				}

				// Notify the player
				gm.notifyPlayerOfMapReset(characterID, "invalid map data")
				return nil
			}
		}

		// If safe spawn fails, use emergency chunk
		return gm.FallbackManager.moveToEmergencyChunk(character)
	}

	// Perform normal safety check for characters with valid maps
	return gm.CharacterSafetyCheck(characterID)
}

// notifyPlayerOfMapReset notifies a player that their position was reset
func (gm *GameManager) notifyPlayerOfMapReset(playerID, reason string) {
	message := fmt.Sprintf("Your character position has been reset to a safe location due to %s. Welcome back!", reason)

	response := messages.ConnectionManagerMessage{
		Type:               messages.ConnectManager_Message_GameCommandResponse,
		RecipientConsoleID: playerID,
		Data: messages.GameMessage{
			Type: messages.GM_SystemMessage,
			Data: messages.GameMessageData{
				CharacterID: playerID,
				Data:        message,
			},
		},
	}

	gm.SendChannel <- response
	gm.Log.Println(logging.LogInfo, fmt.Sprintf("Notified player %s of position reset: %s", playerID, reason))
}

// ResetCharacterToDefaultSpawn resets a character to the default spawn position
// This is useful for fixing corrupted character data or admin commands
func (gm *GameManager) ResetCharacterToDefaultSpawn(characterID string) error {
	if gm.FallbackManager == nil {
		return fmt.Errorf("fallback manager not available")
	}

	character, err := gm.GetCharacter(characterID)
	if err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	// Try to place at first safe spawn point
	if len(gm.FallbackManager.safeSpawnPoints) > 0 {
		spawn := gm.FallbackManager.safeSpawnPoints[0]
		if gm.FallbackManager.tryMoveToSafeSpawn(character, spawn) {
			gm.Log.Println(logging.LogInfo, fmt.Sprintf("Reset character %s to default spawn %s",
				characterID, spawn.Description))

			// Mark character as initialized
			character.Initialized = true

			// Update the character in the database
			err = gm.DB.UpdateRecord("Characters", character)
			if err != nil {
				gm.Log.Println(logging.LogWarn, fmt.Sprintf("Failed to update character %s in database: %v",
					characterID, err))
			}

			// Notify the player
			gm.notifyPlayerOfMapReset(characterID, "admin reset")
			return nil
		}
	}

	// If safe spawn fails, use emergency chunk
	err = gm.FallbackManager.moveToEmergencyChunk(character)
	if err != nil {
		return fmt.Errorf("failed to move to emergency chunk: %w", err)
	}

	// Notify the player
	gm.notifyPlayerOfMapReset(characterID, "admin reset to emergency area")
	return nil
}
