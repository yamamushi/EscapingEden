package accounts

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (am *AccountManager) handleLogin(email, password string) (response messages.AccountLoginResponse) {

	response = messages.AccountLoginResponse{}
	account := messages.Account{}

	err := am.DB.One("Accounts", "Email", email, &account)
	if err != nil {
		am.Log.Println(logging.LogWarn, "Attempted login to nonexistent account:", email)
		response.Error = messages.AMError_AccountDoesNotExist
		return response
	}

	success := am.ComparePasswords(account.HashedPassword, password)
	if !success {
		am.Log.Println(logging.LogWarn, "Attempted login with incorrect password for account:", email)
		response.Error = messages.AMError_InvalidPassword
		return response
	}

	account.HashedPassword = ""
	response.Account = account

	return response
}
