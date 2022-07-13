package character

import (
	"github.com/yamamushi/EscapingEden/messages"
	"time"
)

func (cm *CharacterManager) UpdateLoginHistory(info messages.CharacterInfo) messages.CharacterInfo {
	charInfo := messages.CharacterInfo{}

	err := cm.DB.One("Characters", "ID", info.ID, &charInfo)
	if err != nil {
		charInfo.Error = messages.CMError_DBError.Error()
		return charInfo
	}

	if charInfo.UserID != info.UserID {
		charInfo.Error = messages.CMError_HistoryUpdatePermissionError.Error()
		return charInfo
	}

	charInfo.LastLoginTime = time.Now()
	charInfo.FirstLogin = info.FirstLogin

	err = cm.DB.UpdateRecord("Characters", &charInfo)
	if err != nil {
		charInfo.Error = messages.CMError_DBError.Error()
		return charInfo
	}

	charInfo.Error = messages.CMError_Null.Error()
	return charInfo
}
