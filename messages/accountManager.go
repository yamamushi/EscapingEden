package messages

// Account is used for registration/db records/etc to store player accounts.
// It has no knowledge of any other data.
type Account struct {
	ID                  string `storm:"index"`  // Indexed Unique ID for the account, we use UUID, not the auto-increment, ever.
	Username            string `storm:"unique"` // Username of the account, must be unique.
	DiscordTag          string `storm:"unique"` // Discord Tag of the account, must be unique in username#0000 format.
	DiscordID           string `storm:"unique"` // Discord ID of the account, must be unique.
	HashedPassword      string // Hashed password of the account.
	ValidationStatus    int    // 0 = pending, 1 = validated
	ValidationCode      string // The validation code for the account.
	PasswordResetStatus int    // 0 = no reset requested, 1 = reset requested
	PasswordResetCode   string // The temporary reset code for the account.
}

type AccountRegistrationRequest struct {
	Username  string
	Password  string // Plaintext here, we hash it later down the chain
	DiscordID string
}

type AccountRegistrationResponse struct {
	Error          AMErrorType
	ValidationCode string
}

type AccountLoginRequest struct {
	Username string
	Password string // Plaintext here, we hash it later down the chain
}

type AccountLoginResponse struct {
	Account Account
	Error   AMErrorType
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
	Type            AccountManagerMessageType
	SenderSessionID string
	Data            interface{}
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
	AMError_UserNotInServer
	AMError_UsernameAlreadyExists
	AMError_DiscordAlreadyExists
	AMError_DiscordMessageError
	AMError_PendingValidation
	AMError_DBError
	AMError_InvalidPassword
	AMError_InvalidDiscordID
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
	case AMError_DiscordAlreadyExists:
		return "DiscordID Already Exists"
	case AMError_DiscordMessageError:
		return "Discord Message Error, please check your private message settings for the server."
	case AMError_InvalidPassword:
		return "Invalid Password"
	case AMError_InvalidDiscordID:
		return "Invalid DiscordID"
	case AMError_InvalidUsername:
		return "Invalid Username"
	case AMError_PendingValidation:
		return "This Account is currently pending validation, please check your discord messages."
	case AMError_UserNotInServer:
		return "User is not in the discord server, please join the server and try registering again."
	default:
		return "Unknown Error"
	}
}
