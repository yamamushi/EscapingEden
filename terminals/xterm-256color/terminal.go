package xterm_256color

import (
	"github.com/yamamushi/EscapingEden/terminals"
	"strconv"
)

type XTerm256Color struct {
	terminals.Terminal
}

func New() *XTerm256Color {
	term := &XTerm256Color{}
	term.Type = terminals.TermTypeXTerm256Color
	return term
}

func (t *XTerm256Color) HideCursor() string {
	return "\033[?25l"
}

func (t *XTerm256Color) MoveCursor(x, y int) string {
	return "\033[" + strconv.Itoa(y) + ";" + strconv.Itoa(x) + "H"
}

func (t *XTerm256Color) RepeatChar(count int) string {
	return "\033[" + strconv.Itoa(count) + "b"
}

func (t *XTerm256Color) ClearTerminal() string {
	// We send ESC[2J to clear the screen and NOT \033c because that will break our cursor hiding
	return "\033[2J"
}