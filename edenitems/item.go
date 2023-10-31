package edenitems

type Item struct {
	ID          string
	Name        string
	Description string
	Weight      float64
	Type        ItemType
	Stackable   bool
	Hotkey      string
	Attributes  map[string]bool
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
