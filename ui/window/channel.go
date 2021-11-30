package window

/*
Functions related to go channel routines from the window back to the console
*/
import (
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"log"
	"strconv"
)

// Error sends a message with type error to the console
func (w *Window) Error(err string) {
	consoleMessage := &types.ConsoleMessage{Type: "error", Message: err}
	w.ConsoleSend <- consoleMessage.String()
}

// Quit sends a message with type quit to the console
func (w *Window) Quit() {
	consoleMessage := &types.ConsoleMessage{Type: "quit"}
	w.ConsoleSend <- consoleMessage.String()
}

// ForceConsoleRefresh sends a message with type console and message field refresh to the console
func (w *Window) ForceConsoleRefresh() {
	w.ResetWindowDrawings()
	message := &types.ConsoleMessage{Type: "console", Message: "refresh", WindowID: int(w.GetID())}
	w.ConsoleSend <- message.String()
}

// RequestPopupFromConsole sends a message with type popup to the console requesting for a popup to be displayed
func (w *Window) RequestPopupFromConsole(x, y, width, height int, content string) {
	log.Println("Requesting popup from console")
	config := config.NewWindowConfig(x, y, width, height, content)
	request := types.ConsoleMessage{Type: "console", Message: "popup", Options: config.String()}
	w.ConsoleSend <- request.String()
}

// RequestHelpFromConsole sends a message with type help to the console requesting for the help menu to be displayed
func (w *Window) RequestHelpFromConsole(page types.HelpPage) {
	log.Println("Requesting help from console")
	config := config.NewWindowConfig(w.ConsoleWidth/2-40, w.ConsoleHeight/2-10, 100, 20, "")
	config.Page = page
	request := types.ConsoleMessage{Type: "console", Message: "help", Options: config.String()}
	w.ConsoleSend <- request.String()
}

// RequestFlushFromConsole this tells the console we want our window to be flushed before the next draw
// What this actually does is sets a flag in the console, which is checked at draw time. It tells the console
// We want to flush the window area in its last sent buffer, which forces redrawing of the window area only
func (w *Window) RequestFlushFromConsole() {
	log.Println("Requesting help from console")
	request := types.ConsoleMessage{Type: "console", Message: "flush", Options: strconv.Itoa(int(w.GetID()))}
	w.ConsoleSend <- request.String()
}
