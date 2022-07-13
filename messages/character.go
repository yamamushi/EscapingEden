package messages

type CharacterInfo struct {
	ID          string `storm:"index"`
	Name        string `storm:"unique"`
	InventoryID string `storm:"unique"`
}

func (c *CharacterInfo) GetID() string {
	return c.ID
}

func (c *CharacterInfo) GetName() string {
	return c.Name
}

func (c *CharacterInfo) GetInventoryID() string {
	return c.InventoryID
}
