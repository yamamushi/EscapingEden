package accounts

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (am *AccountManager) AddCharacter(userID string, characterID string) (messages.Account, messages.AMErrorType) {
	userAccount := messages.Account{}
	err := am.DB.One("Accounts", "ID", userID, &userAccount)
	if err != nil {
		am.Log.Println(logging.LogError, "Error getting user account:", err)
		return userAccount, messages.AMError_DBError
	}

	userAccount.Characters = append(userAccount.Characters, characterID)
	userAccount.LastCharacterID = characterID

	err = am.DB.UpdateRecord("Accounts", userAccount)
	if err != nil {
		am.Log.Println(logging.LogError, "Error updating user account:", err)
		return userAccount, messages.AMError_DBError
	}

	return userAccount, messages.AMError_Null
}
