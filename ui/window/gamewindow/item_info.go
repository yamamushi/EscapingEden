package gamewindow

import (
	"fmt"
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

	// Item name and basic info
	info.Content = append(info.Content, fmt.Sprintf("Name: %s", info.Item.Name))
	info.Content = append(info.Content, fmt.Sprintf("Type: %s", info.Item.Type.String()))
	info.Content = append(info.Content, fmt.Sprintf("Category: %s", info.Item.Category))
	info.Content = append(info.Content, "")

	// Description (word wrap)
	wrappedDesc := info.wrapText(info.Item.Description, 45)
	info.Content = append(info.Content, "Description:")
	for _, line := range wrappedDesc {
		info.Content = append(info.Content, "  "+line)
	}
	info.Content = append(info.Content, "")

	// Physical properties
	info.Content = append(info.Content, "Physical Properties:")
	info.Content = append(info.Content, fmt.Sprintf("  Weight: %.2f kg", info.Item.Weight))
	if info.Item.Symbol != "" {
		info.Content = append(info.Content, fmt.Sprintf("  Symbol: %s", info.Item.Symbol))
	}
	// Show foreground color
	if info.Item.FGColor.R != 0 || info.Item.FGColor.G != 0 || info.Item.FGColor.B != 0 {
		info.Content = append(info.Content, fmt.Sprintf("  FG Color: RGB(%d, %d, %d)", info.Item.FGColor.R, info.Item.FGColor.G, info.Item.FGColor.B))
	}
	// Show background color if not black
	if info.Item.BGColor.R != 0 || info.Item.BGColor.G != 0 || info.Item.BGColor.B != 0 {
		info.Content = append(info.Content, fmt.Sprintf("  BG Color: RGB(%d, %d, %d)", info.Item.BGColor.R, info.Item.BGColor.G, info.Item.BGColor.B))
	}
	info.Content = append(info.Content, "")

	// Game properties
	info.Content = append(info.Content, "Game Properties:")
	if info.Item.Stackable {
		if info.Item.MaxStack > 0 {
			info.Content = append(info.Content, fmt.Sprintf("  Stackable: Yes (max %d)", info.Item.MaxStack))
		} else {
			info.Content = append(info.Content, "  Stackable: Yes")
		}
	} else {
		info.Content = append(info.Content, "  Stackable: No")
	}

	if info.Item.Value > 0 {
		info.Content = append(info.Content, fmt.Sprintf("  Value: %d coins", info.Item.Value))
	}

	if info.Item.Durability > 0 {
		info.Content = append(info.Content, fmt.Sprintf("  Durability: %d", info.Item.Durability))
	} else if info.Item.Durability == -1 {
		info.Content = append(info.Content, "  Durability: Infinite")
	}

	if info.Item.Rarity != "" {
		info.Content = append(info.Content, fmt.Sprintf("  Rarity: %s", strings.Title(info.Item.Rarity)))
	}
	info.Content = append(info.Content, "")

	// Attributes and tags
	if len(info.Item.Attributes) > 0 {
		info.Content = append(info.Content, "Attributes:")
		for attr, enabled := range info.Item.Attributes {
			if enabled {
				info.Content = append(info.Content, fmt.Sprintf("  • %s", attr))
			}
		}
		info.Content = append(info.Content, "")
	}

	if len(info.Item.Tags) > 0 {
		info.Content = append(info.Content, "Tags:")
		tagLine := "  " + strings.Join(info.Item.Tags, ", ")
		wrappedTags := info.wrapText(tagLine, 47)
		for _, line := range wrappedTags {
			info.Content = append(info.Content, line)
		}
		info.Content = append(info.Content, "")
	}

	// Controls
	info.Content = append(info.Content, "Controls:")
	info.Content = append(info.Content, "  ! - Close")

	// Calculate dimensions
	maxWidth := 0
	for _, line := range info.Content {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	info.Width = maxWidth + 4
	info.Height = len(info.Content) + 4

	// Center the window
	info.X = (info.GW.Width - info.Width) / 2
	info.Y = (info.GW.Height - info.Height) / 2
}

func (info *ItemInfoDisplay) DrawMenuItems(gw *GameWindow) {
	for index, content := range info.Content {
		// Special handling for the symbol line in Physical Properties section
		if info.Item.Symbol != "" && strings.HasPrefix(content, "  Symbol: ") && strings.Contains(content, info.Item.Symbol) {
			// This is the line with the item symbol - draw it with colors
			symbolIndex := strings.Index(content, info.Item.Symbol)

			// Draw the part before the symbol
			beforeSymbol := content[:symbolIndex]
			info.PrintToMenu(gw, 2, index+2, beforeSymbol, "")

			// Draw the colored symbol
			itemColorCode := info.Item.FGColor.FG() + info.Item.BGColor.BG()
			info.PrintToMenu(gw, 2+len(beforeSymbol), index+2, info.Item.Symbol, itemColorCode)

			// Draw any remaining part after the symbol
			afterSymbol := content[symbolIndex+len(info.Item.Symbol):]
			info.PrintToMenu(gw, 2+len(beforeSymbol)+len(info.Item.Symbol), index+2, afterSymbol, "")
		} else {
			// Normal drawing for all other lines
			info.PrintToMenu(gw, 2, index+2, content, "")
		}
	}
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
