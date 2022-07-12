package ui

import (
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/window/help"
)

// ToggleSettings opens a new settings window using the options
func (c *Console) ToggleSettings(options *config.WindowConfig) {
	if !c.IsHelpOpen() {
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

// CloseSettings closes the help window
func (c *Console) CloseSettings() {
	// Loop through windows and remove the popup
	for _, w := range c.Windows {
		if w.GetID() == config.WindowSettingsBox {
			c.RemoveWindow(w.GetID())
			c.SetActiveWindow(c.LastActiveWindow)
			break
		}
	}
	// When we close the popup, our window is all garbage so, we're going to force a re-draw on everything
	c.ForceRedraw()
}

// HandleSettingsMessage handles messages for the help window
func (c *Console) HandleSettingsMessage(message messages.WindowMessage) {
	switch message.Data.(string) {
	case "close":
		c.CloseSettings()
	}
}

// IsSettingsOpen returns true if the help window is open
func (c *Console) IsSettingsOpen() bool {
	for _, w := range c.Windows {
		if w.GetID() == config.WindowSettingsBox {
			return true
		}
	}
	return false
}

// GetSettingsWindowConfig returns the settings window config
func (c *Console) GetSettingsWindowConfig() *config.WindowConfig {
	for _, w := range c.Windows {
		if w.GetID() == config.WindowSettingsBox {
			return w.GetConfig()
		}
	}
	return nil
}
