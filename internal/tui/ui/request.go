package ui

import (
	"image/color"
	"restclient/internal/tui/helper"
	"restclient/internal/tui/styles"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"
	lv2 "charm.land/lipgloss/v2"
)

const (
	methodBtnWidth = 10
	sendBtnWidth   = 10
	gapWidth       = 1
	urlPlaceholder = "https://api.example.com/endpoint"
)

var boxStyle = lv2.NewStyle().
	Border(lv2.RoundedBorder()).
	Padding(0, 1)

type RequestModel struct {
	width  int
	height int
	Active bool

	method         string
	methodDropdown *huh.Form
	dropdownOpen   bool

	urlInput textinput.Model

	OffsetX int
	OffsetY int
}

type SendRequestMsg struct {
	Method string
	URL    string
}

type rowLayout struct {
	btnStart, btnEnd   int
	urlStart, urlEnd   int
	sendStart, sendEnd int
	rowHeight          int
}

func (m RequestModel) innerContentWidth() int {
	w := m.width - 6
	if w < 1 {
		w = 1
	}
	return w
}

func (m RequestModel) layout() rowLayout {
	btnW := lv2.Width(m.methodButtonView())
	sendW := lv2.Width(m.sendButtonView())

	frame := boxStyle.GetHorizontalFrameSize()
	minInner := lv2.Width(urlPlaceholder)
	urlOuterWidth := m.innerContentWidth() - btnW - sendW - gapWidth*2
	if urlOuterWidth < minInner+frame {
		urlOuterWidth = minInner + frame
	}

	btnStart := 0
	btnEnd := btnStart + btnW
	urlStart := btnEnd + gapWidth
	urlEnd := urlStart + urlOuterWidth
	sendStart := urlEnd + gapWidth
	sendEnd := sendStart + sendW

	return rowLayout{
		btnStart: btnStart, btnEnd: btnEnd,
		urlStart: urlStart, urlEnd: urlEnd,
		sendStart: sendStart, sendEnd: sendEnd,
		rowHeight: 3,
	}
}

func (m RequestModel) urlInnerWidth() int {
	l := m.layout()
	frame := boxStyle.GetHorizontalFrameSize()
	inner := (l.urlEnd - l.urlStart) - frame
	if inner < 1 {
		inner = 1
	}
	return inner
}

func newURLInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = urlPlaceholder
	ti.Prompt = "> "
	ti.SetVirtualCursor(true)
	return ti
}

func NewRequest(width, height int) RequestModel {
	m := RequestModel{
		width:  width,
		height: height,
		method: "GET",
	}
	m.methodDropdown = newMethodDropdown()
	m.urlInput = newURLInput()
	m.urlInput.SetWidth(m.urlInnerWidth())
	return m
}

func (m RequestModel) SetSize(width, height int) RequestModel {
	m.width = width
	m.height = height
	m.urlInput.SetWidth(m.urlInnerWidth())
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

func (m RequestModel) Init() tea.Cmd {
	return nil
}

func (m RequestModel) Update(msg tea.Msg) (RequestModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		mouse := msg.Mouse()

		relX := mouse.X - m.OffsetX
		relY := mouse.Y - m.OffsetY

		l := m.layout()

		inRow := relY >= 0 && relY < l.rowHeight

		switch {
		case inRow && relX >= l.btnStart && relX < l.btnEnd:
			m.urlInput.Blur()
			if m.dropdownOpen {
				m.dropdownOpen = false
			} else {
				m.dropdownOpen = true
				m.methodDropdown = newMethodDropdown()
				cmds = append(cmds, m.methodDropdown.Init())
			}
			return m, tea.Batch(cmds...)

		case m.dropdownOpen:
			m.urlInput.Blur()
			m.dropdownOpen = false
			return m, nil

		case inRow && relX >= l.sendStart && relX < l.sendEnd:
			m.urlInput.Blur()
			cmds = append(cmds, func() tea.Msg {
				return SendRequestMsg{Method: m.method, URL: m.urlInput.Value()}
			})
			return m, tea.Batch(cmds...)

		case inRow && relX >= l.urlStart && relX < l.urlEnd:
			focusCmd := m.urlInput.Focus()
			adjusted := msg
			adjusted.X = relX - l.urlStart
			adjusted.Y = 0

			var cmd tea.Cmd
			m.urlInput, cmd = m.urlInput.Update(adjusted)
			return m, tea.Batch(focusCmd, cmd)

		default:
			m.urlInput.Blur()
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
			if val := m.methodDropdown.GetString("method"); val != "" {
				m.method = val
			}
			m.dropdownOpen = false
		}

		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	m.urlInput, cmd = m.urlInput.Update(msg)
	cmds = append(cmds, cmd)

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
		Width(methodBtnWidth).
		Padding(0, 1).
		Border(lv2.RoundedBorder()).
		Render(m.method + " " + arrow)
}

func (m RequestModel) sendButtonView() string {
	return lv2.NewStyle().
		Foreground(styles.TextNormal).
		Bold(true).
		Align(lv2.Center).
		Width(sendBtnWidth).
		Padding(0, 1).
		Border(lv2.RoundedBorder()).
		BorderForeground(lv2.Color("42")).
		Render("Send")
}

func (m RequestModel) urlBoxView() string {
	borderColor := styles.BorderNormal
	if m.urlInput.Focused() {
		borderColor = styles.BorderActive
	}
	return boxStyle.
		BorderForeground(borderColor).
		Render(m.urlInput.View())
}

func (m RequestModel) View() string {
	panelBorderColor := styles.BorderNormal
	if m.Active {
		panelBorderColor = styles.BorderActive
	}

	btn := m.methodButtonView()
	send := m.sendButtonView()
	urlBox := m.urlBoxView()

	row := lv2.JoinHorizontal(lv2.Center, btn, " ", urlBox, " ", send)

	var content string
	if m.dropdownOpen {
		dropdownView := lv2.NewStyle().
			Width(lv2.Width(btn)).
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
		panelBorderColor,
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
