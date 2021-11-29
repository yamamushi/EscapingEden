package config

import (
	"encoding/json"
	"github.com/yamamushi/EscapingEden/ui/types"
)

type WindowID int

// These are used as WindowID's for tracking drawing and other behavior.
// They are unique, and only one window of any type can be open at any time.
const (
	WindowDebugBox WindowID = iota
	WindowHelpBox
	WindowChatBox
	WindowLoginMenu
	WindowToolBox
	WindowPopupBox
)

type WindowConfig struct {
	X       int
	Y       int
	Width   int
	Height  int
	Content string
	Page    types.HelpPage
}

func NewWindowConfig(x, y, width, height int, content string) *WindowConfig {
	return &WindowConfig{X: x, Y: y, Width: width, Height: height, Content: content}
}

func (c *WindowConfig) String() string {
	output, _ := json.Marshal(c)
	return string(output)
}
