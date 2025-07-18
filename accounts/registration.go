package accounts

import (
	"strings"

	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/edenbot"
	"github.com/yamamushi/EscapingEden/edenutil"
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

	// Sanitize and validate input
	username = edenutil.SanitizeInput(strings.TrimSpace(username))
	password = edenutil.SanitizeInput(password)
	discordTag = edenutil.SanitizeInput(strings.TrimSpace(discordTag))

	response := messages.AccountRegistrationResponse{}

	// Validate username format
	if !edenutil.ValidateUsername(username) {
		response.Error = messages.AMError_InvalidUsername
		return response
	}

	// Validate password strength
	if !edenutil.ValidatePassword(password) {
		response.Error = messages.AMError_InvalidPassword
		return response
	}

	// Check if Discord is required based on configuration
	if am.Config.Discord.RequireDiscord {
		// Validate Discord ID format (if it's not a tag format)
		//if !strings.Contains(discordTag, "#") && !edenutil.ValidateDiscordID(discordTag) {
		//	response.Error = messages.AMError_InvalidDiscordID
		//	return response
		//}

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
	} else {
		// Discord not required - set discordTag to empty string to skip Discord validation
		discordTag = ""
	}

	usernameStatus, usernameError := am.UsernameExists(username)
	if usernameError != nil || usernameStatus {
		response.Error = messages.AMError_UsernameAlreadyExists
		return response
	}

	hash, hashPassError := am.HashPassword(password)
	if hashPassError != nil {
		response.Error = messages.AMError_SystemError
		return response
	}

	var discordUser *edenbot.DiscordUser
	var registrationCode uuid.UUID
	var uuidError error

	if am.Config.Discord.RequireDiscord {
		// Discord validation required
		// Validate if user is in discord server
		var serverError error
		discordUser, serverError = am.EB.IsUserInDiscordServer(discordTag)
		if serverError != nil {
			response.Error = messages.AMError_UserNotInServer
			return response
		}

		// generate a random uuid for the account registration
		registrationCode, uuidError = uuid.NewUUID()
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
	} else {
		// Discord validation not required - account is immediately validated
		registrationCode = uuid.New()
	}

	// Now we can add the account to the database
	account := messages.Account{
		ID:             uuid.New().String(),
		Username:       username,
		HashedPassword: hash,
		DiscordTag:     discordTag,
		ValidationCode: registrationCode.String(),
	}

	if am.Config.Discord.RequireDiscord && discordUser != nil {
		// Set Discord ID and validation status for Discord-required accounts
		account.DiscordID = discordUser.ID
		account.ValidationStatus = 0 // 0 = pending validation
	} else {
		// No Discord required - account is immediately validated
		account.DiscordID = ""
		account.ValidationStatus = 1 // 1 = validated
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
