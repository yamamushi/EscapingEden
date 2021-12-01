package window

/*
Functions related to go channel routines from the window back to the console
*/
import (
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"log"
)

// Error sends a message with type error to the console
func (w *Window) Error(err string) {
	request := messages.WindowMessage{Type: messages.WM_Error, MessageContent: err}
	w.ConsoleSend <- request
}

// Quit sends a message with type quit to the console
func (w *Window) Quit() {
	request := messages.WindowMessage{Type: messages.WM_QuitConsole}
	w.ConsoleSend <- request
}

// ForceConsoleRefresh sends a message with type console and message field refresh to the console
func (w *Window) ForceConsoleRefresh() {
	w.ResetWindowDrawings()
	request := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_RefreshConsole, TargetID: w.GetID()}
	w.ConsoleSend <- request
}

// RequestPopupFromConsole sends a message with type popup to the console requesting for a popup to be displayed
func (w *Window) RequestPopupFromConsole(x, y, width, height int, content string) {
	log.Println("Requesting popup from console")
	popupConfig := config.NewWindowConfig(x, y, width, height, content)
	request := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_NewPopup, PopupOptions: *popupConfig}
	w.ConsoleSend <- request
}

// RequestHelpFromConsole sends a message with type help to the console requesting for the help menu to be displayed
func (w *Window) RequestHelpFromConsole(page types.HelpPage) {
	log.Println("Requesting help from console")
	helpConfig := config.NewWindowConfig(w.ConsoleWidth/2-40, w.ConsoleHeight/2-10, 100, 20, "")
	helpConfig.Page = page
	request := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_ToggleHelp, HelpOptions: *helpConfig}
	w.ConsoleSend <- request
}

// RequestFlushFromConsole this tells the console we want our window to be flushed before the next draw
// What this actually does is sets a flag in the console, which is checked at draw time. It tells the console
// We want to flush the window area in its last sent buffer, which forces redrawing of the window area only
func (w *Window) RequestFlushFromConsole() {
	log.Println("Requesting flush of window area from console")
	request := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_FlushConsoleBuffer, TargetID: w.GetID()}
	w.ConsoleSend <- request
}
