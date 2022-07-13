package dashboard

import "github.com/yamamushi/EscapingEden/logging"

func (dw *DashboardWindow) loginCharacterByID(charID string) {
	dw.Log.Println(logging.LogInfo, "Logging in character by ID:", charID)
	dw.GetCharacterByID(charID)
	go dw.HandleReceiveChannel()
}
