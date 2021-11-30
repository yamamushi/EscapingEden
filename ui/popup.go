package ui

import (
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/window/popupbox"
)

// OpenPopup opens a new popup window using the options
func (c *Console) OpenPopup(options *config.WindowConfig) {
	//log.Println(options)
	popupBox := popupbox.NewPopupBox(options.X, options.Y, options.Width, options.Height, c.Width, c.Height, c.PopupBoxMessages, c.WindowMessages)
	popupBox.Init()
	popupBox.SetContents(options.Content)
	c.LastActiveWindow = c.GetActiveWindow() // Save the last active window
	c.AddWindow(popupBox)                    // Add the popup to the list of windows
	c.SetActiveWindow(popupBox)              // Set the popup as the active window
	//popupBox.FlushLastSent()
	popupBox.FlushLastSent()
}

// ClosePopup closes the popup window
func (c *Console) ClosePopup() {
	// Loop through windows and remove the popup
	for _, w := range c.Windows {
		if w.GetID() == config.WindowPopupBox {
			c.RemoveWindow(w.GetID())
			c.SetActiveWindow(c.LastActiveWindow)
			break
		}
	}
	// When we close the popup, our window is all garbage so, we're going to force a re-draw on everything
	c.ForceRedraw()
}

// HandlePopupMessage handles the messages sent to the popup window
func (c *Console) HandlePopupMessage(message *types.ConsoleMessage) {
	switch message.Message {
	case "close":
		c.ClosePopup()
	}
}

// IsPopupOpen returns true if a popup window is open
func (c *Console) IsPopupOpen() bool {
	for _, w := range c.Windows {
		if w.GetID() == config.WindowPopupBox {
			return true
		}
	}
	return false
}

// GetPopupWindowConfig returns the config of the popup window
func (c *Console) GetPopupWindowConfig() *config.WindowConfig {
	for _, w := range c.Windows {
		if w.GetID() == config.WindowPopupBox {
			return w.GetConfig()
		}
	}
	return nil
}
