package components

import lv2 "charm.land/lipgloss/v2"

type ButtonModel struct {
	height  int
	width   int
	Text    string
	padding [2]int
}

func NewButton(text string, width, height int, padding ...[2]int) ButtonModel {
	p := [2]int{0, 1}
	if len(padding) > 0 {
		p = padding[0]
	}
	return ButtonModel{
		Text:    text,
		height:  height,
		width:   width,
		padding: p,
	}
}

func (m ButtonModel) View() string {
	innerWidth := m.width - m.padding[1]*2
	if innerWidth < 0 {
		innerWidth = 0
	}
	return lv2.NewStyle().
		Foreground(lv2.Color("#FFF")).
		Background(lv2.Color("#6C91BF")).
		Padding(m.padding[0], m.padding[1]).
		Width(innerWidth).
		AlignHorizontal(lv2.Center).
		Render(m.Text)
}
