package dashboard

import "github.com/yamamushi/EscapingEden/messages"

func (dw *DashboardWindow) drawCharacterLoginPending() {

	dw.PrintLn(dw.X+2, dw.Y+2, "Login Pending...", dw.Terminal.Bold())

}

func (dw *DashboardWindow) finalizeLogin(charInfo messages.CharacterInfo) {
	for {
		dw.pendingLoginMutex.Lock()
		if dw.characterManagerValidated && dw.accountManagerValidated {
			dw.stopHandler = true
			dw.pendingLoginMutex.Unlock()
			dw.LoginCharacter(charInfo)
			dw.RequestFlushFromConsole()
			return
		}
		dw.pendingLoginMutex.Unlock()
	}
}
