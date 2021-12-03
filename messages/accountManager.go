package messages

// Account is used for registration/db records/etc to store player accounts.
// It has no knowledge of any other data.
type Account struct {
	ID             string `storm:"index"`  // Indexed Unique ID for the account, we use UUID, not the auto-increment, ever.
	Username       string `storm:"unique"` // Username of the account, must be unique.
	Email          string `storm:"unique"` // Email of the account, must be unique.
	HashedPassword string // Hashed password of the account.
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
