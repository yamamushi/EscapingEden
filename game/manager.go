package game

import (
	"sync"

	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/edendb"
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
}

func NewGameManager(receiveChannel chan messages.GameManagerMessage, sendChannel chan messages.ConnectionManagerMessage, db edendb.DatabaseType, log logging.LoggerType, conf *edenconfig.Config) *GameManager {
	manager := &GameManager{ReceiveChannel: receiveChannel, SendChannel: sendChannel, DB: db, Log: log, Config: conf}
	manager.Init()
	return manager
}

// Init initializes the database for the game manager if needed
func (gm *GameManager) Init() error {

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

	// Start performance monitoring
	monitor := NewChunkMonitor(gm)
	monitor.StartPerformanceMonitoring()

	// Schedule periodic character safety checks
	gm.ScheduleCharacterSafetyChecks()

	return nil
}
