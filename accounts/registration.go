package accounts

import (
	"github.com/yamamushi/EscapingEden/messages"
	"golang.org/x/crypto/bcrypt"
)

/*
AccountManager functions related to account registration
*/

func (am *AccountManager) UsernameExists(username string) bool {
	result := messages.Account{}
	err := am.DB.One("Characters", "Username", username, &result)
	if err != nil {
		if err.Error() == "not found" {
			return false
		}
	}
	return true
}

func (am *AccountManager) EmailExists(email string) bool {
	result := messages.Account{}
	err := am.DB.One("Characters", "Email", email, &result)
	if err != nil {
		if err.Error() == "not found" {
			return false
		}
	}
	return true
}

// HashPassword uses bcrypt.GenerateFromPassword to hash a password
func (am *AccountManager) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (am *AccountManager) ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
