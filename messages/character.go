package messages

import "time"

type CharacterInfo struct {
	ID            string `storm:"index"`
	Name          string `storm:"unique"`
	Color         string // The escape code of the color of the character, can include both FG and BG
	InventoryID   string `storm:"unique"`
	LastLoginTime time.Time
}

func (c *CharacterInfo) GetID() string {
	return c.ID
}

func (c *CharacterInfo) GetColor() string {
	return c.Color
}

func (c *CharacterInfo) GetName() string {
	return c.Name
}

func (c *CharacterInfo) GetInventoryID() string {
	return c.InventoryID
}
