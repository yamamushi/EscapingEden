package popupbox

import (
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
		message := types.ConsoleMessage{Type: "popupbox", Message: "close"}
		pb.ConsoleSend <- message.String()
		log.Println("PopupBox sent close message to console")
	}

}
