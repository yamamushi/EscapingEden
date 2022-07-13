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
	// Create a console message with type Console_Message_LoginUser, we don't pack any data with this message (yet, TBD)
	msg := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_SetCharacterLoggedIn, TargetID: dw.GetID(), Data: messages.CharacterInfo{Name: dw.charCreatorName, ID: charID}}
	// Send the message to the console so that we can enable the full dashboard control
	dw.SendToConsole(msg)
}

// NotifyConsoleLoggedOut is called when the user logs out
func (dw *DashboardWindow) LogoutCharacter() {
	// Create a console message with type Console_Message_LoginUser, we don't pack any data with this message (yet, TBD)
	msg := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_SetCharacterLoggedOut, TargetID: dw.GetID()}
	// Send the message to the console so that we can enable the full dashboard control
	dw.SendToConsole(msg)
}
