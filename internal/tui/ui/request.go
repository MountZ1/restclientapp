package ui

import (
	"restclient/internal/tui/helper"
	"restclient/internal/tui/styles"

	tea "charm.land/bubbletea/v2"
)

type RequestModel struct {
	width  int
	height int
	Active bool
}

func NewRequest(width, height int) RequestModel {
	return RequestModel{
		width:  width,
		height: height,
	}
}

func (m RequestModel) Update(msg tea.Msg) (RequestModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sidebarWidth := msg.Width / 4
		m.width = msg.Width - sidebarWidth
		m.height = msg.Height / 2
	}

	return m, nil
}

func (m RequestModel) View() string {
	borderColor := styles.BorderNormal
	if m.Active {
		borderColor = styles.BorderActive
	}
	return helper.RenderWithTitle(
		"Request Content",
		"Request",
		m.width,
		m.height,
		borderColor,
	)
}
