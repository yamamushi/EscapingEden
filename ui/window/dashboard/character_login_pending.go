package dashboard

func (dw *DashboardWindow) drawCharacterLoginPending() {

	dw.PrintLn(dw.X+2, dw.Y+2, "Login Pending...", dw.Terminal.Bold())

}
