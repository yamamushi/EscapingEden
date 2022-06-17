package login

import (
	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/messages"
)

func (lw *LoginWindow) RegistrationSubmit() *RegistrationError {
	lw.registrationSubmitMutex.Lock()
	defer lw.registrationSubmitMutex.Unlock()

	// We don't need to lock our lw.registrationErrorMutex because we're generating a new one here, not changing it.

	// We need to lock to make sure no requests are happening twice, and that
	// Our registration data is locked
	regError := &RegistrationError{
		usernameError:        "",
		passwordError:        "",
		passwordConfirmError: "",
		discordError:         "",
		rulesError:           "",
		errorRequest:         "",
	}

	if lw.registrationSubmitData.Username == "" {
		regError.usernameError = "You must enter a username."
	}
	if edenutil.CheckBlacklist(lw.registrationSubmitData.Username, edenutil.BlackListUsernames) {
		regError.usernameError = "That username is not allowed."
	}

	if lw.registrationSubmitData.Password == "" {
		regError.passwordError = "You must enter a password."
	}
	if lw.registrationSubmitData.Password != lw.registrationSubmitData.PasswordConfirm ||
		lw.registrationSubmitData.PasswordConfirm == "" {
		regError.passwordConfirmError = "Your passwords do not match."
	}

	// discord verification
	if lw.registrationSubmitData.DiscordID == "" {
		regError.discordError = "You must enter a Discord ID."
	}

	if edenutil.CheckBlacklist(lw.registrationSubmitData.DiscordID, edenutil.BlackListDiscordIDS) {
		regError.discordError = "Discord User is not allowed to register."
	}

	if !lw.registrationAgreeRules {
		regError.rulesError = "You must agree to the rules before you can register."
	}
	if !regError.Empty() {
		return regError
	}

	registrationData := messages.AccountRegistrationRequest{
		Username:  lw.registrationSubmitData.Username,
		Password:  lw.registrationSubmitData.Password,
		DiscordID: lw.registrationSubmitData.DiscordID,
	}
	windowMessage := messages.WindowMessage{Type: messages.WM_RequestRegistration, Data: registrationData}
	lw.SendToConsole(windowMessage)

	go lw.HandleReceiveChannel() // We're going to start listening for responses now

	return nil
}
