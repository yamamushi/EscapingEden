package gamewindow

import (
	"fmt"
	"github.com/yamamushi/EscapingEden/edenitems"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
	"log"
	"sort"
)

func (gw *GameWindow) ItemForHotkey(hotkey string) *edenitems.Item {
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()
	for _, item := range gw.Inventory {
		if item.Hotkey == hotkey {
			return &item
		}
	}
	return nil // Return nil if we don't find an item
}

func (gw *GameWindow) IsInventoryOpen() bool {
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()
	for _, window := range gw.Menus {
		if window.GetType() == MenuTypeInventory {
			return true
		}
	}
	return false
}

func (gw *GameWindow) CloseInventory() {
	for _, menu := range gw.Menus {
		if menu.GetType() == MenuTypeInventory {
			gw.RemoveMenuBox(menu)
			gw.InventoryDisplayType = edenitems.ItemTypeNull
			gw.MenuCallback = nil
			gw.InventoryCallbackPrompt = ""
			gw.SetStatusBarMessage("")
		}
	}
}

func (gw *GameWindow) DisplayInventoryAfterReceive(toggle bool) {
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()
	gw.DisplayInventoryPostReceive = toggle
}

func (gw *GameWindow) RequestInventoryUpdate(callback interface{}, callbackPrompt string) {
	inventoryRequest := messages.WindowMessage{Data: messages.GameManagerMessage{Data: messages.GameMessageData{CharacterID: gw.GetCharacterInfoField("id")}, Type: messages.GameManager_RequestInventory}, Type: messages.WM_GameCommand}
	gw.SendToConsole(inventoryRequest)
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()
	gw.MenuCallback = callback
	gw.InventoryCallbackPrompt = callbackPrompt
}

func (gw *GameWindow) UpdateInventory(inventory []edenitems.Item) {
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()
	gw.Inventory = inventory
	//gw.Log.Println(logging.LogInfo, "Inventory updated - ", len(gw.Inventory))
}

func (gw *GameWindow) UpdateInventoryDisplayType(itemType edenitems.ItemType) {
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()
	gw.InventoryDisplayType = itemType
}

func (gw *GameWindow) DisplayInventory() {
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()

	gw.DisplayInventoryPostReceive = false

	inventoryWindow := InventoryDisplay{Inventory: gw.Inventory,
		DisplayType: gw.InventoryDisplayType, GW: gw}
	inventoryWindow.ResponseCallback = gw.MenuCallback
	inventoryWindow.CallbackStatusBarMessage = gw.InventoryCallbackPrompt
	inventoryWindow.X = gw.Width - 30
	inventoryWindow.Y = gw.Height/2 - 10
	inventoryWindow.Type = MenuTypeInventory
	//inventoryWindow.Width = 29
	//inventoryWindow.Height = len(gw.Inventory) + 4

	//gw.Log.Println(logging.LogInfo, "Inventory Display Type: ", gw.InventoryDisplayType)
	//gw.Log.Println(logging.LogInfo, "Inventory Display Internal: ", inventoryWindow.DisplayType)

	if gw.InventoryDisplayType != edenitems.ItemTypeNull {
		gw.Log.Println(logging.LogInfo, fmt.Sprintf("%ss in Inventory", gw.InventoryDisplayType.String()))
		inventoryWindow.Title = fmt.Sprintf("%ss", gw.InventoryDisplayType.String())
	} else {
		inventoryWindow.Title = "Inventory"
	}

	inventoryWindow.PrepareContent()

	gw.AddMenuBox(&inventoryWindow)
}

func (gw *GameWindow) BuildHotKeys() {
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()
	for _, item := range gw.Inventory {
		if item.Type == gw.InventoryDisplayType || gw.InventoryDisplayType == edenitems.ItemTypeNull {
			gw.Hotkeys[item.Hotkey] = item
		}
	}
}

type InventoryDisplay struct {
	MenuBox
	Inventory   []edenitems.Item
	DisplayType edenitems.ItemType
	Content     []string
	Hotkeys     map[string]edenitems.Item
	GW          *GameWindow
}

func (inv *InventoryDisplay) GetType() MenuType {
	return MenuTypeInventory
}

func (inv *InventoryDisplay) UpdateCallbackFunction() {
	inv.ResponseCallback = inv.GW.MenuCallback
}

func (inv *InventoryDisplay) HandleInput(gw *GameWindow, inputType types.InputType, input string) {
	//inv.UpdateCallbackFunction()
	// Handle input for the menu box
	switch inputType {
	case types.InputEscape:
		//gw.Log.Println(logging.LogInfo, "InventoryDisplay received close")
		gw.CloseMenus = true
		return
	case types.InputCharacter:
		inv.HandleCharInput(input)
	}
}

func (inv *InventoryDisplay) HandleCharInput(input string) {
	// Check if the input is a hotkey in inv.Hotkeys
	if inv.CheckHotkeys {
		if hotkeyItem, ok := inv.Hotkeys[input]; ok {
			log.Println("Hotkey item: ", hotkeyItem)
			if inv.ResponseCallback != nil {
				//inv.CallbackData = hotkeyItem.Name
				switch inv.ResponseCallback.(type) {
				case func(*MenuBox, string):
					inv.ResponseCallback.(func(box *MenuBox, item string))(&inv.MenuBox, input)
					log.Println("Callback is func(*MenuBox, string)")
				case func(string):
					inv.ResponseCallback.(func(string))(input)
					log.Println("Callback is func(string)")
				}
				return
			}

		} else {
			inv.SetCallbackStatusBarMessage("Invalid item selected, please select an item from the list")
			//inv.GW.CloseMenus = true
		}
	} else {
		switch inv.ResponseCallback.(type) {
		case func(*MenuBox, string):
			inv.ResponseCallback.(func(box *MenuBox, item string))(&inv.MenuBox, input)
			log.Println("Callback is func(*MenuBox, string)")
		case func(string):
			inv.ResponseCallback.(func(string))(input)
			log.Println("Callback is func(string)")
		}
	}

}

func (inv *InventoryDisplay) Draw(gw *GameWindow) {
	inv.Clear(gw)
	inv.DrawMenuItems(gw)
	inv.DrawBorder(gw)
	inv.DrawTitle(gw)
	inv.DrawPopupMenu(gw)
	inv.StatusBarMessageMutex.Lock()
	defer inv.StatusBarMessageMutex.Unlock()
	gw.SetStatusBarMessage(inv.CallbackStatusBarMessage)
}

func (inv *InventoryDisplay) PrepareContent() {
	stackableCounts := make(map[string]int) // Storing the count for each stackable item
	weightMap := make(map[string]float64)   // Storing the weight for each item/stack of items
	countMap := make(map[string]string)     // Storing the output strings for each item/stack of items
	inv.Hotkeys = make(map[string]edenitems.Item)

	weight := 0.0
	//gw.Log.Println(logging.LogInfo, "Inventory Display Type: ", inv.DisplayType)

	for _, item := range inv.Inventory {
		if item.Type == inv.DisplayType || inv.DisplayType == edenitems.ItemTypeNull {
			//itemInfo := fmt.Sprintf("%s) %-*s", item.Hotkey, maxNameWidth, item.Name)
			itemInfo := fmt.Sprintf("%s) %s", item.Hotkey, item.Name)
			if item.Stackable {
				// Check if we've encountered this stackable item before
				if count, ok := stackableCounts[item.Name]; ok {
					stackableCounts[item.Name]++
					itemInfo += fmt.Sprintf(" (%d)", count+1)
				} else {
					stackableCounts[item.Name] = 1
				}
			}
			// Add the weight to the item info to two decimal places

			//itemInfo += fmt.Sprintf(" - %.2fkg", item.Weight)
			weightMap[item.Name] += item.Weight
			countMap[item.Name] = itemInfo
			inv.Hotkeys[item.Hotkey] = item
			weight += item.Weight
		}
	}

	// Create a slice to hold the item names for sorting
	itemNames := make([]string, 0, len(inv.Hotkeys))
	for _, item := range inv.Hotkeys {
		itemNames = append(itemNames, item.Name)
	}
	// Define a custom sorting function for item names
	sort.Slice(itemNames, func(i, j int) bool {
		return itemNames[i] > itemNames[j]
	}) // If we want to fix the sorting, because this will appear as a-A-b-B most likely, we'll replace this with a custom sort function
	// This just hasn't been tested with a large enough inventory to see if it's a problem yet

	maxNameWidth := 0
	for _, itemEntry := range countMap {
		if len(itemEntry) > maxNameWidth {
			maxNameWidth = len(itemEntry)
		}
	}

	// Now iterate over the inventory map and print the items to the screen
	//linecount := 0
	for _, itemName := range itemNames {
		itemEntry := countMap[itemName]
		//weightInfo := fmt.Sprintf("- %.2fkg", weightMap[itemName])
		itemInfo := fmt.Sprintf("%-*s", maxNameWidth, itemEntry)
		itemInfo += fmt.Sprintf(" - %.2fkg", weightMap[itemName])
		//gw.Log.Println(logging.LogInfo, itemInfo)
		//inventoryContents += itemInfo + "\n"
		//gw.PrintStringToMap(inv.X+2, inv.Y+index+2, itemInfo, "")
		inv.Content = append(inv.Content, itemInfo)
		//linecount++
	}

	weightLine := ""
	if inv.DisplayType != edenitems.ItemTypeNull {
		//gw.Log.Println(logging.LogInfo, fmt.Sprintf("%ss Weight: %.2fkg", gw.InventoryDisplayType.String(), weight))
		//inventoryContents += fmt.Sprintf("%ss Weight: %.2fkg", inv.DisplayType.String(), weight)
		weightLine = fmt.Sprintf("%ss Weight: %.2fkg", inv.DisplayType.String(), weight)
		//gw.PrintStringToMap(inv.X+2, inv.Y+linecount+3, weightLine, "")
		//inv.Content = append(inv.Content, weightLine+"\n")
	} else {
		//gw.Log.Println(logging.LogInfo, fmt.Sprintf("Inventory Weight: %.2fkg", weight))
		weightLine = fmt.Sprintf("Weight: %.2fkg", weight)
	}

	inv.Content = append(inv.Content, "")
	//gw.PrintStringToMap(inv.X+2, inv.Y+linecount+3, weightLine, "")
	inv.Content = append(inv.Content, weightLine)

	widthAdjustment := 0
	for _, line := range inv.Content {
		if len(line) > widthAdjustment {
			widthAdjustment = len(line)
		}
	}

	inv.Width = widthAdjustment + 4
	inv.Height = len(inv.Content) + 4

}

func (inv *InventoryDisplay) DrawMenuItems(gw *GameWindow) {
	//gw.Log.Println(logging.LogInfo, "Drawing inventory items")
	// Create a map to count stackable items by their names

	// Reset the inventory display type
	for index, content := range inv.Content {
		inv.PrintToMenu(gw, 2, index+2, content, "")
	}
}
