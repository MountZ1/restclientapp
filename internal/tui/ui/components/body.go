package components

import (
	"restclient/internal/tui/styles"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	lv2 "charm.land/lipgloss/v2"
)

type BodyType int

const (
	BodyJSON BodyType = iota
	BodyFormData
	BodyText
)

var bodyTypeLabels = []string{"JSON", "Form Data", "Text"}

var bodyBoxStyle = lv2.NewStyle().
	Border(lv2.RoundedBorder()).
	Padding(0, 1)

// BodyModel wraps a textarea used for editing the request body, plus a
// clickable row for selecting the body's content type. The textarea
// scrolls internally when content exceeds its visible height, so no
// extra scroll wrapper is needed here.
type BodyModel struct {
	Type  BodyType
	Input textarea.Model

	width int
}

func placeholderFor(t BodyType) string {
	switch t {
	case BodyFormData:
		return "key=value&another=value"
	case BodyText:
		return "Plain text body..."
	default:
		return `{ "key": "value" }`
	}
}

func NewBody() BodyModel {
	ta := textarea.New()
	ta.Placeholder = placeholderFor(BodyJSON)
	ta.SetVirtualCursor(true)
	ta.ShowLineNumbers = false
	return BodyModel{Input: ta}
}

// typeLabelView renders one content-type label, highlighted if active.
func (m BodyModel) typeLabelView(t BodyType) string {
	style := lv2.NewStyle().Foreground(styles.TextSubtle).Padding(0, 1)
	if m.Type == t {
		style = style.Foreground(styles.TextNormal).Bold(true).Underline(true)
	}
	return style.Render(bodyTypeLabels[t])
}

func (m BodyModel) typeRowView() string {
	parts := make([]string, 0, len(bodyTypeLabels)*2-1)
	for i := range bodyTypeLabels {
		if i > 0 {
			parts = append(parts, " ")
		}
		parts = append(parts, m.typeLabelView(BodyType(i)))
	}
	return lv2.JoinHorizontal(lv2.Top, parts...)
}

func (m BodyModel) typeBounds() []tabBounds {
	bounds := make([]tabBounds, len(bodyTypeLabels))
	cursor := 0
	for i := range bodyTypeLabels {
		w := lv2.Width(m.typeLabelView(BodyType(i)))
		bounds[i] = tabBounds{start: cursor, end: cursor + w}
		cursor += w + 1
	}
	return bounds
}

// TypeRowHeight reports how many rows the content-type selector row
// occupies.
func (m BodyModel) TypeRowHeight() int {
	return lv2.Height(m.typeRowView())
}

// SetSize resizes the editor to fit inside the given outer width and
// height, reserving space for the content-type row above it.
func (m BodyModel) SetSize(width, height int) BodyModel {
	m.width = width

	editorHeight := height - m.TypeRowHeight() - 1 // -1 for the gap row
	innerWidth := width - bodyBoxStyle.GetHorizontalFrameSize()
	innerHeight := editorHeight - bodyBoxStyle.GetVerticalFrameSize()
	if innerWidth < 1 {
		innerWidth = 1
	}
	if innerHeight < 1 {
		innerHeight = 1
	}
	m.Input.SetWidth(innerWidth)
	m.Input.SetHeight(innerHeight)
	return m
}

func (m BodyModel) Focus() (BodyModel, tea.Cmd) {
	cmd := m.Input.Focus()
	return m, cmd
}

func (m BodyModel) Blur() BodyModel {
	m.Input.Blur()
	return m
}

func (m BodyModel) Focused() bool {
	return m.Input.Focused()
}

// HandleTypeClick processes a click at coordinates relative to the
// top-left of the content-type row.
func (m BodyModel) HandleTypeClick(relX, relY int) BodyModel {
	if relY != 0 {
		return m
	}
	for i, b := range m.typeBounds() {
		if relX >= b.start && relX < b.end {
			m.Type = BodyType(i)
			m.Input.Placeholder = placeholderFor(m.Type)
		}
	}
	return m
}

func (m BodyModel) Update(msg tea.Msg) (BodyModel, tea.Cmd) {
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m BodyModel) View() string {
	borderColor := styles.BorderNormal
	if m.Input.Focused() {
		borderColor = styles.BorderActive
	}

	editor := bodyBoxStyle.
		BorderForeground(borderColor).
		Render(m.Input.View())

	gap := lv2.NewStyle().Height(1).Render("")

	return lv2.JoinVertical(lv2.Left, m.typeRowView(), gap, editor)
}
