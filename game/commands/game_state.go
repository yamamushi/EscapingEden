package commands

import (
	"fmt"

	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

// GameStateProvider provides access to game state and implements GameStateProviderInterface
type GameStateProvider struct {
	log              logging.LoggerType
	sendChannel      chan messages.ConnectionManagerMessage
	getCharacterFunc func(string) (interface{}, error)
	getMapChunkFunc  func(string) interface{}
	getTileFunc      func(int, int, int, string) (interface{}, error)
}

// NewGameStateProvider creates a new game state provider
func NewGameStateProvider(
	log logging.LoggerType,
	sendChannel chan messages.ConnectionManagerMessage,
	getCharacterFunc func(string) (interface{}, error),
	getMapChunkFunc func(string) interface{},
	getTileFunc func(int, int, int, string) (interface{}, error),
) *GameStateProvider {
	return &GameStateProvider{
		log:              log,
		sendChannel:      sendChannel,
		getCharacterFunc: getCharacterFunc,
		getMapChunkFunc:  getMapChunkFunc,
		getTileFunc:      getTileFunc,
	}
}

// GetCharacter returns character information
func (p *GameStateProvider) GetCharacter(playerID string) (interface{}, error) {
	return p.getCharacterFunc(playerID)
}

// GetMapChunk returns map chunk information
func (p *GameStateProvider) GetMapChunk(chunkID string) (interface{}, error) {
	return p.getMapChunkFunc(chunkID), nil
}

// GetTile returns tile information at specific coordinates
func (p *GameStateProvider) GetTile(x, y, z int, chunkID string) (interface{}, error) {
	return p.getTileFunc(x, y, z, chunkID)
}

// UpdateGameState updates the game state based on command results
func (p *GameStateProvider) UpdateGameState(result *CommandResult) error {
	if result == nil || !result.Success {
		return nil
	}

	// Process side effects
	for _, effect := range result.SideEffects {
		switch effect.Type {
		case "tile_changed":
			p.log.Println(logging.LogInfo, fmt.Sprintf("Tile changed: %v", effect.Data))
			// Additional tile update logic could go here

		case "player_moved":
			p.log.Println(logging.LogInfo, fmt.Sprintf("Player moved: %v", effect.Data))
			// Additional movement update logic could go here

		case "inventory_changed":
			p.log.Println(logging.LogInfo, fmt.Sprintf("Inventory changed: %v", effect.Data))
			// Additional inventory update logic could go here

		case "notification":
			// Send notification to player
			if playerID, ok := effect.Data["player_id"].(string); ok {
				if message, ok := effect.Data["message"].(string); ok {
					p.sendNotification(playerID, message)
				}
			}
		}
	}

	return nil
}

// sendNotification sends a notification to a player
func (p *GameStateProvider) sendNotification(playerID, message string) {
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

	p.sendChannel <- response
}

// DefaultPermissionChecker provides basic permission checking
type DefaultPermissionChecker struct {
	log logging.LoggerType
}

// NewDefaultPermissionChecker creates a new permission checker
func NewDefaultPermissionChecker(log logging.LoggerType) *DefaultPermissionChecker {
	return &DefaultPermissionChecker{
		log: log,
	}
}

// HasPermission checks if a player has the required permission
func (p *DefaultPermissionChecker) HasPermission(playerID string, permission Permission) bool {
	// For now, grant all permissions to all players
	// In a real implementation, this would check against the role system
	return true
}

// GetPlayerPermissions returns all permissions for a player
func (p *DefaultPermissionChecker) GetPlayerPermissions(playerID string) []Permission {
	// For now, return all permissions
	// In a real implementation, this would check against the role system
	return []Permission{
		PermissionMove,
		PermissionDig,
		PermissionBuild,
		PermissionCraft,
		PermissionAttack,
	}
}
