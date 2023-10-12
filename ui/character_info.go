package ui

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (c *Console) UpdateCharacterInfo(charInfo messages.CharacterInfo) {
	c.characterInfoMutex.Lock()
	defer c.characterInfoMutex.Unlock()
	c.CharacterInfo = charInfo
	c.Log.Println(logging.LogInfo, "Updated character info for ", c.CharacterInfo.ID)
	c.UpdateWindowsCharacterInfo()
	c.Log.Println(logging.LogInfo, "Updated windows character info")
}

func (c *Console) UpdateWindowsCharacterInfo() {
	for _, w := range c.Windows {
		w.SetCharacterInfo(c.CharacterInfo)
	}
}

func (c *Console) GetCharacterName() string {
	c.userInfoMutex.Lock()
	defer c.userInfoMutex.Unlock()
	return c.CharacterInfo.Name
}
