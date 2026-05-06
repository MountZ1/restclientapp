package components

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ButtonModel struct {
	Text    string
	Focused bool
}

func Button(text string) Component {
	return &ButtonModel{
		Text:    text,
		Focused: false,
	}
}

func (s *ButtonModel) Focus() {
	s.Focused = true
}

func (s *ButtonModel) Blur() {
	s.Focused = false
}

func (s *ButtonModel) IsFocus() bool {
	return s.Focused
}

func (b *ButtonModel) Update(msg tea.Msg) (Component, tea.Cmd) {
	if !b.Focused {
		return b, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			// Button clicked
			// TODO: Trigger action
		}
	}

	return b, nil
}

func (b *ButtonModel) View() string {
	Button := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Background(lipgloss.Color("#333")).
		Padding(0, 3)

	return Button.Render(b.Text)
}
