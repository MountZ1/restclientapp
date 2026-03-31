package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func sidebar(width, height int, data string) string {
	titleStyle := lipgloss.NewStyle().
		Padding(0, 0)
	title := titleStyle.Render("[ C Collections ]")
	widthLenght := countDimentionLength(width, 20)

	borderColor := lipgloss.Color("#526D82")
	remainWidth := widthLenght - lipgloss.Width(title)

	borderStyle := lipgloss.NewStyle().Foreground(borderColor)
	leftBorder := borderStyle.Render("╭" + strings.Repeat("─", 1))
	rightBorder := borderStyle.Render(strings.Repeat("─", remainWidth-1) + "╮")

	topBorder := lipgloss.JoinHorizontal(lipgloss.Left, leftBorder, title, rightBorder)

	contentStyle := lipgloss.NewStyle().
		Width(widthLenght).
		Height(height - 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		BorderTop(false).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		PaddingLeft(1)

	content := contentStyle.Render(data)
	return lipgloss.JoinVertical(lipgloss.Left, topBorder, content)
}

func mainbar(width, height int) string {
	title := lipgloss.NewStyle().Render("[ M Main Bar ]")

	topRemainWidth := countDimentionLength(width, 80) - lipgloss.Width(title) - 3
	leftBorder := lipgloss.NewStyle().Foreground(lipgloss.Color("#526D82")).Render("╭" + strings.Repeat("─", 1))
	rightBorder := lipgloss.NewStyle().Foreground(lipgloss.Color("#526D82")).Render(strings.Repeat("─", topRemainWidth-1) + "╮")
	topBorder := lipgloss.JoinHorizontal(lipgloss.Top, leftBorder, title, rightBorder)

	contentStyle := lipgloss.NewStyle().
		Width(countDimentionLength(width, 80) - 3).
		Height(countDimentionLength(height, 55) - 3).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#526D82")).
		BorderTop(false).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		PaddingLeft(1)

	content := contentStyle.Render("Main")

	return lipgloss.JoinVertical(lipgloss.Top, topBorder, content)
}

func result(width, height int) string {
	title := lipgloss.NewStyle().Render("[ R Result ]")

	topRemainWidth := countDimentionLength(width, 80) - lipgloss.Width(title) - 3
	leftBorder := lipgloss.NewStyle().Foreground(lipgloss.Color("#526D82")).Render("╭" + strings.Repeat("─", 1))
	rightBorder := lipgloss.NewStyle().Foreground(lipgloss.Color("#526D82")).Render(strings.Repeat("─", topRemainWidth-1) + "╮")
	topBorder := lipgloss.JoinHorizontal(lipgloss.Top, leftBorder, title, rightBorder)

	contentStyle := lipgloss.NewStyle().
		Width(countDimentionLength(width, 80) - 3).
		Height(countDimentionLength(height, 50) - 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#526D82")).
		BorderTop(false).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		PaddingLeft(1)

	content := contentStyle.Render("Main")

	return lipgloss.JoinVertical(lipgloss.Top, topBorder, content)
}
