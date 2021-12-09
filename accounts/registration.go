package accounts

import (
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

/*
AccountManager functions related to account registration
*/

// UsernameExists checks if a username is already in use
func (am *AccountManager) UsernameExists(username string) (bool, error) {
	result := messages.Account{}
	err := am.DB.One("Accounts", "Username", username, &result)
	if err != nil {
		if err.Error() == "not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// EmailExists checks if an email is already registered to an account
func (am *AccountManager) EmailExists(email string) (bool, error) {
	result := messages.Account{}
	err := am.DB.One("Accounts", "Email", email, &result)
	if err != nil {
		if err.Error() == "not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateAccount creates a new account, returns nil on success or error on failure
func (am *AccountManager) CreateAccount(username, password, email string) messages.AccountRegistrationResponse {
	am.Log.Println(logging.LogInfo, "Creating account:", username, email, password)

	response := messages.AccountRegistrationResponse{}
	// Before we work, lets make sure the username and email are not already taken
	usernameStatus, err := am.UsernameExists(username)
	if err != nil || usernameStatus {
		response.Error = messages.AMError_UsernameAlreadyExists
		return response
	}

	emailStatus, err := am.EmailExists(email)
	if err != nil || emailStatus {
		response.Error = messages.AMError_EmailAlreadyExists
		return response
	}

	hash, err := am.HashPassword(password)
	if err != nil {
		response.Error = messages.AMError_SystemError
		return response
	}

	account := messages.Account{
		ID:             uuid.New().String(),
		Username:       username,
		HashedPassword: hash,
		Email:          email,
	}
	err = am.DB.AddRecord("Accounts", &account)
	if err != nil {
		am.Log.Println(logging.LogError, "Error creating account:", err)
		response.Error = messages.AMError_AccountAlreadyExists
		return response
	}

	return response
}
