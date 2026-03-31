package components

import tea "github.com/charmbracelet/bubbletea"

type InputModel struct {
	Placeholder string
	Focused     bool
}

func Input(placeholder string) Component {
	return &InputModel{
		Placeholder: placeholder,
		Focused:     false,
	}
}

func (i *InputModel) Focus() {
	i.Focused = true
}

func (i *InputModel) Blur() {
	i.Focused = false
}

func (i *InputModel) IsFocus() bool {
	return i.Focused
}

func (i *InputModel) Update(msg tea.Msg) (Component, tea.Cmd) {
	if !i.Focused {
		return i, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			// Button clicked
			// TODO: Trigger action
		}
	}

	return i, nil
}

func (i *InputModel) View() string {
	return ""
}
