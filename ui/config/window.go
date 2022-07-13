package config

import (
	"encoding/json"
	"github.com/yamamushi/EscapingEden/ui/types"
)

// WindowID is the ID type of the window
// Only one window of any given ID type can exist in the console
type WindowID int

// These are used as WindowID's for tracking drawing and other behavior.
// They are unique, and only one window of any type can be open at any time.
const (
	WindowDebugBox WindowID = iota
	WindowHelpBox
	WindowChatBox
	WindowLoginMenu
	WindowUserDashboard
	WindowToolBox
	WindowPopupBox
	WindowSettingsBox
)

// String returns the string representation of the WindowID
func (w WindowID) String() string {
	switch w {
	case WindowDebugBox:
		return "WindowDebugBox"
	case WindowHelpBox:
		return "WindowHelpBox"
	case WindowChatBox:
		return "WindowChatBox"
	case WindowLoginMenu:
		return "WindowLoginMenu"
	case WindowToolBox:
		return "WindowToolBox"
	case WindowPopupBox:
		return "WindowPopupBox"
	default:
		return "WindowID"
	}
}

// WindowConfig is the configuration for a window
type WindowConfig struct {
	X       int
	Y       int
	Width   int
	Height  int
	Content string
	Page    types.HelpPage
}

// NewWindowConfig creates a new WindowConfig
func NewWindowConfig(x, y, width, height int, content string) *WindowConfig {
	return &WindowConfig{X: x, Y: y, Width: width, Height: height, Content: content}
}

// String returns the string representation of the WindowConfig
func (c *WindowConfig) String() string {
	output, _ := json.Marshal(c)
	return string(output)
}
