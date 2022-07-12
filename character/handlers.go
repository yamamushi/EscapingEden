package character

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"log"
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
				log.Println("Character Manager: Create Character")
				continue
			case messages.CharManager_DeleteCharacter:
				log.Println("Character Manager: Delete Character")
				continue
			case messages.CharManager_ListCharacters:
				log.Println("Character Manager: List Characters")
				continue
			case messages.CharManager_UpdateCharacter:
				log.Println("Character Manager: Update Character")
				continue
			case messages.CharManager_GetCharacter:
				log.Println("Character Manager: Get Character")
				continue
			case messages.CharManager_GetCharacterInfo:
				log.Println("Character Manager: Get Character Info")
				continue
			default:
				continue
			}
		}
	}
}
