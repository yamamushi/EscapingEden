package xterm_256color

import (
	"github.com/yamamushi/EscapingEden/terminals"
)

type XTerm256Color struct {
	terminals.Terminal
}

func New() *XTerm256Color {
	return &XTerm256Color{}
}

func (t *XTerm256Color) Init() {
	t.Type = terminals.TermTypeXTerm256Color
}
