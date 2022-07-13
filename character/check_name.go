package character

import (
	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/messages"
)

// Checks to see if a character with the given name exists in the database

func (cm *CharacterManager) CheckCharNameInUse(name string) (bool, error) {

	// Get Character
	// Look for character in db with the same name
	result := messages.CharacterInfo{}
	err := cm.DB.One("Characters", "Name", name, &result)
	if err != nil {
		if err.Error() == "not found" {
			return false, messages.CMError_Null
		}
		return true, messages.CMError_DBError
	} else {
		return true, messages.CMError_NameAlreadyExists
	}
}

func (cm *CharacterManager) CheckCharNameAllowed(name string) (bool, messages.CMErrorType) {

	if edenutil.CheckBlacklist(name, edenutil.BlackListUsernames) {
		return false, messages.CMError_InvalidName
	}
	return true, messages.CMError_Null // Name is allowed
}
