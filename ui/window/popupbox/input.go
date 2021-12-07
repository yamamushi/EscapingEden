package popupbox

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
)

// HandleInput is a function that handles input events.
func (pb *PopupBox) HandleInput(input types.Input) {
	pb.pbMutex.Lock()
	defer pb.pbMutex.Unlock()

	if pb.GetActive() {
		pb.Log.Println(logging.LogInfo, "PopupBox Handling input")
	}

	switch input.Type {
	case types.InputUp:
		pb.Log.Println(logging.LogInfo, "PopupBox Up")
		pb.DecreaseContentPos()
		return
	case types.InputDown:
		pb.Log.Println(logging.LogInfo, "PopupBox Down")
		pb.IncreaseContentPos()
		return
	case types.InputReturn:
		pb.Log.Println(logging.LogInfo, "PopupBox Handling input return - attempting to close popup")
		message := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_ClosePopup}
		pb.ConsoleSend <- message
		pb.Log.Println(logging.LogInfo, "PopupBox sent close message to console")
	}

}
