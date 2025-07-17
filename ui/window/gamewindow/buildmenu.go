package gamewindow

import (
	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (gw *GameWindow) BuildMenu() {
	// Build our options
	options := []MenuBoxOption{{Name: "wall", Keybind: "w", Callback: gw.BuildWall}, {Name: "stairs", Keybind: "s"}, {Name: "floor", Keybind: "f"}, {Name: "door", Keybind: "d"}}
	// Create a new menu box
	mb := &MenuBox{X: gw.Width - 25, Y: gw.Height/2 - 10, Width: 21, Height: len(options) + 4, Title: "Build", Options: options, CallbackStatusBarMessage: "Build what?"}
	// Add the menu box to the game window
	gw.AddMenuBox(mb)
}

func (gw *GameWindow) BuildWall(box *MenuBox) {
	gw.StatusBarMutex.Lock()
	defer gw.StatusBarMutex.Unlock()
	box.SetCallbackStatusBarMessage("Build with what?")
	box.ResponseCallback = gw.BuildWallSend
	gw.RequestInventoryUpdate(nil, "")
	gw.DisplayInventoryPostReceive = false
	gw.LockPendingInventory()
}

func (gw *GameWindow) BuildWallConfirmDirection(box *MenuBox, input string) {
	gw.StatusBarMutex.Lock()
	//gw.StatusBarMessage = "Building wall with " + box.CallbackData.(string) + " in " + input + " direction"
	defer gw.StatusBarMutex.Unlock()
	//log.Println("BuildWallConfirmDirection received input: ", input)

	item := gw.ItemForHotkey(box.CallbackData.(string))

	// Check vi movement keys
	if input != "y" || input != "u" || input != "h" || input != "j" || input != "k" || input != "l" || input != "b" || input != "n" {
		box.SetCallbackStatusBarMessage("Building wall with " + item.Name + " in " + input + " direction")
		gw.StatusBarMessage = "Building wall with " + item.Name + " in " + input + " direction"
		gw.SendBuildWallRequest(item, nil, input)
	} else {
		box.SetCallbackStatusBarMessage("Invalid direction selected.")
		gw.StatusBarMessage = "Invalid direction selected."
	}

	// Cleanup
	gw.CloseMenus = true
}

func (gw *GameWindow) BuildWallSend(box *MenuBox, input string) {
	gw.StatusBarMutex.Lock()
	defer gw.StatusBarMutex.Unlock()
	//log.Println("BuildWallSend received input: ", input)
	if input == "?" {
		gw.InventoryDisplayType = edentypes.ItemMaterial
		gw.RequestInventoryUpdate(gw.BuildWallSend, "Build with what?")
		gw.DisplayInventoryAfterReceive(true)
		return
	}

	// Right now we don't have a way of parsing the material
	// So later we'll have to retrieve the material from the inventory
	item := gw.ItemForHotkey(input)
	if item == nil {
		//log.Println("BuildWallSend received nil item")
		box.SetCallbackStatusBarMessage("You don't have that item.")
		//gw.MenusMutex.Lock()
		//defer gw.MenusMutex.Unlock()
		//gw.CloseMenus = true
		return
	}
	if item.Type != edentypes.ItemMaterial {
		//log.Println("BuildWallSend received non-material item")
		box.SetCallbackStatusBarMessage("That item is not a material suitable for building.")
		//gw.MenusMutex.Lock()
		//defer gw.MenusMutex.Unlock()
		//gw.CloseMenus = true
		return
	}
	//log.Println("BuildWallSend received item: ", item)

	box.SetCallbackStatusBarMessage("Building wall with " + item.Name + " in which direction?")
	box.ResponseCallback = gw.BuildWallConfirmDirection
	box.CallbackData = input
	box.ToggleHotkeyCheck(false)
}

func (gw *GameWindow) SendBuildWallRequest(item *edentypes.Item, tool *edentypes.Item, dir string) {
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

	// Note that toolID is "" right now
	message := messages.WindowMessage{Type: messages.WM_GameCommand, Data: messages.GameManagerMessage{Type: messages.GameManager_BuildWallCommand, Data: messages.GameMessageData{CharacterID: gw.GetCharacterInfoField("id"), Data: messages.GameCharBuildWall{DeltaX: deltaX, DeltaY: deltaY, ItemID: item.ID, ToolID: ""}}}}
	gw.SendToConsole(message)

}
