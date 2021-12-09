package messages

type CharacterManagerMessageType int

const (
	CharManager_None CharacterManagerMessageType = iota
	CharManager_CreateCharacter
	CharManager_DeleteCharacter
	CharManager_ListCharacters
	CharManager_UpdateCharacter
	CharManager_GetCharacter
	CharManager_GetCharacterInfo
)

type CharacterManagerMessage struct {
	Type               CharacterManagerMessageType `json:"type"`
	SenderConsoleID    string                      `json:"sender_id"`
	RecipientConsoleID string                      `json:"recipient_id"`
	Data               interface{}                 `json:"data"`
}
