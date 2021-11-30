package types

// InputType is an enum storing the types of input
type InputType int

// A list of input types that can be used
const (
	InputUnknown InputType = iota
	InputCharacter
	InputUp
	InputDown
	InputLeft
	InputRight
	InputEscape
	InputBackspace
	InputTab
	InputReturn
	InputNewline
)

// Input is a struct that stores the type of input and the data associated with it
type Input struct {
	Type InputType
	Data string
}
