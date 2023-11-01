package gamewindow

import (
	"github.com/yamamushi/EscapingEden/edentypes"
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
				if gw.DisplayInventoryPostReceive {
					gw.DisplayInventory()
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
			}
		}
	}
}
