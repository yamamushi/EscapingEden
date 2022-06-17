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

	if account.PasswordResetStatus > 0 {
		am.Log.Println(logging.LogWarn, "Attempted login to account that is pending password reset:", username)
		response.Error = messages.AMError_PendingPasswordReset
		return response
	}

	account.HashedPassword = ""
	response.Account = account

	return response
}

func (am *AccountManager) handleValidatePasswordReset(resetRequest messages.AccountProcessForgotPasswordData) (response bool) {
	account := messages.Account{}

	err := am.DB.One("Accounts", "DiscordTag", resetRequest.DiscordTag, &account)
	if err != nil {
		am.Log.Println(logging.LogWarn, "Attempted validate a password reset for a nonexistent account:", resetRequest.Username)
		return false
	}
	if account.DiscordTag != resetRequest.DiscordTag {
		am.Log.Println(logging.LogWarn, "Attempted validate a password reset for a discord tag that does not match:", resetRequest.DiscordTag)
		return false
	}
	if account.PasswordResetCode != resetRequest.Code {
		am.Log.Println(logging.LogWarn, "Attempted validate a password reset for a password reset code that does not match:", resetRequest.Username)
		return false
	}

	account.PasswordResetStatus = 1 // Reset requested, now a reset can be performed on the password. Without this, a password reset will fail.
	err = am.UpdateAccount(account)
	if err != nil {
		// We fail here because without the reset status being set, a password reset will fail anyway.
		am.Log.Println(logging.LogWarn, "Failed to update account after password reset request:", resetRequest.Username)
		return false
	}
	return true
}

func (am *AccountManager) handleChangePassword(resetRequest messages.AccountProcessForgotPasswordData) (response bool) {
	account := messages.Account{}

	err := am.DB.One("Accounts", "DiscordTag", resetRequest.DiscordTag, &account)
	if err != nil {
		am.Log.Println(logging.LogWarn, "Attempted a password reset for a nonexistent account:", resetRequest.Username)
		return false
	}
	if account.DiscordTag != resetRequest.DiscordTag {
		am.Log.Println(logging.LogWarn, "Attempted a password reset for a discord tag that does not match:", resetRequest.DiscordTag)
		return false
	}
	if account.PasswordResetCode != resetRequest.Code {
		am.Log.Println(logging.LogWarn, "Attempted a password reset for a password reset code that does not match:", resetRequest.Username)
		return false
	}

	if account.PasswordResetStatus > 0 {
		account.PasswordResetStatus = 0 // We reset the status to 0 so that logins are allowed again.
		account.PasswordResetCode = ""
		account.HashedPassword, err = am.HashPassword(resetRequest.NewPassword)
		if err != nil {
			am.Log.Println(logging.LogWarn, "Failed to hash password for password reset:", resetRequest.Username)
			return false
		}

		err = am.UpdateAccountField("PasswordResetStatus", 0, account)
		if err != nil {
			am.Log.Println(logging.LogWarn, "Failed to update account PasswordResetStatus after password reset:", resetRequest.Username)
			return false
		}

		err = am.UpdateAccountField("PasswordResetCode", "", account)
		if err != nil {
			am.Log.Println(logging.LogWarn, "Failed to update account PasswordResetCode after password reset:", resetRequest.Username)
		}

		err = am.UpdateAccountField("HashedPassword", account.HashedPassword, account)
		if err != nil {
			am.Log.Println(logging.LogWarn, "Failed to update account HashedPassword after password reset:", resetRequest.Username)
			return false
		}

		return true
	} else {
		return false // Password reset status is not set, so we can't change the password.
	}
}
