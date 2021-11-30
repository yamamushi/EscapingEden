package window

import "github.com/yamamushi/EscapingEden/ui/types"

// UpdateParams is used when handling resize events to update the various window parameters in a safe state
func (w *Window) UpdateParams(x, y, width, height, consoleWidth, consoleHeight int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()

	// This function can probably also be used later for window moving

	// if x or y are less than 1 set them to 1
	if x < 1 {
		x = 1
	}
	if y < 1 {
		y = 1
	}
	w.X = x
	w.Y = y

	// if w or h are less than 1 set them to 1
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	w.Width = width
	w.Height = height
	w.ConsoleWidth = consoleWidth
	w.ConsoleHeight = consoleHeight

	w.pointMap = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
	w.lastSentContents = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
	w.pointMapInitialized = true

	for i := w.Y + 1; i < w.Y+w.Height+1; i++ {
		for j := w.X + 1; j < w.X+w.Width; j++ {
			if j > len(w.pointMap)-1 {
				return
			}
			if i > len(w.pointMap[i])-1 {
				return
			}
			w.pointMap[j][i] = types.Point{X: j, Y: i, Character: " ", EscapeCode: ""}
		}
	}
	w.lastSentContents = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
}
