package gamewindow

type GameMenuType int

// A list of menu types
const (
	MenuType_Default GameMenuType = iota
	MenuType_Inventory
	MenuType_Build
)

func (gw *GameWindow) CreateMenu(menuType GameMenuType) {
	switch menuType {
	case MenuType_Build:
		gw.BuildMenu()

	}
}

func (gw *GameWindow) AddMenuBox(mb MenuBoxType) {
	gw.MenusMutex.Lock()
	defer gw.MenusMutex.Unlock()
	gw.Menus = append(gw.Menus, mb)
}

// not sure if this actually works or not
func (gw *GameWindow) RemoveMenuBox(mb MenuBoxType) {
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
