package tui

import (
	services "restclient/internal/services/fs"
	"restclient/internal/services/logger"
	"restclient/internal/tui/custommodel"
	"restclient/internal/tui/ui"
	"restclient/internal/tui/ui/components"
	"restclient/internal/tui/ui/components/dialog"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	sidebar  ui.SidebarModel
	request  ui.RequestModel
	response ui.ResponseModel
	dialog   *dialog.DialogModel
	counter  int
	height   int
	width    int
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.sidebar.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.dialog != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.dialog = nil
				return m, nil
			}
		case dialog.DialogSubmitMsg:
			m.dialog = nil
			err := services.CreateRequestCollection(services.Create{
				Location: msg.Location,
				Type:     msg.Type,
				Name:     msg.Name,
			})
			if err != nil {
				logger.Error("Failed to create %s: %v", msg.Type, err)
			}
			return m, components.LoadCollection()
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
		}

		*m.dialog, cmd = m.dialog.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			m.counter++
		case "ctrl+s":
			return m, func() tea.Msg { return custommodel.SetActivityMSG{Target: "sidebar"} }
		case "ctrl+r":
			return m, func() tea.Msg { return custommodel.SetActivityMSG{Target: "request"} }
		case "ctrl+d":
			return m, func() tea.Msg { return custommodel.SetActivityMSG{Target: "response"} }
		}
	case tea.MouseClickMsg:
		mouse := msg.Mouse()
		if mouse.X < m.sidebar.Width {
			m.sidebar.Active = true
			m.request.Active = false
			m.response.Active = false
		} else if mouse.Y < m.height/2 {
			m.sidebar.Active = false
			m.request.Active = true
			m.response.Active = false
		} else {
			m.sidebar.Active = false
			m.request.Active = false
			m.response.Active = true
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case ui.OpenDialogMsg:
		dialog := dialog.NewDialog(
			msg.Title,
			msg.Message,
			msg.Location,
			msg.Content,
		)
		m.dialog = &dialog
		return m, m.dialog.Init()
	case custommodel.SetActivityMSG:
		m.sidebar.Active = msg.Target == "sidebar"
		m.request.Active = msg.Target == "request"
		m.response.Active = msg.Target == "response"

	case dialog.DialogSubmitMsg:
		m.dialog = nil
		err := services.CreateRequestCollection(services.Create{
			Location: msg.Location,
			Type:     msg.Type,
			Name:     msg.Name,
		})
		if err != nil {
			logger.Error("Failed to create %s: %v", msg.Type, err, msg.Location)
		}

		return m, components.LoadCollection()
	}

	m.sidebar, cmd = m.sidebar.Update(msg)
	cmds = append(cmds, cmd)
	m.request, cmd = m.request.Update(msg)
	cmds = append(cmds, cmd)
	m.response, cmd = m.response.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	sidebar := m.sidebar.View()
	right := lipgloss.JoinVertical(lipgloss.Top,
		m.request.View(),
		m.response.View(),
	)
	main := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, right)

	var rendered string

	if m.dialog != nil {
		w := lipgloss.Width(main)
		h := lipgloss.Height(main)

		dialogStr := m.dialog.View()
		dw := lipgloss.Width(dialogStr)
		dh := lipgloss.Height(dialogStr)

		x := (w - dw) / 2
		y := (h - dh) / 2

		comp := lipgloss.NewCompositor(
			lipgloss.NewLayer(main),
			lipgloss.NewLayer(dialogStr).X(x).Y(y).Z(1),
		)
		rendered = comp.Render()
	} else {
		rendered = main
	}

	v := tea.NewView(rendered)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func Run() error {
	p := tea.NewProgram(
		model{
			sidebar: ui.NewSidebar(30, 20),
		},
	)
	_, err := p.Run()
	return err
}
