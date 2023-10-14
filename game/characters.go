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
}

type ActiveCharacters []*ActiveCharacter

func (gm *GameManager) LoadCharacter(id string) (err error) {
	gm.Log.Println(logging.LogInfo, "Loading character:", id)
	// First load the character's info from the database
	character := messages.CharacterInfo{}
	err = gm.DB.One("Characters", "ID", id, &character)
	if err != nil {
		gm.Log.Println(logging.LogError, "Failed to load character:", err.Error())
		return err
	}
	gm.AddToLiveCharacterList(character)
	return nil
}

func (gm *GameManager) AddToLiveCharacterList(character messages.CharacterInfo) {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	gm.ActiveCharacters = append(gm.ActiveCharacters, &ActiveCharacter{ID: character.ID, Name: character.Name, Record: &character})
}

func (gm *GameManager) RemoveFromLiveCharacterList(characterID string) {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	for i, character := range gm.ActiveCharacters {
		if character.ID == characterID {
			gm.ActiveCharacters = append(gm.ActiveCharacters[:i], gm.ActiveCharacters[i+1:]...)
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

func (gm *GameManager) GetCharacterAt(X, Y int) (character *messages.CharacterInfo) {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	for i, character := range gm.ActiveCharacters {
		if character.Record.Position.X == X && character.Record.Position.Y == Y {
			return gm.ActiveCharacters[i].Record
		}
	}
	return &messages.CharacterInfo{}
}
