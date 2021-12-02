package ui

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/window"
	"strconv"
)

// Draw returns the console as a byte array.
func (c *Console) Draw() []byte {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var s string
	//s = s + c.ConsoleCommands
	//c.ConsoleCommands = ""
	if !c.IsConsoleValidSize() {
		s = s + "\033[2J"
		s = s + "Invalid console size, Escaping Eden requires a terminal size of" + strconv.Itoa(MINWIDTH) + "x" + strconv.Itoa(MINHEIGHT) + "or greater.\r\n"
		s = s + "Please resize your terminal, or press q to disconnect.\n"
		s = s + "If your terminal is empty after resizing, you can press ctrl-r to force a screen refresh.\n"

		if c.LastSentOutput != s {
			c.LastSentOutput = s
			return []byte(s)
		} else {
			return []byte("")
		}
	}

	if !c.consoleInitialized {
		c.consoleInitialized = true
		return []byte(s)
	}

	if c.forceScreenRefresh {
		//log.Println("force screen refresh")
		c.forceScreenRefresh = false
		s = s + c.ResetTerminal()
		return []byte(s)
	}

	if c.resizeActive {
		//log.Println("Handling resize in buffer")
		for _, w := range c.Windows {
			w.FlushLastSent()
		}
		c.resizeActive = false
	}

	//log.Println("Drawing console")
	for _, target := range c.Windows {

		c.UpdateWindow(target)
	}
	consoleDraw := c.GenerateScreenFromPointMap()
	if consoleDraw != "" {
		s = s + consoleDraw
	}

	// We do a last minute check for aborting sending messages
	// This is useful when a screen has asked for a screen refresh
	// And we don't want to send data that is going to get overwritten immediately
	c.abortSync.Lock()
	defer c.abortSync.Unlock()
	// If the last output was not the same as the current output, we send it to the client and update the last output.
	if c.LastSentOutput != s && s != "" && !c.abortSend {
		//log.Println("Sending new output to client, length:", len(s))
		c.LastSentOutput = s
		return []byte(s)
	} else {
		c.abortSend = false
		return []byte("")
	}
}

// UpdateWindow takes a WindowType as an argument and updates
// This will result in the window's pointmap being generated
func (c *Console) UpdateWindow(window window.WindowType) {
	// Get Window Attrs
	winX, winY, _, _ := c.GetWindowAttrs(window)

	// First we want to clear the window for any new content coming in
	window.ResetWindowDrawings()

	// Now we want to check for any new content updates
	window.UpdateContents()

	// Draw the contents of the window
	window.Draw(winX, winY)

	// Now we want to draw the window border
	window.DrawBorder(winX, winY)
}

// GenerateScreenFromPointMap generates the screen from the point map
func (c *Console) GenerateScreenFromPointMap() string {
	for _, target := range c.Windows {
		c.FlushWindowArea(target.GetID())
		if target.GetID() == config.WindowHelpBox || target.GetID() == config.WindowPopupBox {
			// If this is a hovered window and is active, we're going to update that window now
			c.UpdateWindow(target)
		}
		if !target.GetHidden() {
			c.AddToPointMap(target.GetPointMap(), target.GetConfig(), target.GetID())
		}
	}
	return c.PrintPointMap()
}

// ClearPointMap clears the point map
func (c *Console) ClearPointMap() {
	// First clear the console before we redraw it
	for i := 0; i < c.Height+1; i++ {
		for j := 0; j < c.Width; j++ {
			//if w.GetCharAt(i, j) != " " { // && w.GetEscapeCodeAt(i, j) != "" {
			//log.Println("Blank point found: ", i, j)
			c.PrintChar(j, i, "", "")
			//}
		}
	}
}

// FlushLastSent flushes the last sent output
func (c *Console) FlushLastSent() {
	c.pmapMutex.Lock()
	defer c.pmapMutex.Unlock()
	//log.Println("Flushing last sent console")
	c.LastSentPointMap = types.NewPointMap(c.Width, c.Height)
}

// ResetWindowDrawings resets the pointmap
func (c *Console) ResetWindowDrawings() {
	//log.Println("console reset called")
	c.FlushLastSent()
	c.ClearPointMap()
}

// PrintLn prints a line of text to the pointmap
func (c *Console) PrintLn(X int, Y int, text string, escapeCode string) {
	c.pmapMutex.Lock()
	defer c.pmapMutex.Unlock()
	if X > len(c.PointMap)-1 {
		return
	}
	if Y > len(c.PointMap[X])-1 {
		return
	}

	for i, character := range text {
		// For the point at X+1, Y, set the character to the character at the current index of the text string
		c.PointMap[X+i][Y] = types.Point{X: X + i, Y: Y, Character: string(character), EscapeCode: escapeCode}
	}
}

// PrintChar prints a character to the pointmap
func (c *Console) PrintChar(X int, Y int, text string, escapeCode string) {
	c.pmapMutex.Lock()
	defer c.pmapMutex.Unlock()
	if X > len(c.PointMap)-1 || X < 0 {
		return
	}
	if Y > len(c.PointMap[X])-1 || Y < 0 {
		return
	}
	c.PointMap[X][Y] = types.Point{X: X, Y: Y, Character: text, EscapeCode: escapeCode}
}

// GetHoveredWindowConfig gets the config of the hovered window
func (c *Console) GetHoveredWindowConfig() *config.WindowConfig {
	var hoveredWindowConfig *config.WindowConfig
	if c.IsHelpOpen() {
		hoveredWindowConfig = c.GetHelpWindowConfig()
	} else if c.IsPopupOpen() {
		hoveredWindowConfig = c.GetPopupWindowConfig()
	}
	return hoveredWindowConfig
}

// DoesWindowNeedFlush returns whether or not a given Window ID needs to be flushed from the last sent buffer
func (c *Console) DoesWindowNeedFlush(winID config.WindowID) bool {
	for _, id := range c.flushWindowList {
		if id == winID {
			return true
		}
	}
	return false
}

// RemoveWindowFromFlushList removes a window from the flush list
func (c *Console) RemoveWindowFromFlushList(winID config.WindowID) {
	for i, id := range c.flushWindowList {
		if id == winID {
			c.flushWindowList = append(c.flushWindowList[:i], c.flushWindowList[i+1:]...)
			return
		}
	}
}

// FlushWindowArea removes the window area from our last sent point map buffer
func (c *Console) FlushWindowArea(winID config.WindowID) {
	c.pmapMutex.Lock()
	defer c.pmapMutex.Unlock()
	if c.DoesWindowNeedFlush(winID) {
		c.RemoveWindowFromFlushList(winID)

		win := c.GetWindowByID(winID)
		if win == nil {
			c.Log.Println(logging.LogWarn, "Invalid request to flush window area, window ID does not exist")
			return
		}

		//log.Println("Flushing window area: ", winID.String())
		//log.Println("GetX, GetY: ", win.GetX(), win.GetY())
		//log.Println("GetWidth, GetHeight: ", win.GetWidth(), win.GetHeight())
		for j := win.GetY(); j < win.GetY()+win.GetHeight()+2; j++ {
			for i := win.GetX(); i < win.GetX()+win.GetWidth()+1; i++ {
				//log.Println("Flushing point: ", i, j)
				//if w.GetCharAt(i, j) != " " { // && w.GetEscapeCodeAt(i, j) != "" {
				//log.Println("Blank point found: ", i, j)'
				c.PointMap[i][j] = types.Point{X: i, Y: i, Character: " ", EscapeCode: ""}
				c.LastSentPointMap[i][j] = types.Point{X: i, Y: j, Character: "\033[0m", EscapeCode: ""}
				//}
			}
		}
		//c.LastSentPointMap = types.NewPointMap(c.Width, c.Height)
	}
}

// AddToPointMap adds a window's pointmap to the console's PointMap
func (c *Console) AddToPointMap(input types.PointMap, inputConfig *config.WindowConfig, windowID config.WindowID) {
	c.pmapMutex.Lock()
	defer c.pmapMutex.Unlock()

	// iterate through entire w.pointMap and print out the character at each point
	hoveredWindowConfig := c.GetHoveredWindowConfig()

	for y := inputConfig.Y; y < inputConfig.Y+inputConfig.Height+2; y++ {
		for x := inputConfig.X; x < inputConfig.X+inputConfig.Width+1; x++ {

			// This is how we avoid drawing over open windows, of which we only support one at a time right now
			if hoveredWindowConfig != nil {
				if x >= hoveredWindowConfig.X && x < hoveredWindowConfig.X+hoveredWindowConfig.Width+1 && y >= hoveredWindowConfig.Y && y < hoveredWindowConfig.Y+hoveredWindowConfig.Height+2 {
					if windowID != config.WindowHelpBox && windowID != config.WindowPopupBox {
						continue
					}
				}
			}
			c.PointMap[x][y] = input[x][y]
		}
	}
}

// PrintPointMap prints the points to the console of the content that has changed.
// Only new content should ever be leaving this function. If we see the length of output increase
// Dramatically, this is a good place to start debugging from.
func (c *Console) PrintPointMap() string {
	c.pmapMutex.Lock()
	defer c.pmapMutex.Unlock()

	// iterate through entire w.pointMap and print out the character at each point
	output := ""
	lastSentChar := ""
	lastSentEscape := ""
	lastY := 0
	lastX := 0
	bufferCount := 0

	for y := 0; y < c.Height+1; y++ {
		for x := 0; x < c.Width+1; x++ {
			if c.PointMap[x][y].Character != "" || c.PointMap[x][y].EscapeCode != "" {
				if c.LastSentPointMap[x][y].Print() != c.PointMap[x][y].Print() {
					//log.Println("LastSentPointMap: ", c.LastSentPointMap[x][y].Character)
					//log.Println("Currently Read PointMap: ", target[x][y].Character)

					pointMapChar := c.PointMap[x][y].Character
					pointMapEscape := c.PointMap[x][y].EscapeCode

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
							lastY = c.PointMap[x][y].Y
							lastX = c.PointMap[x][y].X
							bufferCount = 0
						}
						// Now that we have dealt with the buffer count, we can print the new character
						output += c.PointMap[x][y].Print()
					}

					// Finally, no matter what we do with the character, we still append it to
					// The last sent contents, as printing it will still take up column spaces
					c.LastSentPointMap[x][y] = c.PointMap[x][y]
				}
			}
		}
	}

	return output
}
