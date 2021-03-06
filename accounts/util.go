package accounts

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"golang.org/x/crypto/bcrypt"
)

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

// UpdateAccount updates the provided account in the database
func (am *AccountManager) UpdateAccount(account messages.Account) error {
	return am.DB.UpdateRecord("Accounts", &account)
}

func (am *AccountManager) UpdateAccountField(field string, value interface{}, account messages.Account) error {
	return am.DB.UpdateField("Accounts", field, value, &account)
}
