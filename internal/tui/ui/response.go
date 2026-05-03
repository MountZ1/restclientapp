package ui

import (
	"restclient/internal/tui/helper"
	"restclient/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type ResponseModel struct {
	width  int
	height int
	Active bool
}

func NewResponse(width, height int) ResponseModel {
	return ResponseModel{
		width:  width,
		height: height,
	}
}

func (m ResponseModel) Update(msg tea.Msg) (ResponseModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sidebarWidth := msg.Width / 4
		m.width = msg.Width - sidebarWidth
		m.height = msg.Height/2 - 1
	}

	return m, nil
}

func (m ResponseModel) View() string {
	borderColor := styles.BorderNormal
	if m.Active {
		borderColor = styles.BorderActive
	}

	return helper.RenderWithTitle(
		"Response Content",
		"[ Response ]",
		m.width,
		m.height,
		borderColor,
	)
}
