package game

import (
	"errors"
	"sync"

	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/edenutil"
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

	// Debug logging to see what map IDs are loaded from database
	gm.Log.Println(logging.LogInfo, "Loaded character", character.Name, "from database: position", character.Position.X, character.Position.Y, "CurrentMapID", character.CurrentMapID, "Position.MapChunkID", character.Position.MapChunkID)

	// Automatically sync character inventory with JSON definitions
	if len(character.Inventory) > 0 {
		character.Inventory = gm.SyncInventoryWithJSON(character.Inventory)
		// Save the updated character back to the database
		err = gm.DB.UpdateRecord("Characters", &character)
		if err != nil {
			gm.Log.Println(logging.LogWarn, "Failed to save synced character inventory:", err.Error())
		}
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
			return // Prevents continuing iteration on modified slice
		}
	}
}

func (gm *GameManager) MovePlayerToMap(characterID string, currentMapID string) error {
	character, err := gm.GetCharacter(characterID)
	if err != nil {
		gm.Log.Println(logging.LogError, "Failed to get character:", err.Error())
		return err
	}

	character.CurrentMapID = currentMapID
	err = gm.DB.UpdateRecord("Characters", character)
	if err != nil {
		gm.Log.Println(logging.LogError, "Failed to update character after moving to new mapchunk", err.Error())
		return err
	}
	//gm.Log.Println(logging.LogInfo, "Map Transfer", character.CurrentMapID)
	return nil
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
func (gm *GameManager) GetCharacterAt(chunk *MapChunk, X, Y int) (character *messages.CharacterInfo) {
	for i, character := range gm.ActiveCharacters {
		if character.Record.Position.X == X && character.Record.Position.Y == Y && character.Record.CurrentMapID == chunk.ID {
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

func (gm *GameManager) GetCharacterInventory(characterID string) ([]edentypes.Item, error) {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()
	for i, character := range gm.ActiveCharacters {
		if character.ID == characterID {
			if len(gm.ActiveCharacters[i].Record.Inventory) == 0 {
				// Create default inventory using JSON-based item system
				inventory := gm.CreateDefaultInventory()
				gm.AssignItemHotkeys(inventory)
				gm.ActiveCharacters[i].Record.Inventory = inventory
			} else {
				// Automatically sync existing inventory items with JSON definitions
				gm.ActiveCharacters[i].Record.Inventory = gm.SyncInventoryWithJSON(gm.ActiveCharacters[i].Record.Inventory)
			}
			return gm.ActiveCharacters[i].Record.Inventory, nil
		}
	}
	return []edentypes.Item{}, errors.New("character not found")
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
func (gm *GameManager) AssignItemHotkeys(inventory []edentypes.Item) {
	// Assign hotkeys to items in the inventory, starting with lowercase a, and going up to z, then A to Z, then 0 to 9
	// Stackable items should be assigned the same hotkey
	// Non-stackable items should be assigned the next available hotkey
	// If there are no more hotkeys available, the item should not be assigned a hotkey
	hotkey := 'a'                        // Start with lowercase 'a'
	hotkeyMap := make(map[string]string) // Map to store assigned hotkeys for stackable items

	for i := range inventory {
		// Check if the item is stackable
		if inventory[i].Stackable {
			// Check if an item with the same name has been assigned a hotkey
			if stackableHotkey, ok := hotkeyMap[inventory[i].Name]; ok {
				// Assign the same hotkey to the current stackable item
				inventory[i].Hotkey = stackableHotkey
			} else {
				// Assign the next available hotkey to the item
				if hotkey <= 'z' || (hotkey >= 'A' && hotkey <= 'Z') || (hotkey >= '0' && hotkey <= '9') {
					inventory[i].Hotkey = string(hotkey)
					hotkeyMap[inventory[i].Name] = string(hotkey) // Update the hotkey map for stackable items
					hotkey++
				} else {
					// No more hotkeys available, don't assign a hotkey to this item
					inventory[i].Hotkey = ""
				}
			}
		} else {
			// Non-stackable item, assign a new hotkey
			if hotkey <= 'z' || (hotkey >= 'A' && hotkey <= 'Z') || (hotkey >= '0' && hotkey <= '9') {
				inventory[i].Hotkey = string(hotkey)
				hotkey++
			} else {
				// No more hotkeys available, don't assign a hotkey to this item
				inventory[i].Hotkey = ""
			}
		}
	}
}

func (gm *GameManager) AddToCharacterInventory(characterID string, item edentypes.Item) error {
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

// SendLoadingMessage sends a loading message to a specific player
func (gm *GameManager) SendLoadingMessage(characterID, message string) {
	// Find the connection ID for this character
	gm.activeCharactersMutex.Lock()
	var connectionID string
	for _, character := range gm.ActiveCharacters {
		if character.ID == characterID {
			connectionID = character.ConnectionID
			break
		}
	}
	gm.activeCharactersMutex.Unlock()

	if connectionID == "" {
		gm.Log.Println(logging.LogWarn, "Could not find connection ID for character", characterID, "- cannot send loading message")
		return
	}

	// Send the loading message to the specific player
	response := messages.ConnectionManagerMessage{
		Type:               messages.ConnectManager_Message_GameCommandResponse,
		RecipientConsoleID: connectionID,
		Data: messages.GameMessage{
			Type:    messages.GM_LoadingMessage,
			Message: message,
		},
	}

	gm.SendChannel <- response
	gm.Log.Println(logging.LogInfo, "Sent loading message to player", characterID, "(connection:", connectionID, "):", message)
}

// CreateDefaultInventory creates a default starting inventory for new characters
// Only uses JSON-based item definitions - no hardcoded fallbacks
func (gm *GameManager) CreateDefaultInventory() []edentypes.Item {
	var inventory []edentypes.Item

	// Only create items from JSON definitions
	if gm.ItemRegistry == nil {
		gm.Log.Println(logging.LogError, "ItemRegistry is nil - cannot create default inventory")
		return inventory
	}

	// Add 10 wood items
	for j := 0; j < 10; j++ {
		if item, err := gm.ItemRegistry.CreateItemFromDefinition("wood", ""); err == nil {
			item.ID = edenutil.GenerateID() // Generate unique ID for this instance
			if err := gm.DB.AddRecord("Items", item); err != nil {
				gm.Log.Println(logging.LogError, "Failed to add wood item to DB:", err.Error())
			}
			inventory = append(inventory, *item)
		} else {
			gm.Log.Println(logging.LogError, "Failed to create wood item from registry:", err.Error())
		}
	}

	// Add 10 stone items
	for j := 0; j < 10; j++ {
		if item, err := gm.ItemRegistry.CreateItemFromDefinition("stone", ""); err == nil {
			item.ID = edenutil.GenerateID() // Generate unique ID for this instance
			if err := gm.DB.AddRecord("Items", item); err != nil {
				gm.Log.Println(logging.LogError, "Failed to add stone item to DB:", err.Error())
			}
			inventory = append(inventory, *item)
		} else {
			gm.Log.Println(logging.LogError, "Failed to create stone item from registry:", err.Error())
		}
	}

	// Add 1 pickaxe
	if item, err := gm.ItemRegistry.CreateItemFromDefinition("pickaxe", ""); err == nil {
		item.ID = edenutil.GenerateID() // Generate unique ID for this instance
		if err := gm.DB.AddRecord("Items", item); err != nil {
			gm.Log.Println(logging.LogError, "Failed to add pickaxe item to DB:", err.Error())
		}
		inventory = append(inventory, *item)
	} else {
		gm.Log.Println(logging.LogError, "Failed to create pickaxe item from registry:", err.Error())
	}

	gm.Log.Println(logging.LogInfo, "Created default inventory with", len(inventory), "items from JSON definitions")
	return inventory
}
