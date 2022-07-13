package character

import (
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/messages"
	"log"
	"time"
)

func (cm *CharacterManager) CreateCharacter(info messages.CharacterInfo) messages.CharacterInfo {

	charID := uuid.New().String()

	newCharInfo := messages.CharacterInfo{
		ID:            charID,
		UserID:        info.UserID,
		Name:          info.Name,
		FGColor:       info.FGColor,
		BGColor:       info.BGColor,
		InventoryID:   "",          // TODO - Generate inventory ID
		LastLoginTime: time.Time{}, // Not yet initialized, we update this when we also update the account login history details.
		FirstLogin:    1,           // 1 = true, 0 = false
	}

	err := cm.DB.AddRecord("Characters", &newCharInfo)
	if err != nil {
		log.Println(err.Error())
		newCharInfo.Error = messages.CMError_DBError.Error()
		return newCharInfo
	}

	newCharInfo.Error = messages.CMError_Null.Error()
	return newCharInfo
}
