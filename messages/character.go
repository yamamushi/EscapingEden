package messages

import (
	"github.com/yamamushi/EscapingEden/ui/util"
	"time"
)

type CharacterInfo struct {
	ID            string `storm:"index"`
	UserID        string
	Name          string         `storm:"unique"`
	FGColor       util.ColorCode // The escape code of the FG color of the character
	BGColor       util.ColorCode
	InventoryID   string `storm:"unique"`
	LastLoginTime time.Time
	FirstLogin    int // 0 = false, 1 = true
	Error         string
}

func (c *CharacterInfo) GetID() string {
	return c.ID
}

func (c *CharacterInfo) GetColorFG() util.ColorCode {
	return c.FGColor
}

func (c *CharacterInfo) GetColorBG() util.ColorCode {
	return c.BGColor
}

func (c *CharacterInfo) GetName() string {
	return c.Name
}

func (c *CharacterInfo) GetInventoryID() string {
	return c.InventoryID
}
