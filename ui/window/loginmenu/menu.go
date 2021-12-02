package login

// drawMenu draws the default login menu
func (lw *LoginWindow) drawMenu() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()
	lw.LockMutex()
	defer lw.UnlockMutex()
	//lw.FlushLastSent()

	// First we are going to setup our default login screen
	lw.PrintLn(lw.X+10, lw.Y+5, "Welcome to Escaping Eden!", "")
	lw.PrintLn(lw.X+10, lw.Y+6, "Please select a menu option from below:", "")

	lw.PrintLn(lw.X+11, lw.Y+8, "(", "")
	lw.PrintLn(lw.X+12, lw.Y+8, "l", lw.Terminal.Bold())
	lw.PrintLn(lw.X+13, lw.Y+8, ")ogin", "")

	lw.PrintLn(lw.X+11, lw.Y+9, "(", "")
	lw.PrintLn(lw.X+12, lw.Y+9, "r", lw.Terminal.Bold())
	lw.PrintLn(lw.X+13, lw.Y+9, ")egister", "")

	lw.PrintLn(lw.X+11, lw.Y+10, "(", "")
	lw.PrintLn(lw.X+12, lw.Y+10, "q", lw.Terminal.Bold())
	lw.PrintLn(lw.X+13, lw.Y+10, ")uit", "")
	/* Disabling this testing logo for now!
	artFile, err := util.OpenASCIIArtFile("assets/ascii/logo.txt")
	if err != nil {
		log.Println(err)
	} else {
		for y, line := range artFile.Data {
			for x, point := range line {
				printX := x + lw.Width - artFile.Width
				printY := y + lw.Height - artFile.Height + 3
				if printX < lw.Width+1 && printY < lw.Height+3 && printY < lw.ConsoleHeight && printX < lw.ConsoleWidth {
					lw.PrintChar(printX, printY, point.Character, point.EscapeCode)
				}
			}
		}
	}
	*/
	return
}
