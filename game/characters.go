package game

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
)

type ActiveCharacter struct {
	// The character's ID
	ID string
	// The character's name
	Name string
	// The character's current view (this is constantly updated)
	View types.PointMap
	// The character's DB Record which we will be managing from now on
	Record messages.CharacterInfo
}

type ActiveCharacters []ActiveCharacter

func (gw *GameManager) LoadCharacter(id string) {
	gw.Log.Println(logging.LogInfo, "Loading character:", id)
	// First load the character's info from the database
	character := messages.CharacterInfo{}
	err := gw.DB.One("Characters", "ID", id, &character)
	if err != nil {
		gw.Log.Println(logging.LogError, "Failed to load character:", err.Error())
		return
	}

}
