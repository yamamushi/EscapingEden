package edentypes

import "github.com/yamamushi/EscapingEden/ui/util"

type Item struct {
	ID          string
	Name        string
	Description string
	Weight      float64
	Type        ItemType
	Stackable   bool
	Hotkey      string
	Attributes  map[string]bool

	// Visual properties
	Symbol  string         `json:"symbol"`   // Character/symbol when dropped on ground
	FGColor util.ColorCode `json:"fg_color"` // Foreground color
	BGColor util.ColorCode `json:"bg_color"` // Background color

	// Extended properties
	Value      int    `json:"value"`      // Economic value
	Durability int    `json:"durability"` // For tools/equipment (-1 = infinite)
	MaxStack   int    `json:"max_stack"`  // Maximum stack size (0 = not stackable)
	Rarity     string `json:"rarity"`     // common, uncommon, rare, epic, legendary

	// Metadata
	Category string   `json:"category"` // More specific than Type (e.g., "wood", "stone", "metal")
	Tags     []string `json:"tags"`     // Searchable tags
}

type ItemType int

const (
	ItemTypeNull ItemType = iota
	ItemMaterial
	ItemTool
)

// Print item type as a string
func (itemType ItemType) String() string {
	return [...]string{"Null", "Material", "Tool"}[itemType]
}

func GetInventoryWeight(inventory []Item) float64 {

	weight := 0.0
	for _, item := range inventory {
		weight += item.Weight
	}
	return weight

	return 0
}

func GetInventoryHotkeyMap(inventory []Item) map[string]Item {
	hotkeyMap := make(map[string]Item)
	for _, item := range inventory {
		hotkeyMap[item.Hotkey] = item // this will ensure that the last item with a given hotkey is the one that is used in a stack
	}
	return hotkeyMap
}
