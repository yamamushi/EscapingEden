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
				cm.Log.Println(logging.LogInfo, "Character Manager: Create Character")
				// TODO - Replace with actual character creation and storage in database
				cm.OutputChannel <- messages.ConnectionManagerMessage{Type: messages.ConnectManager_Message_CharacterCreationResponse, Data: managerMessage.Data, RecipientConsoleID: managerMessage.SenderConsoleID}
				cm.Log.Println(logging.LogInfo, "Character Manager: Create Character response sent")
				continue
			case messages.CharManager_DeleteCharacter:
				cm.Log.Println(logging.LogInfo, "Character Manager: Delete Character")
				continue
			case messages.CharManager_ListCharacters:
				cm.Log.Println(logging.LogInfo, "Character Manager: List Characters")
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
				inUse, err := cm.CheckCharNameInUse(managerMessage.Data.(string))
				cm.OutputChannel <- messages.ConnectionManagerMessage{Type: messages.ConnectManager_Message_CharNameValidationResponse, Data: messages.CharManagerNameCheckResponse{NameInUse: inUse, Error: err.Error()}, RecipientConsoleID: managerMessage.SenderConsoleID}

			default:
				continue
			}
		}
	}
}
