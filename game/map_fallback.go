package game

import (
	"fmt"

	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

// MapFallbackManager handles map loading failures and ensures players always see valid maps
type MapFallbackManager struct {
	gm               *GameManager
	safeSpawnPoints  []SafeSpawnPoint
	emergencyChunk   *MapChunk
	fallbackAttempts map[string]int // Track fallback attempts per character
}

// SafeSpawnPoint represents a known good spawn location
type SafeSpawnPoint struct {
	ChunkX, ChunkY, ChunkZ int
	LocalX, LocalY, LocalZ int
	Description            string
	Priority               int // Lower numbers = higher priority
}

// NewMapFallbackManager creates a new map fallback manager
func NewMapFallbackManager(gm *GameManager) *MapFallbackManager {
	mfm := &MapFallbackManager{
		gm:               gm,
		fallbackAttempts: make(map[string]int),
	}

	// Initialize safe spawn points
	mfm.initializeSafeSpawnPoints()

	// Create emergency chunk
	mfm.createEmergencyChunk()

	return mfm
}

// initializeSafeSpawnPoints sets up known safe locations
func (mfm *MapFallbackManager) initializeSafeSpawnPoints() {
	mfm.safeSpawnPoints = []SafeSpawnPoint{
		{
			ChunkX: 0, ChunkY: 0, ChunkZ: 0,
			LocalX: 127, LocalY: 127, LocalZ: 0, // Center of chunk
			Description: "Origin Center",
			Priority:    1,
		},
		{
			ChunkX: 0, ChunkY: 0, ChunkZ: 0,
			LocalX: 50, LocalY: 50, LocalZ: 0,
			Description: "Origin Safe Zone",
			Priority:    2,
		},
		{
			ChunkX: 1, ChunkY: 1, ChunkZ: 0,
			LocalX: 127, LocalY: 127, LocalZ: 0,
			Description: "Secondary Spawn",
			Priority:    3,
		},
	}
}

// createEmergencyChunk creates a minimal safe chunk for absolute emergencies
func (mfm *MapFallbackManager) createEmergencyChunk() {
	chunkID := "emergency-safe-chunk"

	// Create a simple, safe chunk
	mfm.emergencyChunk = &MapChunk{
		ID: chunkID,
		GlobalPosition: struct {
			X int
			Y int
			Z int
		}{X: -1, Y: -1, Z: -1}, // Special coordinates for emergency chunk
		TileMap: make([][][]Tile, 255),
	}

	// Initialize with safe floor tiles
	for x := 0; x < 255; x++ {
		mfm.emergencyChunk.TileMap[x] = make([][]Tile, 255)
		for y := 0; y < 255; y++ {
			mfm.emergencyChunk.TileMap[x][y] = make([]Tile, 3)
			for z := 0; z < 3; z++ {
				mfm.emergencyChunk.TileMap[x][y][z] = Tile{
					TileType: "floor",
				}
			}
		}
	}

	// Add some basic features to make it less boring
	mfm.addEmergencyChunkFeatures()
}

// addEmergencyChunkFeatures adds basic features to the emergency chunk
func (mfm *MapFallbackManager) addEmergencyChunkFeatures() {
	// Create a safe room in the center
	centerX, centerY := 127, 127

	// Add walls around the safe area
	for x := centerX - 10; x <= centerX+10; x++ {
		for y := centerY - 10; y <= centerY+10; y++ {
			if x == centerX-10 || x == centerX+10 || y == centerY-10 || y == centerY+10 {
				if x >= 0 && x < 255 && y >= 0 && y < 255 {
					mfm.emergencyChunk.TileMap[x][y][0] = Tile{TileType: "wall"}
				}
			}
		}
	}

	// Add an exit
	if centerX >= 0 && centerX < 255 && centerY+10 >= 0 && centerY+10 < 255 {
		mfm.emergencyChunk.TileMap[centerX][centerY+10][0] = Tile{TileType: "floor"}
	}
}

// HandleMapLoadFailure handles when a character's map fails to load
func (mfm *MapFallbackManager) HandleMapLoadFailure(characterID string, failedChunkID string, reason error) error {
	mfm.gm.Log.Println(logging.LogWarn, fmt.Sprintf("Map load failure for character %s, chunk %s: %v",
		characterID, failedChunkID, reason))

	// Track fallback attempts
	mfm.fallbackAttempts[characterID]++

	// Get character
	character, err := mfm.gm.GetCharacter(characterID)
	if err != nil {
		return fmt.Errorf("failed to get character for fallback: %w", err)
	}

	// Try safe spawn points in order of priority
	for _, spawnPoint := range mfm.safeSpawnPoints {
		if mfm.tryMoveToSafeSpawn(character, spawnPoint) {
			mfm.gm.Log.Println(logging.LogInfo, fmt.Sprintf("Successfully moved character %s to safe spawn: %s",
				characterID, spawnPoint.Description))

			// Notify player about the teleport
			mfm.notifyPlayerOfTeleport(characterID, spawnPoint.Description, reason)
			return nil
		}
	}

	// If all safe spawns fail, use emergency chunk
	return mfm.moveToEmergencyChunk(character)
}

// tryMoveToSafeSpawn attempts to move a character to a safe spawn point
func (mfm *MapFallbackManager) tryMoveToSafeSpawn(character *messages.CharacterInfo, spawn SafeSpawnPoint) bool {
	// Try to load the target chunk
	chunk, err := mfm.gm.GetMapChunk(spawn.ChunkX, spawn.ChunkY, spawn.ChunkZ)
	if err != nil {
		mfm.gm.Log.Println(logging.LogWarn, fmt.Sprintf("Failed to load safe spawn chunk %d,%d,%d: %v",
			spawn.ChunkX, spawn.ChunkY, spawn.ChunkZ, err))
		return false
	}

	// Validate the spawn location
	if !mfm.isLocationSafe(chunk, spawn.LocalX, spawn.LocalY, spawn.LocalZ) {
		mfm.gm.Log.Println(logging.LogWarn, fmt.Sprintf("Safe spawn location %s is not actually safe",
			spawn.Description))
		return false
	}

	// Move character to safe location
	mfm.gm.activeCharactersMutex.Lock()
	// Use coordinate-based chunk ID instead of UUID
	coordinateBasedID := fmt.Sprintf("%d-%d-%d", chunk.GlobalPosition.X, chunk.GlobalPosition.Y, chunk.GlobalPosition.Z)

	character.Position.X = spawn.LocalX
	character.Position.Y = spawn.LocalY
	character.Position.Z = spawn.LocalZ
	character.Position.MapChunkID = coordinateBasedID
	character.CurrentMapID = coordinateBasedID
	character.Initialized = true
	mfm.gm.activeCharactersMutex.Unlock()

	return true
}

// moveToEmergencyChunk moves character to the emergency safe chunk
func (mfm *MapFallbackManager) moveToEmergencyChunk(character *messages.CharacterInfo) error {
	mfm.gm.Log.Println(logging.LogWarn, fmt.Sprintf("Moving character %s to emergency chunk", character.ID))

	// Add emergency chunk to cache if not already there
	mfm.gm.ChunkCache["emergency-safe-chunk"] = mfm.emergencyChunk

	// Move character to center of emergency chunk
	mfm.gm.activeCharactersMutex.Lock()
	character.Position.X = 127
	character.Position.Y = 127
	character.Position.Z = 0
	character.Position.MapChunkID = "emergency-safe-chunk"
	character.CurrentMapID = "emergency-safe-chunk"
	character.Initialized = true
	mfm.gm.activeCharactersMutex.Unlock()

	// Notify player
	mfm.notifyPlayerOfEmergencyTeleport(character.ID)

	return nil
}

// isLocationSafe checks if a location is safe for a character
func (mfm *MapFallbackManager) isLocationSafe(chunk *MapChunk, x, y, z int) bool {
	// Check bounds
	if x < 0 || x >= len(chunk.TileMap) ||
		y < 0 || y >= len(chunk.TileMap[0]) ||
		z < 0 || z >= len(chunk.TileMap[0][0]) {
		return false
	}

	// Check if tile is walkable
	tile := &chunk.TileMap[x][y][z]
	tileInfo := mfm.gm.GetTileInfo(tile)

	// Floor tiles are generally safe
	return tile.TileType == "floor" || tileInfo.Character == "."
}

// notifyPlayerOfTeleport sends a message to the player about the teleport
func (mfm *MapFallbackManager) notifyPlayerOfTeleport(characterID, location string, reason error) {
	message := fmt.Sprintf("Map loading failed (%v). You have been teleported to a safe location: %s",
		reason, location)

	// Send message to player via connection manager
	response := messages.ConnectionManagerMessage{
		Type:               messages.ConnectManager_Message_GameCommandResponse,
		RecipientConsoleID: characterID, // Assuming character ID matches console ID
		Data: messages.GameMessage{
			Type: messages.GM_SystemMessage,
			Data: messages.GameMessageData{
				CharacterID: characterID,
				Data:        message,
			},
		},
	}

	mfm.gm.SendChannel <- response
}

// notifyPlayerOfEmergencyTeleport notifies player of emergency teleport
func (mfm *MapFallbackManager) notifyPlayerOfEmergencyTeleport(characterID string) {
	message := "Critical map loading failure. You have been moved to an emergency safe area. " +
		"Please contact an administrator if this problem persists."

	response := messages.ConnectionManagerMessage{
		Type:               messages.ConnectManager_Message_GameCommandResponse,
		RecipientConsoleID: characterID,
		Data: messages.GameMessage{
			Type: messages.GM_SystemMessage,
			Data: messages.GameMessageData{
				CharacterID: characterID,
				Data:        message,
			},
		},
	}

	mfm.gm.SendChannel <- response
}

// ValidateCharacterPosition ensures a character's position is valid and safe
func (mfm *MapFallbackManager) ValidateCharacterPosition(characterID string) error {
	character, err := mfm.gm.GetCharacter(characterID)
	if err != nil {
		return err
	}

	// Try to load the character's current chunk
	chunk, err := mfm.gm.GetMapChunk(0, 0, 0) // This would need to be derived from character position
	if err != nil {
		// If we can't load their current chunk, trigger fallback
		return mfm.HandleMapLoadFailure(characterID, character.CurrentMapID, err)
	}

	// Check if their position is safe
	if !mfm.isLocationSafe(chunk, character.Position.X, character.Position.Y, character.Position.Z) {
		mfm.gm.Log.Println(logging.LogWarn, fmt.Sprintf("Character %s is in unsafe location, moving to safety", characterID))
		return mfm.HandleMapLoadFailure(characterID, character.CurrentMapID,
			fmt.Errorf("character in unsafe location"))
	}

	return nil
}

// GetFallbackStats returns statistics about fallback usage
func (mfm *MapFallbackManager) GetFallbackStats() map[string]interface{} {
	totalFallbacks := 0
	for _, count := range mfm.fallbackAttempts {
		totalFallbacks += count
	}

	return map[string]interface{}{
		"total_fallbacks":       totalFallbacks,
		"characters_affected":   len(mfm.fallbackAttempts),
		"safe_spawn_points":     len(mfm.safeSpawnPoints),
		"emergency_chunk_ready": mfm.emergencyChunk != nil,
	}
}

// ResetFallbackAttempts clears the fallback attempt counter for a character
func (mfm *MapFallbackManager) ResetFallbackAttempts(characterID string) {
	delete(mfm.fallbackAttempts, characterID)
}

// AddSafeSpawnPoint adds a new safe spawn point
func (mfm *MapFallbackManager) AddSafeSpawnPoint(spawn SafeSpawnPoint) {
	mfm.safeSpawnPoints = append(mfm.safeSpawnPoints, spawn)

	// Sort by priority
	for i := 0; i < len(mfm.safeSpawnPoints)-1; i++ {
		for j := i + 1; j < len(mfm.safeSpawnPoints); j++ {
			if mfm.safeSpawnPoints[i].Priority > mfm.safeSpawnPoints[j].Priority {
				mfm.safeSpawnPoints[i], mfm.safeSpawnPoints[j] = mfm.safeSpawnPoints[j], mfm.safeSpawnPoints[i]
			}
		}
	}
}
