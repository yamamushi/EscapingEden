package game

import (
	"errors"
	"github.com/yamamushi/EscapingEden/edenitems"
	"github.com/yamamushi/EscapingEden/edenutil"
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
	//gm.Log.Println(logging.LogInfo, "Loading character:", id)
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

// RemoveFromLiveCharacterList removes a character from the live character list, and broadcasts a message to all connected consoles
func (gm *GameManager) RemoveFromLiveCharacterList(ID string) {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	for i, character := range gm.ActiveCharacters {
		//gm.Log.Println(logging.LogInfo, "Checking character:", character.ConnectionID)
		if character.ID == ID || character.ConnectionID == ID {
			characterName := character.Name
			err := gm.DB.UpdateRecord("Characters", character.Record)
			if err != nil {
				gm.Log.Println(logging.LogError, "Failed to update character after removing from game manager:", err.Error())
			}
			gm.ActiveCharacters = append(gm.ActiveCharacters[:i], gm.ActiveCharacters[i+1:]...)
			response := messages.ConnectionManagerMessage{
				Type: messages.ConnectManager_Message_Broadcast,
				Data: edenutil.EdenTime{}.CurrentTimeString() + " - " + characterName + " left the world.",
			}
			gm.SendChannel <- response
			//gm.Log.Println(logging.LogInfo, "Removed character from game manager:", ID)
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

func (gm *GameManager) GetCharacterName(characterID string) (name string) {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	for i, character := range gm.ActiveCharacters {
		if character.ID == characterID {
			return gm.ActiveCharacters[i].Name
		}
	}
	return ""
}

func (gm *GameManager) GetCharacterInventory(characterID string) ([]edenitems.Item, error) {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	for i, character := range gm.ActiveCharacters {
		if character.ID == characterID {
			if len(gm.ActiveCharacters[i].Record.Inventory) == 0 {
				// Add some default items to the inventory
				inventory := []edenitems.Item{}
				for j := 0; j < 10; j++ {
					id := edenutil.GenerateID()
					wood := edenitems.Item{ID: id, Name: "Wood", Description: "A piece of wood.", Type: edenitems.ItemMaterial, Weight: 1, Stackable: true}
					err := gm.DB.AddRecord("Items", &wood)
					if err != nil {
						gm.Log.Println(logging.LogError, "Failed to add item to DB:", err.Error())
					}
					inventory = append(inventory, wood)
				}
				for j := 0; j < 10; j++ {
					id := edenutil.GenerateID()
					stone := edenitems.Item{ID: id, Name: "Stone", Description: "A piece of stone.", Type: edenitems.ItemMaterial, Weight: 1, Stackable: true}
					err := gm.DB.AddRecord("Items", &stone)
					if err != nil {
						gm.Log.Println(logging.LogError, "Failed to add item to DB:", err.Error())
					}
					inventory = append(inventory, stone)
				}
				pickid := edenutil.GenerateID()
				pickaxeAttributes := make(map[string]bool)
				pickaxeAttributes["digging"] = true
				pickaxe := edenitems.Item{ID: pickid, Name: "Pickaxe", Description: "A pickaxe that looks like it should be suitable for digging through stone", Type: edenitems.ItemTool, Weight: 5, Stackable: false, Attributes: pickaxeAttributes}
				err := gm.DB.AddRecord("Items", &pickaxe)
				if err != nil {
					gm.Log.Println(logging.LogError, "Failed to add item to DB:", err.Error())
				}
				inventory = append(inventory, pickaxe)
				gm.AssignItemHotkeys(inventory)
				gm.ActiveCharacters[i].Record.Inventory = inventory
			}
			return gm.ActiveCharacters[i].Record.Inventory, nil
		}
	}
	return []edenitems.Item{}, errors.New("character not found")
}

func (gm *GameManager) RemoveFromCharacterInventory(characterID string, itemID string) error {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	for i, character := range gm.ActiveCharacters {
		if character.ID == characterID {
			for j, item := range gm.ActiveCharacters[i].Record.Inventory {
				if item.ID == itemID {
					gm.ActiveCharacters[i].Record.Inventory = append(gm.ActiveCharacters[i].Record.Inventory[:j], gm.ActiveCharacters[i].Record.Inventory[j+1:]...)
					return nil
				}
			}
		}
	}
	return errors.New("character not found")
}

// AssignItemHotkeys is not thread safe, it is assumed that the caller has already locked the mutex!
func (gm *GameManager) AssignItemHotkeys(inventory []edenitems.Item) {
	// Assign hotkeys to items in the inventory, starting with lowercase a, and going up to z, then A to Z, then 0 to 9
	// Stackable items should be assigned the same hotkey
	// Non-stackable items should be assigned the next available hotkey
	// If there are no more hotkeys available, the item should not be assigned a hotkey
	hotkey := 'a'                        // Start with lowercase 'a'
	hotkeyMap := make(map[string]string) // Map to store assigned hotkeys for stackable items

	for i, item := range inventory {
		// Check if the item is stackable
		if item.Stackable {
			// Check if an item with the same name has been assigned a hotkey
			if stackableHotkey, ok := hotkeyMap[item.Name]; ok {
				// Assign the same hotkey to the current stackable item
				item.Hotkey = stackableHotkey
			} else {
				// Assign the next available hotkey to the item
				if hotkey <= 'z' || (hotkey >= 'A' && hotkey <= 'Z') || (hotkey >= '0' && hotkey <= '9') {
					item.Hotkey = string(hotkey)
					hotkeyMap[item.Name] = string(hotkey) // Update the hotkey map for stackable items
					hotkey++
					//log.Println("Assigned hotkey", item.Hotkey, "to item", item.Name)
				} else {
					// No more hotkeys available, don't assign a hotkey to this item
					item.Hotkey = ""
				}
			}
		} else {
			// Non-stackable item, assign a new hotkey
			if hotkey <= 'z' || (hotkey >= 'A' && hotkey <= 'Z') || (hotkey >= '0' && hotkey <= '9') {
				item.Hotkey = string(hotkey)
				hotkey++
				//log.Println("Assigned hotkey", item.Hotkey, "to item", item.Name)
			} else {
				// No more hotkeys available, don't assign a hotkey to this item
				item.Hotkey = ""
			}
		}

		// Update the item in the inventory
		inventory[i] = item
	}
}

func (gm *GameManager) AddToCharacterInventory(characterID string, item edenitems.Item) error {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	for i, character := range gm.ActiveCharacters {
		if character.ID == characterID {
			gm.ActiveCharacters[i].Record.Inventory = append(gm.ActiveCharacters[i].Record.Inventory, item)
			return nil
		}
	}
	return errors.New("character not found")
}
