package messages

// Account is used during registration to create a new player account.
// It is also used during logins to retrieve the character list for an account.
type Account struct {
	ID             string `storm:"index"`
	Username       string `storm:"unique"`
	Email          string `storm:"unique"`
	HashedPassword string

	Characters map[string]string `storm:"unique"` // Character ID -> Character Name
}

type AccountRegistrationRequest struct {
	Username string
	Password string // Plaintext here, we hash it later down the chain
	Email    string
}

type AccountRegistrationResponse struct {
	Success bool
	Message string
}

type AccountManagerMessageType int

const (
	AccountManager_Message_Null AccountManagerMessageType = iota
	AccountManager_Message_Register
	AccountManager_Message_Login
	AccountManager_Message_Logout
	AccountManager_Message_GetCharacters
)

type AccountManagerMessage struct {
	Type                 AccountManagerMessageType
	SenderSessionID      string
	AccountResult        Account
	RegistrationRequest  AccountRegistrationRequest
	RegistrationResponse AccountRegistrationResponse
}
