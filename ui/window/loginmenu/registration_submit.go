package login

type RegistrationErrorData struct {
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

func (lw *LoginWindow) RegistrationSubmit(RegistrationSubmitData) RegistrationErrorData {
	lw.registrationSubmitMutex.Lock()
	defer lw.registrationSubmitMutex.Unlock()
	// We need to lock to make sure no requests are happening twice, and that
	// Our registration data is locked

	return RegistrationErrorData{
		UsernameError:        "This is a test error",
		PasswordError:        "",
		PasswordConfirmError: "",
		EmailError:           "",
		Error:                nil,
	}
}
