package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResponseModel struct {
	width  int
	height int
}

func NewResponse(width, height int) ResponseModel {
	return ResponseModel{
		width:  width,
		height: height,
	}
}

func (m ResponseModel) Update(msg tea.Msg) (ResponseModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sidebarWidth := msg.Width / 4
		m.width = msg.Width - sidebarWidth
		m.height = msg.Height/2 - 2
	}

	return m, nil
}

func (m ResponseModel) View() string {
	return lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height).
		Border(lipgloss.RoundedBorder()).
		Render("Response Content")
}
