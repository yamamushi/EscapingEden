package xterm_256color

import (
	"github.com/yamamushi/EscapingEden/terminals"
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
