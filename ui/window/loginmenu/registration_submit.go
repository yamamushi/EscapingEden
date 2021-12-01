package login

type RegistrationError struct {
	UsernameError        string
	PasswordError        string
	PasswordConfirmError string
	EmailError           string
	Error                error
}

type RegistrationSubmitData struct {
	Username        string
	Password        string
	PasswordConfirm string
	Email           string
}

func (lw *LoginWindow) RegistrationSubmit(RegistrationSubmitData) *RegistrationError {
	lw.registrationSubmitMutex.Lock()
	defer lw.registrationSubmitMutex.Unlock()
	// We need to lock to make sure no requests are happening twice, and that
	// Our registration data is locked

	return &RegistrationError{
		UsernameError:        "This is a test error",
		PasswordError:        "",
		PasswordConfirmError: "",
		EmailError:           "",
		Error:                nil,
	}
}
