package login

import (
	"github.com/yamamushi/EscapingEden/ui/types"
	"log"
)

// handleRegistrationInput handles input for the registration screen of the login window
func (lw *LoginWindow) handleRegistrationInput(input types.Input) {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	if !lw.GetActive() {
		return
	}
	log.Println("Handling registration input")

	switch lw.registrationState {
	case RegistrationMain:
		lw.handleRegistrationMainInput(input)
	case RegistrationUserInfo:
		lw.handleRegistrationUserInfo(input)
	}
}

func (lw *LoginWindow) handleRegistrationMainInput(input types.Input) {
	switch input.Type {
	case types.InputCharacter:
		switch input.Data {
		case "b":
			log.Println("Opening controls help page")
			lw.RequestHelpFromConsole(types.HelpPageControls)
			return
		case "d":
			log.Println("Opening death help page")
			lw.RequestHelpFromConsole(types.HelpPageDeath)
			return
		case "r":
			log.Println("Opening rules help page")
			lw.RequestHelpFromConsole(types.HelpPageRules)
			return
		default:
			lw.Error("Invalid input received")
			return
		}
	case types.InputLeft:
		log.Println("Left arrow pressed")
		lw.optionSelected = 1
		return
	case types.InputRight:
		log.Println("Right arrow pressed")
		lw.optionSelected = 2
		return
	case types.InputReturn:
		log.Println("Return pressed")
		if lw.optionSelected == 1 {
			lw.windowState = LoginWindowMenu
		}
		if lw.optionSelected == 2 {
			lw.registrationState = RegistrationUserInfo
		}
		lw.optionSelected = 0
		//lw.ResetWindowDrawings()
		lw.ForceConsoleRefresh() // Whenever we switch to a different window state, we need to reset the console
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
			lw.Error("Invalid input received")
			return
		}
	case types.InputLeft:
		log.Println("Left arrow pressed")
		lw.optionSelected = 1
	case types.InputRight:
		log.Println("Right arrow pressed")
		lw.optionSelected = 2
	case types.InputUp:
		log.Println("Up arrow pressed")
	case types.InputDown:
		log.Println("Down arrow pressed")
	case types.InputReturn:
		log.Println("Return pressed")
		if lw.optionSelected == 1 {
			lw.registrationState = RegistrationMain
		}
		if lw.optionSelected == 2 {
			log.Println("Unimplemented state")
		}
		lw.optionSelected = 0
		lw.ForceConsoleRefresh() // Whenever we switch to a different window state, we need to reset the console

	}
}
