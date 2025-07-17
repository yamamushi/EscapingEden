package gamewindow

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/ui/types"
)

type ItemInfoDisplay struct {
	MenuBox
	Item    edentypes.Item
	Content []string
	GW      *GameWindow
}

func (info *ItemInfoDisplay) GetType() MenuType {
	return MenuTypeItemInfo
}

// formatItemColor converts integer color codes to escape sequences (same format as tiles)
func formatItemColorInfo(fgColor, bgColor int) string {
	fg := strconv.Itoa(fgColor)
	bg := strconv.Itoa(bgColor)
	return "\033[38;5;" + fg + "m" + "\033[48;5;" + bg + "m"
}

func (info *ItemInfoDisplay) HandleInput(gw *GameWindow, inputType types.InputType, input string) {
	switch inputType {
	case types.InputEscape:
		// Close the item info popup by removing it from the menu stack (using unsafe version since mutex is already held by input handler)
		gw.RemoveMenuBoxUnsafe(info)
		return
	case types.InputCharacter:
		// Close on ! key
		if input == "!" {
			gw.RemoveMenuBoxUnsafe(info)
			return
		}
	}
}

func (info *ItemInfoDisplay) Draw(gw *GameWindow) {
	info.Clear(gw)
	info.DrawMenuItems(gw)
	info.DrawBorder(gw)
	info.DrawTitle(gw)
	info.DrawPopupMenu(gw)
	gw.SetStatusBarMessage("Press ! to close item info")
}

func (info *ItemInfoDisplay) PrepareContent() {
	info.Content = []string{}

	// Header with name and symbol
	headerLine := fmt.Sprintf("Name: %s", info.Item.Name)
	if info.Item.Symbol != "" {
		headerLine += fmt.Sprintf("  Symbol: %s", info.Item.Symbol)
	}
	info.Content = append(info.Content, headerLine)

	// Basic info line - Type, Category, Weight
	basicLine := fmt.Sprintf("Type: %s  Category: %s  Weight: %.2f kg",
		info.Item.Type.String(), info.Item.Category, info.Item.Weight)
	info.Content = append(info.Content, basicLine)

	// Game properties line
	var gameProps []string

	// Show stackable status
	if info.Item.Stackable {
		if info.Item.MaxStack > 0 {
			gameProps = append(gameProps, fmt.Sprintf("Stackable: yes (max %d)", info.Item.MaxStack))
		} else {
			gameProps = append(gameProps, "Stackable: yes")
		}
	} else {
		gameProps = append(gameProps, "Stackable: no")
	}

	// Show equippable status
	if info.Item.Equippable {
		gameProps = append(gameProps, "Equippable: yes")
	} else {
		gameProps = append(gameProps, "Equippable: no")
	}

	// Only show durability if it's a positive value (not infinite/unset)
	if info.Item.Durability > 0 {
		gameProps = append(gameProps, fmt.Sprintf("Durability: %d", info.Item.Durability))
	}

	if len(gameProps) > 0 {
		info.Content = append(info.Content, strings.Join(gameProps, "  "))
	}

	info.Content = append(info.Content, "")

	// Description (word wrap to wider width for notecard style)
	if info.Item.Description != "" {
		wrappedDesc := info.wrapText(info.Item.Description, 70)
		info.Content = append(info.Content, "Description:")
		for _, line := range wrappedDesc {
			info.Content = append(info.Content, "  "+line)
		}
		info.Content = append(info.Content, "")
	}

	// Attributes and tags on same line if they exist
	var extraInfo []string
	if len(info.Item.Attributes) > 0 {
		var attrs []string
		for attr, enabled := range info.Item.Attributes {
			if enabled {
				attrs = append(attrs, attr)
			}
		}
		if len(attrs) > 0 {
			extraInfo = append(extraInfo, "Attributes: "+strings.Join(attrs, ", "))
		}
	}

	if len(info.Item.Tags) > 0 {
		extraInfo = append(extraInfo, "Tags: "+strings.Join(info.Item.Tags, ", "))
	}

	for _, line := range extraInfo {
		wrappedLine := info.wrapText(line, 70)
		for _, wrapped := range wrappedLine {
			info.Content = append(info.Content, wrapped)
		}
	}

	if len(extraInfo) > 0 {
		info.Content = append(info.Content, "")
	}

	// Controls
	info.Content = append(info.Content, "Press ! to close")

	// Calculate dimensions for notecard style (wider, shorter)
	maxWidth := 0
	for _, line := range info.Content {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Ensure minimum width for notecard appearance
	if maxWidth < 60 {
		maxWidth = 60
	}

	info.Width = maxWidth + 4
	info.Height = len(info.Content) + 4

	// Position in upper left area to avoid overlapping inventory
	info.X = 5
	info.Y = 3
}

func (info *ItemInfoDisplay) DrawMenuItems(gw *GameWindow) {
	for index, content := range info.Content {
		// Special handling for the header line that contains the symbol
		if info.Item.Symbol != "" && strings.Contains(content, "Symbol: "+info.Item.Symbol) {
			// This is the header line with the item symbol - draw it with colors and bold field names
			info.drawLineWithBoldFields(gw, index, content, true)
		} else {
			// Draw other lines with bold field names
			info.drawLineWithBoldFields(gw, index, content, false)
		}
	}
}

// drawLineWithBoldFields draws a line with field names in bold
func (info *ItemInfoDisplay) drawLineWithBoldFields(gw *GameWindow, lineIndex int, content string, hasSymbol bool) {
	// Field names that should be bold
	fieldNames := []string{"Name:", "Type:", "Category:", "Weight:", "Symbol:", "Durability:", "Stackable", "Equippable", "Description:", "Attributes:", "Tags:"}

	x := 2
	y := lineIndex + 2

	// If this line has a symbol, handle it specially
	if hasSymbol && info.Item.Symbol != "" && strings.Contains(content, "Symbol: "+info.Item.Symbol) {
		symbolIndex := strings.Index(content, info.Item.Symbol)

		// Draw everything before the symbol (with bold field names)
		beforeSymbol := content[:symbolIndex]
		x = info.drawTextWithBoldFields(gw, x, y, beforeSymbol, fieldNames)

		// Draw the colored symbol
		itemColorCode := formatItemColorInfo(info.Item.FGColor, info.Item.BGColor)
		info.PrintToMenu(gw, x, y, info.Item.Symbol, itemColorCode)
		x += len([]rune(info.Item.Symbol))

		// Draw everything after the symbol
		afterSymbol := content[symbolIndex+len(info.Item.Symbol):]
		info.drawTextWithBoldFields(gw, x, y, afterSymbol, fieldNames)
	} else {
		// Normal line - just draw with bold field names
		info.drawTextWithBoldFields(gw, x, y, content, fieldNames)
	}
}

// drawTextWithBoldFields draws text with specified field names in bold
func (info *ItemInfoDisplay) drawTextWithBoldFields(gw *GameWindow, startX, y int, text string, fieldNames []string) int {
	if text == "" {
		return startX
	}

	x := startX
	remaining := text

	for len(remaining) > 0 {
		// Find the next field name
		nextFieldIndex := -1
		nextFieldName := ""

		for _, field := range fieldNames {
			if index := strings.Index(remaining, field); index != -1 {
				if nextFieldIndex == -1 || index < nextFieldIndex {
					nextFieldIndex = index
					nextFieldName = field
				}
			}
		}

		if nextFieldIndex == -1 {
			// No more field names, draw the rest normally
			info.PrintToMenu(gw, x, y, remaining, "")
			x += len([]rune(remaining))
			break
		}

		// Draw text before the field name
		if nextFieldIndex > 0 {
			beforeField := remaining[:nextFieldIndex]
			info.PrintToMenu(gw, x, y, beforeField, "")
			x += len([]rune(beforeField))
		}

		// Draw the field name in bold
		info.PrintToMenu(gw, x, y, nextFieldName, gw.Terminal.Bold())
		x += len([]rune(nextFieldName))

		// Continue with the rest of the text
		remaining = remaining[nextFieldIndex+len(nextFieldName):]
	}

	return x
}

// Implement missing MenuBoxType interface methods
func (info *ItemInfoDisplay) CloseMenus(gw *GameWindow) {
	// Close this item info popup
	gw.CloseMenus = true
}

func (info *ItemInfoDisplay) SetCallbackStatusBarMessage(message string) {
	// Not used for item info, but required by interface
}

func (info *ItemInfoDisplay) GetCallbackDataString() string {
	// Not used for item info, but required by interface
	return ""
}

func (info *ItemInfoDisplay) ToggleHotkeyCheck(toggle bool) {
	// Not used for item info, but required by interface
}

// wrapText wraps text to fit within the specified width
func (info *ItemInfoDisplay) wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	var lines []string
	currentLine := ""

	for _, word := range words {
		if len(currentLine) == 0 {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}
