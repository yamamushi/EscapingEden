package types

type HelpPage int

const (
	HelpPageMain HelpPage = iota
	HelpPageRules
	HelpPageDeath
	HelpPageControls
	HelpPageCredits
	HelpPageAbout
	HelpPageIndex
)
