package dashboard

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/ui/types"
)

// handleMenuInput handles input for the login window
func (dw *DashboardWindow) handleMenuInput(input types.Input) {
	dw.dwMutex.Lock()
	defer dw.dwMutex.Unlock()

	if !dw.GetActive() {
		return
	}

	switch input.Type {
	case types.InputCharacter:
		switch input.Data {
		case "a":
			//log.Println("login as last selected")
			return
		case "b":
			//log.Println("manage characters selected")
			//dw.ResetWindowDrawings() // needed?
			//dw.ForceConsoleRefresh()
			//dw.RequestFlushFromConsole()
			return
		case "c":
			//log.Println("manage settings selected")
			return
		case "d":
			//log.Println("logout selected")
			dw.NotifyConsoleLoggedOut()
			dw.Log.Println(logging.LogInfo, "User: "+dw.GetUserInfoField("username")+" logged out")
			return
		default:
			dw.Error("Invalid input received")
			return
		}
	}
}
