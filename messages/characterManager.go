package messages

type CMErrorType int

const (
	CMError_Null CMErrorType = iota
	CMError_DBError
	CMError_NameAlreadyExists
	CMError_DeleteCharacterError
	CMError_CreateCharacterError_Default
	CMError_InvalidName
	CMError_CreateCharacterError_InvalidColor
	CMError_CreateCharacterError_InventoryCreateError
)

func (cme CMErrorType) Error() string {
	switch cme {
	case CMError_Null:
		return "Null Error"
	case CMError_DBError:
		return "Database Error"
	case CMError_NameAlreadyExists:
		return "Name Already Exists"
	}
	return "Unknown Error"
}

type CharacterManagerMessageType int

const (
	CharManager_None CharacterManagerMessageType = iota
	CharManager_CreateCharacter
	CharManager_DeleteCharacter
	CharManager_ListCharacters
	CharManager_UpdateLoginHistory
	CharManager_UpdateCharacter
	CharManager_GetCharacter
	CharManager_GetCharacterInfo
	CharManager_CheckName
)

type CharacterManagerMessage struct {
	Type               CharacterManagerMessageType `json:"type"`
	SenderConsoleID    string                      `json:"sender_id"`
	RecipientConsoleID string                      `json:"recipient_id"`
	Data               interface{}                 `json:"data"`
}

type CharManagerNameCheckResponse struct {
	NameInUse bool   `json:"name_in_use"`
	Error     string `json:"error"`
}
