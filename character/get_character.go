package character

import "github.com/yamamushi/EscapingEden/messages"

func (cm *CharacterManager) GetCharacter(ID string) (messages.CharacterInfo, messages.CMErrorType) {
	charInfo := messages.CharacterInfo{}

	err := cm.DB.One("Characters", "ID", ID, &charInfo)
	if err != nil {
		return charInfo, messages.CMError_DBError
	}

	return charInfo, messages.CMError_Null
}
