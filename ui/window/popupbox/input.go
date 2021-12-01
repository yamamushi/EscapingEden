package popupbox

import (
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
	"log"
)

// HandleInput is a function that handles input events.
func (pb *PopupBox) HandleInput(input types.Input) {
	pb.pbMutex.Lock()
	defer pb.pbMutex.Unlock()

	if pb.GetActive() {
		log.Println("PopupBox Handling input")
	}

	switch input.Type {
	case types.InputUp:
		log.Println("PopupBox Up")
		pb.DecreaseContentPos()
		return
	case types.InputDown:
		log.Println("PopupBox Down")
		pb.IncreaseContentPos()
		return
	case types.InputReturn:
		log.Println("PopupBox Handling input return - attempting to close popup")
		message := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_ClosePopup}
		pb.ConsoleSend <- message
		log.Println("PopupBox sent close message to console")
	}

}
