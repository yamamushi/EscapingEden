package gamewindow

import (
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
			gw.HandleCommand(types.InputCharacter, input.Data)
		case types.InputEscape:
			gw.HandleCommand(types.InputEscape, "")
		}
	}
}

func (gw *GameWindow) HandleCommand(inputType types.InputType, input string) {
	gw.commandMutex.Lock()
	defer gw.commandMutex.Unlock()
	//gw.log.Println(logging.LogInfo, "GameWindow Command: ", input)
	gw.MenusMutex.Lock()
	if len(gw.Menus) > 0 {
		gw.Menus[len(gw.Menus)-1].HandleInput(gw, inputType, input) // Handle input for the top menu
		gw.MenusMutex.Unlock()
		return
	}
	gw.MenusMutex.Unlock()
	//gw.Log.Println(logging.LogInfo, "GameWindow Input: ", strconv.Itoa(int(input[0])))
	// convert input to an int and send the value to the console
	if inputType == types.InputEscape {
		return // Do nothing
	}
	if int(input[0]) == 4 {
		// ^D
		//gw.Log.Println(logging.LogInfo, "GameWindow received ^D, handling dig")
		gw.StatusBarMutex.Lock()
		gw.StatusBarMessage = "Dig in which direction?"
		gw.StatusBarMutex.Unlock()
		return
	} else if int(input[0]) == 2 {
		// ctrl-b
		//gw.Log.Println(logging.LogInfo, "GameWindow received ^B, handling build")
		gw.StatusBarMutex.Lock()
		gw.StatusBarMessage = "Build in which direction?"
		gw.StatusBarMutex.Unlock()
		return
	}
	gw.StatusBarMutex.Lock()
	gw.StatusBarMessage = ""
	gw.StatusBarMutex.Unlock()
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
	// Up and down movement
	case "<":
		gw.StatusBarMutex.Lock()
		gw.StatusBarMessage = "You cannot go up here."
		gw.StatusBarMutex.Unlock()
	case ">":
		gw.StatusBarMutex.Lock()
		gw.StatusBarMessage = "You cannot go down here."
		gw.StatusBarMutex.Unlock()
	// Other commands
	case "d":
		gw.StatusBarMutex.Lock()
		gw.StatusBarMessage = "What do you want to drop?"
		gw.StatusBarMutex.Unlock()
	case ",":
		gw.StatusBarMutex.Lock()
		gw.StatusBarMessage = "There is nothing here to pick up."
		gw.StatusBarMutex.Unlock()
	case "i":
		gw.RequestInventoryDisplay(nil, "")
		return
	case "t":
		if len(gw.Menus) > 0 {
			gw.RemoveMenuBox(gw.Menus[0])
			return
		} else {
			gw.CreateMenu(MenuType_Build)
		}

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
