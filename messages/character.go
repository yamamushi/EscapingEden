package messages

type CharacterInfo struct {
	ID          string `storm:"index"`
	Name        string `storm:"unique"`
	InventoryID string `storm:"unique"`
}
