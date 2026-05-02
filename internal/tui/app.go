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
	return tea.Batch(
		m.sidebar.Init(),
		tea.EnterAltScreen,
	)
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
		case "S":
			m.sidebar.Active = !m.sidebar.Active
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width + 2
		m.height = msg.Height + 2
	}

	m.sidebar, cmd = m.sidebar.Update(msg)
	m.response, cmd = m.response.Update(msg)
	m.request, cmd = m.request.Update(msg)

	return m, cmd
}

func (m model) View() string {
	sidebar := m.sidebar.View()
	right := lipgloss.JoinVertical(
		lipgloss.Top,
		m.request.View(),
		m.response.View(),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, right)
}

func Run() error {
	tea := tea.NewProgram(
		model{
			sidebar: ui.NewSidebar(30, 20),
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err := tea.Run()

	return err
}
