package terminals

type TermTypeID int

const (
	TermTypeUnknown TermTypeID = iota
	TermTypeXTerm256Color
)

type TerminalType interface {
	Init()
	GetType() TermTypeID

	// Reset or normal
	// int 0
	// All attributes off
	Reset() string

	// Bold or increased intensity
	// int 1
	// As with faint, the color change is a PC (SCO / CGA) invention.
	Bold() string

	// Faint or decreased intensity or dim
	// int 2
	// May be implemented as a light font weight like bold.
	Faint() string

	// Italic or slanted
	// int 3
	// Not widely supported. Sometimes treated as inverse or blink
	Italic() string

	// Underline
	// int 4
	// Style extensions exist for Kitty, VTE, mintty and iTerm2.
	Underline() string

	// SlowBlink
	// int 5
	// Sets blinking to less than 150 times per minute
	SlowBlink() string

	// RapidBlink
	// int 6
	// MS-DOS ANSI.SYS, 150+ per minute; not widely supported
	RapidBlink() string

	// Reverse video or swap foreground and background (invert)
	// int 7
	// Swap foreground and background colors; inconsistent emulation
	Reverse() string

	// Conceal or hide
	// int 8
	// Not widely supported.
	Conceal() string

	// CrossedOut or strike-through
	// int 9
	// Characters legible, but marked for deletion.
	CrossedOut() string

	// PrimaryFont or default font
	// int 10
	// Select default font
	PrimaryFont() string

	// 11-19: Reserved for font change

	// Fraktur select fraktur font
	// int 20
	// Rarely supported
	Fraktur() string

	// BoldOff Doubly underlined; or: not bold
	// int 21
	// Double-underline per ECMA-48 but instead disables bold
	// intensity on several terminals, including in the Linux kernel's console before version 4.17.[33]
	BoldOff() string

	// NormalIntensity
	// int 22
	// Neither bold nor faint; color changes where intensity is implemented as such.
	NormalIntensity() string

	// NeitherItalicOrBlackLetter off
	// int 23
	// Neither italic, nor blackletter
	NeitherItalicOrBlackLetter() string

	// UnderlineOff Not underlined
	// int 24
	// Neither singly nor doubly underlined
	UnderlineOff() string

	// BlinkOff not blinking
	// int 25
	// Turn blinking off
	BlinkOff() string

	// ProportionalSpacingOff Proportional spacing off
	// int 26
	// ITU T.61 and T.416, not known to be used on terminals
	ProportionalSpacingOff() string

	// ReverseOff Normal video
	// int 27
	// Turn reverse video off
	ReverseOff() string

	// Reveal Not concealed
	// int 28
	// Turn concealed mode off
	Reveal() string

	// NotCrossedOut Not crossed out
	// int 29
	// Turn strike-through mode off
	NotCrossedOut() string

	// 30-37: Set foreground color
	// 38: Set foreground color to default
	// 39: Set foreground color to default
	// 40-47: Set background color
	// 48: Set background color to default
	// 49: Set background color to default

	// DisableProportionalSpacing
	// int 50
	// Turn off proportional spacing (T.61 and T.416)
	DisableProportionalSpacing() string

	// Framed
	// int 51
	// Implemented as "emoji variation selector" in mintty.[34]
	Framed() string

	// Encircled
	// int 52
	// Implemented as "emoji variation selector" in mintty.[34]
	Encircled() string

	// Overlined
	// int 53
	Overlined() string

	// NeitherFrameNorEncircled Neither framed nor encircled
	// int 54
	NeitherFrameNorEncircled() string

	// NotOverlined Not overlined
	// int 55
	NotOverlined() string

	// 58 - 107 : Reserved for future standardization, few terminals support things in this range.
}

type Terminal struct {
	Type TermTypeID
}

func NewTerminal(TypeID TermTypeID) *Terminal {
	return &Terminal{
		Type: TypeID,
	}
}

// Init runs any terminal initialization code needed
// Generally, this is just a no-op
func (t *Terminal) Init() {
	return
}

// GetType returns the terminal type as a TermTypeID
func (t *Terminal) GetType() TermTypeID {
	return t.Type
}

// Reset or normal
// int 0
// All attributes off
func (t *Terminal) Reset() string {
	return "\033[0m"
}

// Bold or increased intensity
// int 1
// As with faint, the color change is a PC (SCO / CGA) invention.
func (t *Terminal) Bold() string {
	return "\033[1m"
}

// Faint or decreased intensity or dim
// int 2
// May be implemented as a light font weight like bold.
func (t *Terminal) Faint() string {
	return "\033[2m"
}

// Italic or slanted
// int 3
// Not widely supported. Sometimes treated as inverse or blink
func (t *Terminal) Italic() string {
	return "\033[3m"
}

// Underline
// int 4
// Style extensions exist for Kitty, VTE, mintty and iTerm2.
func (t *Terminal) Underline() string {
	return "\033[4m"
}

// SlowBlink
// int 5
// Sets blinking to less than 150 times per minute
func (t *Terminal) SlowBlink() string {
	return "\033[5m"
}

// RapidBlink
// int 6
// MS-DOS ANSI.SYS, 150+ per minute; not widely supported
func (t *Terminal) RapidBlink() string {
	return "\033[6m"
}

// Reverse video or swap foreground and background (invert)
// int 7
// Swap foreground and background colors; inconsistent emulation
func (t *Terminal) Reverse() string {
	return "\033[7m"
}

// Conceal or hide
// int 8
// Not widely supported.
func (t *Terminal) Conceal() string {
	return "\033[8m"
}

// CrossedOut or strike-through
// int 9
// Characters legible, but marked for deletion.
func (t *Terminal) CrossedOut() string {
	return "\033[9m"
}

// PrimaryFont or default font
// int 10
// Select default font
func (t *Terminal) PrimaryFont() string {
	return "\033[10m"
}

// Fraktur select fraktur font
// int 20
// Rarely supported
func (t *Terminal) Fraktur() string {
	return "\033[20m"
}

// BoldOff Doubly underlined; or: not bold
// int 21
// Double-underline per ECMA-48 but instead disables bold
// intensity on several terminals, including in the Linux kernel's console before version 4.17.[33]
func (t *Terminal) BoldOff() string {
	return "\033[21m"
}

// NormalIntensity
// int 22
// Neither bold nor faint; color changes where intensity is implemented as such.
func (t *Terminal) NormalIntensity() string {
	return "\033[22m"
}

// NeitherItalicOrBlackLetter off
// int 23
// Neither italic, nor blackletter
func (t *Terminal) NeitherItalicOrBlackLetter() string {
	return "\033[23m"
}

// UnderlineOff Not underlined
// int 24
// Neither singly nor doubly underlined
func (t *Terminal) UnderlineOff() string {
	return "\033[24m"
}

// BlinkOff not blinking
// int 25
// Turn blinking off
func (t *Terminal) BlinkOff() string {
	return "\033[25m"
}

// ProportionalSpacingOff Proportional spacing off
// int 26
// ITU T.61 and T.416, not known to be used on terminals
func (t *Terminal) ProportionalSpacingOff() string {
	return "\033[26m"
}

// ReverseOff Normal video
// int 27
// Turn reverse video off
func (t *Terminal) ReverseOff() string {
	return "\033[27m"
}

// Reveal Not concealed
// int 28
// Turn concealed mode off
func (t *Terminal) Reveal() string {
	return "\033[28m"
}

// NotCrossedOut Not crossed out
// int 29
// Turn strike-through mode off
func (t *Terminal) NotCrossedOut() string {
	return "\033[29m"
}

// DisableProportionalSpacing
// int 50
// Turn off proportional spacing (T.61 and T.416)
func (t *Terminal) DisableProportionalSpacing() string {
	return "\033[50m"
}

// Framed
// int 51
// Implemented as "emoji variation selector" in mintty.[34]
func (t *Terminal) Framed() string {
	return "\033[51m"
}

// Encircled
// int 52
// Implemented as "emoji variation selector" in mintty.[34]
func (t *Terminal) Encircled() string {
	return "\033[52m"
}

// Overlined
// int 53
func (t *Terminal) Overlined() string {
	return "\033[53m"
}

// NeitherFrameNorEncircled Neither framed nor encircled
// int 54
func (t *Terminal) NeitherFrameNorEncircled() string {
	return "\033[54m"
}

// NotOverlined Not overlined
// int 55
func (t *Terminal) NotOverlined() string {
	return "\033[55m"
}
