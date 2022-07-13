package character

import "github.com/yamamushi/EscapingEden/messages"

func (cm *CharacterManager) GetCharacterByID(ID string) messages.CharacterInfo {
	charInfo := messages.CharacterInfo{}

	err := cm.DB.One("Characters", "ID", ID, &charInfo)
	if err != nil {
		charInfo.Error = messages.CMError_DBError.Error()
		return charInfo
	}

	charInfo.Error = messages.CMError_Null.Error()
	return charInfo
}
