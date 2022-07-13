package ui

import (
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/window/popupbox"
)

// OpenPopup opens a new popup window using the options
func (c *Console) OpenPopup(options *config.WindowConfig) {
	//log.Println(options)
	popupBox := popupbox.NewPopupBox(options.X, options.Y, options.Width, options.Height, c.Width, c.Height,
		c.PopupBoxWindowMessages, c.WindowMessages, c.Log, c.Terminal)
	popupBox.Init()
	popupBox.SetContents(options.Content)
	c.LastActiveWindow = c.GetActiveWindow() // Save the last active window
	c.AddWindow(popupBox)                    // Add the popup to the list of windows
	c.SetActiveWindow(popupBox)              // Set the popup as the active window
	c.ForceRedraw()
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
func (c *Console) HandlePopupMessage(message messages.WindowMessage) {
	switch message.Data.(string) {
	case "close":
		c.ClosePopup()
	}
}

// IsPopupOpen returns true if a popup window is open, including the help and settings popups
func (c *Console) IsPopupOpen() bool {
	for _, w := range c.Windows {
		if w.GetID() == config.WindowPopupBox || w.GetID() == config.WindowHelpBox || w.GetID() == config.WindowSettingsBox {
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
