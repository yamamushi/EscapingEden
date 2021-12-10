package login

type RegistrationError struct {
	usernameError        string
	passwordError        string
	passwordConfirmError string
	discordError         string
	rulesError           string
	errorRequest         string
}

func (r RegistrationError) IsEmpty() bool {
	return r.usernameError == "" && r.passwordError == "" && r.passwordConfirmError == "" && r.discordError == "" && r.rulesError == ""
}

func (r RegistrationError) UsernameError() string {
	if r.usernameError != "" {
		return "Error: " + r.usernameError
	} else {
		return ""
	}
}

func (r RegistrationError) PasswordError() string {
	if r.passwordError != "" {
		return "Error: " + r.passwordError
	} else {
		return ""
	}
}

func (r RegistrationError) PasswordConfirmError() string {
	if r.passwordConfirmError != "" {
		return "Error: " + r.passwordConfirmError
	} else {
		return ""
	}
}

func (r RegistrationError) DiscordError() string {
	if r.discordError != "" {
		return "Error: " + r.discordError
	} else {
		return ""
	}
}

func (r RegistrationError) RulesError() string {
	if r.rulesError != "" {
		return "Error: " + r.rulesError
	} else {
		return ""
	}
}

func (r RegistrationError) ErrorRequest() string {
	if r.errorRequest != "" {
		return r.errorRequest
	} else {
		return ""
	}
}

func (r RegistrationError) Empty() bool {
	if r.usernameError != "" || r.passwordError != "" || r.passwordConfirmError != "" || r.discordError != "" || r.rulesError != "" {
		return false
	} else {
		return true
	}
}

type RegistrationSubmitData struct {
	Username        string
	Password        string
	PasswordConfirm string
	DiscordID       string
}
