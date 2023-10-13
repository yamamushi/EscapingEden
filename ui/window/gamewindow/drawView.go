package gamewindow

import (
	"github.com/yamamushi/EscapingEden/messages"
)

func (gw *GameWindow) drawView(view messages.GameCharView) {
	//gw.log.Println(logging.LogInfo, "Game Window received view from game manager, drawing")

	receivedPoints := view.View

	// For each point, x, y, print the character to the visible map
	for i := 0; i < len(receivedPoints); i++ {
		for j := 0; j < len(receivedPoints[i]); j++ {
			//gw.PrintStringToMap(i, j, receivedPoints[i][j].Character, receivedPoints[i][j].EscapeCode)
			//gw.log.Println(logging.LogInfo, "Game Window received view from game manager, drawing", receivedPoints[i][j].Character)
			if j < gw.Height-3 {
				gw.DrawToVisibleMap(i, j, receivedPoints[i][j].Character, receivedPoints[i][j].EscapeCode)
			}
		}
	}

}
