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
				case messages.AMError_EmailAlreadyExists:
					lw.registrationErrorData.emailError = lw.registrationResponse.Error.Error()
				case messages.AMError_AccountAlreadyExists:
					lw.registrationErrorData.usernameError = "This account already exists"
					//lw.registrationErrorData.emailError = messages.AMError_EmailAlreadyExists.Error()
				case messages.AMError_SystemError:
					lw.registrationErrorData.errorRequest = lw.registrationResponse.Error.Error()
				default:
					lw.registrationErrorData.errorRequest = "Unhandled Error - Please report this issue."
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

				lw.loginSubmitData.Error = "" // Flush any previous errors
				lw.loginResponse = windowMessage.Data.(messages.AccountLoginResponse)

				if lw.loginResponse.Error != messages.AMError_Null {
					lw.loginSubmitData.Error = lw.loginResponse.Error.Error()
					if lw.loginResponse.Error == messages.AMError_InvalidPassword {
						// Increment our login attempts if we have an invalid password
						lw.loginAttempts += 1
					}
				}
				lw.loginResponseReceived = true
				//lw.RequestFlushFromConsole()
				return
			}
		}
	}
}
