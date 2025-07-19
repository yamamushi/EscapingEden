package gamewindow

import (
	"fmt"
	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
)

func (gw *GameWindow) BuildMenu() {
	// Build our options
	options := []MenuBoxOption{
		{Name: "wall", Keybind: "w", Callback: gw.BuildWall},
		{Name: "stairs", Keybind: "s", Callback: gw.BuildStairs},
		{Name: "floor", Keybind: "f", Callback: gw.BuildFloor},
		{Name: "door", Keybind: "d", Callback: gw.BuildDoor},
	}
	// Create a new menu box
	mb := &MenuBox{X: gw.Width - 25, Y: gw.Height/2 - 10, Width: 21, Height: len(options) + 4, Title: "Build", Options: options, CallbackStatusBarMessage: "Build what?"}

	// Debug: Log that build menu is being created
	gw.Log.Println(logging.LogInfo, "Creating build menu with", len(options), "options")

	// Add the menu box to the game window
	gw.AddMenuBox(mb)
}

func (gw *GameWindow) BuildWall(box *MenuBox) {
	// Debug: Log that BuildWall was called
	gw.Log.Println(logging.LogInfo, "BuildWall function called!")

	// Show custom material selection popup first
	gw.ShowMaterialSelectionPopup("Select material to build wall with:", gw.BuildWallMaterialSelected)

	// Don't close menus immediately - let the material selection popup handle it
	// The build menu will be closed when the material selection popup is displayed
}

// BuildWallMaterialSelected handles when a material is selected from the material selection popup
func (gw *GameWindow) BuildWallMaterialSelected(box *MenuBox, input string) {
	// Debug: Log that this function was called
	gw.Log.Println(logging.LogInfo, fmt.Sprintf("BuildWallMaterialSelected called with input: '%s'", input))

	// Get the selected material item using the stored material hotkeys
	item := gw.GetMaterialForHotkey(input)
	if item == nil {
		gw.Log.Println(logging.LogError, fmt.Sprintf("No material found for hotkey '%s'", input))
		gw.SetStatusBarMessage("Invalid material selected.")
		gw.CloseMenus = true
		return
	}

	gw.Log.Println(logging.LogInfo, fmt.Sprintf("Selected material: %s", item.Name))

	// Close all menus and enter direction selection mode
	gw.CloseMenus = true

	// Store the selected material for direction selection
	gw.SetBuildWallMaterial(item)

	// Set status bar message and enter direction selection mode
	gw.SetStatusBarMessage("Building wall with " + item.Name + ". Select direction (hjklyubn) or ! to cancel:")
	gw.Log.Println(logging.LogInfo, "Entering direction selection mode - no popup needed")
}

// BuildWallWithDirection handles the direction selection and sends the build request
func (gw *GameWindow) BuildWallWithDirection(direction string, item *edentypes.Item) {
	// Validate direction input
	validDirections := map[string]bool{
		"y": true, "u": true, "h": true, "j": true,
		"k": true, "l": true, "b": true, "n": true,
	}

	if !validDirections[direction] {
		gw.SetStatusBarMessage("Invalid direction selected.")
		gw.CloseMenus = true
		return
	}

	// Send the build request
	gw.SetStatusBarMessage("Building wall with " + item.Name + " in " + direction + " direction")
	gw.SendBuildWallRequest(item, nil, direction)

	// Close the direction menu
	gw.CloseMenus = true
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

// BuildStairs handles building stairs
func (gw *GameWindow) BuildStairs(box *MenuBox) {
	// Show custom material selection popup first
	gw.ShowMaterialSelectionPopup("Select material to build stairs with:", gw.BuildStairsMaterialSelected)

	// Don't close menus immediately - let the material selection popup handle it
}

// BuildStairsMaterialSelected handles when a material is selected for stairs
func (gw *GameWindow) BuildStairsMaterialSelected(box *MenuBox, input string) {
	// Get the selected material item using the stored material hotkeys
	item := gw.GetMaterialForHotkey(input)
	if item == nil {
		gw.SetStatusBarMessage("Invalid material selected.")
		gw.CloseMenus = true
		return
	}

	// For stairs, we might not need direction - just build at current location
	gw.SetStatusBarMessage("Building stairs with " + item.Name)
	// TODO: Implement stairs building command when available
	gw.CloseMenus = true
}

// BuildFloor handles building floors
func (gw *GameWindow) BuildFloor(box *MenuBox) {
	// Show custom material selection popup first
	gw.ShowMaterialSelectionPopup("Select material to build floor with:", gw.BuildFloorMaterialSelected)

	// Don't close menus immediately - let the material selection popup handle it
}

// BuildFloorMaterialSelected handles when a material is selected for floor
func (gw *GameWindow) BuildFloorMaterialSelected(box *MenuBox, input string) {
	// Get the selected material item using the stored material hotkeys
	item := gw.GetMaterialForHotkey(input)
	if item == nil {
		gw.SetStatusBarMessage("Invalid material selected.")
		gw.CloseMenus = true
		return
	}

	// For floor, we might not need direction - just build at current location
	gw.SetStatusBarMessage("Building floor with " + item.Name)
	// TODO: Implement floor building command when available
	gw.CloseMenus = true
}

// BuildDoor handles building doors
func (gw *GameWindow) BuildDoor(box *MenuBox) {
	// Show custom material selection popup first
	gw.ShowMaterialSelectionPopup("Select material to build door with:", gw.BuildDoorMaterialSelected)

	// Don't close menus immediately - let the material selection popup handle it
}

// BuildDoorMaterialSelected handles when a material is selected for door
func (gw *GameWindow) BuildDoorMaterialSelected(box *MenuBox, input string) {
	// Get the selected material item using the stored material hotkeys
	item := gw.GetMaterialForHotkey(input)
	if item == nil {
		gw.SetStatusBarMessage("Invalid material selected.")
		gw.CloseMenus = true
		return
	}

	// Close the material selection popup
	gw.CloseMenus = true

	// Create a new menu for direction selection (doors need direction)
	gw.SetStatusBarMessage("Building door with " + item.Name + ". Select direction (hjklyubn) or ! to cancel:")

	// Create a simple direction selection menu
	options := []MenuBoxOption{
		{Name: "↖ y", Keybind: "y"},
		{Name: "↑ k", Keybind: "k"},
		{Name: "↗ u", Keybind: "u"},
		{Name: "← h", Keybind: "h"},
		{Name: "→ l", Keybind: "l"},
		{Name: "↙ b", Keybind: "b"},
		{Name: "↓ j", Keybind: "j"},
		{Name: "↘ n", Keybind: "n"},
	}

	directionMenu := &MenuBox{
		X:                        gw.Width/2 - 10,
		Y:                        gw.Height/2 - 6,
		Width:                    20,
		Height:                   len(options) + 4,
		Title:                    "Select Direction",
		Options:                  options,
		CallbackStatusBarMessage: "Building door with " + item.Name + ". Select direction:",
		ResponseCallback: func(menuBox *MenuBox, direction string) {
			gw.BuildDoorWithDirection(direction, item)
		},
	}
	directionMenu.ToggleHotkeyCheck(true)
	gw.AddMenuBox(directionMenu)
}

// BuildDoorWithDirection handles the direction selection for door building
func (gw *GameWindow) BuildDoorWithDirection(direction string, item *edentypes.Item) {
	// Validate direction input
	validDirections := map[string]bool{
		"y": true, "u": true, "h": true, "j": true,
		"k": true, "l": true, "b": true, "n": true,
	}

	if !validDirections[direction] {
		gw.SetStatusBarMessage("Invalid direction selected.")
		gw.CloseMenus = true
		return
	}

	// Send the build request
	gw.SetStatusBarMessage("Building door with " + item.Name + " in " + direction + " direction")
	// TODO: Implement door building command when available
	gw.CloseMenus = true
}

// MaterialStack represents a stack of identical materials
type MaterialStack struct {
	Item     edentypes.Item
	Quantity int
}

// MaterialSelectionDisplay represents a custom popup for selecting building materials
type MaterialSelectionDisplay struct {
	*MenuBox
	Materials []edentypes.Item
	Stacks    []MaterialStack
	Hotkeys   map[string]edentypes.Item
	GW        *GameWindow
}

func (msd *MaterialSelectionDisplay) GetType() MenuType {
	return MenuTypeInventory // Reuse inventory type for now
}

func (msd *MaterialSelectionDisplay) HandleInput(gw *GameWindow, inputType types.InputType, input string) {
	// Debug: Log all input received by material selection popup
	gw.Log.Println(logging.LogInfo, fmt.Sprintf("MaterialSelectionDisplay.HandleInput called: type=%v, input='%s'", inputType, input))

	// Handle input for the material selection
	switch inputType {
	case types.InputEscape:
		// Close material selection
		gw.Log.Println(logging.LogInfo, "Material selection closing due to Escape key")
		gw.CloseMenus = true
		return
	case types.InputCharacter:
		// Check if this is the close key
		if input == "!" {
			gw.Log.Println(logging.LogInfo, "Material selection closing due to ! key")
			gw.CloseMenus = true
			return
		}

		// Check if this is a hotkey for a material
		if _, exists := msd.Hotkeys[input]; exists {
			gw.Log.Println(logging.LogInfo, fmt.Sprintf("Valid material hotkey '%s' selected", input))
			// Call the response callback with the selected material
			if msd.ResponseCallback != nil {
				gw.Log.Println(logging.LogInfo, "Calling response callback for material selection")
				switch callback := msd.ResponseCallback.(type) {
				case func(*MenuBox, string):
					callback(msd.MenuBox, input)
				default:
					gw.Log.Println(logging.LogError, fmt.Sprintf("Unexpected callback type: %T", msd.ResponseCallback))
				}
			} else {
				gw.Log.Println(logging.LogError, "No response callback set for material selection")
			}
			return
		}

		// Invalid selection
		gw.Log.Println(logging.LogInfo, fmt.Sprintf("Invalid material hotkey '%s' selected", input))
		gw.SetStatusBarMessage("Invalid material selected. Press ! to cancel.")
	default:
		gw.Log.Println(logging.LogInfo, fmt.Sprintf("Material selection received unhandled input type: %v", inputType))
	}
}

func (msd *MaterialSelectionDisplay) PrepareContent() {
	msd.Hotkeys = make(map[string]edentypes.Item)

	// Get all building materials from inventory
	materials := []edentypes.Item{}
	for _, item := range msd.GW.Inventory {
		if item.Type == edentypes.ItemMaterial {
			materials = append(materials, item)
		}
	}

	// Debug: Log how many materials were found
	msd.GW.Log.Println(logging.LogInfo, fmt.Sprintf("PrepareContent called, found %d materials in inventory", len(materials)))
	if len(materials) == 0 {
		msd.GW.Log.Println(logging.LogInfo, fmt.Sprintf("Total inventory items: %d", len(msd.GW.Inventory)))
		for i, item := range msd.GW.Inventory {
			msd.GW.Log.Println(logging.LogInfo, fmt.Sprintf("Item %d: %s (Type: %v)", i, item.Name, item.Type))
		}

		// If no materials found, show all items as potential materials for debugging
		msd.GW.Log.Println(logging.LogInfo, "No materials found, showing all items for debugging")
		for _, item := range msd.GW.Inventory {
			materials = append(materials, item)
		}
	}

	// Group materials by name to create stacks
	stackMap := make(map[string]*MaterialStack)
	for _, material := range materials {
		if stack, exists := stackMap[material.Name]; exists {
			stack.Quantity++
		} else {
			stackMap[material.Name] = &MaterialStack{
				Item:     material,
				Quantity: 1,
			}
		}
	}

	// Convert map to slice for consistent ordering
	msd.Stacks = []MaterialStack{}
	for _, stack := range stackMap {
		msd.Stacks = append(msd.Stacks, *stack)
	}

	msd.GW.Log.Println(logging.LogInfo, fmt.Sprintf("Created %d material stacks from %d individual items", len(msd.Stacks), len(materials)))

	// Assign alphabetical hotkeys (a, b, c, d, etc.) to stacks
	hotkey := 'a'
	options := []MenuBoxOption{}

	for _, stack := range msd.Stacks {
		if hotkey > 'z' {
			break // Don't go beyond 'z' for now
		}

		hotkeyStr := string(hotkey)
		msd.Hotkeys[hotkeyStr] = stack.Item

		// Create display name with symbol and quantity
		displayName := fmt.Sprintf("%s) %s x%d", hotkeyStr, stack.Item.Name, stack.Quantity)
		if stack.Item.Symbol != "" {
			displayName = fmt.Sprintf("%s) %s %s x%d", hotkeyStr, stack.Item.Symbol, stack.Item.Name, stack.Quantity)
		}

		// Add type info for debugging (can be removed later)
		displayName += fmt.Sprintf(" (Type: %v)", stack.Item.Type)

		options = append(options, MenuBoxOption{
			Name:    displayName,
			Keybind: hotkeyStr,
		})

		hotkey++
	}

	// If still no materials found, show a message
	if len(msd.Stacks) == 0 {
		options = append(options, MenuBoxOption{
			Name:    "No items found in inventory",
			Keybind: "",
		})
		msd.GW.SetStatusBarMessage("No items found in inventory. Press ! to close.")
	} else {
		msd.GW.SetStatusBarMessage("Select material or press ! to cancel")
	}

	msd.Options = options
	msd.Materials = materials

	// Set dimensions based on content
	msd.Width = 45 // Made wider to accommodate quantity info
	msd.Height = len(options) + 4
	if msd.Height < 8 {
		msd.Height = 8
	}

	// Position in center-left area
	msd.X = msd.GW.Width / 4
	msd.Y = msd.GW.Height/2 - msd.Height/2

	// Store the material hotkeys in the GameWindow for later access
	msd.GW.InventoryMutex.Lock()
	msd.GW.Hotkeys = msd.Hotkeys
	msd.GW.InventoryMutex.Unlock()
}

// GetMaterialForHotkey returns the material item for the given hotkey
func (gw *GameWindow) GetMaterialForHotkey(hotkey string) *edentypes.Item {
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()

	if item, exists := gw.Hotkeys[hotkey]; exists {
		return &item
	}
	return nil
}

// SetBuildWallMaterial stores the selected material for wall building
func (gw *GameWindow) SetBuildWallMaterial(item *edentypes.Item) {
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()
	// Store in a simple field - we'll add this to the GameWindow struct
	gw.buildWallMaterial = item
	gw.buildWallMode = true
	gw.Log.Println(logging.LogInfo, fmt.Sprintf("Stored build wall material: %s", item.Name))
}

// HandleBuildWallDirection handles direction input when in build wall mode
func (gw *GameWindow) HandleBuildWallDirection(input string) bool {
	// Check for cancel
	if input == "!" {
		gw.Log.Println(logging.LogInfo, "Build wall cancelled by user")
		gw.buildWallMode = false
		gw.buildWallMaterial = nil
		gw.SetStatusBarMessage("Build cancelled.")
		return true
	}

	// Check for valid direction
	validDirections := map[string]bool{
		"y": true, "u": true, "h": true, "j": true,
		"k": true, "l": true, "b": true, "n": true,
	}

	if validDirections[input] {
		gw.Log.Println(logging.LogInfo, fmt.Sprintf("Valid direction '%s' selected for wall building", input))

		// Get the stored material
		gw.InventoryMutex.Lock()
		material := gw.buildWallMaterial
		gw.InventoryMutex.Unlock()

		if material != nil {
			// Build the wall
			gw.BuildWallWithDirection(input, material)

			// Reset build wall mode
			gw.buildWallMode = false
			gw.buildWallMaterial = nil
		} else {
			gw.Log.Println(logging.LogError, "No material stored for wall building")
			gw.SetStatusBarMessage("Error: No material selected.")
			gw.buildWallMode = false
		}
		return true
	}

	// Invalid input - show message but stay in build wall mode
	gw.SetStatusBarMessage("Invalid direction. Use hjklyubn or ! to cancel.")
	return true
}

// DisplayMaterialSelection creates and displays the material selection popup after inventory is received
func (gw *GameWindow) DisplayMaterialSelection() {
	gw.InventoryMutex.Lock()
	callback := gw.MaterialSelectionCallback
	prompt := gw.MaterialSelectionPrompt
	gw.DisplayMaterialSelectionAfterReceive = false
	inventorySize := len(gw.Inventory)
	gw.InventoryMutex.Unlock()

	// Debug: Log that DisplayMaterialSelection was called
	gw.Log.Println(logging.LogInfo, "DisplayMaterialSelection called, inventory size:", inventorySize)

	// Close existing menus (like the build menu) before showing material selection
	gw.MenusMutex.Lock()
	currentMenus := make([]MenuBoxType, len(gw.Menus))
	copy(currentMenus, gw.Menus)
	gw.MenusMutex.Unlock()

	// Remove existing menus
	for _, menu := range currentMenus {
		gw.RemoveMenuBox(menu)
	}

	// Create material selection using standard MenuBox
	materialSelection := gw.CreateMaterialSelectionMenu(callback, prompt)

	gw.Log.Println(logging.LogInfo, "Material selection popup created and added to menus")
	gw.AddMenuBox(materialSelection)
}

// CreateMaterialSelectionMenu creates a standard MenuBox for material selection
func (gw *GameWindow) CreateMaterialSelectionMenu(callback func(*MenuBox, string), prompt string) *MenuBox {
	gw.Log.Println(logging.LogInfo, "CreateMaterialSelectionMenu called")

	// Get all building materials from inventory and group them by name
	materials := []edentypes.Item{}
	for _, item := range gw.Inventory {
		if item.Type == edentypes.ItemMaterial {
			materials = append(materials, item)
		}
	}

	gw.Log.Println(logging.LogInfo, fmt.Sprintf("Found %d materials in inventory", len(materials)))
	if len(materials) == 0 {
		gw.Log.Println(logging.LogInfo, "No materials found, showing all items for debugging")
		for _, item := range gw.Inventory {
			materials = append(materials, item)
		}
	}

	// Group materials by name to create stacks
	stackMap := make(map[string]*MaterialStack)
	for _, material := range materials {
		if stack, exists := stackMap[material.Name]; exists {
			stack.Quantity++
		} else {
			stackMap[material.Name] = &MaterialStack{
				Item:     material,
				Quantity: 1,
			}
		}
	}

	// Convert map to slice for consistent ordering
	stacks := []MaterialStack{}
	for _, stack := range stackMap {
		stacks = append(stacks, *stack)
	}

	gw.Log.Println(logging.LogInfo, fmt.Sprintf("Created %d material stacks", len(stacks)))

	// Create menu options and hotkey mapping
	hotkeys := make(map[string]edentypes.Item)
	options := []MenuBoxOption{}
	hotkey := 'a'

	for _, stack := range stacks {
		if hotkey > 'z' {
			break
		}

		hotkeyStr := string(hotkey)
		hotkeys[hotkeyStr] = stack.Item

		// Create display name with symbol and quantity
		displayName := fmt.Sprintf("%s x%d", stack.Item.Name, stack.Quantity)
		if stack.Item.Symbol != "" {
			displayName = fmt.Sprintf("%s %s x%d", stack.Item.Symbol, stack.Item.Name, stack.Quantity)
		}

		options = append(options, MenuBoxOption{
			Name:    displayName,
			Keybind: hotkeyStr,
		})

		hotkey++
	}

	if len(stacks) == 0 {
		options = append(options, MenuBoxOption{
			Name:    "No materials found in inventory",
			Keybind: "",
		})
	}

	// Store hotkeys for later access
	gw.InventoryMutex.Lock()
	gw.Hotkeys = hotkeys
	gw.InventoryMutex.Unlock()

	// Create the menu box
	menuBox := &MenuBox{
		X:                        gw.Width / 4,
		Y:                        gw.Height/2 - (len(options)+4)/2,
		Width:                    40,
		Height:                   len(options) + 4,
		Title:                    "Select Material",
		Options:                  options,
		CallbackStatusBarMessage: prompt,
		ResponseCallback:         callback,
		Type:                     MenuTypeInventory,
		CheckHotkeys:             true,
	}

	// Debug: Log callback information
	gw.Log.Println(logging.LogInfo, fmt.Sprintf("Created material selection menu with %d options", len(options)))
	if callback != nil {
		gw.Log.Println(logging.LogInfo, fmt.Sprintf("ResponseCallback set: %T", callback))
	} else {
		gw.Log.Println(logging.LogError, "ResponseCallback is nil!")
	}

	// Debug: Log the options to see if keybinds are set correctly
	for i, option := range options {
		gw.Log.Println(logging.LogInfo, fmt.Sprintf("Option %d: Name='%s', Keybind='%s'", i, option.Name, option.Keybind))
	}

	return menuBox
}

// ShowMaterialSelectionPopup creates and displays a custom material selection popup
func (gw *GameWindow) ShowMaterialSelectionPopup(prompt string, callback func(*MenuBox, string)) {
	// Debug: Log that ShowMaterialSelectionPopup was called
	gw.Log.Println(logging.LogInfo, "ShowMaterialSelectionPopup called with prompt:", prompt)

	// Store the callback for when inventory is received
	gw.InventoryMutex.Lock()
	gw.MaterialSelectionCallback = callback
	gw.MaterialSelectionPrompt = prompt
	gw.InventoryMutex.Unlock()

	// Request fresh inventory data
	gw.Log.Println(logging.LogInfo, "Requesting inventory update for material selection")
	gw.RequestInventoryUpdate(nil, "")
	gw.DisplayMaterialSelectionAfterReceive = true
}
