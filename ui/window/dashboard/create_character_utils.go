package dashboard

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/util"
	"time"
)

func (dw *DashboardWindow) ValidateCharacterName() {
	dw.Log.Println(logging.LogInfo, "Sending create character request to the upstream managers")
	message := messages.WindowMessage{Type: messages.WM_RequestCharNameValidation, Data: dw.charCreatorName}
	dw.SendToConsole(message)
	go dw.HandleReceiveChannel()
}

func (dw *DashboardWindow) CreateCharacter() {
	colorRed := util.ColorCode{R: 255, G: 0, B: 0}
	colorGreen := util.ColorCode{R: 0, G: 255, B: 0}
	colorBlue := util.ColorCode{R: 0, G: 0, B: 255}
	colorCode := util.ColorCode{}
	if dw.charColorOption == 0 {
		colorCode = colorRed
	}
	if dw.charColorOption == 1 {
		colorCode = colorBlue
	}
	if dw.charColorOption == 2 {
		colorCode = colorGreen
	}

	charInfo := messages.CharacterInfo{
		UserID:        dw.UserInfo.ID,
		Name:          dw.charCreatorName,
		FGColor:       colorCode,
		FirstLogin:    1,
		LastLoginTime: time.Now(),
	}

	msg := messages.WindowMessage{Type: messages.WM_RequestCharacterCreation, Data: charInfo}
	// Send the message to the console so that we can enable the full dashboard control
	dw.SendToConsole(msg)
}

func (dw *DashboardWindow) GetCharacterByID(charID string) {
	//dw.Log.Println(logging.LogInfo, "GetCharacterByID:", charID)
	charInfo := messages.CharacterInfo{ID: charID, UserID: dw.UserInfo.ID}
	msg := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_RequestCharacterByID, Data: charInfo}
	dw.SendToConsole(msg)
}

func (dw *DashboardWindow) LoginCharacter(charInfo messages.CharacterInfo) {
	// Create a console message with type Console_Message_LoginUser, we don't pack any data with this message (yet, TBD)
	msg := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_SetCharacterLoggedIn, TargetID: dw.GetID(), Data: charInfo}
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

func (dw *DashboardWindow) ValidateCharacterLogin(charInfo messages.CharacterInfo) {
	// Create a console message with type WMC_RequestCharacterHistoryUpdate
	//dw.Log.Println(logging.LogInfo, "ValidateCharacterLogin")
	msg := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_RequestCharacterHistoryUpdate, TargetID: dw.GetID(), Data: charInfo}
	// Send the message to the console so that we can enable the full dashboard control
	dw.SendToConsole(msg)
	//dw.Log.Println(logging.LogInfo, "ValidateCharacterLogin message sent")

}
