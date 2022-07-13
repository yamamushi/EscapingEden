package dashboard

func (dw *DashboardWindow) drawCreateCharacterPending() {

	dw.PrintLn(dw.X+2, dw.Y+2, "Character Creation Pending...", dw.Terminal.Bold())

}
