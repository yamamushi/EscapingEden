package gamewindow

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
)

// cleanUnicodeWhitespace removes all types of Unicode whitespace characters
// including regular spaces, non-breaking spaces, zero-width spaces, etc.
func cleanUnicodeWhitespace(s string) string {
	// Convert to runes for proper Unicode handling
	runes := []rune(s)
	result := make([]rune, 0, len(runes))

	for _, r := range runes {
		// Skip all Unicode whitespace characters
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' &&
			r != '\u00A0' && // Non-breaking space
			r != '\u2000' && r != '\u2001' && r != '\u2002' && r != '\u2003' && // En quad, Em quad, En space, Em space
			r != '\u2004' && r != '\u2005' && r != '\u2006' && r != '\u2007' && // Three-per-em space, Four-per-em space, Six-per-em space, Figure space
			r != '\u2008' && r != '\u2009' && r != '\u200A' && r != '\u200B' && // Punctuation space, Thin space, Hair space, Zero-width space
			r != '\u200C' && r != '\u200D' && r != '\u2028' && r != '\u2029' && // Zero-width non-joiner, Zero-width joiner, Line separator, Paragraph separator
			r != '\u202F' && r != '\u205F' && r != '\u3000' && r != '\uFEFF' { // Narrow no-break space, Medium mathematical space, Ideographic space, Zero-width no-break space
			result = append(result, r)
		}
	}

	return string(result)
}

// formatItemColor converts integer color codes to escape sequences (same format as tiles)
func formatItemColor(fgColor, bgColor int) string {
	fg := strconv.Itoa(fgColor)
	bg := strconv.Itoa(bgColor)
	return "\033[38;5;" + fg + "m" + "\033[48;5;" + bg + "m"
}

func (gw *GameWindow) ItemForHotkey(hotkey string) *edentypes.Item {
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
			gw.InventoryDisplayType = edentypes.ItemTypeNull
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

func (gw *GameWindow) UpdateInventory(inventory []edentypes.Item) {
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()
	gw.Inventory = inventory
	//gw.Log.Println(logging.LogInfo, "Inventory updated - ", len(gw.Inventory))
}

func (gw *GameWindow) UpdateInventoryDisplayType(itemType edentypes.ItemType) {
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
	inventoryWindow.X = gw.Width - 35 // Move further left to avoid border overlap
	inventoryWindow.Y = 3             // Position near top with some margin
	inventoryWindow.Type = MenuTypeInventory
	inventoryWindow.ToggleHotkeyCheck(true) // Enable hotkey checking for interactive inventory
	//inventoryWindow.Width = 29
	//inventoryWindow.Height = len(gw.Inventory) + 4

	//gw.Log.Println(logging.LogInfo, "Inventory Display Type: ", gw.InventoryDisplayType)
	//gw.Log.Println(logging.LogInfo, "Inventory Display Internal: ", inventoryWindow.DisplayType)

	if gw.InventoryDisplayType != edentypes.ItemTypeNull {
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
		if item.Type == gw.InventoryDisplayType || gw.InventoryDisplayType == edentypes.ItemTypeNull {
			gw.Hotkeys[item.Hotkey] = item
		}
	}
}

type InventoryDisplay struct {
	MenuBox
	Inventory   []edentypes.Item
	DisplayType edentypes.ItemType
	Content     []string
	Hotkeys     map[string]edentypes.Item
	GW          *GameWindow
}

func (inv *InventoryDisplay) GetType() MenuType {
	return MenuTypeInventory
}

func (inv *InventoryDisplay) UpdateCallbackFunction() {
	inv.ResponseCallback = inv.GW.MenuCallback
}

func (inv *InventoryDisplay) HandleInput(gw *GameWindow, inputType types.InputType, input string) {
	// Handle input for the menu box
	switch inputType {
	case types.InputEscape:
		// Close inventory
		gw.CloseMenus = true
		return
	case types.InputCharacter:
		// Check if this is the close key
		if input == "!" {
			gw.CloseMenus = true
			return
		}

		// Check if this is a hotkey for an item
		if item, exists := inv.Hotkeys[input]; exists {
			// Show item info popup
			inv.ShowItemInfo(gw, item)
			return
		}

		// Handle other character input if needed
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
				case func(string):
					inv.ResponseCallback.(func(string))(input)
				}
				return
			}

		} else {
			inv.SetCallbackStatusBarMessage("Invalid item selected, please select an item from the list")
			return
			//inv.GW.CloseMenus = true
		}
	} else {
		switch inv.ResponseCallback.(type) {
		case func(*MenuBox, string):
			inv.ResponseCallback.(func(box *MenuBox, item string))(&inv.MenuBox, input)
		case func(string):
			inv.ResponseCallback.(func(string))(input)
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

	// Show instructions for inventory interaction
	if inv.CallbackStatusBarMessage != "" {
		gw.SetStatusBarMessage(inv.CallbackStatusBarMessage)
	} else {
		gw.SetStatusBarMessage("Press item letter (a, b, c...) for info, ! to close")
	}
}

func (inv *InventoryDisplay) PrepareContent() {
	stackableCounts := make(map[string]int) // Storing the count for each stackable item
	weightMap := make(map[string]float64)   // Storing the weight for each item/stack of items
	countMap := make(map[string]string)     // Storing the output strings for each item/stack of items
	inv.Hotkeys = make(map[string]edentypes.Item)

	weight := 0.0
	//gw.Log.Println(logging.LogInfo, "Inventory Display Type: ", inv.DisplayType)

	for _, item := range inv.Inventory {
		if item.Type == inv.DisplayType || inv.DisplayType == edentypes.ItemTypeNull {
			// Create clean item display with exactly one space between symbol and name
			// Clean any Unicode whitespace from symbol and name
			cleanSymbol := cleanUnicodeWhitespace(item.Symbol)
			cleanName := cleanUnicodeWhitespace(item.Name)

			var itemInfo string
			if cleanSymbol != "" {
				itemInfo = fmt.Sprintf("%s) %s %s", item.Hotkey, cleanSymbol, cleanName)
			} else {
				itemInfo = fmt.Sprintf("%s) %s", item.Hotkey, cleanName)
			}
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
		itemInfo := fmt.Sprintf("%s - %.2fkg", itemEntry, weightMap[itemName])

		//gw.Log.Println(logging.LogInfo, itemInfo)
		//inventoryContents += itemInfo + "\n"
		//gw.PrintStringToMap(inv.X+2, inv.Y+index+2, itemInfo, "")
		inv.Content = append(inv.Content, itemInfo)
		//linecount++
	}

	weightLine := ""
	if inv.DisplayType != edentypes.ItemTypeNull {
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

	// Add close instruction at the bottom
	inv.Content = append(inv.Content, "")
	inv.Content = append(inv.Content, "Press ! to close")

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

	for index, content := range inv.Content {
		// Apply colored symbols to inventory items
		parts := strings.Split(content, ")")
		if len(parts) > 1 {
			hotkey := strings.TrimSpace(parts[0])
			if item, exists := inv.Hotkeys[hotkey]; exists && item.Symbol != "" {
				// Find the symbol in the content and apply color
				symbolIndex := strings.Index(content, item.Symbol)
				if symbolIndex != -1 {
					// Draw everything before the symbol
					beforeSymbol := content[:symbolIndex]
					inv.PrintToMenu(gw, 2, index+2, beforeSymbol, "")

					// Draw the colored symbol
					itemColorCode := formatItemColor(item.FGColor, item.BGColor)
					inv.PrintToMenu(gw, 2+len(beforeSymbol), index+2, item.Symbol, itemColorCode)

					// Draw everything after the symbol - use proper Unicode-aware slicing
					// Convert to runes for proper Unicode handling
					contentRunes := []rune(content)
					symbolRunes := []rune(item.Symbol)

					// Find the symbol position in runes
					symbolStartRune := len([]rune(beforeSymbol))
					symbolEndRune := symbolStartRune + len(symbolRunes)

					// Get the part after the symbol
					if symbolEndRune < len(contentRunes) {
						afterSymbol := string(contentRunes[symbolEndRune:])
						// Use rune length for proper positioning with Unicode characters
						beforeSymbolRuneLen := len([]rune(beforeSymbol))
						symbolRuneLen := len(symbolRunes)
						inv.PrintToMenu(gw, 2+beforeSymbolRuneLen+symbolRuneLen, index+2, afterSymbol, "")
					}
					continue
				}
			}
		}

		// Normal drawing for lines without symbols
		inv.PrintToMenu(gw, 2, index+2, content, "")
	}
}

// ShowItemInfo displays detailed information about an item
func (inv *InventoryDisplay) ShowItemInfo(gw *GameWindow, item edentypes.Item) {
	itemInfo := ItemInfoDisplay{
		Item: item,
		GW:   gw,
	}
	itemInfo.Type = MenuTypeItemInfo
	itemInfo.Title = fmt.Sprintf("Item Info: %s", item.Name)
	itemInfo.PrepareContent()

	// Add the item info popup on top of the inventory (using unsafe version since mutex is already held by input handler)
	gw.AddMenuBoxUnsafe(&itemInfo)
}
