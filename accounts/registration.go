package accounts

import (
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"golang.org/x/crypto/bcrypt"
)

/*
AccountManager functions related to account registration
*/

// UsernameExists checks if a username is already in use
func (am *AccountManager) UsernameExists(username string) (bool, error) {
	result := messages.Account{}
	err := am.DB.One("Characters", "Username", username, &result)
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
	err := am.DB.One("Characters", "Email", email, &result)
	if err != nil {
		if err.Error() == "not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// HashPassword uses bcrypt.GenerateFromPassword to hash a password
func (am *AccountManager) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		am.Log.Println(logging.LogError, "Error hashing password:", err)
		return "", messages.AMError_SystemError
	}
	return string(bytes), nil
}

// ComparePasswords compares a password with a hash using bcrypt.CompareHashAndPassword
func (am *AccountManager) ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// CreateAccount creates a new account, returns nil on success or error on failure
func (am *AccountManager) CreateAccount(username, password, email string) messages.AMErrorType {
	am.Log.Println(logging.LogInfo, "Creating account:", username, email, password)
	// Before we work, lets make sure the username and email are not already taken
	usernameStatus, err := am.UsernameExists(username)
	if err != nil || usernameStatus {
		return messages.AMError_UsernameAlreadyExists
	}

	emailStatus, err := am.EmailExists(email)
	if err != nil || emailStatus {
		return messages.AMError_EmailAlreadyExists
	}

	hash, err := am.HashPassword(password)
	if err != nil {
		return messages.AMError_Null
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
		return messages.AMError_AccountAlreadyExists
	}

	return messages.AMError_Null
}
