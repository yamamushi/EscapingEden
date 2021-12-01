package terminals

type TermTypeID int

const (
	TermTypeUnknown TermTypeID = iota
	TermTypeXTerm256Color
)

type TerminalType interface {
	GetType() TermTypeID
	Init()
}

type Terminal struct {
	Type TermTypeID
}

func NewTerminal(TypeID TermTypeID) *Terminal {
	return &Terminal{
		Type: TypeID,
	}
}

func (t *Terminal) Init() {
	return
}

func (t *Terminal) GetType() TermTypeID {
	return t.Type
}
