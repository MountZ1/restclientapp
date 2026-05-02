package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SidebarModel struct {
	lastKey string
	width   int
	height  int
}

func NewSidebar(width, height int) SidebarModel {
	return SidebarModel{
		width:  width,
		height: height,
	}
}

func (m SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.lastKey = msg.String()
	case tea.WindowSizeMsg:
		m.width = msg.Width / 4
		m.height = msg.Height
	}
	return m, nil
}

func (m SidebarModel) View() string {
	return lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - 2).
		Border(lipgloss.RoundedBorder()).
		Render("Last key: " + m.lastKey)
}
