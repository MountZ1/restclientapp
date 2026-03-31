package components

import tea "github.com/charmbracelet/bubbletea"

type Component interface {
	Update(tea.Msg) (Component, tea.Cmd)
	View() string
	Focus()
	Blur()
	IsFocus() bool
}
