package window

import (
	"github.com/yamamushi/EscapingEden/ui/types"
	"strconv"
)

func (w *Window) PrintLn(X int, Y int, text string, escapeCode string) {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	if X > len(w.pointMap)-1 {
		return
	}
	if Y > len(w.pointMap[X])-1 {
		return
	}

	for i, character := range text {
		//log.Println("inserting character:", string(character))
		// For the point at X, Y+1, set the character to the character at the current index of the text string
		w.pointMap[X+i][Y] = types.Point{X: X + i, Y: Y, Character: string(character), EscapeCode: escapeCode}
		//log.Println("pointMap:", X, Y+i, w.pointMap[X][Y+i].Character)
	}
}

func (w *Window) PrintChar(X int, Y int, text string, escapeCode string) {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	if X > len(w.pointMap)-1 || X < 0 {
		return
	}
	if Y > len(w.pointMap[X])-1 || Y < 0 {
		return
	}
	w.pointMap[X][Y] = types.Point{X: X, Y: Y, Character: text, EscapeCode: escapeCode}
}

func (w *Window) GetCharAt(X, Y int) string {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	if X > len(w.pointMap)-1 {
		return ""
	}
	if Y > len(w.pointMap[X])-1 {
		return ""
	}
	return w.pointMap[X][Y].Character
}

func (w *Window) GetEscapeCodeAt(X, Y int) string {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	if X > len(w.pointMap)-1 {
		return ""
	}
	if Y > len(w.pointMap[X])-1 {
		return ""
	}
	return w.pointMap[X][Y].EscapeCode
}

func (w *Window) IsPointBlank(X, Y int) bool {
	if w.GetEscapeCodeAt(X, Y) == "" && w.GetCharAt(X, Y) == "" {
		return true
	} else {
		return false
	}
}

func (w *Window) GetPointMap() types.PointMap {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	return w.pointMap
}

func (w *Window) PointMapToString() string {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()

	// iterate through entire w.pointMap and print out the character at each point
	output := ""
	lastSentChar := ""
	lastSentEscape := ""
	lastY := 0
	lastX := 0
	bufferCount := 0

	for y := 0; y < len(w.pointMap[0]); y++ {
		for x := 0; x < len(w.pointMap); x++ {
			if w.pointMap[x][y].Character != "" || w.pointMap[x][y].EscapeCode != "" {

				if w.lastSentContents[x][y].Print() != w.pointMap[x][y].Print() {
					pointMapChar := w.pointMap[x][y].Character
					pointMapEscape := w.pointMap[x][y].EscapeCode
					// If this character is the last one sent, then we increase the buffer count
					// and repeat
					if pointMapChar == lastSentChar && pointMapEscape == lastSentEscape &&
						y == lastY && (x)-lastX == 1 {
						//log.Println(pointMapChar)
						bufferCount++
						lastX = x
					} else {
						// If we reached a new character, and the buffer count is greater than 0
						// We need to print the repeated last character bufferCount times
						if bufferCount > 0 {
							repeatCode := lastSentEscape + "\033[" + strconv.Itoa(bufferCount) + "b" + "\033[0m"
							output += repeatCode
							// Finally Reset the buffer count
							bufferCount = 0
						} else {
							// If the buffer count was already 0, we update the last sent character
							// And reset the buffer count for verbosity
							lastSentChar = pointMapChar
							lastSentEscape = pointMapEscape
							lastY = w.pointMap[x][y].Y
							lastX = w.pointMap[x][y].X
							bufferCount = 0
						}
						// Now that we have dealt with the buffer count, we can print the new character
						output += w.pointMap[x][y].Print()

					}

					// Finally, no matter what we do with the character, we still append it to
					// The last sent contents, as printing it will still take up column spaces
					w.lastSentContents[x][y] = w.pointMap[x][y]
				}
			}
		}
	}
	return output
}

func (w *Window) ClearMap(winX int, winY int, visibleLength, visibleHeight int, startX, startY int) {
	// First clear the window before we redraw it
	for i := w.Y + 1; i < w.Y+w.Height+1; i++ {
		for j := w.X + 1; j < w.X+w.Width; j++ {
			//if w.GetCharAt(i, j) != " " { // && w.GetEscapeCodeAt(i, j) != "" {
			//log.Println("Blank point found: ", i, j)
			w.PrintChar(j, i, " ", "")
			//}
		}
	}
}

func (w *Window) FlushLastSent() {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	w.lastSentContents = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
}

func (w *Window) ResetWindowDrawings() {
	w.ClearMap(w.X, w.Y, w.Width, w.Height, 0, 0)
	w.FlushLastSent()
	w.SetContents("")
}
