package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RequestModel struct {
	width  int
	height int
}

func NewRequest(width, height int) RequestModel {
	return RequestModel{
		width:  width,
		height: height,
	}
}

func (m RequestModel) Update(msg tea.Msg) (RequestModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sidebarWidth := msg.Width / 4
		m.width = msg.Width - sidebarWidth
		m.height = msg.Height / 2
	}

	return m, nil
}

func (m RequestModel) View() string {
	return lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - 2).
		Border(lipgloss.RoundedBorder()).
		Render("Request Content")
}
