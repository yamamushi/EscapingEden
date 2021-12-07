package login

import (
	"strings"
)

// handleMenuInput handles input for the login window
func (lw *LoginWindow) handleMenuInput(input string) {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	if !lw.GetActive() {
		return
	}
	if len(input) == 0 {
		return
	}
	input = strings.ToLower(input)

	input = input[:1]

	switch input {
	case "l":
		//log.Println("Login selected")
		lw.windowState = LoginWindowLogin
		lw.loginState = LoginUsername
		//lw.ResetWindowDrawings()
		//lw.ForceConsoleRefresh()
		lw.RequestFlushFromConsole()

		return
	case "r":
		//log.Println("Register selected")
		lw.windowState = LoginWindowRegister
		lw.registrationState = RegistrationMain
		//lw.ResetWindowDrawings()
		//lw.ForceConsoleRefresh()
		lw.RequestFlushFromConsole()
		return
	case "q":
		//log.Println("Quit selected")
		lw.Quit()
		return
	default:
		lw.Error("Invalid input received")
		return
	}
}
