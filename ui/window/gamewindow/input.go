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
			consoleMessage := messages.WindowMessage{Data: messages.GameMessage{Data: messages.GameMessageData{CharacterID: gw.GetCharacterInfoField("id")}}, Type: messages.WM_GameCommand}
			gw.ConsoleSend <- consoleMessage
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
}
