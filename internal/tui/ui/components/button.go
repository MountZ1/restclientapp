package components

import lv2 "charm.land/lipgloss/v2"

type ButtonModel struct {
	height  int
	width   int
	Text    string
	padding [2]int
}

func NewButton(text string, width, height int, padding ...[2]int) ButtonModel {
	p := [2]int{0, 2}
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
	return lv2.NewStyle().
		Foreground(lv2.Color("#FFF")).
		Background(lv2.Color("#6C91BF")).
		Padding(m.padding[0], m.padding[1]).
		Width(m.width).
		AlignHorizontal(lv2.Center).
		Render(m.Text)
}
