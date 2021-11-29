package window

import "github.com/yamamushi/EscapingEden/ui/types"

func (w *Window) IncreaseContentPos() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.DirectionInput = types.InputUp
}

func (w *Window) DecreaseContentPos() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.DirectionInput = types.InputDown
}

func (w *Window) SupportsScrolling() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.ScrollingSupported
}

func (w *Window) CheckScrollBufferNew() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.ScrollBufferHasNew
}

func (w *Window) SetScrollBufferNew(new bool) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.ScrollBufferHasNew = new
}
