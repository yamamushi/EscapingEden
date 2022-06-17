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
			case messages.Edenbot_Message_NewRegistration:
				//
			case messages.Edenbot_Message_ForgotPassword:
				//eb.Log.Println(logging.LogInfo, "Edenbot Manager received a forgot password request")
				eb.ForgotPassword(edenbotMessage.Data.(messages.AccountForgotPasswordData))
			case messages.Edenbot_Message_PlayerLoggedIn:
				//
			case messages.Edenbot_Message_PlayerLoggedOut:
				//
			case messages.Edenbot_Message_CharacterCreated:
				//
			case messages.Edenbot_Message_Shutdown:
				//
				eb.Log.Println(logging.LogInfo, "Edenbot received shutdown message, shutting down")
				err := eb.Shutdown()
				if err != nil {
					eb.Log.Println(logging.LogError, "Error shutting down Edenbot: ", err)
				}
				eb.Log.Println(logging.LogInfo, "Edenbot shut down successfully")
			default:
				//
			}
		}
	}
}
