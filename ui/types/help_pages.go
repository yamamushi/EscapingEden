package types

// HelpPage is a help page enum
type HelpPage int

// A list of our help pages

const (
	HelpPageMain HelpPage = iota
	HelpPageAbout
	HelpPageCredits
	HelpPageRules
	HelpPageControls
	HelpPageDeath
	HelpPageIndex
)

// String returns the string representation of a HelpPage
func (hp HelpPage) String() string {
	switch hp {
	case HelpPageMain:
		return "main"
	case HelpPageRules:
		return "rules"
	case HelpPageDeath:
		return "death"
	case HelpPageControls:
		return "controls"
	case HelpPageCredits:
		return "credits"
	case HelpPageAbout:
		return "about"
	default:
		return "unknown"
	}
}
