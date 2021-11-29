package types

type HelpPage int

const (
	HelpPageMain HelpPage = iota
	HelpPageAbout
	HelpPageCredits
	HelpPageRules
	HelpPageControls
	HelpPageDeath
	HelpPageIndex
)

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
