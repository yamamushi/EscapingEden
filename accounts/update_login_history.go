package accounts

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (am *AccountManager) UpdateLoginCharacterHistory(info messages.CharacterInfo) messages.Account {
	userAccount := messages.Account{}
	err := am.DB.One("Accounts", "ID", info.UserID, &userAccount)
	if err != nil {
		am.Log.Println(logging.LogError, "Error getting user account:", err)
		userAccount.Error = messages.AMError_DBError.Error()
		return userAccount
	}

	userAccount.LastCharacterName = info.Name
	userAccount.LastCharacterID = info.ID
	err = am.DB.UpdateRecord("Accounts", &userAccount)
	if err != nil {
		am.Log.Println(logging.LogError, "Error updating user account:", err)
		userAccount.Error = messages.AMError_DBError.Error()
		return userAccount
	}
	return userAccount
}
