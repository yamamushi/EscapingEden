package accounts

import (
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"strings"
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
func (am *AccountManager) DiscordExists(discordID string) (*messages.Account, messages.AMErrorType) {
	result := messages.Account{}
	err := am.DB.One("Accounts", "DiscordID", discordID, &result)
	if err != nil {
		if err.Error() == "not found" {
			return nil, messages.AMError_Null // no error, no account
		}
		return nil, messages.AMError_DBError
	}
	return &result, messages.AMError_Null
}

// DiscordTagExists checks if a discord tag is already registered to an account
func (am *AccountManager) DiscordTagExists(discordTag string) (*messages.Account, messages.AMErrorType) {
	result := messages.Account{}
	err := am.DB.One("Accounts", "DiscordTag", discordTag, &result)
	if err != nil {
		if err.Error() == "not found" {
			return nil, messages.AMError_Null // no error, no account
		}
		return nil, messages.AMError_DBError
	}
	return &result, messages.AMError_Null
}

// CreateAccount creates a new account, returns nil on success or error on failure
func (am *AccountManager) CreateAccount(username, password, discordTag string) messages.AccountRegistrationResponse {

	username = strings.TrimSpace(username)
	discordTag = strings.TrimSpace(discordTag)

	response := messages.AccountRegistrationResponse{}
	// Before we work, lets make sure the username and discord are not already taken
	foundAccount, err := am.DiscordTagExists(discordTag)
	if err == messages.AMError_DBError {
		response.Error = messages.AMError_DBError
		return response
	}

	if foundAccount != nil {
		if foundAccount.ValidationStatus == 0 && foundAccount.Username == username {
			response.Error = messages.AMError_PendingValidation
			response.ValidationCode = foundAccount.ValidationCode
			return response
		}
		response.Error = messages.AMError_DiscordAlreadyExists
		return response
	}

	usernameStatus, usernameError := am.UsernameExists(username)
	if usernameError != nil || usernameStatus {
		response.Error = messages.AMError_UsernameAlreadyExists
		return response
	}

	// Validate if user is in discord server
	discordUser, serverError := am.EB.IsUserInDiscordServer(discordTag)
	if serverError != nil {
		response.Error = messages.AMError_UserNotInServer
		return response
	}

	hash, hashPassError := am.HashPassword(password)
	if hashPassError != nil {
		response.Error = messages.AMError_SystemError
		return response
	}
	// generate a random uuid for the account registration
	registrationCode, uuidError := uuid.NewUUID()
	if uuidError != nil {
		response.Error = messages.AMError_SystemError
		return response
	}
	// Send a message to the user to validate their account
	// If we can't reach a user, there's no point to adding the account to the database
	pmError := am.EB.ValidateUser(username, discordUser.ID)
	if pmError != nil {
		response.Error = messages.AMError_DiscordMessageError
		return response
	}

	// Now we can add the account to the database
	account := messages.Account{
		ID:               uuid.New().String(),
		Username:         username,
		HashedPassword:   hash,
		DiscordTag:       discordTag,
		DiscordID:        discordUser.ID,
		ValidationCode:   registrationCode.String(),
		ValidationStatus: 0, // 0 = pending validation
	}
	dbError := am.DB.AddRecord("Accounts", &account)
	if dbError != nil {
		// If we hit this, we have a problem because validation above missed something
		am.Log.Println(logging.LogError, "Error creating account:", dbError)
		response.Error = messages.AMError_SystemError
		return response
	}

	response.ValidationCode = registrationCode.String()
	return response
}
