package login

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (lw *LoginWindow) HandleReceiveChannel() {
	for {
		select {
		case windowMessage := <-lw.ConsoleReceive:
			switch windowMessage.Type {
			case messages.WM_RegistrationResponse:
				lw.registrationStatusMutex.Lock()
				defer lw.registrationStatusMutex.Unlock()

				lw.Log.Println(logging.LogInfo, "Login Window received registration response from console")
				lw.registrationResponse = windowMessage.Data.(messages.AccountRegistrationResponse)
				lw.registrationResponseReceived = true
				lw.RequestFlushFromConsole()
				return // We launched when our registration request was submitted, now we can return since we got a response
			}
		}
	}
}
