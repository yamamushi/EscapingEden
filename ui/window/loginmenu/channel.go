package login

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (lw *LoginWindow) HandleReceiveChannel() {
	for {
		select {
		case windowMessage := <-lw.ConsoleReceive:
			switch windowMessage.Type {
			case messages.WM_RegistrationResponse:
				lw.registrationStatusMutex.Lock()
				defer lw.registrationStatusMutex.Unlock()

				lw.Log.Println(logging.LogInfo, "handleLogin Window received registration response from console")
				lw.registrationResponse = windowMessage.Data.(messages.AccountRegistrationResponse)

				lw.registrationErrorMutex.Lock()
				defer lw.registrationErrorMutex.Unlock()

				switch lw.registrationResponse.Error {
				case messages.AMError_UsernameAlreadyExists:
					lw.registrationErrorData.usernameError = lw.registrationResponse.Error.Error()
				case messages.AMError_DiscordAlreadyExists:
					lw.registrationErrorData.discordError = lw.registrationResponse.Error.Error()
				case messages.AMError_AccountAlreadyExists:
					lw.registrationErrorData.usernameError = "This account already exists"
					//lw.registrationErrorData.discordError = messages.AMError_DiscordAlreadyExists.Error()
				case messages.AMError_PendingValidation:
					lw.registrationErrorData.discordError = "This account is pending discord validation."
				case messages.AMError_UserNotInServer:
					lw.registrationErrorData.discordError = "This account is not in the discord server."
				case messages.AMError_DiscordMessageError:
					lw.registrationErrorData.discordError = "Discord message error, check your private message settings."
				case messages.AMError_DBError:
					lw.registrationErrorData.errorRequest = "Unhandled DB error, please report this issue"
				case messages.AMError_SystemError:
					lw.registrationErrorData.errorRequest = lw.registrationResponse.Error.Error()
				default:
					lw.registrationErrorData.errorRequest = "Unhandled error, please report this issue."
				}

				lw.registrationResponseReceived = true
				//lw.RequestFlushFromConsole()
				return // We launched when our registration request was submitted, now we can return since we got a response

			case messages.WM_LoginResponse:
				// Handle the login response here and then return
				lw.loginStatusMutex.Lock()
				defer lw.loginStatusMutex.Unlock()

				lw.loginSubmitMutex.Lock()
				defer lw.loginSubmitMutex.Unlock()

				lw.Log.Println(logging.LogInfo, "handleLogin Window received login response from console")

				lw.loginSubmitData.Error = "" // Flush any previous errors
				lw.loginResponse = windowMessage.Data.(messages.AccountLoginResponse)

				if lw.loginResponse.Error != messages.AMError_Null {
					lw.loginSubmitData.Error = lw.loginResponse.Error.Error()
					lw.loginAttempts += 1
				}
				lw.loginResponseReceived = true
				//lw.RequestFlushFromConsole()
				return
			case messages.WM_PasswordResetValidateResponse:
				success := windowMessage.Data.(bool)
				if success {
					lw.RequestFlushFromConsole()
					lw.loginState = LoginForgotPasswordSuccess
				} else {
					lw.RequestFlushFromConsole()
					lw.loginState = LoginForgotPasswordFailed
				}
			case messages.WM_PasswordResetProcessResponse:
				success := windowMessage.Data.(bool)
				if success {
					lw.Log.Println(logging.LogInfo, "handleLogin Window received password reset success")
					lw.loginMenuMessage = "Your password has been reset. Please login with your new password."
					lw.RequestFlushFromConsole()
					lw.loginProcessForgotPasswordPendingData = messages.AccountProcessForgotPasswordData{}
					lw.loginForgotPasswordData = messages.AccountForgotPasswordData{}
					lw.loginForgotPasswordNewPasswordData = LoginForgotPasswordSuccessData{}
					lw.loginState = LoginNull
					lw.windowState = LoginWindowMenu
				} else {
					lw.Log.Println(logging.LogInfo, "handleLogin Window received password reset failed")
					lw.RequestFlushFromConsole()
					lw.loginForgotPasswordNewPasswordData.Error = "Password reset failed, please try again."
				}
			}
		}
	}
}
