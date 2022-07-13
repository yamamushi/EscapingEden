package dashboard

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/util"
)

func (dw *DashboardWindow) HandleReceiveChannel() {
	for {
		select {
		case windowMessage := <-dw.ConsoleReceive:
			switch windowMessage.Type {
			case messages.WM_RequestCharNameValidationResponse:
				dw.Log.Println(logging.LogInfo, "Received WM_RequestCharNameValidationResponse")
				data := windowMessage.Data.(messages.CharManagerNameCheckResponse)
				if data.NameInUse {
					dw.Log.Println(logging.LogInfo, "Character Name In Use")
					if data.Error != "" {
						dw.Log.Println(logging.LogInfo, data.Error)
					}
					dw.charCreatorUsernameError = "Character Name In Use"
					dw.characterCreatorState = CharacterCreatorCharacterDetails
					dw.RequestFlushFromConsole()
				} else {
					dw.Log.Println(logging.LogInfo, "Character Name Not In Use")
					colorCode := util.ColorCode{0, 255, 0}
					bgColorCode := util.ColorCode{255, 0, 0}
					charInfo := messages.CharacterInfo{Name: dw.charCreatorName, ID: "tmpID", FGColor: colorCode, BGColor: bgColorCode, InventoryID: "tmpInventoryID"}
					dw.LoginCharacter(charInfo)
				}
			}
		}
	}
}
