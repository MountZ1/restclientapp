package ui

import (
	"image/color"
	"restclient/internal/tui/helper"
	"restclient/internal/tui/styles"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"
	lv2 "charm.land/lipgloss/v2"
)

type RequestModel struct {
	width  int
	height int
	Active bool

	method         string
	methodDropdown *huh.Form
	dropdownOpen   bool

	urlInput *huh.Form
	url      string

	OffsetX int
	OffsetY int
}

func NewRequest(width, height int) RequestModel {
	m := RequestModel{
		width:  width,
		height: height,
		method: "GET",
		url:    "",
	}
	m.methodDropdown = newMethodDropdown()
	m.urlInput = newURLInput()
	return m
}

func newMethodDropdown() *huh.Form {
	placeholder := "GET"
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("method").
				Value(&placeholder).
				Options(
					huh.NewOption("GET", "GET"),
					huh.NewOption("POST", "POST"),
					huh.NewOption("PUT", "PUT"),
					huh.NewOption("DELETE", "DELETE"),
					huh.NewOption("PATCH", "PATCH"),
				),
		),
	).WithShowHelp(false)
}

func newURLInput() *huh.Form {
	placeholder := ""
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("url").
				Value(&placeholder).
				Placeholder("https://api.example.com/endpoint"),
		),
	).WithShowHelp(false)
}

func (m RequestModel) Init() tea.Cmd {
	return m.urlInput.Init()
}

func (m RequestModel) Update(msg tea.Msg) (RequestModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sidebarWidth := msg.Width / 4
		m.width = msg.Width - sidebarWidth
		m.height = msg.Height / 2

	case tea.MouseClickMsg:
		mouse := msg.Mouse()

		// Koordinat relatif ke panel Request
		relX := mouse.X - m.OffsetX
		relY := mouse.Y - m.OffsetY

		btnWidth := lv2.Width(m.methodButtonView())
		btnHeight := lv2.Height(m.methodButtonView())

		// Tombol ada di y=1 (setelah border atas), x=1 (setelah border kiri)
		if relY >= 1 && relY <= btnHeight && relX >= 1 && relX <= btnWidth {
			if m.dropdownOpen {
				m.dropdownOpen = false
			} else {
				m.dropdownOpen = true
				m.methodDropdown = newMethodDropdown()
				cmds = append(cmds, m.methodDropdown.Init())
			}
			return m, tea.Batch(cmds...)
		}

		if m.dropdownOpen {
			m.dropdownOpen = false
			return m, nil
		}

	case tea.KeyMsg:
		if msg.String() == "esc" && m.dropdownOpen {
			m.dropdownOpen = false
			return m, nil
		}
	}

	if m.dropdownOpen {
		form, cmd := m.methodDropdown.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.methodDropdown = f
		}
		cmds = append(cmds, cmd)

		if m.methodDropdown.State == huh.StateCompleted {
			// baca hasil pilihan lewat GetString
			if val := m.methodDropdown.GetString("method"); val != "" {
				m.method = val
			}
			m.dropdownOpen = false
		}

		return m, tea.Batch(cmds...)
	}

	form, cmd := m.urlInput.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.urlInput = f
	}
	cmds = append(cmds, cmd)

	if val := m.urlInput.GetString("url"); val != "" {
		m.url = val
	}

	return m, tea.Batch(cmds...)
}

func (m RequestModel) methodButtonView() string {
	mc := methodColor(m.method)
	arrow := "▼"
	if m.dropdownOpen {
		arrow = "▲"
	}
	return lv2.NewStyle().
		Foreground(mc).
		Bold(true).
		Padding(0, 1).
		Border(lv2.RoundedBorder()).
		Render(m.method + " " + arrow)
}

func (m RequestModel) View() string {
	borderColor := styles.BorderNormal
	if m.Active {
		borderColor = styles.BorderActive
	}

	btn := m.methodButtonView()
	btnWidth := lv2.Width(btn)

	urlWidth := m.width - btnWidth - 5
	urlStyle := lv2.NewStyle().
		Width(urlWidth).
		Foreground(styles.TextNormal)

	urlView := urlStyle.Render(m.urlInput.View())
	row := lv2.JoinHorizontal(lv2.Center, btn, " ", urlView)

	var content string
	if m.dropdownOpen {
		dropdownView := lv2.NewStyle().
			Width(btnWidth).
			Render(m.methodDropdown.View())
		content = lv2.JoinVertical(lv2.Left, row, dropdownView)
	} else {
		content = row
	}

	return helper.RenderWithTitle(
		content,
		"Request",
		m.width,
		m.height,
		borderColor,
	)
}

func methodColor(method string) color.Color {
	switch method {
	case "GET":
		return styles.ColorGET
	case "POST":
		return styles.ColorPOST
	case "PUT":
		return styles.ColorPUT
	case "DELETE":
		return styles.ColorDELETE
	case "PATCH":
		return styles.ColorPATCH
	default:
		return styles.ColorUnknown
	}
}
