package dashboard

import logging "github.com/EscapingEden/Logging-Go"

func (dw *DashboardWindow) loginCharacterByID(charID string) {
	dw.Log.Println(logging.LogInfo, "Logging in character by ID:", charID)
	dw.GetCharacterByID(charID)
	go dw.HandleReceiveChannel()
}
