package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/game/commands"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

type GameManager struct {
	ReceiveChannel chan messages.GameManagerMessage       // We only receive player manager messages
	SendChannel    chan messages.ConnectionManagerMessage // We only send connection manager messages

	Log logging.LoggerType
	DB  edendb.DatabaseType

	// Active Characters
	ActiveCharacters      ActiveCharacters
	activeCharactersMutex sync.Mutex

	// ChunkCache is a map of map chunk IDs to map chunks
	ChunkCache map[string]*MapChunk

	// Temporary will remove in a refactor of the map generation
	TileTypes map[string]TileInfo
	Config    *edenconfig.Config

	// Performance optimizations
	ChunkRegistry  *ChunkRegistry
	LazyLoader     *LazyChunkLoader
	WorldValidated bool

	// Map fallback system
	FallbackManager *MapFallbackManager

	// Role management system
	RoleManager *RoleManager

	// Command system
	CommandRegistry *commands.CommandRegistry

	// Monitoring
	StartTime time.Time

	// Item system
	ItemRegistry *edentypes.ItemRegistry
}

func NewGameManager(receiveChannel chan messages.GameManagerMessage, sendChannel chan messages.ConnectionManagerMessage, db edendb.DatabaseType, log logging.LoggerType, conf *edenconfig.Config) *GameManager {
	manager := &GameManager{ReceiveChannel: receiveChannel, SendChannel: sendChannel, DB: db, Log: log, Config: conf}
	manager.Init()
	return manager
}

// Init initializes the database for the game manager if needed
func (gm *GameManager) Init() error {
	gm.Log.Println(logging.LogInfo, "Initializing Game Manager...")

	// Set start time for monitoring
	gm.StartTime = time.Now()

	id, _ := uuid.NewUUID()
	err := gm.DB.AddRecord("Game", &World{ID: id.String()})
	if err != nil {
		if err.Error() == "already exists" {
			gm.Log.Println(logging.LogInfo, "GameManager", "Init", messages.GMError_DBError.Error())
		}
	}
	_ = gm.DB.UpdateRecord("Game", &World{ID: id.String()})

	//	search := messages.Account{Username: "Test"}
	result := World{}
	err = gm.DB.One("Game", "ID", id.String(), &result)
	if err != nil {
		gm.Log.Println(logging.LogError, messages.GMError_DBError.Error(), err)
		return err
	}

	gm.ChunkCache = make(map[string]*MapChunk)

	// Initialize performance optimizations
	gm.ChunkRegistry = NewChunkRegistry("./assets/world")
	gm.LazyLoader = NewLazyChunkLoader(gm, 100) // Cache up to 100 chunks

	// We're going to load the tile types first
	gm.LoadTileTypes()

	// Fast world validation instead of full loading
	gm.FastWorldValidation()

	// Initialize map fallback system
	gm.FallbackManager = NewMapFallbackManager(gm)

	// Initialize role management system
	gm.RoleManager = NewRoleManager(gm)

	// Ensure first user is super admin
	err = gm.RoleManager.EnsureFirstUserIsSuperAdmin()
	if err != nil {
		gm.Log.Println(logging.LogWarn, "Failed to ensure first user is super admin:", err)
	}

	// Initialize command system
	gm.CommandRegistry = commands.NewCommandRegistry(gm.Log)

	// Set up command system components
	gm.initializeCommandSystem()

	// Start performance monitoring
	monitor := NewChunkMonitor(gm)
	monitor.StartPerformanceMonitoring()

	// Schedule periodic character safety checks
	gm.ScheduleCharacterSafetyChecks()

	// Schedule periodic character saves (backup save every 5 minutes)
	gm.SchedulePeriodicCharacterSaves()

	// Initialize item registry and load item definitions
	gm.ItemRegistry = edentypes.NewItemRegistry()
	err = gm.ItemRegistry.LoadItemsFromDirectory("./assets/items")
	if err != nil {
		gm.Log.Println(logging.LogWarn, "Failed to load item definitions:", err)
		gm.Log.Println(logging.LogInfo, "Continuing with empty item registry - items will use fallback system")
	} else {
		itemCount := len(gm.ItemRegistry.Items)
		gm.Log.Println(logging.LogInfo, fmt.Sprintf("Loaded %d item definitions from JSON files", itemCount))

		// Repair existing items in the database with updated color and symbol data
		gm.RepairItemDatabase()
	}

	return nil
}

// initializeCommandSystem sets up the command system components
func (gm *GameManager) initializeCommandSystem() {
	// Create and register command handlers
	digHandler := commands.NewDigCommandHandler(
		gm.DB,
		gm.Log,
		func(playerID string) (interface{}, error) {
			return gm.GetCharacter(playerID)
		},
		gm.GetCharacterInventory,
		gm.HandleDigRequest,
	)
	gm.CommandRegistry.RegisterHandler(digHandler)

	// Set up command logger
	logger := commands.NewDefaultCommandLogger(gm.Log)
	gm.CommandRegistry.SetLogger(logger)

	// Set up metrics collector
	metrics := commands.NewDefaultCommandMetrics()
	gm.CommandRegistry.SetMetrics(metrics)

	// Set up game state provider
	gameState := commands.NewGameStateProvider(
		gm.Log,
		gm.SendChannel,
		func(playerID string) (interface{}, error) {
			return gm.GetCharacter(playerID)
		},
		func(chunkID string) interface{} {
			return gm.GetMapChunkByID(chunkID)
		},
		func(x, y, z int, chunkID string) (interface{}, error) {
			chunk := gm.GetMapChunkByID(chunkID)
			if chunk == nil {
				return nil, fmt.Errorf("chunk not found: %s", chunkID)
			}
			// Add bounds checking and tile access logic here
			return nil, fmt.Errorf("tile access not implemented")
		},
	)
	gm.CommandRegistry.SetGameStateProvider(gameState)

	// Set up permission checker
	permissions := commands.NewDefaultPermissionChecker(gm.Log)
	gm.CommandRegistry.SetPermissionChecker(permissions)

	gm.Log.Println(logging.LogInfo, "Command system initialized")
}

// deltaToDirection converts x,y deltas to a direction string
func (gm *GameManager) deltaToDirection(deltaX, deltaY int) string {
	switch {
	case deltaX == 0 && deltaY == -1:
		return "n" // north
	case deltaX == 0 && deltaY == 1:
		return "s" // south
	case deltaX == 1 && deltaY == 0:
		return "e" // east
	case deltaX == -1 && deltaY == 0:
		return "w" // west
	case deltaX == 1 && deltaY == -1:
		return "ne" // northeast
	case deltaX == -1 && deltaY == -1:
		return "nw" // northwest
	case deltaX == 1 && deltaY == 1:
		return "se" // southeast
	case deltaX == -1 && deltaY == 1:
		return "sw" // southwest
	default:
		return "" // invalid direction
	}
}

// ProcessCommand processes a player command through the command system
func (gm *GameManager) ProcessCommand(cmd *commands.PlayerCommand) (*commands.CommandResult, error) {
	if gm.CommandRegistry == nil {
		return nil, fmt.Errorf("command registry not initialized")
	}

	// Process the command through the registry
	result, err := gm.CommandRegistry.ProcessCommand(cmd)

	// Log command processing
	if err != nil {
		gm.Log.Println(logging.LogError, fmt.Sprintf("Command processing error - Type: %s, Player: %s, Error: %v",
			cmd.Type, cmd.PlayerID, err))
	} else if result != nil {
		gm.Log.Println(logging.LogInfo, fmt.Sprintf("Command processed - Type: %s, Player: %s, Success: %v",
			cmd.Type, cmd.PlayerID, result.Success))
	}

	return result, err
}

// SchedulePeriodicCharacterSaves starts a goroutine that periodically saves all active characters
func (gm *GameManager) SchedulePeriodicCharacterSaves() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute) // Save every 5 minutes
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				gm.SaveAllActiveCharacters()
			}
		}
	}()

	gm.Log.Println(logging.LogInfo, "Periodic character saves scheduled (every 5 minutes)")
}

// SaveAllActiveCharacters saves all currently active characters to the database
func (gm *GameManager) SaveAllActiveCharacters() {
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()

	savedCount := 0
	errorCount := 0

	for _, character := range gm.ActiveCharacters {
		if character.Record != nil {
			err := gm.DB.UpdateRecord("Characters", character.Record)
			if err != nil {
				gm.Log.Println(logging.LogError, "Failed to save character", character.Name, "during periodic save:", err.Error())
				errorCount++
			} else {
				if gm.Config.Logger.DebugChunk {
					gm.Log.Println(logging.LogInfo, "Periodic save: character", character.Name, "at position", character.Record.Position.X, character.Record.Position.Y, "in map", character.Record.CurrentMapID, "chunk", character.Record.Position.MapChunkID)
				}
				savedCount++
			}
		}
	}

	if gm.Config.Logger.DebugChunk {
		gm.Log.Println(logging.LogInfo, fmt.Sprintf("Periodic character save completed: %d saved, %d errors", savedCount, errorCount))
	}
}

// RepairItemDatabase updates existing items in the database with correct color and symbol data from JSON definitions
func (gm *GameManager) RepairItemDatabase() {
	gm.Log.Println(logging.LogInfo, "Starting item database repair...")

	repairedCount := 0
	errorCount := 0

	// Get all characters to repair their inventories
	var characters []messages.CharacterInfo
	err := gm.DB.All("Characters", &characters)
	if err != nil {
		gm.Log.Println(logging.LogError, "Failed to get characters for item repair:", err)
		return
	}

	for _, character := range characters {
		characterUpdated := false

		// Repair items in character inventory
		for i, item := range character.Inventory {
			// Check if this item needs repair (missing symbol or colors)
			needsRepair := item.Symbol == "" || item.FGColor == 0

			if needsRepair {
				// Try to find the item definition by matching name directly
				var def edentypes.ItemDefinition
				var exists bool
				for _, itemDef := range gm.ItemRegistry.Items {
					if itemDef.Name == item.Name {
						def = itemDef
						exists = true
						break
					}
				}
				if exists {
					// Update the item with data from JSON definition
					character.Inventory[i].Symbol = def.Symbol
					character.Inventory[i].FGColor = def.FGColor
					character.Inventory[i].BGColor = def.BGColor
					character.Inventory[i].Category = def.Category
					character.Inventory[i].Tags = def.Tags
					character.Inventory[i].Rarity = def.Rarity
					character.Inventory[i].Value = def.Value
					character.Inventory[i].MaxStack = def.MaxStack

					// Update durability if it's not set
					if character.Inventory[i].Durability == 0 && def.Durability != 0 {
						character.Inventory[i].Durability = def.Durability
					}

					characterUpdated = true
					repairedCount++

					gm.Log.Println(logging.LogInfo, fmt.Sprintf("Repaired item '%s' for character '%s' with symbol '%s'",
						item.Name, character.Name, def.Symbol))
				}
			}
		}

		// Save the character if any items were updated
		if characterUpdated {
			err := gm.DB.UpdateRecord("Characters", &character)
			if err != nil {
				gm.Log.Println(logging.LogError, fmt.Sprintf("Failed to save repaired character '%s': %v", character.Name, err))
				errorCount++
			}
		}
	}

	// Also repair standalone items in the database
	var items []edentypes.Item
	err = gm.DB.All("Items", &items)
	if err != nil {
		gm.Log.Println(logging.LogWarn, "Failed to get standalone items for repair:", err)
	} else {
		for _, item := range items {
			needsRepair := item.Symbol == "" || item.FGColor == 0

			if needsRepair {
				// Try to find the item definition by matching name directly
				var def edentypes.ItemDefinition
				var exists bool
				for _, itemDef := range gm.ItemRegistry.Items {
					if itemDef.Name == item.Name {
						def = itemDef
						exists = true
						break
					}
				}
				if exists {
					// Update the item with data from JSON definition
					item.Symbol = def.Symbol
					item.FGColor = def.FGColor
					item.BGColor = def.BGColor
					item.Category = def.Category
					item.Tags = def.Tags
					item.Rarity = def.Rarity
					item.Value = def.Value
					item.MaxStack = def.MaxStack

					if item.Durability == 0 && def.Durability != 0 {
						item.Durability = def.Durability
					}

					err := gm.DB.UpdateRecord("Items", &item)
					if err != nil {
						gm.Log.Println(logging.LogError, fmt.Sprintf("Failed to save repaired item '%s': %v", item.Name, err))
						errorCount++
					} else {
						repairedCount++
						gm.Log.Println(logging.LogInfo, fmt.Sprintf("Repaired standalone item '%s' with symbol '%s'",
							item.Name, def.Symbol))
					}
				}
			}
		}
	}
	gm.Log.Println(logging.LogInfo, fmt.Sprintf("Item database repair completed: %d items repaired, %d errors", repairedCount, errorCount))

}

// SyncItemWithJSON synchronizes an item's immutable properties with its JSON definition
// while preserving mutable fields like durability, ID, and hotkey
func (gm *GameManager) SyncItemWithJSON(item *edentypes.Item) bool {
	if gm.ItemRegistry == nil {
		return false
	}

	// Try to find the item definition by matching name directly
	var def edentypes.ItemDefinition
	var exists bool
	for _, itemDef := range gm.ItemRegistry.Items {
		if itemDef.Name == item.Name {
			def = itemDef
			exists = true
			break
		}
	}
	if !exists {
		return false
	}

	// Preserve mutable fields
	originalID := item.ID
	originalHotkey := item.Hotkey
	originalDurability := item.Durability

	// Sync immutable properties from JSON
	item.Symbol = def.Symbol
	item.FGColor = def.FGColor
	item.BGColor = def.BGColor
	item.Category = def.Category
	item.Tags = def.Tags
	item.MaxStack = def.MaxStack
	item.Description = def.Description
	item.Weight = def.Weight
	item.Type = edentypes.ItemType(def.Type)
	item.Stackable = def.Stackable
	item.Equippable = def.Equippable

	// Restore mutable fields
	item.ID = originalID
	item.Hotkey = originalHotkey
	// Only restore durability if it was previously set (not 0)
	if originalDurability != 0 {
		item.Durability = originalDurability
	} else if def.Durability != 0 {
		item.Durability = def.Durability
	}

	return true
}

// SyncInventoryWithJSON synchronizes all items in an inventory with their JSON definitions
func (gm *GameManager) SyncInventoryWithJSON(inventory []edentypes.Item) []edentypes.Item {
	syncedCount := 0
	for i := range inventory {
		if gm.SyncItemWithJSON(&inventory[i]) {
			syncedCount++
		}
	}

	if syncedCount > 0 {
		gm.Log.Println(logging.LogInfo, fmt.Sprintf("Synced %d items with JSON definitions", syncedCount))
	}

	return inventory
}
