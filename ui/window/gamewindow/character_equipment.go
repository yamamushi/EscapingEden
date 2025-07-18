package gamewindow

import (
	"fmt"
	"strings"

	"github.com/yamamushi/EscapingEden/ui/types"
)

type CharacterEquipmentDisplay struct {
	MenuBox
	GW *GameWindow
	// Navigation state
	selectedRow    int
	selectedCol    int
	equipmentSlots [][]EquipmentSlot // 2D array of equipment slots for navigation
}

// EquipmentSlot represents a single equipment slot with its metadata
type EquipmentSlot struct {
	slot        string
	displayName string
	itemID      string
	isEmpty     bool
}

func (eq *CharacterEquipmentDisplay) GetType() MenuType {
	return MenuTypeCharacterEquipment
}

func (eq *CharacterEquipmentDisplay) HandleInput(gw *GameWindow, inputType types.InputType, input string) {
	switch inputType {
	case types.InputEscape:
		// Close equipment window
		gw.CloseMenus = true
		return
	case types.InputCharacter:
		// Check if this is the close key
		if input == "!" {
			gw.CloseMenus = true
			return
		}
	case types.InputUp:
		eq.navigateUp()
		eq.updateStatusBar()
	case types.InputDown:
		eq.navigateDown()
		eq.updateStatusBar()
	case types.InputLeft:
		eq.navigateLeft()
		eq.updateStatusBar()
	case types.InputRight:
		eq.navigateRight()
		eq.updateStatusBar()
	}
}

func (eq *CharacterEquipmentDisplay) Draw(gw *GameWindow) {
	eq.Clear(gw)
	eq.DrawMenuItems(gw)
	eq.DrawBorder(gw)
	eq.DrawTitle(gw)
	eq.DrawPopupMenu(gw)
	// Draw custom status bar on the left side for equipment display
	eq.drawLeftStatusBar(gw)
}

// Custom status bar drawing for equipment display (positioned on the left)
func (eq *CharacterEquipmentDisplay) drawLeftStatusBar(gw *GameWindow) {
	gw.StatusBarMutex.Lock()
	defer gw.StatusBarMutex.Unlock()

	// Clear the status bar first
	gw.ClearStatusBar()

	// Split the status message into equipment info and help text
	if gw.StatusBarMessage != "" {
		// Look for the pattern " - Use arrows to navigate, ! to close"
		parts := strings.Split(gw.StatusBarMessage, " - Use arrows to navigate, ! to close")

		if len(parts) == 2 {
			// Equipment info on the left
			equipmentInfo := parts[0]
			gw.PrintStringToStatusBar(2, 0, equipmentInfo, gw.Terminal.Bold())

			// Help text on the right
			helpText := "Use arrows to navigate, ! to close"
			helpTextLen := len(helpText)
			gw.PrintStringToStatusBar(gw.Width-helpTextLen-2, 0, helpText, gw.Terminal.Bold())
		} else {
			// Fallback: if the message doesn't match expected pattern, show it on the left
			gw.PrintStringToStatusBar(2, 0, gw.StatusBarMessage, gw.Terminal.Bold())
		}
	}
}

func (eq *CharacterEquipmentDisplay) PrepareContent() {
	// Get character equipment
	equipment := eq.GW.CharacterInfo.Equipment

	// Define all equipment slots organized in groups
	bodySlots := []struct {
		slot        string
		displayName string
	}{
		{"head", "Head"},
		{"face", "Face"},
		{"neck", "Neck"},
		{"shoulders", "Shoulders"},
		{"back", "Back"},
		{"body", "Body"},
		{"waist", "Waist"},
	}

	armSlots := []struct {
		slot        string
		displayName string
	}{
		{"left_arm", "Left Arm"},
		{"right_arm", "Right Arm"},
		{"left_hand", "Left Hand"},
		{"right_hand", "Right Hand"},
	}

	legSlots := []struct {
		slot        string
		displayName string
	}{
		{"left_leg", "Left Leg"},
		{"right_leg", "Right Leg"},
		{"left_foot", "Left Foot"},
		{"right_foot", "Right Foot"},
	}

	leftRingSlots := []struct {
		slot        string
		displayName string
	}{
		{"left_index_ring", "L.Index"},
		{"left_middle_ring", "L.Middle"},
		{"left_ring_ring", "L.Ring"},
		{"left_pinky_ring", "L.Pinky"},
	}

	rightRingSlots := []struct {
		slot        string
		displayName string
	}{
		{"right_index_ring", "R.Index"},
		{"right_middle_ring", "R.Middle"},
		{"right_ring_ring", "R.Ring"},
		{"right_pinky_ring", "R.Pinky"},
	}

	// Build equipment slots grid for navigation
	maxRows := 7 // Body slots are the longest group
	eq.equipmentSlots = make([][]EquipmentSlot, maxRows)

	for i := 0; i < maxRows; i++ {
		eq.equipmentSlots[i] = make([]EquipmentSlot, 5) // 5 columns

		// Body column (0)
		if i < len(bodySlots) {
			itemID := equipment.GetEquippedItemID(bodySlots[i].slot)
			eq.equipmentSlots[i][0] = EquipmentSlot{
				slot:        bodySlots[i].slot,
				displayName: bodySlots[i].displayName,
				itemID:      itemID,
				isEmpty:     itemID == "",
			}
		}

		// Arms column (1)
		if i < len(armSlots) {
			itemID := equipment.GetEquippedItemID(armSlots[i].slot)
			eq.equipmentSlots[i][1] = EquipmentSlot{
				slot:        armSlots[i].slot,
				displayName: armSlots[i].displayName,
				itemID:      itemID,
				isEmpty:     itemID == "",
			}
		}

		// Legs column (2)
		if i < len(legSlots) {
			itemID := equipment.GetEquippedItemID(legSlots[i].slot)
			eq.equipmentSlots[i][2] = EquipmentSlot{
				slot:        legSlots[i].slot,
				displayName: legSlots[i].displayName,
				itemID:      itemID,
				isEmpty:     itemID == "",
			}
		}

		// Left rings column (3)
		if i < len(leftRingSlots) {
			itemID := equipment.GetEquippedItemID(leftRingSlots[i].slot)
			eq.equipmentSlots[i][3] = EquipmentSlot{
				slot:        leftRingSlots[i].slot,
				displayName: leftRingSlots[i].displayName,
				itemID:      itemID,
				isEmpty:     itemID == "",
			}
		}

		// Right rings column (4)
		if i < len(rightRingSlots) {
			itemID := equipment.GetEquippedItemID(rightRingSlots[i].slot)
			eq.equipmentSlots[i][4] = EquipmentSlot{
				slot:        rightRingSlots[i].slot,
				displayName: rightRingSlots[i].displayName,
				itemID:      itemID,
				isEmpty:     itemID == "",
			}
		}
	}

	// Create wide columnar layout for equipment
	content := []string{}

	// Build header row with proper spacing
	content = append(content, fmt.Sprintf("%-30s %-30s %-30s %-25s %-25s", "BODY", "ARMS", "LEGS", "LEFT RINGS", "RIGHT RINGS"))

	// Build rows with equipment from each category
	for i := 0; i < maxRows; i++ {
		var bodySlot, armSlot, legSlot, leftRingSlot, rightRingSlot string

		// Body column - limit to fit in 25 chars
		if i < len(bodySlots) {
			itemID := equipment.GetEquippedItemID(bodySlots[i].slot)
			item := eq.getItemDisplay(itemID)
			// Truncate if too long
			if len(item) > 12 {
				item = item[:9] + "..."
			}
			bodySlot = fmt.Sprintf("%-10s: %s", bodySlots[i].displayName, item)
		}

		// Arms column - limit to fit in 25 chars
		if i < len(armSlots) {
			itemID := equipment.GetEquippedItemID(armSlots[i].slot)
			item := eq.getItemDisplay(itemID)
			// Truncate if too long
			if len(item) > 12 {
				item = item[:9] + "..."
			}
			armSlot = fmt.Sprintf("%-10s: %s", armSlots[i].displayName, item)
		}

		// Legs column - limit to fit in 25 chars
		if i < len(legSlots) {
			itemID := equipment.GetEquippedItemID(legSlots[i].slot)
			item := eq.getItemDisplay(itemID)
			// Truncate if too long
			if len(item) > 12 {
				item = item[:9] + "..."
			}
			legSlot = fmt.Sprintf("%-10s: %s", legSlots[i].displayName, item)
		}

		// Left rings column - limit to fit in 20 chars
		if i < len(leftRingSlots) {
			itemID := equipment.GetEquippedItemID(leftRingSlots[i].slot)
			item := eq.getItemDisplay(itemID)
			// Truncate if too long
			if len(item) > 10 {
				item = item[:7] + "..."
			}
			leftRingSlot = fmt.Sprintf("%-8s: %s", leftRingSlots[i].displayName, item)
		}

		// Right rings column - limit to fit in 20 chars
		if i < len(rightRingSlots) {
			itemID := equipment.GetEquippedItemID(rightRingSlots[i].slot)
			item := eq.getItemDisplay(itemID)
			// Truncate if too long
			if len(item) > 10 {
				item = item[:7] + "..."
			}
			rightRingSlot = fmt.Sprintf("%-8s: %s", rightRingSlots[i].displayName, item)
		}

		// Only add row if at least one column has content
		if bodySlot != "" || armSlot != "" || legSlot != "" || leftRingSlot != "" || rightRingSlot != "" {
			line := fmt.Sprintf("%-30s %-30s %-30s %-25s %-25s", bodySlot, armSlot, legSlot, leftRingSlot, rightRingSlot)
			content = append(content, line)
		}
	}

	content = append(content, "")
	content = append(content, fmt.Sprintf("Total equipped: %d/23                                        Press ! to close", equipment.GetEquippedItemCount()))

	// Make the window very wide to accommodate columns
	eq.Width = 145 // Wide enough for 5 columns with increased spacing
	eq.Height = len(content) + 2

	// Position higher up on screen
	eq.X = 2
	eq.Y = 1

	// Convert content to menu options for display
	eq.Options = []MenuBoxOption{}
	for i, line := range content {
		eq.Options = append(eq.Options, MenuBoxOption{
			Name:     line,
			SkipDraw: false,
			Keybind:  fmt.Sprintf("%d", i),
		})
	}

	// Update status bar with current selection
	eq.updateStatusBar()
}

func (eq *CharacterEquipmentDisplay) DrawMenuItems(gw *GameWindow) {
	// Custom drawing for the wide columnar equipment display
	for index, option := range eq.Options {
		line := option.Name

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			eq.PrintToMenu(gw, 2, index+2, line, "")
			continue
		}

		// Column headers - draw in bold
		if strings.Contains(line, "BODY") && strings.Contains(line, "ARMS") && strings.Contains(line, "LEGS") {
			eq.PrintToMenu(gw, 2, index+2, line, gw.Terminal.Bold())
			continue
		}

		// Separator line - draw in dim color
		if strings.Contains(line, "─") {
			eq.PrintToMenu(gw, 2, index+2, line, "\033[90m")
			continue
		}

		// Check if this line contains equipment information
		if strings.Contains(line, ":") && (strings.Contains(line, "[") || strings.Contains(line, "(nothing)")) {
			// This line contains equipment slots - parse and colorize
			eq.drawEquipmentLine(gw, 2, index+2, line)
			continue
		}

		// Default drawing for other lines (like "Total equipped" and "Press ! to close")
		eq.PrintToMenu(gw, 2, index+2, line, "")
	}
}

// Helper function to draw equipment lines with proper coloring and highlighting
func (eq *CharacterEquipmentDisplay) drawEquipmentLine(gw *GameWindow, x, y int, line string) {
	// Calculate which row this line represents
	// The header is at index 0, so equipment rows start at index 1
	// y starts at 2 (window offset), so equipment line 0 is at y=3
	lineRow := y - 3 // Adjust for header row

	// Check if any slot in this row is selected
	selectedInThisRow := false
	selectedColInRow := -1
	if lineRow >= 0 && lineRow < len(eq.equipmentSlots) {
		for col := 0; col < 5; col++ {
			if eq.isValidSlot(lineRow, col) && lineRow == eq.selectedRow && col == eq.selectedCol {
				selectedInThisRow = true
				selectedColInRow = col
				break
			}
		}
	}

	if selectedInThisRow {
		// Highlight the selected equipment slot
		eq.drawHighlightedEquipmentLine(gw, x, y, line, selectedColInRow)
	} else {
		// Draw normally
		eq.PrintToMenu(gw, x, y, line, "")
	}
}

// Helper function to draw equipment line with highlighting for selected slot
func (eq *CharacterEquipmentDisplay) drawHighlightedEquipmentLine(gw *GameWindow, x, y int, line string, selectedCol int) {
	// Column boundaries for highlighting
	columnStarts := []int{0, 30, 60, 90, 115}
	columnWidths := []int{30, 30, 30, 25, 25}

	if selectedCol >= 0 && selectedCol < len(columnStarts) {
		// Draw the line in parts: before highlight, highlight, after highlight
		start := columnStarts[selectedCol]
		width := columnWidths[selectedCol]

		// Draw part before highlight
		if start > 0 {
			beforePart := line[:min(start, len(line))]
			eq.PrintToMenu(gw, x, y, beforePart, "")
		}

		// Draw highlighted part
		if start < len(line) {
			end := min(start+width, len(line))
			highlightPart := line[start:end]
			// Use reverse video for highlighting
			eq.PrintToMenu(gw, x+start, y, highlightPart, "\033[7m") // Reverse video
		}

		// Draw part after highlight
		afterStart := start + width
		if afterStart < len(line) {
			afterPart := line[afterStart:]
			eq.PrintToMenu(gw, x+afterStart, y, afterPart, "")
		}
	} else {
		// Fallback to normal drawing
		eq.PrintToMenu(gw, x, y, line, "")
	}
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function to format item display
func (eq *CharacterEquipmentDisplay) getItemDisplay(itemID string) string {
	if itemID != "" {
		return eq.getItemName(itemID)
	}
	return "(nothing)"
}

// Helper function to get item name from ID using character inventory
func (eq *CharacterEquipmentDisplay) getItemName(itemID string) string {
	if itemID == "" {
		return "(nothing)"
	}

	// Look up the item in the character's inventory
	// The character's inventory contains all items the character owns, including equipped ones
	for _, item := range eq.GW.CharacterInfo.Inventory {
		if item.ID == itemID {
			return item.Name
		}
	}

	// If not found in inventory, return a fallback
	return "Unknown Item"
}

// Implement missing MenuBoxType interface methods
func (eq *CharacterEquipmentDisplay) CloseMenus(gw *GameWindow) {
	gw.CloseMenus = true
}

func (eq *CharacterEquipmentDisplay) SetCallbackStatusBarMessage(message string) {
	// Not used for equipment display, but required by interface
}

func (eq *CharacterEquipmentDisplay) GetCallbackDataString() string {
	// Not used for equipment display, but required by interface
	return ""
}

func (eq *CharacterEquipmentDisplay) ToggleHotkeyCheck(toggle bool) {
	// Not used for equipment display, but required by interface
}

// Navigation methods
func (eq *CharacterEquipmentDisplay) navigateUp() {
	if eq.selectedRow > 0 {
		// Find the next valid slot above
		for row := eq.selectedRow - 1; row >= 0; row-- {
			if eq.isValidSlot(row, eq.selectedCol) {
				eq.selectedRow = row
				break
			}
		}
	}
}

func (eq *CharacterEquipmentDisplay) navigateDown() {
	if eq.selectedRow < len(eq.equipmentSlots)-1 {
		// Find the next valid slot below
		for row := eq.selectedRow + 1; row < len(eq.equipmentSlots); row++ {
			if eq.isValidSlot(row, eq.selectedCol) {
				eq.selectedRow = row
				break
			}
		}
	}
}

func (eq *CharacterEquipmentDisplay) navigateLeft() {
	// Try to find a slot in the same row to the left
	for col := eq.selectedCol - 1; col >= 0; col-- {
		if eq.isValidSlot(eq.selectedRow, col) {
			eq.selectedCol = col
			return
		}
	}

	// If no slot found in same row, try to find the rightmost slot in the previous column
	// that has content, starting from the bottom and working up
	for col := eq.selectedCol - 1; col >= 0; col-- {
		// Find the last valid slot in this column
		for row := len(eq.equipmentSlots) - 1; row >= 0; row-- {
			if eq.isValidSlot(row, col) {
				eq.selectedRow = row
				eq.selectedCol = col
				return
			}
		}
	}

	// If we're at the leftmost column, wrap around to the rightmost column
	// Try to find a slot in the same row starting from the rightmost column
	for col := 4; col > eq.selectedCol; col-- {
		if eq.isValidSlot(eq.selectedRow, col) {
			eq.selectedCol = col
			return
		}
	}

	// If no slot found in same row, find the last valid slot in the rightmost columns
	for col := 4; col > eq.selectedCol; col-- {
		// Find the last valid slot in this column
		for row := len(eq.equipmentSlots) - 1; row >= 0; row-- {
			if eq.isValidSlot(row, col) {
				eq.selectedRow = row
				eq.selectedCol = col
				return
			}
		}
	}
}

func (eq *CharacterEquipmentDisplay) navigateRight() {
	// Try to find a slot in the same row to the right
	for col := eq.selectedCol + 1; col < 5; col++ {
		if eq.isValidSlot(eq.selectedRow, col) {
			eq.selectedCol = col
			return
		}
	}

	// If no slot found in same row, find the closest slot in the next column
	for col := eq.selectedCol + 1; col < 5; col++ {
		// First try to find a slot at or below the current row
		for row := eq.selectedRow; row < len(eq.equipmentSlots); row++ {
			if eq.isValidSlot(row, col) {
				eq.selectedRow = row
				eq.selectedCol = col
				return
			}
		}
		// If no slot found at or below, try above the current row
		for row := eq.selectedRow - 1; row >= 0; row-- {
			if eq.isValidSlot(row, col) {
				eq.selectedRow = row
				eq.selectedCol = col
				return
			}
		}
	}

	// If we're at the rightmost column, wrap around to the leftmost column
	// Try to find a slot in the same row starting from the leftmost column
	for col := 0; col < eq.selectedCol; col++ {
		if eq.isValidSlot(eq.selectedRow, col) {
			eq.selectedCol = col
			return
		}
	}

	// If no slot found in same row, find the first valid slot in the leftmost columns
	for col := 0; col < eq.selectedCol; col++ {
		// Find the first valid slot in this column
		for row := 0; row < len(eq.equipmentSlots); row++ {
			if eq.isValidSlot(row, col) {
				eq.selectedRow = row
				eq.selectedCol = col
				return
			}
		}
	}
}

// Check if a slot position is valid (has equipment slot data)
func (eq *CharacterEquipmentDisplay) isValidSlot(row, col int) bool {
	if row < 0 || row >= len(eq.equipmentSlots) || col < 0 || col >= 5 {
		return false
	}
	return eq.equipmentSlots[row][col].slot != ""
}

// Update status bar with current selection information
func (eq *CharacterEquipmentDisplay) updateStatusBar() {
	if !eq.isValidSlot(eq.selectedRow, eq.selectedCol) {
		eq.GW.SetStatusBarMessage("Press ! to close character equipment")
		return
	}

	currentSlot := eq.equipmentSlots[eq.selectedRow][eq.selectedCol]

	var statusMessage string
	if currentSlot.isEmpty {
		statusMessage = fmt.Sprintf("%s: (empty) - Use arrows to navigate, ! to close", currentSlot.displayName)
	} else {
		// Show the item name instead of ID
		itemName := eq.getItemName(currentSlot.itemID)
		statusMessage = fmt.Sprintf("%s: %s - Use arrows to navigate, ! to close", currentSlot.displayName, itemName)
	}

	eq.GW.SetStatusBarMessage(statusMessage)
}

// DisplayCharacterEquipment creates and displays the character equipment window
func (gw *GameWindow) DisplayCharacterEquipment() {
	equipmentWindow := CharacterEquipmentDisplay{
		GW:          gw,
		selectedRow: 0,
		selectedCol: 0,
	}
	equipmentWindow.Type = MenuTypeCharacterEquipment
	equipmentWindow.Title = "Character Equipment"
	equipmentWindow.PrepareContent()

	gw.AddMenuBox(&equipmentWindow)
}
