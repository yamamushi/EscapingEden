package login

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"time"
)

func (lw *LoginWindow) loginSubmit() {
	lw.loginSubmitMutex.Lock()
	defer lw.loginSubmitMutex.Unlock()

	lw.loginSubmitData.Error = "" // clear error

	// If less than 3 seconds have passed since the last submit, ignore this submit
	if time.Since(lw.loginLastAttempt) < time.Second*3 {
		lw.loginSubmitData.Error = "You must wait before your next login attempt"
		lw.loginResponseReceived = true
		return
	}

	if lw.loginAttempts >= 3 {
		// We notify the console (subsequently the connection manager) that this connection
		// Has made a bad login attempt, we track these as too many will result in a disconnect and
		// Eventually an IP ban (temporarily).
		lw.Log.Println(logging.LogInfo, "Sending bad login record to the connection manager")
		windowMessage := messages.WindowMessage{Type: messages.WM_BadLoginAttempt}
		lw.SendToConsole(windowMessage)
		lw.loginAttempts = 0
		lw.loginSubmitData.Error = "Too many bad login attempts, please check your credentials and try again."
		lw.loginResponseReceived = true
		return
	}

	// check if username is empty
	if lw.loginSubmitData.Username == "" {
		lw.loginSubmitData.Error = "Username cannot be empty"
		lw.loginResponseReceived = true
		return
	}
	// check if password is empty
	if lw.loginSubmitData.Password == "" {
		lw.loginSubmitData.Error = "Password cannot be empty"
		lw.loginResponseReceived = true
		return
	}

	loginData := messages.AccountLoginRequest{
		Username: lw.loginSubmitData.Username,
		Password: lw.loginSubmitData.Password,
	}
	windowMessage := messages.WindowMessage{Type: messages.WM_RequestLogin, Data: loginData}
	lw.SendToConsole(windowMessage)
	lw.loginLastAttempt = time.Now()

	go lw.HandleReceiveChannel() // We're going to start listening for responses now

	return
}
