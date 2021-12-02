package login

import (
	"github.com/yamamushi/EscapingEden/messages"
	"log"
)

func (lw *LoginWindow) HandleReceiveChannel() {
	for {
		select {
		case windowMessage := <-lw.ConsoleReceive:
			switch windowMessage.Type {
			case messages.WM_RegistrationResponse:
				lw.registrationStatusMutex.Lock()
				defer lw.registrationStatusMutex.Unlock()

				log.Println("Login Window received registration response from console")
				lw.registrationResponse = windowMessage.Data.(messages.AccountRegistrationResponse)
				return // We launched when our registration request was submitted, now we can return since we got a response
			}
		}
	}
}
