package accounts

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (am *AccountManager) handleLogin(username, password string) (response messages.AccountLoginResponse) {

	response = messages.AccountLoginResponse{}
	account := messages.Account{}

	err := am.DB.One("Accounts", "Username", username, &account)
	if err != nil {
		am.Log.Println(logging.LogWarn, "Attempted login to nonexistent account:", username)
		response.Error = messages.AMError_AccountDoesNotExist
		return response
	}

	success := am.ComparePasswords(account.HashedPassword, password)
	if !success {
		am.Log.Println(logging.LogWarn, "Attempted login with incorrect password for account:", username)
		response.Error = messages.AMError_InvalidPassword
		return response
	}

	if account.ValidationStatus == 0 {
		am.Log.Println(logging.LogWarn, "Attempted login to account that is pending validation:", username)
		response.Error = messages.AMError_PendingValidation
		return response
	}

	account.HashedPassword = ""
	response.Account = account

	return response
}
