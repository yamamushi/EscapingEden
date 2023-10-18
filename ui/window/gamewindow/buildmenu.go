package gamewindow

func (gw *GameWindow) BuildMenu() {
	// Build our options
	options := []MenuBoxOption{{Name: "wall", Keybind: "w", Callback: gw.BuildWall}, {Name: "stairs", Keybind: "s"}, {Name: "floor", Keybind: "f"}, {Name: "door", Keybind: "d"}}
	// Create a new menu box
	mb := &MenuBox{X: gw.Width - 25, Y: gw.Height/2 - 10, Width: 21, Height: len(options) + 4, Title: "Build", Options: options}
	// Add the menu box to the game window
	gw.AddMenuBox(mb)
}

func (gw *GameWindow) BuildWall(box *MenuBox) {
	gw.StatusBarMutex.Lock()
	gw.StatusBarMessage = "Build with what?"
	gw.StatusBarMutex.Unlock()
	box.ResponseCallback = gw.BuildWallSend

}

func (gw *GameWindow) BuildWallConfirmDirection(box *MenuBox, input string) {
	gw.StatusBarMutex.Lock()
	gw.StatusBarMessage = "Building wall with " + box.CallbackData.(string) + " in " + input + " direction"
	gw.StatusBarMutex.Unlock()

	gw.CloseMenu = true
}

func (gw *GameWindow) BuildWallSend(box *MenuBox, input string) {
	if input == "?" {
		// Now we need to open a popup menu, and set the response callback to the BuildWallSend function
		// Create a new menu box
		//options := []MenuBoxOption{{Name: "wood", Keybind: "w"}, {Name: "stone", Keybind: "s"}, {Name: "brick", Keybind: "b"}}
		//mb := &MenuBox{X: gw.Width - 25, Y: gw.Height/2 - 10, Width: 21, Height: len(options) + 4, Title: "Materials", Options: options, ResponseCallback: gw.BuildWallConfirmDirection}
		//box.PopupMenu = mb
		gw.RequestInventoryDisplay()
		return
	}
	gw.StatusBarMutex.Lock()
	// Right now we don't have a way of parsing the material
	// So later we'll have to retrieve the material from the inventory
	gw.StatusBarMessage = "Building wall with " + input + " in which direction?"
	gw.StatusBarMutex.Unlock()
	box.ResponseCallback = gw.BuildWallConfirmDirection
	box.CallbackData = input

	//gw.CloseMenu = true
}

// Need an inventory menu, all it should do is load the inventory and display it, if an item is selected it should
// Send the selected item to the provided callback, and if no callback is provided it should do something else
// We'll figure that out later
