package util

import (
	"bufio"
	"github.com/yamamushi/EscapingEden/ui/types"
	"os"
	"strings"
)

// These functions are for working with "art" files, ie ascii art files.
// It is assumed that the art files do not have any other ANSI escape codes other than the following:
// [s - save cursor position - discarded
// [u - restore cursor position - discarded
// [*m - Colors

type AsciiArtFile struct {
	Filename string
	Data     types.PointMap
	Height   int
	Width    int
}

func OpenASCIIArtFile(filename string) (*AsciiArtFile, error) {
	// Opens the file at filename
	// Returns the contents of the file as a string
	// Returns an error if the file cannot be opened
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var contents string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		contents += scanner.Text() + "\n"
	}
	lines := strings.Split(contents, "\n")
	height := len(lines)
	var width int
	// First we need to iterate and get the longest line, which serves as our maximum possible width
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	pointMap := types.NewPointMap(height, width)
	// Now we reset width because our actual width is going to be smaller if we start using escape codes
	width = 0

	var newEscapeCode string
	var applyEscapeCode string
	var lineIndex int
	//var lastChar rune
	for y, line := range lines {
		for _, char := range line {
			if char == '\033' && newEscapeCode == "" {
				// If we find an escape code, we store it and continue to the next character
				newEscapeCode = "\033"
				continue
			}

			if newEscapeCode != "" {
				if char == 'u' || char == 's' {
					// ESC[u	restores the cursor to the last saved position
					// ESC[s	save cursor position
					// We're going to discard these and reset the escape code, because we only care about color codes
					// We generally don't care about cursor codes here
					newEscapeCode = ""
					continue
				} else if char == 'm' {
					// If we're in an escape code and we find the 'm'
					// Character, we know we're at the end of a color code
					// And need to reset the last escape code
					// We also need to save the escape code we've saved to this point into the applyEscapeCode sequence
					applyEscapeCode = newEscapeCode + string(char)
					newEscapeCode = ""
					continue
				} else {
					// If we're not at the end of the escape code, we continue adding to it
					newEscapeCode += string(char)
				}
			} else {
				// Only when we apply a character do we increment the line index
				lineIndex++
				// Don't color spaces
				//if char != ' ' {
				// Now that we know we're not IN an escape sequence, we can apply the escape code we have
				point := types.Point{
					X:          lineIndex,
					Y:          y,
					EscapeCode: applyEscapeCode,
					Character:  string(char),
				}
				pointMap[y][lineIndex] = point
				// we don't reset applyEscapeCode because we want to keep the escape code we've applied for the next character
				//}
			}
			//log.Println("Escape code:", applyEscapeCode)
		}
		if width < lineIndex {
			width = lineIndex
		}
		// When we hit a new line, we reset our escape code vars
		newEscapeCode = ""
		applyEscapeCode = ""
		lineIndex = 0

	}
	return &AsciiArtFile{Filename: filename, Data: pointMap, Height: height, Width: width}, nil
}
