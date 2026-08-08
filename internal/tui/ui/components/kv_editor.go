package components

import (
	"restclient/internal/tui/helper"
	"restclient/internal/tui/styles"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	lv2 "charm.land/lipgloss/v2"
)

var fieldBoxStyle = lv2.NewStyle().
	Border(lv2.RoundedBorder())

func fieldFrame() int {
	return fieldBoxStyle.GetHorizontalFrameSize()
}

func renderField(content string, focused bool) string {
	borderColor := styles.BorderNormal
	if focused {
		borderColor = styles.BorderActive
	}
	return fieldBoxStyle.BorderForeground(borderColor).Render(content)
}

type KVRow struct {
	Enabled bool
	Key     textinput.Model
	Value   textinput.Model
}

func newKVRow() KVRow {
	key := textinput.New()
	key.Placeholder = "key"
	key.SetVirtualCursor(true)

	value := textinput.New()
	value.Placeholder = "value"
	value.SetVirtualCursor(true)

	return KVRow{Enabled: true, Key: key, Value: value}
}

const (
	addRowLabel    = "+ Add"
	deleteLabel    = "x"
	deleteColWidth = 3
)

// KeyValueEditor is a scrollable table-like editor for lists of
// key/value pairs, used for request Params and Headers.
type KeyValueEditor struct {
	Rows []KVRow

	checkboxWidth int
	gapWidth      int

	keyOuterWidth   int
	valueOuterWidth int

	viewport viewport.Model
}

type kvHit struct {
	row      int
	col      string
	isAddRow bool
	isDelete bool
}

// KV is a plain key/value pair used to load external data (e.g. parsed
// from a saved request file) into a KeyValueEditor.
type KV struct {
	Key   string
	Value string
}

func NewKeyValueEditor() KeyValueEditor {
	return KeyValueEditor{
		Rows:          []KVRow{newKVRow()},
		checkboxWidth: 4,
		gapWidth:      1,
		viewport:      viewport.New(),
	}
}

func (m KeyValueEditor) rowHeight() int {
	return lv2.Height(renderField("", false))
}

func (m KeyValueEditor) contentLineHeight() int {
	return len(m.Rows)*m.rowHeight() + 1
}

func fieldContentOverhead() int {
	probe := textinput.New()
	probe.SetWidth(20)
	overhead := lv2.Width(probe.View()) - 20
	if overhead < 0 {
		overhead = 0
	}
	return overhead
}

func (m KeyValueEditor) SetSize(width, height int) KeyValueEditor {
	frame := fieldFrame()
	overhead := fieldContentOverhead()

	fixedCols := m.checkboxWidth + m.gapWidth +
		frame + overhead + m.gapWidth +
		frame + overhead + m.gapWidth +
		deleteColWidth

	remaining := width - fixedCols + 2
	if remaining < 8 {
		remaining = 8
	}

	keyInner := remaining / 2
	if keyInner < 4 {
		keyInner = 4
	}
	valueInner := remaining - keyInner
	if valueInner < 4 {
		valueInner = 4
	}

	m.keyOuterWidth = keyInner + frame + overhead
	m.valueOuterWidth = valueInner + frame + overhead

	for i := range m.Rows {
		m.Rows[i].Key.SetWidth(keyInner)
		m.Rows[i].Value.SetWidth(valueInner)
	}

	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)

	content := helper.ClampLines(m.rowsView(), width, m.contentLineHeight())
	m.viewport.SetContent(content)

	return m
}

// SetPairs replaces all rows with the given pairs. If pairs is empty,
// a single blank row is kept so the editor still shows an editable row.
func (m KeyValueEditor) SetPairs(pairs []KV) KeyValueEditor {
	overhead := fieldContentOverhead()
	keyWidth := m.keyOuterWidth - fieldFrame() - overhead
	valueWidth := m.valueOuterWidth - fieldFrame() - overhead

	rows := make([]KVRow, 0, len(pairs))
	for _, p := range pairs {
		row := newKVRow()
		row.Key.SetWidth(keyWidth)
		row.Value.SetWidth(valueWidth)
		row.Key.SetValue(p.Key)
		row.Value.SetValue(p.Value)
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		rows = append(rows, newKVRow())
	}

	m.Rows = rows
	m.viewport.SetContent(helper.ClampLines(m.rowsView(), m.viewport.Width(), m.contentLineHeight()))
	return m
}

func (m KeyValueEditor) hitTest(relX, relY int) (kvHit, bool) {
	if relY < 0 || relX < 0 {
		return kvHit{}, false
	}

	rh := m.rowHeight()
	rowIndex := relY / rh
	rowLocalY := relY % rh

	if rowIndex < len(m.Rows) {
		keyStart := m.checkboxWidth + m.gapWidth
		keyEnd := keyStart + m.keyOuterWidth
		valueStart := keyEnd + m.gapWidth
		valueEnd := valueStart + m.valueOuterWidth
		deleteStart := valueEnd + m.gapWidth
		deleteEnd := deleteStart + deleteColWidth

		switch {
		case relX < m.checkboxWidth:
			return kvHit{row: rowIndex, col: "checkbox"}, true

		case relX >= keyStart && relX < keyEnd:
			if rowLocalY == 0 || rowLocalY == rh-1 {
				return kvHit{}, false
			}
			return kvHit{row: rowIndex, col: "key"}, true

		case relX >= valueStart && relX < valueEnd:
			if rowLocalY == 0 || rowLocalY == rh-1 {
				return kvHit{}, false
			}
			return kvHit{row: rowIndex, col: "value"}, true

		case relX >= deleteStart && relX < deleteEnd:
			return kvHit{row: rowIndex, isDelete: true}, true
		}
		return kvHit{}, false
	}

	if rowIndex == len(m.Rows) {
		if relX < lv2.Width(addRowLabel) {
			return kvHit{isAddRow: true}, true
		}
	}

	return kvHit{}, false
}

func (m KeyValueEditor) blurAll() KeyValueEditor {
	for i := range m.Rows {
		m.Rows[i].Key.Blur()
		m.Rows[i].Value.Blur()
	}
	return m
}

// HandleClick processes a mouse click at coordinates relative to the
// top-left of the editor's visible viewport.
func (m KeyValueEditor) HandleClick(relX, relY int) (KeyValueEditor, tea.Cmd) {
	scrolledY := relY + m.viewport.YOffset()

	hit, ok := m.hitTest(relX, scrolledY)
	m = m.blurAll()

	if !ok {
		m.viewport.SetContent(helper.ClampLines(m.rowsView(), m.viewport.Width(), m.contentLineHeight()))
		return m, nil
	}

	switch {
	case hit.isAddRow:
		row := newKVRow()
		overhead := fieldContentOverhead()
		row.Key.SetWidth(m.keyOuterWidth - fieldFrame() - overhead)
		row.Value.SetWidth(m.valueOuterWidth - fieldFrame() - overhead)
		m.Rows = append(m.Rows, row)

	case hit.isDelete:
		m.Rows = append(m.Rows[:hit.row], m.Rows[hit.row+1:]...)

	case hit.col == "checkbox":
		m.Rows[hit.row].Enabled = !m.Rows[hit.row].Enabled

	case hit.col == "key":
		m.Rows[hit.row].Key.Focus()

	case hit.col == "value":
		m.Rows[hit.row].Value.Focus()
	}

	m.viewport.SetContent(helper.ClampLines(m.rowsView(), m.viewport.Width(), m.contentLineHeight()))
	return m, nil
}

func (m KeyValueEditor) Update(msg tea.Msg) (KeyValueEditor, tea.Cmd) {
	var cmds []tea.Cmd

	for i := range m.Rows {
		var cmd tea.Cmd
		m.Rows[i].Key, cmd = m.Rows[i].Key.Update(msg)
		cmds = append(cmds, cmd)
		m.Rows[i].Value, cmd = m.Rows[i].Value.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.viewport.SetContent(helper.ClampLines(m.rowsView(), m.viewport.Width(), m.contentLineHeight()))

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m KeyValueEditor) rowsView() string {
	lines := make([]string, 0, len(m.Rows)+1)

	for _, row := range m.Rows {
		checkbox := "[ ]"
		if row.Enabled {
			checkbox = "[x]"
		}
		checkboxView := lv2.NewStyle().
			Foreground(styles.TextSubtle).
			Width(m.checkboxWidth).
			Height(m.rowHeight()).
			AlignVertical(lv2.Center).
			Render(checkbox)

		gap := lv2.NewStyle().Width(m.gapWidth).Render("")

		keyView := renderField(row.Key.View(), row.Key.Focused())
		valueView := renderField(row.Value.View(), row.Value.Focused())

		deleteView := lv2.NewStyle().
			Foreground(styles.Danger).
			Width(deleteColWidth).
			Height(m.rowHeight()).
			AlignVertical(lv2.Center).
			Render(deleteLabel)

		line := lv2.JoinHorizontal(lv2.Top,
			checkboxView, gap, keyView, gap, valueView, gap, deleteView)
		lines = append(lines, line)
	}

	addRow := lv2.NewStyle().
		Foreground(styles.TextSubtle).
		Render(addRowLabel)
	lines = append(lines, addRow)

	return lv2.JoinVertical(lv2.Left, lines...)
}

func (m KeyValueEditor) View() string {
	return m.viewport.View()
}
