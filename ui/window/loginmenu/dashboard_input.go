package login

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/ui/types"
)

// handleMenuInput handles input for the login window
func (lw *LoginWindow) handleUserDashboardInput(input types.Input) {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	if !lw.GetActive() {
		return
	}

	switch input.Type {
	case types.InputCharacter:
		switch input.Data {
		case "a":
			//log.Println("handleLogin selected")
			lw.windowState = LoginWindowLogin
			lw.loginState = LoginUserInfo
			lw.loginMenuMessage = ""
			lw.RequestFlushFromConsole()

			return
		case "b":
			//log.Println("Register selected")
			lw.windowState = LoginWindowRegister
			lw.registrationState = RegistrationMain
			lw.loginMenuMessage = ""
			//lw.ResetWindowDrawings()
			//lw.ForceConsoleRefresh()
			lw.RequestFlushFromConsole()
			return
		case "c":
			//log.Println("Register selected")
			lw.windowState = LoginWindowRegister
			lw.registrationState = RegistrationMain
			lw.loginMenuMessage = ""
			//lw.ResetWindowDrawings()
			//lw.ForceConsoleRefresh()
			lw.RequestFlushFromConsole()
			return
		case "d":
			//log.Println("Quit selected")
			lw.windowState = LoginWindowMenu
			lw.NotifyConsoleLoggedOut()
			lw.RequestFlushFromConsole()

			lw.Log.Println(logging.LogInfo, "User logged out")
			return
		default:
			lw.Error("Invalid input received")
			return
		}
	}
}
