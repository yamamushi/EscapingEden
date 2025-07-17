package gamewindow

type GameMenuType int

// A list of menu types
const (
	MenuType_Default GameMenuType = iota
	MenuType_Inventory
	MenuType_Build
	MenuType_Dig
)

func (gw *GameWindow) CreateMenu(menuType GameMenuType) {
	switch menuType {
	case MenuType_Build:
		gw.BuildMenu()
	case MenuType_Dig:
		gw.DigMenu()

	}
}

func (gw *GameWindow) AddMenuBox(mb MenuBoxType) {
	gw.MenusMutex.Lock()
	defer gw.MenusMutex.Unlock()
	gw.Menus = append(gw.Menus, mb)
}

// AddMenuBoxUnsafe adds a menu box without acquiring the mutex (for use when mutex is already held)
func (gw *GameWindow) AddMenuBoxUnsafe(mb MenuBoxType) {
	gw.Menus = append(gw.Menus, mb)
}

// not sure if this actually works or not
func (gw *GameWindow) RemoveMenuBox(mb MenuBoxType) {
	gw.MenusMutex.Lock()
	defer gw.MenusMutex.Unlock()
	for i, menu := range gw.Menus {
		if menu == mb {
			gw.Menus = append(gw.Menus[:i], gw.Menus[i+1:]...)
			return
		}
	}
}

// RemoveMenuBoxUnsafe removes a menu box without acquiring the mutex (for use when mutex is already held)
func (gw *GameWindow) RemoveMenuBoxUnsafe(mb MenuBoxType) {
	for i, menu := range gw.Menus {
		if menu == mb {
			gw.Menus = append(gw.Menus[:i], gw.Menus[i+1:]...)
			return
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
