package gamewindow

import "github.com/yamamushi/EscapingEden/edenitems"

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
	gw.StatusBarMutex.Unlock()
	box.CallbackStatusBarMessage = "Build with what?"
	box.ResponseCallback = gw.BuildWallSend
}

func (gw *GameWindow) BuildWallConfirmDirection(input string) {
	gw.StatusBarMutex.Lock()
	//gw.StatusBarMessage = "Building wall with " + box.CallbackData.(string) + " in " + input + " direction"
	defer gw.StatusBarMutex.Unlock()

	// Check vi movement keys
	if input != "y" || input != "u" || input != "h" || input != "j" || input != "k" || input != "l" || input != "b" || input != "n" {

	}

	// Cleanup
	gw.CloseMenus = true
}

func (gw *GameWindow) BuildWallSend(box *MenuBox, input string) {
	gw.StatusBarMutex.Lock()
	defer gw.StatusBarMutex.Unlock()
	if input == "?" {
		gw.InventoryDisplayType = edenitems.ItemMaterial
		gw.RequestInventoryDisplay(gw.BuildWallConfirmDirection, "Build with what?")
		return
	}

	// Right now we don't have a way of parsing the material
	// So later we'll have to retrieve the material from the inventory
	box.CallbackStatusBarMessage = "Building wall with " + input + " in which direction?"
	box.ResponseCallback = gw.BuildWallConfirmDirection
	box.CallbackData = input

	//gw.CloseMenus = true
}

// Need an inventory menu, all it should do is load the inventory and display it, if an item is selected it should
// Send the selected item to the provided callback, and if no callback is provided it should do something else
// We'll figure that out later
