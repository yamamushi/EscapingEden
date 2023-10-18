package edenitems

type Item struct {
	ID          int
	Name        string
	Description string
	Weight      float64
	Type        ItemType
	Stackable   bool
}

type ItemType int

const (
	ItemUnknown ItemType = iota
	ItemMaterial
	ItemTool
)
