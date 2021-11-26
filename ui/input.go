package ui

type InputType int

const (
	InputUnknown InputType = iota
	InputCharacter
	InputDir
	InputEscape
	InputBackspace
	InputTab
)

type Input struct {
	Type InputType
	Data string
}
