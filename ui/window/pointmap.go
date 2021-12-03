package window

import (
	"github.com/yamamushi/EscapingEden/ui/types"
)

// PrintLn prints a line to the pointmap.
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

// PrintLnColor prints a line to the pointmap taking into account coloring, which means the point escape codes should not be reset.
// Otherwise, buffering at drawing the point map might cause problems on some terminals.
func (w *Window) PrintLnColor(X int, Y int, text string, escapeCode string) {
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
		w.pointMap[X+i][Y] = types.Point{X: X + i, Y: Y, Character: string(character), EscapeCode: escapeCode, NoReset: true}
		//log.Println("pointMap:", X, Y+i, w.pointMap[X][Y+i].Character)
		// if last character in range
		if i == len(text)-1 {
			if (X+i+1) < w.Width && Y < w.Height {
				w.pointMap[X+i+1][Y] = types.Point{X: X + i + 1, Y: Y, Character: " ", EscapeCode: w.Terminal.Reset()}
			}
		}
	}
}

// PrintChar prints a character to the pointmap.
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

// PrintCharColor prints a character to the pointmap taking into account coloring, which means the point escape code should not be reset.
// Otherwise, buffering at drawing the point map might cause problems on some terminals.
func (w *Window) PrintCharColor(X int, Y int, text string, escapeCode string) {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	if X > len(w.pointMap)-1 || X < 0 {
		return
	}
	if Y > len(w.pointMap[X])-1 || Y < 0 {
		return
	}
	w.pointMap[X][Y] = types.Point{X: X, Y: Y, Character: text, EscapeCode: escapeCode, NoReset: true}
}

// GetCharAt returns the character at the given point.
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

// GetEscapeCodeAt returns the escape code at the given point.
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

// IsPointBlank returns true if the point is blank.
func (w *Window) IsPointBlank(X, Y int) bool {
	if w.GetEscapeCodeAt(X, Y) == "" && w.GetCharAt(X, Y) == "" {
		return true
	} else {
		return false
	}
}

// GetPointMap returns the pointmap.
func (w *Window) GetPointMap() types.PointMap {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	return w.pointMap
}

// ClearMap clears the pointmap.
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

// Deprecated: This function is no longer used internally, as it was replaced by
// Rendering at the console layer
// FlushLastSent flushes the last sent pointmap
func (w *Window) FlushLastSent() {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	w.lastSentContents = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
}

// ResetWindowDrawings resets the window's drawings.
func (w *Window) ResetWindowDrawings() {
	w.ClearMap(w.X, w.Y, w.Width, w.Height, 0, 0)
	//w.FlushLastSent()
	w.SetContents("")
}
