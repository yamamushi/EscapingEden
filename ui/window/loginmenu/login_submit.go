package login

import (
	"github.com/yamamushi/EscapingEden/messages"
	"time"
)

func (lw *LoginWindow) LoginSubmit() {
	lw.loginSubmitMutex.Lock()
	defer lw.loginSubmitMutex.Unlock()

	lw.loginSubmitData.Error = "" // clear error

	// If less than 3 seconds have passed since the last submit, ignore this submit
	if time.Since(lw.loginLastAttempt) < time.Second*3 {
		lw.loginSubmitData.Error = "You must wait before your next login attempt"
		return
	}

	if lw.loginAttempts >= 3 {
		// Right now we just hang on the user, but we actually want to track these login attempts
		// And start blacklisting IPs if they are making too many login attempts in a short period of time.
		lw.loginSubmitData.Error = "You have exceeded the maximum number of login attempts."
		return
	}

	// check if username is empty
	if lw.loginSubmitData.Email == "" {
		lw.loginSubmitData.Error = "Email cannot be empty"
		return
	}
	// check if password is empty
	if lw.loginSubmitData.Password == "" {
		lw.loginSubmitData.Error = "Password cannot be empty"
		return
	}

	loginData := messages.AccountLoginRequest{
		Email:    lw.loginSubmitData.Email,
		Password: lw.loginSubmitData.Password,
	}
	windowMessage := messages.WindowMessage{Type: messages.WM_RequestLogin, Data: loginData}
	lw.SendToConsole(windowMessage)
	lw.loginLastAttempt = time.Now()

	go lw.HandleReceiveChannel() // We're going to start listening for responses now

	return
}
