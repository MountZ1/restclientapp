package tui

import (
	services "restclient/internal/services/fs"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	count  int
	Width  int
	Height int
	Ready  bool
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Ready = true

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			m.count++
		case "down":
			m.count--
		}
	}
	return m, nil
}

func (m model) View() string {
	if !m.Ready {
		return "initializing...."
	}

	collections := services.StartColection()

	side := sidebar(m.Width, m.Height, collections)
	main := mainbar(m.Width, m.Height)
	result := result(m.Width, m.Height)

	rightLayout := lipgloss.JoinVertical(lipgloss.Top, main, result)
	return lipgloss.JoinHorizontal(lipgloss.Top, side, rightLayout)
}

func Run() error {
	p := tea.NewProgram(model{}, tea.WithAltScreen(), tea.WithMouseCellMotion())

	_, err := p.Run()

	return err
}
