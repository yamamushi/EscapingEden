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
	if gw.InventoryDisplayType != edenitems.ItemTypeNull {
		gw.Log.Println(logging.LogInfo, fmt.Sprintf("%ss in Inventory", gw.InventoryDisplayType.String()))
	} else {
		gw.Log.Println(logging.LogInfo, "Inventory")
	}

	// Create a map to count stackable items by their names
	stackableCounts := make(map[string]int)
	weightMap := make(map[string]float64)
	countMap := make(map[string]string)
	keyMap := make(map[string]string)

	weight := 0.0

	for _, item := range gw.Inventory {
		if item.Type != gw.InventoryDisplayType && gw.InventoryDisplayType != edenitems.ItemTypeNull {
			continue
		}
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
	for _, itemName := range itemNames {
		itemEntry := countMap[itemName]
		//weightInfo := fmt.Sprintf("- %.2fkg", weightMap[itemName])
		itemInfo := fmt.Sprintf("%-*s", maxNameWidth, itemEntry)
		itemInfo += fmt.Sprintf(" - %.2fkg", weightMap[itemName])
		gw.Log.Println(logging.LogInfo, itemInfo)
	}
	if gw.InventoryDisplayType != edenitems.ItemTypeNull {
		gw.Log.Println(logging.LogInfo, fmt.Sprintf("%ss Weight: %.2fkg", gw.InventoryDisplayType.String(), weight))
	} else {
		gw.Log.Println(logging.LogInfo, fmt.Sprintf("Inventory Weight: %.2fkg", weight))
	}
	// Reset the inventory display type
	gw.InventoryDisplayType = edenitems.ItemTypeNull
}
