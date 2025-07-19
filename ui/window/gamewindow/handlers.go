package gamewindow

import (
	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

// Listen listens for any messages on cw.ReceiveMessages Chan and handles them
func (gw *GameWindow) Listen() {
	for {
		select {
		case receivedMessage := <-gw.ConsoleReceive:
			message := receivedMessage.Data.(messages.GameMessage).Type
			switch message {
			case messages.GM_CharacterPosition:
				//gw.log.Println(logging.LogInfo, "Game Window received message from console ", receivedMessage.Data.(messages.GameMessage).Data.Data)
				continue
			case messages.GM_CharacterView:
				//gw.log.Println(logging.LogInfo, "Game Window received view from console")
				gw.drawView(receivedMessage.Data.(messages.GameMessage).Data.Data.(messages.GameCharView))
			case messages.GM_Inventory:
				//gw.log.Println(logging.LogInfo, "Game Window received inventory from console")
				inventory := receivedMessage.Data.(messages.GameMessage).Data.Data.([]edentypes.Item)
				gw.UnlockPendingInventory()
				// cast the data to []*edenitems.Item
				gw.UpdateInventory(inventory)
				gw.Log.Println(logging.LogInfo, "Inventory received, total items:", len(inventory))
				if gw.DisplayInventoryPostReceive {
					gw.Log.Println(logging.LogInfo, "Displaying inventory after receive")
					gw.DisplayInventory()
				}
				if gw.DisplayMaterialSelectionAfterReceive {
					gw.Log.Println(logging.LogInfo, "Displaying material selection after receive")
					gw.DisplayMaterialSelection()
				}
			case messages.GM_FailedDig:
				//gw.Log.Println(logging.LogInfo, "Game Window received failed dig message from console")
				if len(gw.Menus) > 0 {
					gw.Menus[0].SetCallbackStatusBarMessage("You can't dig there.")
				}
				gw.SetStatusBarMessage("You can't dig there.")

			case messages.GM_FailedBuildWall:
				//gw.Log.Println(logging.LogInfo, "Game Window received failed build wall message from console")
				if len(gw.Menus) > 0 {
					gw.Menus[0].SetCallbackStatusBarMessage("You can't build there.")
				}
				gw.SetStatusBarMessage("You can't build there.")

			case messages.GM_LoadingMessage:
				// Handle loading messages by displaying them in the status bar
				loadingMessage := receivedMessage.Data.(messages.GameMessage).Message
				if len(gw.Menus) > 0 {
					gw.Menus[0].SetCallbackStatusBarMessage(loadingMessage)
				}
				gw.SetStatusBarMessage(loadingMessage)
			}
		}
	}
}
