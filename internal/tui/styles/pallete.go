package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Base
	Bg      = lipgloss.Color("#1a1b2e")
	Surface = lipgloss.Color("#16213e")

	// Border
	BorderNormal = lipgloss.Color("#2d3149")
	BorderActive = lipgloss.Color("#7c7dd9")

	// Text
	TextNormal = lipgloss.Color("#c8d3f5")
	TextSubtle = lipgloss.Color("#636da6")
	TextMuted  = lipgloss.Color("#444b6a")

	// HTTP Method colors
	ColorGET    = lipgloss.Color("#9ece6a")
	ColorPOST   = lipgloss.Color("#e0af68")
	ColorPUT    = lipgloss.Color("#7aa2f7")
	ColorDELETE = lipgloss.Color("#f7768e")
	ColorPATCH  = lipgloss.Color("#bb9af7")

	// Status
	Success = lipgloss.Color("#9ece6a")
	Warning = lipgloss.Color("#e0af68")
	Danger  = lipgloss.Color("#f7768e")

	// Selected/Active item
	SelectedBg = lipgloss.Color("#2d3149")
	SelectedFg = lipgloss.Color("#7c7dd9")
)
