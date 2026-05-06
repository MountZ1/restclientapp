package components

import lv2 "charm.land/lipgloss/v2"

type ButtonModel struct {
	height int
	width  int
	text   string
}

func NewButton(text string, width, height int) ButtonModel {
	return ButtonModel{
		text:   text,
		height: height,
		width:  width,
	}
}

func (m ButtonModel) View() string {
	return lv2.NewStyle().
		Foreground(lv2.Color("#FFF")).
		Background(lv2.Color("#6C91BF")).
		Padding(0, 3).
		Render(m.text)
}
