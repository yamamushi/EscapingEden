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
			lastCharacter := dw.GetUserInfoField("lastcharacter")
			if lastCharacter == "" {
				if lastCharacter == "" {
					dw.characterCreatorState = CharacterCreatorFirstTimeLoginWelcome
				} else {
					dw.characterCreatorState = CharacterCreatorCharacterDetails
				}
				dw.windowState = DashboardCreateCharacter
				//dw.ForceConsoleRefresh() // don't use these both together
				dw.RequestFlushFromConsole()
				return
			} else {
				// TODO - login to last character
				dw.LoginCharacter(lastCharacter)
				return
			}
		case "b":
			lastCharacter := dw.GetUserInfoField("lastcharacter")
			if lastCharacter != "" {
				dw.Log.Println(logging.LogInfo, "manage characters selected")
				//dw.ResetWindowDrawings() // needed?
				//dw.ForceConsoleRefresh()
				//dw.RequestFlushFromConsole()
				return
			}
		case "c":
			dw.Log.Println(logging.LogInfo, "manage user settings selected")
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
