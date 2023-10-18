package gamewindow

import (
	"fmt"
	"github.com/yamamushi/EscapingEden/edenitems"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
	"sort"
)

func (gw *GameWindow) RequestInventoryDisplay(callback interface{}, callbackPrompt string) {
	inventoryRequest := messages.WindowMessage{Data: messages.GameManagerMessage{Data: messages.GameMessageData{CharacterID: gw.GetCharacterInfoField("id")}, Type: messages.GameManager_RequestInventory}, Type: messages.WM_GameCommand}
	gw.SendToConsole(inventoryRequest)
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()
	gw.InventoryCallback = callback
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

	inventoryWindow := InventoryDisplay{Inventory: gw.Inventory,
		DisplayType:    gw.InventoryDisplayType,
		Callback:       gw.InventoryCallback,
		CallbackPrompt: gw.InventoryCallbackPrompt}
	inventoryWindow.X = gw.Width - 30
	inventoryWindow.Y = gw.Height/2 - 10
	inventoryWindow.Width = 29
	inventoryWindow.Height = len(gw.Inventory) + 4

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

type InventoryDisplay struct {
	MenuBox
	Inventory      []edenitems.Item
	DisplayType    edenitems.ItemType
	Content        []string
	Hotkeys        map[string]string
	Callback       interface{}
	CallbackPrompt string
}

func (inv *InventoryDisplay) HandleInput(gw *GameWindow, inputType types.InputType, input string) {
	// Handle input for the menu box
	switch inputType {
	case types.InputEscape:
		//gw.Log.Println(logging.LogInfo, "InventoryDisplay received close")
		gw.CloseMenus = true
		return
	case types.InputCharacter:
		if inv.Callback != nil {
			switch inv.Callback.(type) {
			case func(string):
				inv.Callback.(func(string))(input)
			}
			return
		}
	}
}

func (inv *InventoryDisplay) Draw(gw *GameWindow) {
	inv.Clear(gw)
	inv.DrawMenuItems(gw)
	inv.DrawBorder(gw)
	inv.DrawTitle(gw)
	inv.DrawPopupMenu(gw)
}

func (inv *InventoryDisplay) PrepareContent() {
	stackableCounts := make(map[string]int) // Storing the count for each stackable item
	weightMap := make(map[string]float64)   // Storing the weight for each item/stack of items
	countMap := make(map[string]string)     // Storing the output strings for each item/stack of items
	keyMap := make(map[string]string)       // For organizing our output by hotkey order

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
			keyMap[item.Hotkey] = item.Name
			weight += item.Weight
		}
	}

	// Create a slice to hold the item names for sorting
	itemNames := make([]string, 0, len(keyMap))
	for _, itemName := range keyMap {
		itemNames = append(itemNames, itemName)
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
	gw.SetStatusBarMessage(inv.CallbackPrompt)
}
