package gamewindow

import (
	"github.com/yamamushi/EscapingEden/edenitems"
	"log"
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
			log.Println(item.Attributes)
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

func (gw *GameWindow) HandleDig() {
	// Get the inventory
	gw.RequestInventoryUpdate(nil, "")
	gw.DisplayInventoryPostReceive = false
	gw.LockPendingInventory()

	gw.StatusBarMutex.Lock()
	gw.StatusBarMessage = "Dig in which direction?"
	gw.StatusBarMutex.Unlock()

}
