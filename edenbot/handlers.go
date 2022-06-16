package edenbot

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (eb *EdenBot) HandleInput(started chan bool) {
	eb.Log.Println(logging.LogInfo, "Edenbot Manager Input Handler Launched")
	started <- true
	for {
		select {
		case edenbotMessage := <-eb.InputChannel:
			switch edenbotMessage.Type {
			case messages.Edenbot_Message_Null:
				//
			case messages.Edenbot_Message_Error:
				//
			default:
				//
			}
		}
	}
}
