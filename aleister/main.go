package main

// Aleister is a utility for drawing ASCII art and saving it in a raw text file with escape codes.
import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Aleister")
	if err := ebiten.RunGame(&AleisterApp{}); err != nil {
		log.Fatal(err)
	}
}
