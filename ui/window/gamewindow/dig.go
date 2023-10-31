package gamewindow

import (
	"github.com/yamamushi/EscapingEden/edenitems"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (gw *GameWindow) DigMenu() {
	// Get the inventory
	gw.RequestInventoryUpdate(nil, "")
	gw.DisplayInventoryAfterReceive(false)
	gw.LockPendingInventory()

	// Build our options
	options := []MenuBoxOption{}
	for _, item := range gw.Inventory {
		if item.Type == edenitems.ItemTool {
			if item.Attributes["digging"] {
				option := MenuBoxOption{Name: item.Name, Keybind: item.Hotkey, Callback: gw.HandleDig}
				options = append(options, option)
			}
			//log.Println(item.Attributes)
		}
	}
	if len(options) == 0 {
		gw.SetStatusBarMessage("You have nothing to dig with.")
		return
	}

	// Create a new menu box
	mb := &MenuBox{X: gw.Width - 25, Y: gw.Height/2 - 10, Width: 21, Height: len(options) + 4, Title: "Dig", Options: options, CallbackStatusBarMessage: "Dig with what?"}
	// Add the menu box to the game window
	gw.AddMenuBox(mb)
}

func (gw *GameWindow) HandleDig(box *MenuBox, input string) {
	//gw.Log.Println(logging.LogInfo, "HandleDig called")

	box.SetCallbackStatusBarMessage("Dig in which direction?")

	box.ResponseCallback = gw.DigConfirmDirection
	box.CallbackData = input
	box.ToggleHotkeyCheck(false)
}

func (gw *GameWindow) DigConfirmDirection(box *MenuBox, input string) {
	item := gw.ItemForHotkey(box.CallbackData.(string))

	// Check vi movement keys
	if input != "y" || input != "u" || input != "h" || input != "j" || input != "k" || input != "l" || input != "b" || input != "n" {
		box.SetCallbackStatusBarMessage("Digging with " + item.Name + " in " + input + " direction")
		gw.StatusBarMessage = "Digging with " + item.Name + " in " + input + " direction"
		// Now we fire off our dig request and handle the response in the messages loop outside of here
		gw.SendDigRequest(item, input)
	} else {
		box.SetCallbackStatusBarMessage("Invalid direction selected.")
		gw.StatusBarMessage = "Invalid direction selected."
	}

	// Cleanup
	gw.CloseMenus = true
}

func (gw *GameWindow) SendDigRequest(item *edenitems.Item, dir string) {
	deltaX, deltaY := 0, 0
	switch dir {
	case "y":
		deltaX = -1
		deltaY = -1
	case "u":
		deltaX = 1
		deltaY = -1
	case "h":
		deltaX = -1
		deltaY = 0
	case "j":
		deltaX = 0
		deltaY = 1
	case "k":
		deltaX = 0
		deltaY = -1
	case "l":
		deltaX = 1
		deltaY = 0
	case "b":
		deltaX = -1
		deltaY = 1
	case "n":
		deltaX = 1
		deltaY = 1
	default:
		// We should never get here, but just in case
		gw.Log.Println(logging.LogError, "Invalid direction selected")
		return
	}

	message := messages.WindowMessage{Type: messages.WM_GameCommand, Data: messages.GameManagerMessage{Type: messages.GameManager_DigCommand, Data: messages.GameMessageData{CharacterID: gw.GetCharacterInfoField("id"), Data: messages.GameCharDig{DeltaX: deltaX, DeltaY: deltaY, ItemID: item.ID}}}}
	gw.SendToConsole(message)

}
