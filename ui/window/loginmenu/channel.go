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
				log.Println("Registration response received")
				if windowMessage.Data.(*messages.AccountRegistrationResponse).Success {
					log.Println(windowMessage.Data.(*messages.AccountRegistrationResponse).Message)
				} else {
					log.Println("Something went wrong")
				}
			}
		}
	}
}
