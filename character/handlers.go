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
			case messages.CharManager_CreateCharacter:
				//
			default:
				//
			}
		}
	}
}
