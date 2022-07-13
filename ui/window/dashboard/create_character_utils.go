package dashboard

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (dw *DashboardWindow) CreateCharacter() {
	dw.Log.Println(logging.LogInfo, "Sending create character request to the upstream managers")
	message := messages.WindowMessage{Type: messages.WM_RequestCharNameValidation, Data: dw.charCreatorName}
	dw.SendToConsole(message)
	go dw.HandleReceiveChannel()
}

func (dw *DashboardWindow) LoginCharacter(charID string) {

}
