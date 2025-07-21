package gamewindow

import (
	"fmt"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
	"time"
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
		case types.InputUp, types.InputDown, types.InputLeft, types.InputRight:
			// Forward arrow keys to menus if any are open
			gw.HandleCommand(input.Type, input.Data)
		}
	}
}

func (gw *GameWindow) HandleCommand(inputType types.InputType, input string) {
	gw.commandMutex.Lock()
	defer gw.commandMutex.Unlock()
	//gw.log.Println(logging.LogInfo, "GameWindow Command: ", input)
	gw.MenusMutex.Lock()
	if len(input) > 0 && int(input[0]) == 3 {
		for _, menu := range gw.Menus {
			menu.SetCallbackStatusBarMessage("")
		}
		gw.CloseMenus = true
		gw.MenusMutex.Unlock()
		//gw.Log.Println(logging.LogInfo, "GameWindow received ^C, closing menus")
		gw.SetStatusBarMessage("")
		gw.RequestFlushFromConsole()
		return
	}

	// Handle ESC key for menus
	if inputType == types.InputEscape {
		if len(gw.Menus) > 0 {
			gw.Log.Println(logging.LogInfo, "GameWindow: ESC with menus open, forwarding to menu")
			gw.Menus[len(gw.Menus)-1].HandleInput(gw, inputType, "")
			gw.MenusMutex.Unlock()
			return
		}
		// No menus open, try to interrupt current action
		gw.Log.Println(logging.LogInfo, "GameWindow: ESC with no menus, attempting to interrupt current action")
		gw.InterruptCurrentAction()
		gw.MenusMutex.Unlock()
		return
	}

	if len(gw.Menus) > 0 {
		gw.Menus[len(gw.Menus)-1].HandleInput(gw, inputType, input) // Handle input for the top menu
		gw.MenusMutex.Unlock()
		return
	}
	gw.MenusMutex.Unlock()

	//gw.Log.Println(logging.LogInfo, "GameWindow Input: ", strconv.Itoa(int(input[0])))
	// convert input to an int and send the value to the console
	if int(input[0]) == 4 {
		// ^D
		//gw.Log.Println(logging.LogInfo, "GameWindow received ^D, handling dig")
		if len(gw.Menus) > 0 {
			gw.RemoveMenuBox(gw.Menus[0])
			return
		} else {
			gw.CreateMenu(MenuType_Dig)
		}
		return
	} else if int(input[0]) == 2 {
		// ctrl-b
		//gw.Log.Println(logging.LogInfo, "GameWindow received ^B, handling build")
		if len(gw.Menus) > 0 {
			gw.RemoveMenuBox(gw.Menus[0])
			return
		} else {
			gw.CreateMenu(MenuType_Build)
		}
		return
	}

	// Check if we're in build wall mode BEFORE clearing status bar
	if gw.buildWallMode {
		gw.Log.Println(logging.LogInfo, fmt.Sprintf("Build wall mode active, received input: '%s'", input))
		if gw.HandleBuildWallDirection(input) {
			return
		}
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
	case "W": // Wear
		gw.StatusBarMutex.Lock()
		gw.StatusBarMessage = "What do you want to wear?"
		gw.StatusBarMutex.Unlock()
	case "w": // Wield
		gw.StatusBarMutex.Lock()
		gw.StatusBarMessage = "What do you want to wield?"
		gw.StatusBarMutex.Unlock()
	case "i":
		gw.RequestInventoryUpdate(nil, "")
		gw.DisplayInventoryAfterReceive(true)
		return
	case "c":
		gw.DisplayCharacterEquipment()
		return
	default:
		return // Do nothing
	}
}

func (gw *GameWindow) MovePlayer(deltax, deltay int) {
	//gw.log.Println(logging.LogInfo, "GameWindow MovePlayer: ", deltax, deltay)

	// Convert delta to direction string for action queue
	direction := gw.deltaToDirection(deltax, deltay)

	// Provide immediate visual feedback
	directionName := gw.getDirectionDisplayName(direction)
	gw.SetStatusBarMessage(fmt.Sprintf("🚶 Queuing movement %s...", directionName))

	// Queue movement action instead of immediate execution
	actionData := map[string]interface{}{
		"direction": direction,
		"deltaX":    deltax,
		"deltaY":    deltay,
	}

	queueAction := messages.GameQueueAction{
		CharacterID: gw.GetCharacterInfoField("id"),
		ActionType:  "move",
		Data:        actionData,
	}

	consoleMessage := messages.WindowMessage{
		Type: messages.WM_GameCommand,
		Data: messages.GameManagerMessage{
			Type: messages.GameManager_QueueAction,
			Data: messages.GameMessageData{
				CharacterID: gw.GetCharacterInfoField("id"),
				Data:        queueAction,
			},
		},
	}

	gw.SendToConsole(consoleMessage)

	// Clear status message after a short delay
	go func() {
		time.Sleep(800 * time.Millisecond)
		expectedMsg := fmt.Sprintf("🚶 Queuing movement %s...", directionName)

		gw.StatusBarMutex.Lock()
		currentMsg := gw.StatusBarMessage
		gw.StatusBarMutex.Unlock()

		if currentMsg == expectedMsg {
			gw.SetStatusBarMessage("")
		}
	}()
}

// deltaToDirection converts movement deltas to direction strings
func (gw *GameWindow) deltaToDirection(deltaX, deltaY int) string {
	switch {
	case deltaX == 0 && deltaY == -1:
		return "north"
	case deltaX == 0 && deltaY == 1:
		return "south"
	case deltaX == 1 && deltaY == 0:
		return "east"
	case deltaX == -1 && deltaY == 0:
		return "west"
	case deltaX == 1 && deltaY == -1:
		return "northeast"
	case deltaX == -1 && deltaY == -1:
		return "northwest"
	case deltaX == 1 && deltaY == 1:
		return "southeast"
	case deltaX == -1 && deltaY == 1:
		return "southwest"
	default:
		return "unknown"
	}
}

// InterruptCurrentAction attempts to interrupt the current action for the player with enhanced feedback
func (gw *GameWindow) InterruptCurrentAction() {
	characterID := gw.GetCharacterInfoField("id")
	if characterID == "" {
		gw.Log.Println(logging.LogWarn, "Cannot interrupt action: no character ID")
		gw.SetStatusBarMessage("⚠️ Cannot cancel action - no character found")
		return
	}

	// Provide immediate visual feedback
	gw.SetStatusBarMessage("⏹ Attempting to cancel current action...")

	// Send interrupt action message to GameManager
	interruptMessage := messages.WindowMessage{
		Type: messages.WM_GameCommand,
		Data: messages.GameManagerMessage{
			Type: messages.GameManager_InterruptAction,
			Data: messages.GameMessageData{
				CharacterID: characterID,
			},
		},
	}

	gw.SendToConsole(interruptMessage)
	gw.Log.Println(logging.LogInfo, "Sent action interruption request for character:", characterID)

	// Clear the status message after a short delay to show the result
	go func() {
		time.Sleep(500 * time.Millisecond)
		// The actual result will be shown by the response handler
		// This just clears the "attempting" message if no response comes
		gw.StatusBarMutex.Lock()
		currentMsg := gw.StatusBarMessage
		gw.StatusBarMutex.Unlock()

		if currentMsg == "⏹ Attempting to cancel current action..." {
			gw.SetStatusBarMessage("")
		}
	}()
}

// getDirectionDisplayName converts direction strings to user-friendly display names
func (gw *GameWindow) getDirectionDisplayName(direction string) string {
	switch direction {
	case "north":
		return "north ↑"
	case "south":
		return "south ↓"
	case "east":
		return "east →"
	case "west":
		return "west ←"
	case "northeast":
		return "northeast ↗"
	case "northwest":
		return "northwest ↖"
	case "southeast":
		return "southeast ↘"
	case "southwest":
		return "southwest ↙"
	default:
		return direction
	}
}

// ShowActionQueueFeedback displays feedback when actions are queued or fail
func (gw *GameWindow) ShowActionQueueFeedback(actionType string, success bool, message string) {
	var statusMsg string
	var duration time.Duration = 2 * time.Second

	if success {
		switch actionType {
		case "move":
			statusMsg = fmt.Sprintf("✅ Movement queued: %s", message)
		case "dig":
			statusMsg = fmt.Sprintf("⛏️ Digging queued: %s", message)
		case "build_wall":
			statusMsg = fmt.Sprintf("🧱 Building queued: %s", message)
		case "mine":
			statusMsg = fmt.Sprintf("⚒️ Mining queued: %s", message)
		default:
			statusMsg = fmt.Sprintf("✅ Action queued: %s", message)
		}
	} else {
		switch actionType {
		case "queue_full":
			statusMsg = "⚠️ Queue is full! Cancel actions (ESC) or wait for completion"
			duration = 3 * time.Second
		case "invalid_action":
			statusMsg = fmt.Sprintf("❌ Invalid action: %s", message)
			duration = 3 * time.Second
		case "interrupt_failed":
			statusMsg = fmt.Sprintf("⚠️ Cannot cancel: %s", message)
			duration = 2500 * time.Millisecond
		default:
			statusMsg = fmt.Sprintf("❌ Action failed: %s", message)
			duration = 3 * time.Second
		}
	}

	gw.SetStatusBarMessage(statusMsg)

	// Clear message after duration
	go func() {
		time.Sleep(duration)

		gw.StatusBarMutex.Lock()
		currentMsg := gw.StatusBarMessage
		gw.StatusBarMutex.Unlock()

		if currentMsg == statusMsg {
			gw.SetStatusBarMessage("")
		}
	}()
}

// ShowQueueStatusUpdate displays current queue status information
func (gw *GameWindow) ShowQueueStatusUpdate(queueCount int, maxQueue int, currentAction string) {
	var statusMsg string

	if currentAction != "" {
		if queueCount == 0 {
			statusMsg = fmt.Sprintf("⚡ Executing: %s", currentAction)
		} else {
			statusMsg = fmt.Sprintf("⚡ Executing: %s | Queue: %d/%d", currentAction, queueCount, maxQueue)
		}
	} else if queueCount > 0 {
		statusMsg = fmt.Sprintf("📋 Queue ready: %d/%d actions", queueCount, maxQueue)
	} else {
		statusMsg = "💤 Idle - ready for commands"
	}

	gw.SetStatusBarMessage(statusMsg)

	// Clear after a reasonable time
	go func() {
		time.Sleep(1500 * time.Millisecond)

		gw.StatusBarMutex.Lock()
		currentMsg := gw.StatusBarMessage
		gw.StatusBarMutex.Unlock()

		if currentMsg == statusMsg {
			gw.SetStatusBarMessage("")
		}
	}()
}

// ShowActionProgress displays progress information for long-running actions
func (gw *GameWindow) ShowActionProgress(actionType string, progress float64, timeLeft string) {
	progressPercent := progress * 100
	var icon string

	switch actionType {
	case "move":
		icon = "🚶"
	case "dig":
		icon = "⛏️"
	case "build_wall":
		icon = "🧱"
	case "mine":
		icon = "⚒️"
	default:
		icon = "⚡"
	}

	statusMsg := fmt.Sprintf("%s %s: %.0f%% complete (%s left)", icon, actionType, progressPercent, timeLeft)
	gw.SetStatusBarMessage(statusMsg)
}
