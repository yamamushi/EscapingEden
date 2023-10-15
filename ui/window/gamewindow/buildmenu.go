package gamewindow

func (gw *GameWindow) BuildMenu(x, y, width, height int, title string, options []struct{ Data interface{} }) {
	// Create a new menu box
	mb := &MenuBox{X: x, Y: y, Width: width, Height: height, Title: title, Options: options}
	// Add the menu box to the game window
	gw.AddMenuBox(mb)
}

func (gw *GameWindow) AddMenuBox(mb *MenuBox) {
	gw.MenusMutex.Lock()
	defer gw.MenusMutex.Unlock()
	gw.Menus = append(gw.Menus, mb)
}

// not sure if this actually works or not
func (gw *GameWindow) RemoveMenuBox(mb *MenuBox) {
	gw.MenusMutex.Lock()
	defer gw.MenusMutex.Unlock()
	for i, menu := range gw.Menus {
		if menu == mb {
			gw.Menus = append(gw.Menus[:i], gw.Menus[i+1:]...)
		}
	}
}

func (gw *GameWindow) DrawMenus() {
	gw.MenusMutex.Lock()
	defer gw.MenusMutex.Unlock()

	for _, menu := range gw.Menus {
		menu.Draw(gw)
	}
}
