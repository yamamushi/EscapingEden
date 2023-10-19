package window

import (
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/edenutil"
)

// Draw returns a string of the Window's contents
func (w *Window) Draw(X int, Y int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.DrawContents(X, Y)
}

// DrawContents reads a string one character at a time, placing it within the bounds of the window and returns the string
func (w *Window) DrawContents(winX int, winY int) {
	visibleLength := w.Width - 1
	visibleHeight := w.Height - 1 - w.ScrollBufferLimit
	// maxHeight is the maximum height of the window subtracting the border
	maxHeight := visibleHeight
	lines, _ := w.ContentToLines(winX, winY, visibleLength)
	currentLine := winY + 1 + w.StartY // We use this for scrolling windows that want a little buffer at top

	if len(lines) > maxHeight {
		if w.ScrollingSupported {
			if w.DirectionInput == types.InputUp {
				w.ContentStartPos++
				w.DirectionInput = 0
			} else if w.DirectionInput == types.InputDown {
				w.ContentStartPos--
				w.DirectionInput = 0
			}
		}

		if w.ContentStartPos == 0 {
			w.ScrollBufferHasNew = false
		}

		contentStartPos := 0
		if len(lines)-maxHeight+w.ContentStartPos-1 < 0 {
			contentStartPos = 0 - len(lines) + maxHeight + 1
			w.ContentStartPos = 0 - len(lines) + maxHeight + 1
		} else if len(lines)-maxHeight+w.ContentStartPos-1 > len(lines)-maxHeight-1 {
			contentStartPos = 0
			w.ContentStartPos = 0
		} else {
			contentStartPos = w.ContentStartPos
		}

		for i := len(lines) - maxHeight + contentStartPos - 1; i < len(lines); i++ {
			if currentLine > maxHeight+winY+1 {
				break
			}
			// Print current line
			for j, char := range lines[i] {
				if j+winX+1 > w.X+w.Width-1 {
					break
				}
				if w.GetCharAt(j+winX+1, currentLine) == " " || w.GetCharAt(j+winX+1, currentLine) == "" {
					w.PrintChar(j+winX+1, currentLine, string(char), "")
				}
			}

			//w.PrintLn(winX+1, currentLine, lines[i], "")
			// increment currentLine
			currentLine++
		}

		// Draw our arrows last
		if len(lines)-maxHeight+contentStartPos-1 > 0 {
			// draw an up arrow in grey
			w.PrintChar(winX+visibleLength, winY+1+w.StartY, "\u2191", SHWhite)
		}
		if len(lines)-maxHeight+contentStartPos-1 < len(lines)-maxHeight-1 {
			if w.ScrollBufferHasNew {
				// Draw down arrow in red if there is new content
				w.PrintChar(winX+visibleLength, winY+visibleHeight+1, "\u2193", SHRed)
			} else {
				// Draw down arrow in grey if there is no new content
				w.PrintChar(winX+visibleLength, winY+visibleHeight+1, "\u2193", SHWhite)
			}
		}
	} else {
		if w.ScrollingSupported {
			// If the length of content doesn't exceed our visible height, we don't need to scroll
			// And we can discard the DirectionInput
			if w.DirectionInput == types.InputUp {
				w.ContentStartPos = 0
				w.DirectionInput = 0
			} else if w.DirectionInput == types.InputDown {
				w.ContentStartPos = 0
				w.DirectionInput = 0
			}
		}

		for i := 0; i < len(lines); i++ {

			for j, char := range lines[i] {
				if j+winX+1 > w.X+w.Width-1 {
					break
				}
				if w.GetCharAt(j+winX+1, currentLine) == " " || w.GetCharAt(j+winX+1, currentLine) == "" {
					w.PrintChar(j+winX+1, currentLine, string(char), "")
				}
				//w.PrintChar(i+winX+1, currentLine, string(char), "")
			} // increment currentLine
			currentLine++
		}
	}
}

// ContentToLines converts w.Content into lines that fit into the window size.
// It is an extremely important internal function.
func (w *Window) ContentToLines(winX int, winY int, visibleLength int) ([]string, int) {
	// Split the content into lines
	visibleLength = w.Width - 1 - w.ScrollBufferCharMod
	// maxLength is the maximum length of the window subtracting the border
	//maxLength := visibleLength

	currentColumn := winX + 1

	var lines []string
	//parsed := ""

	var currentWord string
	var currentLineOfText string
	//var lastSpacePosition int

	// For every character in the contents
	for i := 0; i < len(w.Contents); i++ {
		if w.Contents[i] == ' ' {
			currentLineOfText += currentWord + " "
			currentWord = ""
		} else if w.Contents[i] == '\n' {
			currentLineOfText += currentWord
			lines = append(lines, currentLineOfText)
			currentLineOfText = ""
			currentWord = ""
			currentColumn = winX + 1
		} else {
			currentWord += string(w.Contents[i])
		}
		if currentColumn > visibleLength+winX-3 {
			currentColumn = winX + 1
			lines = append(lines, currentLineOfText)
			currentLineOfText = ""
		}
		if i == len(w.Contents)-1 {
			currentLineOfText += currentWord
			lines = append(lines, currentLineOfText)
		}
		currentColumn++
	}
	return lines, len(lines)
}
