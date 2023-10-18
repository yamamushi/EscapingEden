package edenitems

type Item struct {
	ID          string
	Name        string
	Description string
	Weight      float64
	Type        ItemType
	Stackable   bool
	Hotkey      string
}

type ItemType int

const (
	ItemUnknown ItemType = iota
	ItemMaterial
	ItemTool
)

func GetInventoryWeight(inventory []Item) float64 {

	weight := 0.0
	for _, item := range inventory {
		weight += item.Weight
	}
	return weight

	return 0
}
