package ui

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// ArtConvert is a utility for converting txt art files to a renderable format.
// It can also store the converted art in a file that contains color mapping information

type ArtConvert struct {
}

func NewArtConvert() *ArtConvert {
	return &ArtConvert{}
}

type AsciiArtFile struct {
	Filename string
	Data     string
}

func (ac *ArtConvert) OpenAt(filename string, window WindowType, X, Y int) (string, error) {
	af, err := ac.Open(filename)
	if err != nil {
		return "", err
	}
	return ac.Move(af, window, X, Y), nil
}

func (ac *ArtConvert) Open(filename string) (*AsciiArtFile, error) {
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
	return &AsciiArtFile{Filename: filename, Data: contents}, nil
}

// Move moves to X and Y relative to the window position
func (ac *ArtConvert) Move(af *AsciiArtFile, window WindowType, X, Y int) string {
	X = X + window.GetX()
	Y = Y + window.GetY()
	output := "\033[" + strconv.Itoa(X) + ";" + strconv.Itoa(Y) + "H"
	artLines := strings.Split(af.Data, "\n")
	for num, line := range artLines {
		// Move cursor to X+num and Y
		output += "\033[" + strconv.Itoa(X+num) + ";" + strconv.Itoa(Y) + "H"
		output += line + "\n"
	}
	// Reset our output color always
	output += ResetStyle()
	return output
}
