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
	Error AMErrorType
}

type AccountLoginResponse struct {
	Error AMErrorType
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

// AMErrorType is used to define various errors that can occur during account management.
// It is used by anything that expects its messages to get to the Account Manager (registration/login specifically)
// In order to parse the errors it receives.
type AMErrorType int

const (
	AMError_Null AMErrorType = iota
	AMError_SystemError
	AMError_AccountAlreadyExists
	AMError_AccountDoesNotExist
	AMError_UsernameAlreadyExists
	AMError_EmailAlreadyExists
	AMError_InvalidPassword
	AMError_InvalidEmail
	AMError_InvalidUsername
)

func (ame AMErrorType) Error() string {
	switch ame {
	case AMError_Null:
		return "Null Error"
	case AMError_SystemError:
		return "System Error"
	case AMError_AccountAlreadyExists:
		return "Account Already Exists"
	case AMError_AccountDoesNotExist:
		return "Account Does Not Exist"
	case AMError_UsernameAlreadyExists:
		return "Username Already Exists"
	case AMError_EmailAlreadyExists:
		return "Email Already Exists"
	case AMError_InvalidPassword:
		return "Invalid Password"
	case AMError_InvalidEmail:
		return "Invalid Email"
	case AMError_InvalidUsername:
		return "Invalid Username"
	default:
		return "Unknown Error"
	}
}
