package window

/*
Functions related to go channel routines from the window back to the console
*/
import (
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"log"
)

func (w *Window) Error(err string) {
	consoleMessage := &types.ConsoleMessage{Type: "error", Message: err}
	w.ConsoleSend <- consoleMessage.String()
}

func (w *Window) Quit() {
	consoleMessage := &types.ConsoleMessage{Type: "quit"}
	w.ConsoleSend <- consoleMessage.String()
}

func (w *Window) HandleReceive(message types.ConsoleMessage) {
	w.ConsoleReceive <- message.String()
}

func (w *Window) ForceConsoleRefresh() {
	w.ResetWindowDrawings()
	message := &types.ConsoleMessage{Type: "console", Message: "refresh", WindowID: int(w.GetID())}
	w.ConsoleSend <- message.String()
}

func (w *Window) RequestPopupFromConsole(x, y, width, height int, content string) {
	log.Println("Requesting popup from console")
	config := config.NewWindowConfig(x, y, width, height, content)
	log.Println(config.String())
	request := types.ConsoleMessage{Type: "console", Message: "popup", Options: config.String()}
	log.Println(request.String())
	w.ConsoleSend <- request.String()
}

func (w *Window) RequestHelpFromConsole(page types.HelpPage) {
	log.Println("Requesting help from console")
	config := config.NewWindowConfig(w.ConsoleWidth/2-40, w.ConsoleHeight/2-10, 100, 20, "")
	log.Println(page)
	config.Page = page
	log.Println(config.String())
	request := types.ConsoleMessage{Type: "console", Message: "help", Options: config.String()}
	log.Println(request.String())
	w.ConsoleSend <- request.String()
}
