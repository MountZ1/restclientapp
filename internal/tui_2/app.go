package tui_2

import (
	"restclient/internal/collections"
	services "restclient/internal/services/fs"
	"restclient/internal/tui_2/components"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Focus int

const (
	FocusSidebar Focus = iota
	FocusMain
	FocusResult
)

type model struct {
	Width           int
	Height          int
	Ready           bool
	CollectionsList components.Component
	collections     []list.Item
	Focused         Focus
	ShowDialog      bool
	DialogList      components.Component
	baseView        string
}

type CollectionsLoadedMsg struct {
	Items []list.Item
}

func initialModel() model {
	l := components.NewList([]list.Item{}, 0, 0)

	l.Focus()

	dialogItems := []list.Item{
		components.Item{Title_: "New Request"},
		components.Item{Title_: "New Folder"},
	}
	dialog := components.NewList(dialogItems, 30, 10)

	return model{
		CollectionsList: l,
		Focused:         FocusSidebar,
		DialogList:      dialog,
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		files := services.GetCollectionItems()
		items := []list.Item{}
		for _, f := range files {
			icon := collections.IconFile
			if f.IsDir() {
				icon = collections.IconFolder
			}
			items = append(items, components.Item{
				Title_: icon + " " + f.Name(),
				Desc_:  "",
			})
		}
		return CollectionsLoadedMsg{Items: items}
	}
}

func (m *model) applyFocus() {
	if m.Focused == FocusSidebar {
		m.CollectionsList.Focus()
	} else {
		m.CollectionsList.Blur()
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// quit selalu diproses duluan
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "ctrl+c" || key.String() == "q" {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case CollectionsLoadedMsg:
		m.collections = msg.Items
		if l, ok := m.CollectionsList.(*components.ListModel); ok {
			l.SetItems(msg.Items)
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Ready = true
		sideWidth := countDimentionLength(msg.Width, 20)
		m.CollectionsList = components.NewList(m.collections, sideWidth-3, msg.Height-2)
		m.applyFocus()

	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			m.Focused = FocusSidebar
			m.applyFocus()
		case "m":
			m.Focused = FocusMain
			m.applyFocus()
		case "r":
			m.Focused = FocusResult
			m.applyFocus()
		case "ctrl+a":
			m.ShowDialog = !m.ShowDialog
			if m.ShowDialog {
				m.DialogList.Focus()
			} else {
				m.DialogList.Blur()
				m.applyFocus()
			}
		case "esc":
			if m.ShowDialog {
				m.ShowDialog = false
				m.DialogList.Blur()
				m.applyFocus()
				return m, nil
			}
		}

	case tea.MouseMsg:
		if msg.Action != tea.MouseActionPress {
			break
		}
		sideWidth := countDimentionLength(m.Width, 20)
		mainHeight := countDimentionLength(m.Height, 55)
		switch {
		case msg.X < sideWidth:
			m.Focused = FocusSidebar
		case msg.X >= sideWidth && msg.Y < mainHeight:
			m.Focused = FocusMain
		case msg.X >= sideWidth && msg.Y >= mainHeight:
			m.Focused = FocusResult
		}
		m.applyFocus()
	}

	var cmd tea.Cmd
	if m.ShowDialog {
		m.DialogList, cmd = m.DialogList.Update(msg)
	} else {
		m.CollectionsList, cmd = m.CollectionsList.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if !m.Ready {
		return "initializing...."
	}
	side := sidebar(m.Width, m.Height, m.CollectionsList.View())
	main := mainbar(m.Width, m.Height)
	result := result(m.Width, m.Height)
	rightLayout := lipgloss.JoinVertical(lipgloss.Top, main, result)
	base := lipgloss.JoinHorizontal(lipgloss.Top, side, rightLayout)

	if m.ShowDialog {
		dialogBox := renderDialog(m.DialogList.View())
		return overlayCenter(m.Width, m.Height, dialogBox, base)
	}

	return base
}

func Run() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
