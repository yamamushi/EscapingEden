package gamewindow

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
)

// HandleInput handles input for the login window
func (gw *GameWindow) HandleInput(input types.Input) {
	switch gw.windowState {
	case GW_DefaultView:
		// Handle input for the main view
		switch input.Type {
		case types.InputReturn:
			// Send a console message to the ConsoleSend channel
			consoleMessage := messages.WindowMessage{Data: messages.GameManagerMessage{Data: messages.GameMessageData{CharacterID: gw.GetCharacterInfoField("id")}, Type: messages.GameManager_GetCharacterPosition}, Type: messages.WM_GameCommand}
			gw.SendToConsole(consoleMessage)
			return
		case types.InputCharacter:
			gw.HandleCommand(input.Data)
		}
	}
}

func (gw *GameWindow) HandleCommand(input string) {
	gw.commandMutex.Lock()
	defer gw.commandMutex.Unlock()
	gw.log.Println(logging.LogInfo, "GameWindow Command: ", input)
	switch input {
	// vi movement
	case "h":
		gw.MovePlayer(-1, 0)
	case "l":
		gw.MovePlayer(1, 0)
	case "j":
		gw.MovePlayer(0, 1)
	case "k":
		gw.MovePlayer(0, -1)
	// Diagonal movement
	case "y":
		gw.MovePlayer(-1, -1)
	case "u":
		gw.MovePlayer(1, -1)
	case "b":
		gw.MovePlayer(-1, 1)
	case "n":
		gw.MovePlayer(1, 1)
		// Other commands
	default:
		return // Do nothing
	}
}

func (gw *GameWindow) MovePlayer(deltax, deltay int) {
	//gw.log.Println(logging.LogInfo, "GameWindow MovePlayer: ", deltax, deltay)
	consoleMessage := messages.WindowMessage{Data: messages.GameManagerMessage{Type: messages.GameManager_MoveCharacter,
		Data: messages.GameMessageData{CharacterID: gw.GetCharacterInfoField("id"), Data: messages.GameCharMove{DeltaX: deltax, DeltaY: deltay}}}, Type: messages.WM_GameCommand}
	gw.SendToConsole(consoleMessage)
}
