package ui

import "github.com/charmbracelet/lipgloss"

type DialogModel struct {
	title   string
	message string
	input   string
}

func NewDialog(title, message string) DialogModel {
	return DialogModel{
		title:   title,
		message: message,
	}
}

func (d DialogModel) View() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		d.message,
		"",
		d.input+"_",
	)

	return lipgloss.NewStyle().
		Width(40).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Render(content)
}
