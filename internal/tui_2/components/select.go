package components

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SelectModel struct {
	Choices  []string
	Cursor   int
	Selected int
	Focused  bool
}

func List(choices []string) Component {
	return &SelectModel{
		Choices:  choices,
		Cursor:   0,
		Selected: -1,
		Focused:  false,
	}
}

func (s *SelectModel) Focus() {
	s.Focused = true
}

func (s *SelectModel) Blur() {
	s.Focused = false
}

func (s *SelectModel) IsFocus() bool {
	return s.Focused
}

func (s *SelectModel) Update(msg tea.Msg) (Component, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !s.Focused {
			return s, nil
		}

		switch msg.String() {
		case "up", "k":
			if s.Cursor > 0 {
				s.Cursor--
			}
		case "down", "j":
			if s.Cursor < len(s.Choices)-1 {
				s.Cursor++
			}
		case "enter":
			s.Selected = s.Cursor
		}

	case tea.MouseMsg:
		if !s.Focused {
			return s, nil
		}

		switch msg.String() {
		case "click":
			s.Selected = s.Cursor
		}
	}

	return s, nil
}

func (s *SelectModel) View() string {
	style := lipgloss.NewStyle()
	focusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)

	var str string
	for i, choice := range s.Choices {
		cursor := " "
		if s.Cursor == i {
			cursor = "▶"
		}

		renderStyle := style
		if s.Cursor == i && s.Focused {
			renderStyle = focusedStyle
		}

		str += cursor + " " + renderStyle.Render(choice) + "\n"
	}

	return str
}

func (s *SelectModel) GetSelected() string {
	if s.Selected >= 0 && s.Selected < len(s.Choices) {
		return s.Choices[s.Selected]
	}
	return ""
}
