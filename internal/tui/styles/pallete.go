package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Base
	Bg      = lipgloss.Color("#1a1b2e")
	Surface = lipgloss.Color("#16213e")

	// Border
	BorderNormal = lipgloss.Color("#fff")
	BorderActive = lipgloss.Color("#7c7dd9")

	// Text
	TextNormal = lipgloss.Color("#c8d3f5")
	TextSubtle = lipgloss.Color("#636da6")
	TextMuted  = lipgloss.Color("#444b6a")

	ColorGET      = lipgloss.Color("#9ece6a")
	ColorPOST     = lipgloss.Color("#e0af68")
	ColorPUT      = lipgloss.Color("#7aa2f7")
	ColorDELETE   = lipgloss.Color("#f7768e")
	ColorPATCH    = lipgloss.Color("#bb9af7")
	ColorUnknown  = lipgloss.Color("#565f89")
	ActiveGET     = lipgloss.Color("#283b22")
	ActivePOST    = lipgloss.Color("#3d3323")
	ActivePUT     = lipgloss.Color("#243248")
	ActiveDELETE  = lipgloss.Color("#41242d")
	ActivePATCH   = lipgloss.Color("#352742")
	ActiveUnknown = lipgloss.Color("#23263a")

	// Status
	Success = lipgloss.Color("#9ece6a")
	Warning = lipgloss.Color("#e0af68")
	Danger  = lipgloss.Color("#f7768e")

	// Selected/Active item
	SelectedBg = lipgloss.Color("#2d3149")
	SelectedFg = lipgloss.Color("#7c7dd9")
)
