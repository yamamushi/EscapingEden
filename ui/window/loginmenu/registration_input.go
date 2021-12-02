package login

import (
	"github.com/yamamushi/EscapingEden/ui/types"
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
			lw.registrationUserInfoCharInput(input.Data)
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
			// This is where we submit our entered user data
			registrationError := lw.RegistrationSubmit(lw.registrationSubmitData)
			if registrationError != nil {
				// If we got an error, we're just going to update our error state
				lw.registrationErrorData = *registrationError
			} else {
				// Otherwise we succeeded (yay!) and we can go to the success screen
				lw.registrationState = RegistrationSuccess
				// Let's also cleanup our registration data
				lw.registrationAgreeRules = false
				lw.registrationErrorData = RegistrationError{}
				lw.registrationSubmitData = RegistrationSubmitData{}
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
		lw.registrationSubmitData.Username += input
	case UserInfoPassword:
		lw.registrationSubmitData.Password += input
	case UserInfoPasswordConfirm:
		lw.registrationSubmitData.PasswordConfirm += input
	case UserInfoEmail:
		lw.registrationSubmitData.Email += input
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
		}
	case UserInfoPassword:
		if lw.registrationSubmitData.Password != "" {
			lw.registrationSubmitData.Password = lw.registrationSubmitData.Password[:len(lw.registrationSubmitData.Password)-1]
		}
	case UserInfoPasswordConfirm:
		if lw.registrationSubmitData.PasswordConfirm != "" {
			lw.registrationSubmitData.PasswordConfirm = lw.registrationSubmitData.PasswordConfirm[:len(lw.registrationSubmitData.PasswordConfirm)-1]
		}
	case UserInfoEmail:
		if lw.registrationSubmitData.Email != "" {
			lw.registrationSubmitData.Email = lw.registrationSubmitData.Email[:len(lw.registrationSubmitData.Email)-1]
		}
	}
}

func (lw *LoginWindow) handleRegistrationSuccess(input types.Input) {

}
