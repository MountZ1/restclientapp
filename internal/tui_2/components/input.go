package components

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type InputModel struct {
	input   textinput.Model
	Focused bool
}

func Input(placeholder string) Component {
	ti := textinput.New()
	ti.Placeholder = placeholder

	return &InputModel{
		input: ti,
	}
}

func (i *InputModel) Focus() {
	i.Focused = true
	i.input.Focus()
}

func (i *InputModel) Blur() {
	i.Focused = false
	i.input.Blur()
}

func (i *InputModel) IsFocus() bool {
	return i.Focused
}

func (i *InputModel) Update(msg tea.Msg) (Component, tea.Cmd) {
	if !i.Focused {
		return i, nil
	}
	var cmd tea.Cmd
	i.input, cmd = i.input.Update(msg)
	return i, cmd
}

func (i *InputModel) View() string {
	return i.input.View()
}

func (i *InputModel) Value() string {
	return i.input.Value()
}
