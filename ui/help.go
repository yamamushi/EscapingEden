package ui

import (
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/window/help"
)

// ToggleHelp opens a new help window using the options
func (c *Console) ToggleHelp(options *config.WindowConfig) {
	if !c.IsPopupOpen() {
		helpWindow := help.NewHelpWindow(options.X, options.Y, options.Width, options.Height, c.Width, c.Height,
			options.Page, c.PopupBoxWindowMessages, c.WindowMessages, c.Log, c.Terminal)
		helpWindow.Init()
		helpWindow.SetContents(options.Content)
		c.LastActiveWindow = c.GetActiveWindow() // Save the last active window
		c.AddWindow(helpWindow)                  // Add the popup to the list of windows
		c.SetActiveWindow(helpWindow)            // Set the popup as the active window
		c.ForceRedraw()
	} else {
		c.CloseHelp()
	}
}

// CloseHelp closes the help window
func (c *Console) CloseHelp() {
	// Loop through windows and remove the popup
	for _, w := range c.Windows {
		if w.GetID() == config.WindowHelpBox {
			c.RemoveWindow(w.GetID())
			c.SetActiveWindow(c.LastActiveWindow)
			break
		}
	}
	// When we close the popup, our window is all garbage so, we're going to force a re-draw on everything
	c.ForceRedraw()
}

// HandleHelpMessage handles messages for the help window
func (c *Console) HandleHelpMessage(message messages.WindowMessage) {
	switch message.Data.(string) {
	case "close":
		c.CloseHelp()
	}
}

// IsHelpOpen returns true if the help window is open
func (c *Console) IsHelpOpen() bool {
	for _, w := range c.Windows {
		if w.GetID() == config.WindowHelpBox {
			return true
		}
	}
	return false
}

// GetHelpWindowConfig returns the help window config
func (c *Console) GetHelpWindowConfig() *config.WindowConfig {
	for _, w := range c.Windows {
		if w.GetID() == config.WindowHelpBox {
			return w.GetConfig()
		}
	}
	return nil
}
