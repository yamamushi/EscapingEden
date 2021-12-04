package login

import "github.com/yamamushi/EscapingEden/messages"

func (lw *LoginWindow) RegistrationSubmit(RegistrationSubmitData) *RegistrationError {
	lw.registrationSubmitMutex.Lock()
	defer lw.registrationSubmitMutex.Unlock()

	// We don't need to lock our lw.registrationErrorMutex because we're generating a new one here, not changing it.

	// We need to lock to make sure no requests are happening twice, and that
	// Our registration data is locked
	regError := &RegistrationError{
		usernameError:        "",
		passwordError:        "",
		passwordConfirmError: "",
		emailError:           "",
		rulesError:           "",
		errorRequest:         "",
	}

	if lw.registrationSubmitData.Username == "" {
		regError.usernameError = "You must enter a username."
	}
	if lw.registrationSubmitData.Password == "" {
		regError.passwordError = "You must enter a password."
	}
	if lw.registrationSubmitData.Password != lw.registrationSubmitData.PasswordConfirm ||
		lw.registrationSubmitData.PasswordConfirm == "" {
		regError.passwordConfirmError = "Your passwords do not match."
	}
	if lw.registrationSubmitData.Email == "" {
		regError.emailError = "You must enter an email."
	}
	if !lw.registrationAgreeRules {
		regError.rulesError = "You must agree to the rules before you can register."
	}
	if !regError.Empty() {
		return regError
	}

	registrationData := messages.AccountRegistrationRequest{
		Username: lw.registrationSubmitData.Username,
		Password: lw.registrationSubmitData.Password,
		Email:    lw.registrationSubmitData.Email,
	}
	windowMessage := messages.WindowMessage{Type: messages.WM_RequestRegistration, Data: registrationData}
	lw.SendToConsole(windowMessage)

	go lw.HandleReceiveChannel() // We're going to start listening for responses now

	return nil
}
