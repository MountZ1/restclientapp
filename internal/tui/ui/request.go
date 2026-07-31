package ui

import (
	"image/color"
	"restclient/internal/tui/helper"
	"restclient/internal/tui/styles"
	"restclient/internal/tui/ui/components"

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

	// Layout constants for the vertical sections below the panel
	// title border. These need to stay in sync with how each section
	// is actually rendered if its own chrome changes.
	requestRowHeight = 3 // border + content + border of the method/url/send row
	urlTabsGap       = 0 // rows between the URL row and the tab bar
	tabsRowHeight    = 0
	tabsContentGap   = 0 // rows between the tab bar and its content

	// rightMargin keeps every section (the method/url/send row AND the
	// tab content below it) from touching the panel's own right border.
	// Shared in one place so row layout and tab content width can never
	// drift apart again.
	rightMargin = 4
)

// requestTab identifies which tab of the request editor is active.
type requestTab int

const (
	tabParams requestTab = iota
	tabAuth
	tabHeaders
	tabBody
)

var requestTabTitles = []string{"Params", "Authorization", "Headers", "Body"}

var boxStyle = lv2.NewStyle().
	Border(lv2.RoundedBorder()).
	Padding(0, 1)

// dropdownBoxStyle is the border drawn around the floating method
// dropdown list.
var dropdownBoxStyle = lv2.NewStyle().
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

	tabs    components.TabsModel
	body    components.BodyModel
	params  components.KeyValueEditor
	headers components.KeyValueEditor
	auth    components.AuthEditor

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

// innerContentWidth returns the width available for the panel's content,
// after subtracting the chrome added by helper.RenderWithTitle and the
// other nested borders. This value was tuned empirically against the
// actual rendered output, not derived analytically, so if the outer
// chrome changes this may need to be re-tuned.
func (m RequestModel) innerContentWidth() int {
	w := m.width - helper.PanelHorizontalFrame()
	if w < 1 {
		w = 1
	}
	return w
}

func (m RequestModel) innerContentHeight() int {
	h := m.height - helper.PanelVerticalFrame()
	if h < 1 {
		h = 1
	}
	return h
}

func (m RequestModel) tabContentWidth() int {
	w := m.innerContentWidth() - rightMargin
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

	availableWidth := m.innerContentWidth() - rightMargin

	urlOuterWidth := availableWidth - btnW - sendW - gapWidth*2
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
		rowHeight: requestRowHeight,
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

// verticalLayout returns the Y positions, relative to the top of the
// panel's content, of the tab bar and the tab content area, plus how
// tall the content area is allowed to be. This layout never changes
// based on whether the method dropdown is open; the dropdown is
// rendered as a floating overlay instead, see View.
func (m RequestModel) verticalLayout() (tabsY, contentY, contentHeight int) {
	tabsY = requestRowHeight + urlTabsGap
	contentY = tabsY + m.tabs.Height() + tabsContentGap
	contentHeight = m.innerContentHeight() - contentY
	if contentHeight < 3 {
		contentHeight = 3
	}
	return
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
		width:   width,
		height:  height,
		method:  "GET",
		tabs:    components.NewTabs(requestTabTitles),
		body:    components.NewBody(),
		params:  components.NewKeyValueEditor(),
		headers: components.NewKeyValueEditor(),
		auth:    components.NewAuthEditor(),
	}
	m.methodDropdown = newMethodDropdown()
	m.urlInput = newURLInput()
	m.urlInput.SetWidth(m.urlInnerWidth())

	_, _, contentHeight := m.verticalLayout()
	contentWidth := m.tabContentWidth()
	m.body = m.body.SetSize(contentWidth, contentHeight)
	m.params = m.params.SetSize(contentWidth, contentHeight)
	m.headers = m.headers.SetSize(contentWidth, contentHeight)
	m.auth = m.auth.SetWidth(contentWidth)

	return m
}

func (m RequestModel) SetSize(width, height int) RequestModel {
	m.width = width
	m.height = height

	m.urlInput.SetWidth(m.urlInnerWidth())

	_, _, contentHeight := m.verticalLayout()
	contentWidth := m.tabContentWidth() // <- berubah
	m.body = m.body.SetSize(contentWidth, contentHeight)
	m.params = m.params.SetSize(contentWidth, contentHeight)
	m.headers = m.headers.SetSize(contentWidth, contentHeight)
	m.auth = m.auth.SetWidth(contentWidth)
	return m
}

func newMethodDropdown() *huh.Form {
	placeholder := "GET"

	width := methodBtnWidth - dropdownBoxStyle.GetHorizontalFrameSize()
	if width < 1 {
		width = 1
	}

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
				).
				WithWidth(width),
		),
	).WithShowHelp(false)
}

func (m RequestModel) Init() tea.Cmd {
	return nil
}

// blurAllFields removes focus from every focusable field across every
// tab. Used whenever a click lands somewhere that should not keep any
// field focused.
func (m RequestModel) blurAllFields() RequestModel {
	m.urlInput.Blur()
	m.body = m.body.Blur()
	return m
}

func (m RequestModel) Update(msg tea.Msg) (RequestModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		mouse := msg.Mouse()
		relX := mouse.X - m.OffsetX - 1
		relY := mouse.Y - m.OffsetY - 1
		l := m.layout()
		tabsY, contentY, contentHeight := m.verticalLayout()

		inRow := relY >= 0 && relY < l.rowHeight

		switch {
		// Method dropdown button.
		case inRow && relX >= l.btnStart && relX < l.btnEnd:
			m = m.blurAllFields()
			if m.dropdownOpen {
				m.dropdownOpen = false
			} else {
				m.dropdownOpen = true
				m.methodDropdown = newMethodDropdown()
				cmds = append(cmds, m.methodDropdown.Init())
			}
			return m, tea.Batch(cmds...)

		// Any other click while the dropdown is open just closes it.
		case m.dropdownOpen:
			m = m.blurAllFields()
			m.dropdownOpen = false
			return m, nil

		// Send button.
		case inRow && relX >= l.sendStart && relX < l.sendEnd:
			m = m.blurAllFields()
			cmds = append(cmds, func() tea.Msg {
				return SendRequestMsg{Method: m.method, URL: m.urlInput.Value()}
			})
			return m, tea.Batch(cmds...)

		// URL input.
		case inRow && relX >= l.urlStart && relX < l.urlEnd:
			m = m.blurAllFields()
			focusCmd := m.urlInput.Focus()
			adjusted := msg
			adjusted.X = relX - l.urlStart
			adjusted.Y = 0

			var cmd tea.Cmd
			m.urlInput, cmd = m.urlInput.Update(adjusted)
			return m, tea.Batch(focusCmd, cmd)

			// Tab bar.
		case relY == tabsY+m.tabs.FrameOffsetY() && relX >= 0 && relX < lv2.Width(m.tabs.View()):
			m = m.blurAllFields()
			tabRelX := relX - m.tabs.FrameOffsetX()
			if idx, ok := m.tabs.HitTest(tabRelX); ok {
				m.tabs.Active = idx
			}
			return m, nil

		case relY >= contentY && relY < contentY+contentHeight:
			m = m.blurAllFields()
			cRelX := relX
			cRelY := relY - contentY

			switch requestTab(m.tabs.Active) {
			case tabBody:
				if cRelY == 0 {
					m.body = m.body.HandleTypeClick(cRelX, cRelY)
					return m, nil
				}
				var focusCmd, updateCmd tea.Cmd
				m.body, focusCmd = m.body.Focus()
				adjusted := msg
				adjusted.X = cRelX - bodyContentInset
				adjusted.Y = cRelY - m.body.TypeRowHeight() - 1 - bodyContentInset
				m.body, updateCmd = m.body.Update(adjusted)
				return m, tea.Batch(focusCmd, updateCmd)

			case tabParams:
				var cmd tea.Cmd
				m.params, cmd = m.params.HandleClick(cRelX, cRelY)
				return m, cmd

			case tabHeaders:
				var cmd tea.Cmd
				m.headers, cmd = m.headers.HandleClick(cRelX, cRelY)
				return m, cmd

			case tabAuth:
				m.auth = m.auth.HandleClick(cRelX, cRelY)
				return m, nil
			}
			return m, nil

		default:
			m = m.blurAllFields()
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

	m.body, cmd = m.body.Update(msg)
	cmds = append(cmds, cmd)

	m.params, cmd = m.params.Update(msg)
	cmds = append(cmds, cmd)

	m.headers, cmd = m.headers.Update(msg)
	cmds = append(cmds, cmd)

	m.auth, cmd = m.auth.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// bodyContentInset accounts for the border drawn around the body editor
// by BodyModel, so translated click coordinates land on the correct
// character inside the textarea rather than on its border.
const bodyContentInset = 1

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

// dropdownListView renders the floating list of HTTP methods shown when
// the method dropdown is open.
func (m RequestModel) dropdownListView() string {
	return dropdownBoxStyle.
		BorderForeground(styles.BorderActive).
		Render(m.methodDropdown.View())
}

// tabContentView renders whatever should appear below the tab bar for
// the currently active tab.
func (m RequestModel) tabContentView(width, height int) string {
	switch requestTab(m.tabs.Active) {
	case tabBody:
		return m.body.View()
	case tabParams:
		return m.params.View()
	case tabAuth:
		return m.auth.View()
	case tabHeaders:
		return m.headers.View()
	default:
		return components.Placeholder("", width, height)
	}
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

	_, _, contentHeight := m.verticalLayout()
	contentWidth := m.tabContentWidth()

	tabsRow := m.tabs.View()
	tabContent := helper.ClampLines(m.tabContentView(contentWidth, contentHeight), contentWidth, contentHeight)

	// Only include a gap block when it's actually >0 rows — Style.Height(0)
	// on an empty string still renders as 1 line (Height only pads UP,
	// never truncates down), so a "0-height" gap was silently adding a
	// real row to the layout that tabsY/contentY never accounted for.
	parts := []string{row}
	if urlTabsGap > 0 {
		parts = append(parts, lv2.NewStyle().Height(urlTabsGap).Render(""))
	}
	parts = append(parts, tabsRow)
	if tabsContentGap > 0 {
		parts = append(parts, lv2.NewStyle().Height(tabsContentGap).Render(""))
	}
	parts = append(parts, tabContent)

	base := lv2.JoinVertical(lv2.Left, parts...)

	var content string
	if m.dropdownOpen {
		dropdown := m.dropdownListView()
		layerBase := lv2.NewLayer(base).X(0).Y(0)
		layerDropdown := lv2.NewLayer(dropdown).X(0).Y(requestRowHeight).Z(1)
		content = lv2.NewCompositor(layerBase, layerDropdown).Render()
	} else {
		content = base
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
