package helper

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderWithTitle(content, title string, width, height int, borderColor lipgloss.Color) string {
	box := lipgloss.NewStyle().
		Width(width - 2).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Render(content)

	lines := strings.Split(box, "\n")
	if len(lines) > 0 {
		firstLineWidth := lipgloss.Width(lines[0])
		titleWidth := lipgloss.Width(title)
		dashCount := firstLineWidth - titleWidth - 2
		if dashCount < 0 {
			dashCount = 0
		}
		lines[0] = lipgloss.NewStyle().
			Foreground(borderColor).
			Render("╭" + title + strings.Repeat("─", dashCount) + "╮")
	}

	return strings.Join(lines, "\n")
}
