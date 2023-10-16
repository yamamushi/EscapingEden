package character

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (cm *CharacterManager) HandleInput(started chan bool) {
	cm.Log.Println(logging.LogInfo, "Character Manager Input Handler Launched")
	started <- true
	for {
		select {
		case managerMessage := <-cm.InputChannel:
			switch managerMessage.Type {
			case messages.CharManager_None:
				cm.Log.Println(logging.LogError, "Character Manager Input Handler Received Message Type: None from Sender:", managerMessage.SenderConsoleID)
				continue
			case messages.CharManager_CreateCharacter:
				//cm.Log.Println(logging.LogInfo, "Character Manager: Create Character")
				newCharacterInfo := cm.CreateCharacter(managerMessage.Data.(messages.CharacterInfo))
				cm.OutputChannel <- messages.ConnectionManagerMessage{Type: messages.ConnectManager_Message_CharacterCreationResponse, Data: newCharacterInfo, RecipientConsoleID: managerMessage.SenderConsoleID}
				//cm.Log.Println(logging.LogInfo, "Character Manager: Create Character response sent")
				continue
			case messages.CharManager_DeleteCharacter:
				cm.Log.Println(logging.LogInfo, "Character Manager: Delete Character")
				continue
			case messages.CharManager_ListCharacters:
				cm.Log.Println(logging.LogInfo, "Character Manager: List Characters")
				continue
			case messages.CharManager_UpdateLoginHistory:
				//cm.Log.Println(logging.LogInfo, "Character Manager: Update Login History")
				charInfo := cm.UpdateLoginHistory(managerMessage.Data.(messages.CharacterInfo))
				if charInfo.Error != messages.CMError_Null.Error() {
					cm.Log.Println(logging.LogError, "Character Manager: Update Login History Error")
				}
				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_UpdateCharacterHistoryResponse,
					RecipientConsoleID: managerMessage.SenderConsoleID,
					Data: messages.CharManagerUpdateHistoryResponse{
						Data:              charInfo,
						RespondingManager: "character"},
				}
				cm.OutputChannel <- response
				continue
			case messages.CharManager_UpdateCharacter:
				cm.Log.Println(logging.LogInfo, "Character Manager: Update Character")
				continue
			case messages.CharManager_GetCharacter:
				cm.Log.Println(logging.LogInfo, "Character Manager: Get Character")
				continue
			case messages.CharManager_GetCharacterInfo:
				cm.Log.Println(logging.LogInfo, "Character Manager: Get Character Info")
				continue
			case messages.CharManager_CheckName:
				//cm.Log.Println(logging.LogInfo, "Character Manager: Check Name")
				inUse, err := cm.CheckCharNameInUse(managerMessage.Data.(string))
				cm.OutputChannel <- messages.ConnectionManagerMessage{
					Type: messages.ConnectManager_Message_CharNameValidationResponse,
					Data: messages.CharManagerNameCheckResponse{
						NameInUse: inUse,
						Error:     err.Error(),
					},
					RecipientConsoleID: managerMessage.SenderConsoleID,
				}
			case messages.CharManager_RequestCharacterByID:
				//cm.Log.Println(logging.LogInfo, "Character Manager: Request Character By ID")
				characterInfo := cm.GetCharacterByID(managerMessage.Data.(messages.CharacterInfo).ID)
				cm.OutputChannel <- messages.ConnectionManagerMessage{Type: messages.ConnectManager_Message_CharacterRequestResponse, Data: characterInfo, RecipientConsoleID: managerMessage.SenderConsoleID}
			default:
				continue
			}
		}
	}
}
