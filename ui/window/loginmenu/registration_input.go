package login

import (
	"github.com/yamamushi/EscapingEden/ui/types"
	"net/mail"
)

// handleRegistrationInput handles input for the registration screen of the login window
func (lw *LoginWindow) handleRegistrationInput(input types.Input) {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	if !lw.GetActive() {
		return
	}
	//log.Println("Handling registration input")

	switch lw.registrationState {
	case RegistrationMain:
		lw.handleRegistrationMainInput(input)
	case RegistrationUserInfo:
		lw.handleRegistrationUserInfo(input)
	case RegistrationPending:
		lw.handleRegistrationPending(input)
	case RegistrationFailure:
		lw.handleRegistrationFailure(input)
	case RegistrationSuccess:
		lw.handleRegistrationSuccess(input)
	}
}

func (lw *LoginWindow) handleRegistrationMainInput(input types.Input) {
	switch input.Type {
	case types.InputCharacter:
		switch input.Data {
		case "b":
			//log.Println("Opening controls help page")
			lw.RequestHelpFromConsole(types.HelpPageControls)
			return
		case "d":
			//log.Println("Opening death help page")
			lw.RequestHelpFromConsole(types.HelpPageDeath)
			return
		case "r":
			//log.Println("Opening rules help page")
			lw.RequestHelpFromConsole(types.HelpPageRules)
			return
		default:
			lw.Error("Invalid input received")
			return
		}
	case types.InputLeft:
		//log.Println("Left arrow pressed")
		lw.registrationNavOptionSelected = 1
		return
	case types.InputRight:
		//log.Println("Right arrow pressed")
		lw.registrationNavOptionSelected = 2
		return
	case types.InputReturn:
		//log.Println("Return pressed")
		if lw.registrationNavOptionSelected == 1 {
			lw.windowState = LoginWindowMenu
		}
		if lw.registrationNavOptionSelected == 2 {
			lw.registrationState = RegistrationUserInfo
		}
		lw.registrationNavOptionSelected = 0
		//lw.ResetWindowDrawings()
		//lw.ForceConsoleRefresh() // Whenever we switch to a different window state, we need to reset the console
		lw.RequestFlushFromConsole()
		return
	default:
		return
	}
}

func (lw *LoginWindow) handleRegistrationUserInfo(input types.Input) {

	switch input.Type {
	case types.InputCharacter:
		switch input.Data {
		default:
			lw.registrationErrorMutex.Lock()
			defer lw.registrationErrorMutex.Unlock()
			lw.registrationUserInfoCharInput(input.Data)
			lw.RequestFlushFromConsole()
			return
		}
	case types.InputLeft:
		//log.Println("Left arrow pressed")
		if lw.registrationNavOptionSelected != 2 {
			lw.registrationUserInfoLastOptionSelected = lw.registrationUserInfoOptionSelected
			lw.registrationUserInfoOptionSelected = UserInfoNULL
		}
		lw.registrationNavOptionSelected = 1

	case types.InputRight:
		//log.Println("Right arrow pressed")
		if lw.registrationNavOptionSelected != 1 {
			lw.registrationUserInfoLastOptionSelected = lw.registrationUserInfoOptionSelected
			lw.registrationUserInfoOptionSelected = UserInfoNULL
		}
		lw.registrationNavOptionSelected = 2
	case types.InputUp:
		//log.Println("Up arrow pressed")
		if lw.registrationNavOptionSelected == 2 {
			lw.registrationUserInfoOptionSelected = lw.registrationUserInfoLastOptionSelected
		} else if lw.registrationNavOptionSelected == 1 {
			lw.registrationUserInfoOptionSelected = lw.registrationUserInfoLastOptionSelected
		} else if lw.registrationUserInfoOptionSelected > 0 {
			lw.registrationUserInfoOptionSelected--
		}
		lw.registrationNavOptionSelected = 0
		lw.RequestFlushFromConsole()
	case types.InputDown:
		//log.Println("Down arrow pressed")
		if lw.registrationUserInfoOptionSelected < UserInfoNULL-1 {
			lw.registrationUserInfoOptionSelected++
		} else {
			lw.registrationUserInfoLastOptionSelected = lw.registrationUserInfoOptionSelected
			lw.registrationUserInfoOptionSelected = UserInfoNULL
			lw.registrationNavOptionSelected = 2
		}
		//lw.registrationNavOptionSelected = 0
		lw.RequestFlushFromConsole()
	case types.InputReturn:
		//log.Println("Return pressed")
		if lw.registrationNavOptionSelected == 0 {
			switch lw.registrationUserInfoOptionSelected {
			case UserInfoUsername:
				lw.registrationUserInfoOptionSelected = UserInfoPassword
			case UserInfoPassword:
				lw.registrationUserInfoOptionSelected = UserInfoPasswordConfirm
			case UserInfoPasswordConfirm:
				lw.registrationUserInfoOptionSelected = UserInfoEmail
			case UserInfoEmail:
				lw.registrationUserInfoOptionSelected = UserInfoAgreeRules
			case UserInfoAgreeRules:
				// If we have no option selected, and we're at the email line, we may as well just go to the submit button
				lw.registrationUserInfoOptionSelected = UserInfoNULL
				lw.registrationNavOptionSelected = 2
			}
			return
		}
		if lw.registrationNavOptionSelected == 1 {
			// If we're going back, we're going to flush our registration error and data
			lw.registrationState = RegistrationMain
			lw.registrationAgreeRules = false
			lw.registrationErrorData = RegistrationError{}
			lw.registrationSubmitData = RegistrationSubmitData{}
		}
		if lw.registrationNavOptionSelected == 2 {
			if !lw.registrationErrorData.IsEmpty() {
				lw.registrationErrorData.errorRequest = "Please correct the errors above and try again."
			}
			// This is where we submit our entered user data
			registrationError := lw.RegistrationSubmit(lw.registrationSubmitData)
			if registrationError != nil {
				// If we got an error, we're just going to update our error state
				lw.registrationErrorData = *registrationError
			} else {
				// Otherwise we succeeded (yay!) and we can go to the pending screen
				lw.registrationState = RegistrationPending
				// Let's also cleanup our registration data
				//lw.registrationAgreeRules = false

				// We need to reset our error data so that we can process it when we get our response
				// From the account manager
				lw.registrationErrorMutex.Lock()
				lw.registrationErrorData = RegistrationError{}
				lw.registrationErrorMutex.Unlock()
			}
		}
		lw.registrationNavOptionSelected = 0
		lw.registrationUserInfoOptionSelected = 0
		lw.RequestFlushFromConsole()
		//lw.ForceConsoleRefresh() // Whenever we switch to a different window state, we need to reset the console
	case types.InputBackspace:
		//log.Println("Backspace pressed")
		lw.registrationUserInfoBackspaceInput()
		lw.RequestFlushFromConsole()
		return
	}
}

func (lw *LoginWindow) registrationUserInfoCharInput(input string) {
	switch lw.registrationUserInfoOptionSelected {
	case UserInfoUsername:
		//lw.registrationErrorMutex.Lock()
		//defer lw.registrationErrorMutex.Unlock()

		if len(lw.registrationSubmitData.Username) < 16 {
			if input == " " {
				lw.registrationErrorData.usernameError = "Username cannot contain spaces"
			} else {
				lw.registrationSubmitData.Username += input
			}
			if len(lw.registrationSubmitData.Username) < 4 {
				lw.registrationErrorData.usernameError = "Username must be at least 4 characters"
			} else {
				lw.registrationErrorData.usernameError = ""
			}
		} else {
			lw.registrationErrorData.usernameError = "Maximum username length is 16 characters"
		}

	case UserInfoPassword:
		if len(lw.registrationSubmitData.Password) < 32 {
			if input == " " {
				lw.registrationErrorData.passwordError = "Password cannot contain spaces"
			} else {
				lw.registrationSubmitData.Password += input
			}
			if len(lw.registrationSubmitData.Password) < 8 {
				lw.registrationErrorData.passwordError = "Password must be at least 8 characters"
			} else {
				lw.registrationErrorData.passwordError = ""
			}
		} else {
			lw.registrationErrorData.passwordError = "Maximum password length is 32 characters"
		}

	case UserInfoPasswordConfirm:
		if len(lw.registrationSubmitData.PasswordConfirm) < 32 {
			if input == " " {
				lw.registrationErrorData.passwordConfirmError = "Password cannot contain spaces"
			} else {
				lw.registrationSubmitData.PasswordConfirm += input
			}
			if len(lw.registrationSubmitData.PasswordConfirm) < 8 {
				lw.registrationErrorData.passwordConfirmError = "Password must be at least 8 characters"
			} else {
				lw.registrationErrorData.passwordConfirmError = ""
			}
		} else {
			lw.registrationErrorData.passwordConfirmError = "Maximum password length is 32 characters"
		}

	case UserInfoEmail:
		lw.registrationSubmitData.Email += input
		_, err := mail.ParseAddress(lw.registrationSubmitData.Email)
		if err != nil {
			lw.registrationErrorData.emailError = "Invalid email address"
		} else {
			lw.registrationErrorData.emailError = ""
		}

	case UserInfoAgreeRules:
		if input == " " {
			lw.registrationAgreeRules = !lw.registrationAgreeRules
		}

	}
}

func (lw *LoginWindow) registrationUserInfoBackspaceInput() {
	switch lw.registrationUserInfoOptionSelected {
	case UserInfoUsername:
		if lw.registrationSubmitData.Username != "" {
			lw.registrationSubmitData.Username = lw.registrationSubmitData.Username[:len(lw.registrationSubmitData.Username)-1]
			lw.registrationErrorData.usernameError = ""
		}
	case UserInfoPassword:
		if lw.registrationSubmitData.Password != "" {
			lw.registrationSubmitData.Password = lw.registrationSubmitData.Password[:len(lw.registrationSubmitData.Password)-1]
			lw.registrationErrorData.passwordError = ""
		}
	case UserInfoPasswordConfirm:
		if lw.registrationSubmitData.PasswordConfirm != "" {
			lw.registrationSubmitData.PasswordConfirm = lw.registrationSubmitData.PasswordConfirm[:len(lw.registrationSubmitData.PasswordConfirm)-1]
			lw.registrationErrorData.passwordConfirmError = ""
		}
	case UserInfoEmail:
		if lw.registrationSubmitData.Email != "" {
			lw.registrationSubmitData.Email = lw.registrationSubmitData.Email[:len(lw.registrationSubmitData.Email)-1]
			_, err := mail.ParseAddress(lw.registrationSubmitData.Email)
			if err != nil {
				lw.registrationErrorData.emailError = "Invalid email address"
			} else {
				lw.registrationErrorData.emailError = ""
			}
		}
	}
}

func (lw *LoginWindow) handleRegistrationPending(input types.Input) {
	// This is a no-op, we just wait for the response from the account manager
}

func (lw *LoginWindow) handleRegistrationFailure(input types.Input) {
	switch input.Type {
	case types.InputLeft:
		//log.Println("Left arrow pressed")
		lw.registrationFailureOptionSelected = 1
		return
	case types.InputRight:
		//log.Println("Right arrow pressed")
		lw.registrationFailureOptionSelected = 2
		return
	case types.InputReturn:
		//log.Println("Return pressed")
		if lw.registrationFailureOptionSelected == 1 {
			lw.registrationState = RegistrationUserInfo
		}
		if lw.registrationNavOptionSelected == 2 {
			return
		}
		lw.registrationNavOptionSelected = 0
		//lw.ResetWindowDrawings()
		//lw.ForceConsoleRefresh() // Whenever we switch to a different window state, we need to reset the console
		lw.RequestFlushFromConsole()
		return
	default:
		return
	}
}

func (lw *LoginWindow) handleRegistrationSuccess(input types.Input) {
	switch input.Type {
	case types.InputLeft:
		//log.Println("Left arrow pressed")
		lw.registrationSuccessOptionSelected = 1
		return
	case types.InputRight:
		//log.Println("Right arrow pressed")
		lw.registrationSuccessOptionSelected = 2
		return
	case types.InputReturn:
		//log.Println("Return pressed")
		if lw.registrationSuccessOptionSelected == 1 {
			return
		}
		if lw.registrationSuccessOptionSelected == 2 {
			lw.windowState = LoginWindowMenu
		}
		lw.registrationSuccessOptionSelected = 0
		//lw.ResetWindowDrawings()
		//lw.ForceConsoleRefresh() // Whenever we switch to a different window state, we need to reset the console
		lw.RequestFlushFromConsole()
		return
	default:
		return
	}
}
