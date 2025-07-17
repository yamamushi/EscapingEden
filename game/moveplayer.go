package game

import (
	"fmt"
	"strings"

	"github.com/yamamushi/EscapingEden/logging"
)

func (gm *GameManager) MovePlayer(charID string, deltax, deltay int) {
	if gm.Config.Logger.DebugChunk {
		gm.Log.Println(logging.LogInfo, "Game Manager MovePlayer: ", deltax, deltay)
	}

	character, err := gm.GetCharacter(charID)
	if err != nil {
		gm.Log.Println(logging.LogInfo, "error getting character", err.Error())
		return
	}

	currentMap := gm.GetMapChunkByIDWithLoading(character.CurrentMapID, charID)
	if currentMap == nil {
		gm.Log.Println(logging.LogInfo, "No map loaded for character")
		return
	}

	// Log current position and chunk info (only if debug enabled)
	if gm.Config.Logger.DebugChunk {
		gm.Log.Println(logging.LogInfo, "Current position:", character.Position.X, character.Position.Y, "in chunk", character.CurrentMapID)
		gm.Log.Println(logging.LogInfo, "Current chunk global pos:", currentMap.GlobalPosition.X, currentMap.GlobalPosition.Y, currentMap.GlobalPosition.Z)
	}

	gX, gY, gZ := gm.LocalToGlobalTile(character.Position.X, character.Position.Y, 0, currentMap)

	if gm.Config.Logger.DebugChunk {
		gm.Log.Println(logging.LogInfo, "Current global coordinates:", gX, gY, gZ)
		gm.Log.Println(logging.LogInfo, "Target global coordinates:", gX+deltax, gY+deltay, gZ)
	}

	tile, mapChunk := gm.GlobalTile(gX+deltax, gY+deltay, gZ)
	if tile == nil {
		if gm.Config.Logger.DebugChunk {
			gm.Log.Println(logging.LogInfo, "Tile not found at global coords:", gX+deltax, gY+deltay, gZ)
		}
		return
	}
	if strings.Contains(tile.TileType, "wall") {
		if gm.Config.Logger.DebugChunk {
			gm.Log.Println(logging.LogInfo, "Tile is wall at:", gX+deltax, gY+deltay, gZ)
		}
		return
	}

	if gm.Config.Logger.DebugChunk {
		gm.Log.Println(logging.LogInfo, "Target chunk ID:", mapChunk.ID, "Global pos:", mapChunk.GlobalPosition.X, mapChunk.GlobalPosition.Y, mapChunk.GlobalPosition.Z)
	}

	// Use coordinate-based chunk ID for consistency
	coordinateBasedID := fmt.Sprintf("%d-%d-%d", mapChunk.GlobalPosition.X, mapChunk.GlobalPosition.Y, mapChunk.GlobalPosition.Z)
	err = gm.MovePlayerToMap(character.ID, coordinateBasedID)
	if err != nil {
		gm.Log.Println(logging.LogError, "Failed to move player to new map", err.Error())
		return
	}

	lX, lY, lZ, newMapChunk := gm.GlobalToLocalTile(gX+deltax, gY+deltay, gZ)

	if gm.Config.Logger.DebugChunk {
		gm.Log.Println(logging.LogInfo, "Target local coordinates:", lX, lY, lZ)
	}

	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()

	// Validate bounds before accessing the tile array
	if newMapChunk != nil && lX >= 0 && lY >= 0 && lZ >= 0 &&
		lX < len(newMapChunk.TileMap) &&
		lY < len(newMapChunk.TileMap[0]) &&
		lZ < len(newMapChunk.TileMap[0][0]) {

		if newMapChunk.TileMap[lX][lY][lZ].TileType == "floor" {
			// Use coordinate-based chunk ID instead of UUID
			coordinateBasedID := fmt.Sprintf("%d-%d-%d", mapChunk.GlobalPosition.X, mapChunk.GlobalPosition.Y, mapChunk.GlobalPosition.Z)

			// Update character position directly (we already hold the mutex)
			character.Position.X = lX
			character.Position.Y = lY
			character.Position.Z = lZ
			character.Position.MapChunkID = coordinateBasedID
			character.CurrentMapID = coordinateBasedID

			// Save character position to database
			err = gm.DB.UpdateRecord("Characters", character)
			if err != nil {
				gm.Log.Println(logging.LogError, "Failed to save character position to database:", err.Error())
				// Note: We don't return here because the movement was successful in memory
				// The character can continue playing, but their position won't persist if server restarts
			} else if gm.Config.Logger.DebugChunk {
				gm.Log.Println(logging.LogInfo, "Character position saved to database:", lX, lY, "in chunk", mapChunk.ID)
			}
		}
	}
}
