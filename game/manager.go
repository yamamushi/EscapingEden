package game

import (
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"sync"
)

type GameManager struct {
	ReceiveChannel chan messages.GameManagerMessage       // We only receive player manager messages
	SendChannel    chan messages.ConnectionManagerMessage // We only send connection manager messages

	Log logging.LoggerType
	DB  edendb.DatabaseType

	// Active Characters
	ActiveCharacters      ActiveCharacters
	activeCharactersMutex sync.Mutex

	// Temporary will remove in a refactor of the map generation
	MapChunks []MapChunk
	TileTypes map[string]TileInfo
	Config    *edenconfig.Config
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

	// We're going to load the tile types first
	gm.LoadTileTypes()

	// Now we load up the world, this could take a while...
	gm.LoadWorld()
	return nil
}
