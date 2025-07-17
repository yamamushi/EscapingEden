package gamewindow

import (
	"time"

	"github.com/yamamushi/EscapingEden/game/commands"
	"github.com/yamamushi/EscapingEden/messages"
)

// SendCommand sends a command to the server
func (gw *GameWindow) SendCommand(cmd *commands.PlayerCommand) {
	// Create a window message with the command
	message := messages.WindowMessage{
		Type: messages.WM_NewGameCommand,
		Data: cmd,
	}

	// Send to console
	gw.SendToConsole(message)
}

// SendDigCommand sends a dig command
func (gw *GameWindow) SendDigCommand(toolID, direction string) {
	// Create the command
	cmd := &commands.PlayerCommand{
		Type:     commands.CommandDig,
		PlayerID: gw.GetCharacterInfoField("id"),
		Parameters: map[string]interface{}{
			"tool_id":   toolID,
			"direction": direction,
		},
		Timestamp: time.Now(),
	}

	// Send the command
	gw.SendCommand(cmd)
}

// SendMoveCommand sends a move command
func (gw *GameWindow) SendMoveCommand(direction string) {
	// Create the command
	cmd := &commands.PlayerCommand{
		Type:     commands.CommandMove,
		PlayerID: gw.GetCharacterInfoField("id"),
		Parameters: map[string]interface{}{
			"direction": direction,
		},
		Timestamp: time.Now(),
	}

	// Send the command
	gw.SendCommand(cmd)
}

// SendBuildCommand sends a build command
func (gw *GameWindow) SendBuildCommand(itemID, toolID, direction string) {
	// Create the command
	cmd := &commands.PlayerCommand{
		Type:     commands.CommandBuild,
		PlayerID: gw.GetCharacterInfoField("id"),
		Parameters: map[string]interface{}{
			"item_id":   itemID,
			"tool_id":   toolID,
			"direction": direction,
		},
		Timestamp: time.Now(),
	}

	// Send the command
	gw.SendCommand(cmd)
}
