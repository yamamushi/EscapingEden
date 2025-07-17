package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

// AdminTeleportCommand handles the admin teleport command
// Format: /teleport <player> <location>
// Where location can be:
// - A safe spawn name (e.g., "origin")
// - "emergency" for emergency chunk
// - Coordinates in format "x,y,z" (global coordinates)
// - "here" to teleport to admin's location
func (gm *GameManager) AdminTeleportCommand(adminID, targetPlayerID, location string) error {
	// Verify admin permissions (this would be more robust in a real implementation)
	// For now, we'll just log the admin ID
	gm.Log.Println(logging.LogInfo, fmt.Sprintf("Admin %s is teleporting player %s to %s",
		adminID, targetPlayerID, location))

	// Get target player
	targetCharacter, err := gm.GetCharacter(targetPlayerID)
	if err != nil {
		return fmt.Errorf("target player not found: %w", err)
	}

	// Check if fallback manager is available
	if gm.FallbackManager == nil {
		return fmt.Errorf("fallback manager not available")
	}

	// Handle different location types
	switch strings.ToLower(location) {
	case "emergency":
		// Teleport to emergency chunk
		err := gm.FallbackManager.moveToEmergencyChunk(targetCharacter)
		if err != nil {
			return fmt.Errorf("failed to teleport to emergency chunk: %w", err)
		}

		// Notify the player
		gm.notifyPlayerOfAdminTeleport(targetPlayerID, adminID, "emergency safe area")
		return nil

	case "here":
		// Teleport to admin's location
		adminCharacter, err := gm.GetCharacter(adminID)
		if err != nil {
			return fmt.Errorf("admin character not found: %w", err)
		}

		// Copy admin's position to target
		gm.activeCharactersMutex.Lock()
		targetCharacter.Position.X = adminCharacter.Position.X
		targetCharacter.Position.Y = adminCharacter.Position.Y
		targetCharacter.Position.Z = adminCharacter.Position.Z
		targetCharacter.Position.MapChunkID = adminCharacter.CurrentMapID
		targetCharacter.CurrentMapID = adminCharacter.CurrentMapID
		gm.activeCharactersMutex.Unlock()

		// Notify the player
		gm.notifyPlayerOfAdminTeleport(targetPlayerID, adminID, "admin's location")
		return nil
	}

	// Check if it's a named safe spawn
	for _, spawn := range gm.FallbackManager.safeSpawnPoints {
		if strings.EqualFold(spawn.Description, location) {
			// Found matching spawn point
			if gm.FallbackManager.tryMoveToSafeSpawn(targetCharacter, spawn) {
				// Notify the player
				gm.notifyPlayerOfAdminTeleport(targetPlayerID, adminID, spawn.Description)
				return nil
			}
			return fmt.Errorf("failed to teleport to safe spawn %s", spawn.Description)
		}
	}

	// Check if it's coordinates
	if strings.Contains(location, ",") {
		parts := strings.Split(location, ",")
		if len(parts) >= 3 {
			// Parse coordinates
			globalX, errX := strconv.Atoi(strings.TrimSpace(parts[0]))
			globalY, errY := strconv.Atoi(strings.TrimSpace(parts[1]))
			globalZ, errZ := strconv.Atoi(strings.TrimSpace(parts[2]))

			if errX == nil && errY == nil && errZ == nil {
				// Get the map chunk
				chunk, err := gm.GetMapChunk(globalX, globalY, globalZ)
				if err != nil {
					return fmt.Errorf("failed to get map chunk at %d,%d,%d: %w", globalX, globalY, globalZ, err)
				}

				// Convert to local coordinates
				localX, localY, localZ, _ := gm.GlobalToLocalTile(globalX, globalY, globalZ)

				// Check if location is safe
				if !gm.FallbackManager.isLocationSafe(chunk, localX, localY, localZ) {
					return fmt.Errorf("destination coordinates are not safe")
				}

				// Teleport player
				gm.activeCharactersMutex.Lock()
				targetCharacter.Position.X = localX
				targetCharacter.Position.Y = localY
				targetCharacter.Position.Z = localZ
				targetCharacter.Position.MapChunkID = chunk.ID
				targetCharacter.CurrentMapID = chunk.ID
				gm.activeCharactersMutex.Unlock()

				// Notify the player
				gm.notifyPlayerOfAdminTeleport(targetPlayerID, adminID, fmt.Sprintf("coordinates %d,%d,%d", globalX, globalY, globalZ))
				return nil
			}
		}
	}

	return fmt.Errorf("invalid location format. Use: emergency, here, spawn name, or x,y,z coordinates")
}

// ListSafeSpawns returns a list of available safe spawn points for admins
func (gm *GameManager) ListSafeSpawns() []string {
	if gm.FallbackManager == nil {
		return []string{}
	}

	spawns := make([]string, len(gm.FallbackManager.safeSpawnPoints))
	for i, spawn := range gm.FallbackManager.safeSpawnPoints {
		spawns[i] = fmt.Sprintf("%s (Priority %d)", spawn.Description, spawn.Priority)
	}

	return spawns
}

// AdminTeleportSelf allows an admin to teleport themselves to a safe location
func (gm *GameManager) AdminTeleportSelf(adminID, location string) error {
	return gm.AdminTeleportCommand(adminID, adminID, location)
}

// notifyPlayerOfAdminTeleport sends a message to the player about the admin teleport
func (gm *GameManager) notifyPlayerOfAdminTeleport(playerID, adminID, destination string) {
	message := fmt.Sprintf("You have been teleported to %s by admin %s.", destination, adminID)

	// Send message to player
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

	// Log the teleport
	gm.Log.Println(logging.LogInfo, fmt.Sprintf("Admin %s teleported player %s to %s",
		adminID, playerID, destination))
}

// AdminGetPlayerLocation returns a player's current location for admin use
func (gm *GameManager) AdminGetPlayerLocation(playerID string) (string, error) {
	character, err := gm.GetCharacter(playerID)
	if err != nil {
		return "", fmt.Errorf("player not found: %w", err)
	}

	// Get global coordinates
	globalX, globalY, globalZ := gm.LocalToGlobalTile(
		character.Position.X,
		character.Position.Y,
		character.Position.Z,
		gm.GetMapChunkByID(character.CurrentMapID))

	return fmt.Sprintf("Player %s is at local (%d,%d,%d) global (%d,%d,%d) in chunk %s",
		playerID,
		character.Position.X, character.Position.Y, character.Position.Z,
		globalX, globalY, globalZ,
		character.CurrentMapID), nil
}
