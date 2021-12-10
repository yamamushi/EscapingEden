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

// DiscordExists checks if a discord id is already registered to an account
func (am *AccountManager) DiscordExists(discordID string) (bool, error) {
	result := messages.Account{}
	err := am.DB.One("Accounts", "DiscordID", discordID, &result)
	if err != nil {
		if err.Error() == "not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateAccount creates a new account, returns nil on success or error on failure
func (am *AccountManager) CreateAccount(username, password, discordID string) messages.AccountRegistrationResponse {
	am.Log.Println(logging.LogInfo, "Creating account:", username, discordID, password)

	response := messages.AccountRegistrationResponse{}
	// Before we work, lets make sure the username and discord are not already taken
	usernameStatus, err := am.UsernameExists(username)
	if err != nil || usernameStatus {
		response.Error = messages.AMError_UsernameAlreadyExists
		return response
	}

	discordStatus, err := am.DiscordExists(discordID)
	if err != nil || discordStatus {
		response.Error = messages.AMError_DiscordAlreadyExists
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
		DiscordID:      discordID,
	}
	err = am.DB.AddRecord("Accounts", &account)
	if err != nil {
		am.Log.Println(logging.LogError, "Error creating account:", err)
		response.Error = messages.AMError_AccountAlreadyExists
		return response
	}

	return response
}
