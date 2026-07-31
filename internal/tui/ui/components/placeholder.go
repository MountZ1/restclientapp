package components

import (
	"restclient/internal/tui/styles"

	lv2 "charm.land/lipgloss/v2"
)

// Placeholder renders a simple centered message. Used for tabs whose
// editor UI has not been built yet, such as Params, Authorization, and
// Headers.
func Placeholder(message string, width, height int) string {
	return lv2.NewStyle().
		Foreground(styles.TextMuted).
		Width(width).
		Height(height).
		Align(lv2.Center, lv2.Center).
		Render(message)
}
