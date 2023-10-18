package gamewindow

import (
	"github.com/yamamushi/EscapingEden/edenitems"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (gw *GameWindow) RequestInventoryUpdate() {
	inventoryRequest := messages.WindowMessage{Data: messages.GameManagerMessage{Data: messages.GameMessageData{CharacterID: gw.GetCharacterInfoField("id")}, Type: messages.GameManager_RequestInventory}, Type: messages.WM_GameCommand}
	gw.SendToConsole(inventoryRequest)
}

func (gw *GameWindow) UpdateInventory(inventory []*edenitems.Item) {
	gw.InventoryMutex.Lock()
	defer gw.InventoryMutex.Unlock()
	gw.Inventory = inventory
	gw.Log.Println(logging.LogInfo, "Inventory updated - ", len(gw.Inventory))
}
