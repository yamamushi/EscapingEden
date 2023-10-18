package gamewindow

import (
	"fmt"
	"github.com/yamamushi/EscapingEden/edenitems"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"sort"
)

func (gw *GameWindow) RequestInventoryDisplay() {
	inventoryRequest := messages.WindowMessage{Data: messages.GameManagerMessage{Data: messages.GameMessageData{CharacterID: gw.GetCharacterInfoField("id")}, Type: messages.GameManager_RequestInventory}, Type: messages.WM_GameCommand}
	gw.SendToConsole(inventoryRequest)
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

	inventoryWindow := InventoryDisplay{Inventory: gw.Inventory, DisplayType: gw.InventoryDisplayType}
	inventoryWindow.X = gw.Width - 30
	inventoryWindow.Y = gw.Height/2 - 10
	inventoryWindow.Width = 27
	inventoryWindow.Height = len(gw.Inventory) + 4

	if gw.InventoryDisplayType != edenitems.ItemTypeNull {
		gw.Log.Println(logging.LogInfo, fmt.Sprintf("%ss in Inventory", gw.InventoryDisplayType.String()))
		inventoryWindow.Title = fmt.Sprintf("%ss", gw.InventoryDisplayType.String())
	} else {
		inventoryWindow.Title = "Inventory"
	}

	gw.AddMenuBox(&inventoryWindow)
}

type InventoryDisplay struct {
	MenuBox
	Inventory   []edenitems.Item
	DisplayType edenitems.ItemType
}

func (inv *InventoryDisplay) Draw(gw *GameWindow) {
	inv.Clear(gw)
	inv.DrawBorder(gw)
	inv.DrawTitle(gw)
	inv.DrawMenuItems(gw)
	inv.DrawPopupMenu(gw)
}

func (inv *InventoryDisplay) DrawMenuItems(gw *GameWindow) {
	gw.Log.Println(logging.LogInfo, "Drawing inventory items")
	// Create a map to count stackable items by their names
	stackableCounts := make(map[string]int) // Storing the count for each stackable item
	weightMap := make(map[string]float64)   // Storing the weight for each item/stack of items
	countMap := make(map[string]string)     // Storing the output strings for each item/stack of items
	keyMap := make(map[string]string)       // For organizing our output by hotkey order

	weight := 0.0

	Two things, non display types are still being displayed for some reason
	This info should only be generated at the window initialization, instead it's happening ever redraw which is bad

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
	linecount := 0
	for index, itemName := range itemNames {
		itemEntry := countMap[itemName]
		//weightInfo := fmt.Sprintf("- %.2fkg", weightMap[itemName])
		itemInfo := fmt.Sprintf("%-*s", maxNameWidth, itemEntry)
		itemInfo += fmt.Sprintf(" - %.2fkg", weightMap[itemName])
		//gw.Log.Println(logging.LogInfo, itemInfo)
		//inventoryContents += itemInfo + "\n"
		gw.PrintStringToMap(inv.X+2, inv.Y+index+2, itemInfo, "")
		linecount++
	}

	weightLine := ""
	if inv.DisplayType != edenitems.ItemTypeNull {
		//gw.Log.Println(logging.LogInfo, fmt.Sprintf("%ss Weight: %.2fkg", gw.InventoryDisplayType.String(), weight))
		//inventoryContents += fmt.Sprintf("%ss Weight: %.2fkg", inv.DisplayType.String(), weight)
		weightLine = fmt.Sprintf("%ss Weight: %.2fkg", inv.DisplayType.String(), weight)
		gw.PrintStringToMap(inv.X+2, inv.Y+linecount+3, weightLine, "")
	} else {
		//gw.Log.Println(logging.LogInfo, fmt.Sprintf("Inventory Weight: %.2fkg", weight))
		weightLine = fmt.Sprintf("Weight: %.2fkg", weight)
	}

	if len(weightLine) > inv.Width-5 {
		weightLine = weightLine[:inv.Width-5]
	}
	gw.PrintStringToMap(inv.X+2, inv.Y+linecount+3, weightLine, "")

	// Reset the inventory display type
	inv.DisplayType = edenitems.ItemTypeNull
}
