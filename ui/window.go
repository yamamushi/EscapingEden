package ui

import (
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/window"
	"log"
)

func (c *Console) ForceRedrawOn(windowType config.WindowID) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Println("Forcing redraw on: ", windowType)
	for _, w := range c.Windows {
		if w.GetID() == windowType {
			log.Println("Flushing: ", w.GetID())
			w.FlushLastSent()
		}
	}
	//c.forceScreenRefresh = true
}

func (c *Console) SetNextActiveWindow() {
	// Set the active window to the first one in the list, because we know the last one is
	// Always the active one
	c.SetActiveWindowNoThread(c.Windows[0])
}

func (c *Console) SetPrevActiveWindow() {
	// Set the active window to the second to last one in the last, because the last one is
	// Always the active one

	// If the index is less than 2, then we only have one window in the list
	// In which case we don't want to do anything
	if len(c.Windows) < 2 {
		return
	}
	c.SetActiveWindowNoThread(c.Windows[len(c.Windows)-2])
}

// SetActiveWindowNoThread sets the active window and sets all other windows to inactive without locking
func (c *Console) SetActiveWindowNoThread(window window.WindowType) {
	for i, w := range c.Windows {
		if w.GetID() == window.GetID() {
			//log.Println("Active window set to: ", w.GetID())
			w.SetActive(true)
			c.Windows = append(c.Windows[:i], c.Windows[i+1:]...)
			c.Windows = append(c.Windows, window)
		} else {
			w.SetActive(false)
		}
	}
}

// SetActiveWindow sets the active window and sets all other windows to inactive
func (c *Console) SetActiveWindow(window window.WindowType) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, w := range c.Windows {
		if w.GetID() == window.GetID() {
			//log.Println("Active window set to: ", w.GetID())
			w.SetActive(true)
			// We do this to move the window to the end of the slice
			// Since the last one will always be drawn last, ensuring it will be on top
			// of all other drawn windows
			c.Windows = append(c.Windows[:i], c.Windows[i+1:]...)
			c.Windows = append(c.Windows, window)
			//log.Println("Active Window: ", c.Windows[len(c.Windows)-1].GetID())
		} else {
			w.SetActive(false)
		}
	}
}

// GetActiveWindow returns the current active window
func (c *Console) GetActiveWindow() window.WindowType {
	for _, w := range c.Windows {
		if w.GetActive() {
			return w
		}
	}
	return nil
}

// AddWindow adds a window to the console if it is not already in the console by ID.
func (c *Console) AddWindow(w window.WindowType) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, target := range c.Windows {
		if target.GetID() == w.GetID() {
			log.Println("duplicate window: ", target.GetID())
			return
		}
	}
	c.Windows = append(c.Windows, w)
}

// RemoveWindow removes a window from the console by ID.
func (c *Console) RemoveWindow(id config.WindowID) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, target := range c.Windows {
		if target.GetID() == id {
			c.Windows = append(c.Windows[:i], c.Windows[i+1:]...)
			return
		}
	}
}

func (c *Console) GetAllButActiveWindow() (inactive []window.WindowType) {
	for _, win := range c.Windows {
		if !win.GetActive() {
			inactive = append(inactive, win)
		}
	}
	return inactive
}

// InputToActiveWindow sends an input to the active window.
func (c *Console) InputToActiveWindow(input types.Input) {
	for _, target := range c.Windows {
		if target.GetActive() {
			target.HandleInput(input)
			//log.Println("Input Handled on window: ", target.GetID())
			return
		}
	}
}

// GetWindowByID returns a window by ID.
func (c *Console) GetWindowByID(id config.WindowID) window.WindowType {
	for _, target := range c.Windows {
		if target.GetID() == id {
			return target
		}
	}
	return nil
}
