package messages

/*
Window messages are sent to/from the console from/to windows.
*/

import "github.com/yamamushi/EscapingEden/ui/config"

type WindowMessageType int

const (
	WM_Null WindowMessageType = iota
	// These will send messages to the connection manager ultimately
	WM_Error
	WM_QuitConsole

	// This sends a message to the console
	WM_ConsoleCommand

	// These messages are sent to the connection manager for processing upstream
	// We can't do anything with them here, so we just pass them along
	WM_ParseChat
	WM_RequestRegistration
	WM_RequestLogin
	WM_RequestForgotPassword
	WM_ValidateForgotPassword
	WM_ProcessForgotPassword
	WM_BadLoginAttempt

	// These are messages that are parsed by the windows themselves if they receive an event
	WM_RegistrationResponse
	WM_LoginResponse
	WM_PasswordResetValidateResponse
	WM_PasswordResetProcessResponse
)

type WindowMessageCommand int

const (
	WMC_Null WindowMessageCommand = iota
	WMC_NewPopup
	WMC_ClosePopup
	WMC_ToggleHelp
	WMC_RefreshConsole
	WMC_FlushConsoleBuffer
)

type WindowMessage struct {
	Type    WindowMessageType
	Command WindowMessageCommand

	Data         interface{}
	TargetID     config.WindowID
	PopupOptions config.WindowConfig
	HelpOptions  config.WindowConfig
}
