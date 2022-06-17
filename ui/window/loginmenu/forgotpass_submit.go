package login

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (lw *LoginWindow) forgotPasswordGenerateToken() {
	lw.loginForgotPasswordMutex.Lock()
	defer lw.loginForgotPasswordMutex.Unlock()

	lw.Log.Println(logging.LogInfo, "Submitting forgot password request to Eden Bot")
	forgotPassRequest := messages.WindowMessage{Type: messages.WM_RequestForgotPassword, Data: lw.loginForgotPasswordData}
	lw.SendToConsole(forgotPassRequest)
}

func (lw *LoginWindow) forgotPasswordValidate() {
	lw.loginForgotPasswordPendingMutex.Lock()
	defer lw.loginForgotPasswordPendingMutex.Unlock()

	lw.loginProcessForgotPasswordPendingData.DiscordTag = lw.loginForgotPasswordData.DiscordTag
	lw.loginProcessForgotPasswordPendingData.Username = lw.loginForgotPasswordData.Username
	// Now we take the submitted code and send it to the account manager for processing
	forgotPassRequest := messages.WindowMessage{Type: messages.WM_ValidateForgotPassword, Data: lw.loginProcessForgotPasswordPendingData}
	lw.SendToConsole(forgotPassRequest)

	go lw.HandleReceiveChannel() // We're going to start listening for responses now, the channel will handle the rest
}

func (lw *LoginWindow) forgotPasswordSubmitNewPassword() {
	lw.loginForgotPasswordSuccessMutex.Lock()
	defer lw.loginForgotPasswordSuccessMutex.Unlock()

	if lw.loginForgotPasswordNewPasswordData.Password != lw.loginForgotPasswordNewPasswordData.PasswordConfirm {
		lw.Log.Println(logging.LogInfo, "Passwords do not match")
		lw.loginForgotPasswordNewPasswordData.Error = "Passwords do not match"
		return
	}

	lw.loginProcessForgotPasswordPendingData.NewPassword = lw.loginForgotPasswordNewPasswordData.Password
	// Now we take the submitted code and send it to the account manager for processing
	newPassRequest := messages.WindowMessage{Type: messages.WM_ProcessForgotPassword, Data: lw.loginProcessForgotPasswordPendingData}
	lw.Log.Println(logging.LogInfo, "Submitting new password to Account Manager")
	lw.SendToConsole(newPassRequest)
}
