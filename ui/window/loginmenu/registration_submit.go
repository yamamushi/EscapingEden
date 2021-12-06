package login

import (
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/yamamushi/EscapingEden/messages"
	"io/ioutil"
	"strings"
)

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
	if lw.CheckBlacklist(lw.registrationSubmitData.Username) {
		regError.usernameError = "That username is not allowed."
	}

	if lw.registrationSubmitData.Password == "" {
		regError.passwordError = "You must enter a password."
	}
	if lw.registrationSubmitData.Password != lw.registrationSubmitData.PasswordConfirm ||
		lw.registrationSubmitData.PasswordConfirm == "" {
		regError.passwordConfirmError = "Your passwords do not match."
	}

	// email verification
	if lw.registrationSubmitData.Email == "" {
		regError.emailError = "You must enter an email."
	}

	verifier := emailverifier.NewVerifier().EnableAutoUpdateDisposable()
	ret, err := verifier.Verify(lw.registrationSubmitData.Email)
	if err != nil {
		regError.emailError = "Invalid email."
	}
	if ret.Disposable {
		regError.emailError = "Disposable emails are not allowed."
	}
	if !ret.Syntax.Valid {
		regError.emailError = "Invalid email."
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

func (lw *LoginWindow) CheckBlacklist(username string) bool {
	// Open our blacklist file - hardcoded for now, I want to change this later
	blacklist, err := ioutil.ReadFile("assets/blacklist/usernames.txt")
	if err != nil {
		return false // If we can't open the file, we don't have a blacklist
	}
	// Split the file into lines
	blacklistLines := strings.Split(string(blacklist), "\n")
	// Check if the username is in the blacklist
	for _, line := range blacklistLines {
		// if line begins with a # we ignore it
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.TrimSpace(line) == username {
			return true
		}
	}
	return false
}
