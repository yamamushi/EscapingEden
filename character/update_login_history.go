package character

import (
	"github.com/yamamushi/EscapingEden/messages"
	"time"
)

func (cm *CharacterManager) UpdateLoginHistory(info messages.CharacterInfo) messages.CMErrorType {
	charInfo := messages.CharacterInfo{}

	err := cm.DB.One("Characters", "ID", info.ID, &charInfo)
	if err != nil {
		return messages.CMError_DBError
	}

	charInfo.LastLoginTime = time.Now()
	charInfo.FirstLogin = info.FirstLogin

	err = cm.DB.UpdateRecord("Characters", &charInfo)
	if err != nil {
		return messages.CMError_DBError
	}

	return messages.CMError_Null
}
