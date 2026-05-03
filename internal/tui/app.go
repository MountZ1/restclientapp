package tui

import (
	"fmt"
	"os"
	"restclient/internal/tui/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	sidebar  ui.SidebarModel
	request  ui.RequestModel
	response ui.ResponseModel
	dialog   *ui.DialogModel
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			m.counter++
		case "ctrl+s":
			m.request.Active = false
			m.response.Active = false
			m.sidebar.Active = true
		case "ctrl+r":
			m.sidebar.Active = false
			m.response.Active = false
			m.request.Active = true
		case "ctrl+d":
			m.sidebar.Active = false
			m.request.Active = false
			m.response.Active = true
		}
		if m.dialog != nil {
			switch msg.String() {
			case "esc":
				m.dialog = nil
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case ui.OpenDialogMsg:
		dialog := ui.NewDialog(msg.Title, msg.Message)
		m.dialog = &dialog
	}

	m.sidebar, cmd = m.sidebar.Update(msg)
	cmds = append(cmds, cmd)

	m.request, cmd = m.request.Update(msg)
	cmds = append(cmds, cmd)

	m.response, cmd = m.response.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	sidebar := m.sidebar.View()
	right := lipgloss.JoinVertical(lipgloss.Top,
		m.request.View(),
		m.response.View(),
	)
	main := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, right)

	if m.dialog != nil {
		w := lipgloss.Width(main)
		h := lipgloss.Height(main)
		os.WriteFile("debug.txt", []byte(fmt.Sprintf("main w:%d h:%d, screen w:%d h:%d", w, h, m.width, m.height)), 0644)

		return lipgloss.Place(
			w, h,
			lipgloss.Center, lipgloss.Center,
			m.dialog.View(),
		)
	}

	return main
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
