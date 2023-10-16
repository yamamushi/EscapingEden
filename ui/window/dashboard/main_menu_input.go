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
	//dw.Log.Println(logging.LogInfo, "last character:", dw.UserInfo.LastCharacterID)

	switch input.Type {
	case types.InputCharacter:
		switch input.Data {
		case "a":
			lastCharacter := dw.GetUserInfoField("lastcharacter")
			if lastCharacter == "" || dw.UserInfo.LastCharacterID == "" {
				if dw.firstTimeLogin {
					dw.characterCreatorState = CharacterCreatorFirstTimeLoginWelcome
				} else {
					dw.characterCreatorState = CharacterCreatorCharacterDetails
				}
				dw.windowState = DashboardCreateCharacter
				//dw.ForceConsoleRefresh() // don't use these both together
				dw.RequestFlushFromConsole()
				return
			} else {
				// Login by last character ID
				dw.loginCharacterByID(dw.UserInfo.LastCharacterID)
				// Need to use this ID to request the character info struct
				// Need to add the messages necessary for this, and the outbound message as well as listener on the window here
				dw.windowState = DashboardCharacterLoginPending
				dw.RequestFlushFromConsole()
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
