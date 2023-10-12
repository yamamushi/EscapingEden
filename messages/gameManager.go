package messages

type GMErrorType int

const (
	GMError_Null GMErrorType = iota
	GMError_DBError
	GMError_NameAlreadyExists
	GMError_DeleteGameacterError
	GMError_CreateGameacterError_Default
	GMError_HistoryUpdatePermissionError
	GMError_InvalidName
	GMError_CreateGameacterError_InvalidColor
	GMError_CreateGameacterError_InventoryCreateError
)

func (GMe GMErrorType) Error() string {
	switch GMe {
	case GMError_Null:
		return "Null Error"
	case GMError_DBError:
		return "Database Error"
	case GMError_NameAlreadyExists:
		return "Name Already Exists"
	case GMError_InvalidName:
		return "Invalid Name"
	case GMError_HistoryUpdatePermissionError:
		return "History Update Permission Error"
	}
	return "Unknown Error"
}

type GameManagerMessageType int

const (
	GameManager_None GameManagerMessageType = iota
	GameManager_GetCharacterPosition
)

type GameManagerMessage struct {
	Type               GameManagerMessageType `json:"type"`
	SenderConsoleID    string                 `json:"sender_id"`
	RecipientConsoleID string                 `json:"recipient_id"`
	Data               interface{}            `json:"data"`
}
