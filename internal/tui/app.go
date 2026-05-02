package tui

import (
	"restclient/internal/tui/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	sidebar  ui.SidebarModel
	request  ui.RequestModel
	response ui.ResponseModel
	counter  int
	height   int
	width    int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			m.counter++
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.sidebar, cmd = m.sidebar.Update(msg)
	m.response, cmd = m.response.Update(msg)
	m.request, cmd = m.request.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.sidebar.View(),
		lipgloss.JoinVertical(
			lipgloss.Top,
			m.request.View(),
			m.response.View(),
		),
	)
}

func Run() error {
	tea := tea.NewProgram(
		model{
			sidebar: ui.NewSidebar(30, 20),
		},
		tea.WithAltScreen(),
	)
	_, err := tea.Run()

	return err
}
