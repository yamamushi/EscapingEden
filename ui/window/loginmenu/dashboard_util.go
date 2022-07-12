package login

import "github.com/yamamushi/EscapingEden/messages"

// NotifyConsoleLoggedOut is called when the user logs out
func (lw *LoginWindow) NotifyConsoleLoggedOut() {
	// Create a console message with type Console_Message_LoginUser, we don't pack any data with this message (yet, TBD)
	msg := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_SetLoggedOut, TargetID: lw.GetID()}
	// Send the message to the console so that we can enable the full dashboard control
	lw.SendToConsole(msg)
}

// NotifyConsoleLoggedIn is called when the user logs in
func (lw *LoginWindow) NotifyConsoleLoggedIn(info messages.UserInfo) {
	// Create a console message with type Console_Message_LoginUser, we don't pack any data with this message (yet, TBD)
	msg := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_SetLoggedIn, TargetID: lw.GetID(), Data: info}
	// Send the message to the console so that we can enable the full dashboard control
	lw.SendToConsole(msg)
}
