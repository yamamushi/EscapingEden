package window

import "github.com/yamamushi/EscapingEden/ui/types"

// IncreaseContentPos tells the window that a direction input is pressed
func (w *Window) IncreaseContentPos() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.DirectionInput = types.InputUp
}

// DecreaseContentPos tells the window that a direction input is pressed
func (w *Window) DecreaseContentPos() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.DirectionInput = types.InputDown
}

// SupportsScrolling indicates whether or not the window supports scrolling
func (w *Window) SupportsScrolling() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.ScrollingSupported
}

// CheckScrollBufferNew tells the window that a new scroll buffer is available
func (w *Window) CheckScrollBufferNew() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.ScrollBufferHasNew
}

// SetScrollBufferNew sets the scroll buffer new status to the given value
func (w *Window) SetScrollBufferNew(new bool) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.ScrollBufferHasNew = new
}
