package ui

type InputType int

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

type Input struct {
	Type InputType
	Data string
}
