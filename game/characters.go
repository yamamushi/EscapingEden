package game

import (
	"errors"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
	"sync"
)

type ActiveCharacter struct {
	// The character's ID
	ID string
	// The character's name
	Name string
	// The character's current view (this is constantly updated)
	View types.PointMap
	// The character's DB Record which we will be managing from now on
	Record *messages.CharacterInfo

	// Connection ID for the character
	ConnectionID string
	Lock         sync.Mutex
}

type ActiveCharacters []*ActiveCharacter

func (gm *GameManager) LoadCharacter(id string, consoleID string) (err error) {
	gm.Log.Println(logging.LogInfo, "Loading character:", id)
	// First load the character's info from the database
	character := messages.CharacterInfo{}
	err = gm.DB.One("Characters", "ID", id, &character)
	if err != nil {
		gm.Log.Println(logging.LogError, "Failed to load character:", err.Error())
		return err
	}
	gm.AddToLiveCharacterList(character, consoleID)
	return nil
}

func (gm *GameManager) AddToLiveCharacterList(character messages.CharacterInfo, consoleID string) {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	gm.ActiveCharacters = append(gm.ActiveCharacters, &ActiveCharacter{ID: character.ID, Name: character.Name, Record: &character, ConnectionID: consoleID})
}

func (gm *GameManager) RemoveFromLiveCharacterList(ID string) {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	for i, character := range gm.ActiveCharacters {
		gm.Log.Println(logging.LogInfo, "Checking character:", character.ConnectionID)
		if character.ID == ID || character.ConnectionID == ID {
			err := gm.DB.UpdateRecord("Characters", character.Record)
			if err != nil {
				gm.Log.Println(logging.LogError, "Failed to update character after removing from game manager:", err.Error())
			}
			gm.ActiveCharacters = append(gm.ActiveCharacters[:i], gm.ActiveCharacters[i+1:]...)
			gm.Log.Println(logging.LogInfo, "Removed character from game manager:", ID)
		}
	}
}

func (gm *GameManager) GetCharacter(characterID string) (character *messages.CharacterInfo, err error) {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	for i, character := range gm.ActiveCharacters {
		if character.ID == characterID {
			return gm.ActiveCharacters[i].Record, nil
		}
	}
	return &messages.CharacterInfo{}, errors.New("character not found")
}

// Note this does not lock the mutex, it is assumed that the caller has already locked it!
func (gm *GameManager) GetCharacterAt(X, Y int) (character *messages.CharacterInfo) {
	for i, character := range gm.ActiveCharacters {
		if character.Record.Position.X == X && character.Record.Position.Y == Y {
			return gm.ActiveCharacters[i].Record
		}
	}
	return nil
}
