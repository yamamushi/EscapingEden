package dashboard

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (dw *DashboardWindow) HandleReceiveChannel() {
	for {
		select {
		case windowMessage := <-dw.ConsoleReceive:
			switch windowMessage.Type {
			case messages.WM_RequestCharNameValidationResponse:
				//dw.Log.Println(logging.LogInfo, "Received WM_RequestCharNameValidationResponse")
				data := windowMessage.Data.(messages.CharManagerNameCheckResponse)
				if data.NameInUse {
					dw.Log.Println(logging.LogInfo, "Character Name In Use")
					if data.Error != "" {
						dw.Log.Println(logging.LogInfo, data.Error)
					}
					dw.charCreatorUsernameError = "Character Name In Use"
					dw.characterCreatorState = CharacterCreatorCharacterDetails
					dw.RequestFlushFromConsole()
				} else {
					dw.Log.Println(logging.LogInfo, "Character Name Not In Use")
					dw.CreateCharacter()
					dw.RequestFlushFromConsole()
				}

			case messages.WM_RequestCharacterCreationResponse:
				//dw.Log.Println(logging.LogInfo, "Received WM_RequestCharacterCreationResponse")
				data := windowMessage.Data.(messages.CharacterInfo)
				if data.Error != "Null Error" {
					dw.Log.Println(logging.LogInfo, data.Error)
					dw.charCreatorUsernameError = data.Error // TODO - Send these errors through the chat window
					dw.characterCreatorState = CharacterCreatorCharacterDetails
					dw.RequestFlushFromConsole()
				} else {
					dw.Log.Println(logging.LogInfo, "Character Created")
					dw.LoginCharacter(data)
					dw.RequestFlushFromConsole()
				}

			case messages.WM_RequestCharacterResponse:
				//dw.Log.Println(logging.LogInfo, "Received WM_RequestCharacterResponse")
				data := windowMessage.Data.(messages.CharacterInfo)
				if data.Error != "Null Error" {
					dw.Error("Something went wrong when trying to get your character, please try again later.")
					dw.windowState = DashboardMainMenu
					dw.RequestFlushFromConsole()
				} else {
					dw.Log.Println(logging.LogInfo, "Character Found")
					// TODO - update pending screen that character was found and a login is being attempted.
					dw.pendingLoginMessage = "Character found, attempting character timestamp update..."
					go dw.ValidateCharacterLogin(data)
					go dw.HandleReceiveChannel()
					dw.RequestFlushFromConsole()
				}

			case messages.WM_RequestCharacterHistoryAccountUpdateResponse:
				//dw.Log.Println(logging.LogInfo, "Received WM_RequestCharacterHistoryAccountUpdateResponse")
				data := windowMessage.Data.(messages.CharManagerUpdateHistoryResponse)
				account := data.Data.(messages.Account)
				if account.Error != "Null Error" {
					// TODO - This error needs to be handled better - right now we basically ignore it, and logins can break.
					dw.Log.Println(logging.LogError, "Something went wrong when trying to update your account, please try again later.")
					dw.windowState = DashboardMainMenu
					dw.RequestFlushFromConsole()
				} else {
					dw.pendingLoginMutex.Lock()
					dw.accountManagerValidated = true
					dw.pendingLoginMutex.Unlock()
					//dw.Log.Println(logging.LogInfo, "Account Login Time Updated")
					dw.accountManagerValidated = true
				}

			case messages.WM_RequestCharacterHistoryCharacterUpdateResponse:
				//dw.Log.Println(logging.LogInfo, "Received WM_RequestCharacterHistoryCharacterUpdateResponse")
				data := windowMessage.Data.(messages.CharManagerUpdateHistoryResponse)
				charInfo := data.Data.(messages.CharacterInfo)
				if charInfo.Error != "Null Error" {
					// TODO - This error needs to be handled better - right now we basically ignore it, and logins can break.
					dw.Error("Something went wrong when trying to update your character, please try again later.")
					dw.windowState = DashboardMainMenu
					dw.RequestFlushFromConsole()
				} else {
					dw.pendingLoginMutex.Lock()
					dw.characterManagerValidated = true
					dw.pendingLoginMutex.Unlock()
					dw.Log.Println(logging.LogInfo, "Character Login Time Updated")
					dw.RequestFlushFromConsole()
					go dw.finalizeLogin(charInfo)
				}

			}
		}
		dw.pendingLoginMutex.Lock()
		if dw.stopHandler {
			dw.pendingLoginMutex.Unlock()
			return
		}
		dw.pendingLoginMutex.Unlock()
	}
}
