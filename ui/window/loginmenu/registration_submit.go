package login

func (lw *LoginWindow) RegistrationSubmit(RegistrationSubmitData) *RegistrationError {
	lw.registrationSubmitMutex.Lock()
	defer lw.registrationSubmitMutex.Unlock()
	// We need to lock to make sure no requests are happening twice, and that
	// Our registration data is locked
	regError := &RegistrationError{
		usernameError:        "Everything went fine, this is a test",
		passwordError:        "",
		passwordConfirmError: "",
		emailError:           "",
		rulesError:           "",
		errorRequest:         "",
	}

	var inputError bool
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
	if inputError {
		return regError
	}

	return regError
}
