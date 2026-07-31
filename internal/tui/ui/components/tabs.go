package components

import (
	"restclient/internal/tui/styles"

	lv2 "charm.land/lipgloss/v2"
)

var tabsBoxStyle = lv2.NewStyle().
	Border(lv2.RoundedBorder()).
	Padding(0, 1)

// TabsModel renders a bordered row of clickable tab labels and keeps
// track of which tab is currently active.
type TabsModel struct {
	Titles []string
	Active int

	gapWidth int
}

func NewTabs(titles []string) TabsModel {
	return TabsModel{
		Titles:   titles,
		Active:   0,
		gapWidth: 2,
	}
}

func (m TabsModel) tabView(index int) string {
	title := m.Titles[index]

	style := lv2.NewStyle().
		Foreground(styles.TextSubtle).
		Padding(0, 1)

	if index == m.Active {
		style = style.
			Foreground(styles.TextNormal).
			Bold(true).
			Underline(true)
	}

	return style.Render(title)
}

type tabBounds struct {
	start, end int
}

// segmentBounds computes bounds relative to the INTERIOR of the tab
// box, i.e. after the border and padding added by tabsBoxStyle. Callers
// must subtract that same frame from raw coordinates before calling
// HitTest, see FrameOffset.
func (m TabsModel) segmentBounds() []tabBounds {
	bounds := make([]tabBounds, len(m.Titles))

	cursor := 0
	for i := range m.Titles {
		w := lv2.Width(m.tabView(i))
		bounds[i] = tabBounds{start: cursor, end: cursor + w}
		cursor += w + m.gapWidth
	}

	return bounds
}

// FrameOffsetX reports how many columns the border and padding add to
// the left of the tab content, so callers can translate raw click
// coordinates into HitTest's coordinate space.
func (m TabsModel) FrameOffsetX() int {
	return 1 + 1 // left border + left padding
}

// FrameOffsetY reports how many rows the border adds above the tab
// content.
func (m TabsModel) FrameOffsetY() int {
	return 1 // top border
}

func (m TabsModel) innerView() string {
	parts := make([]string, 0, len(m.Titles)*2-1)
	for i := range m.Titles {
		if i > 0 {
			parts = append(parts, lv2.NewStyle().Width(m.gapWidth).Render(""))
		}
		parts = append(parts, m.tabView(i))
	}
	return lv2.JoinHorizontal(lv2.Top, parts...)
}

// View renders the full bordered tab bar.
func (m TabsModel) View() string {
	return tabsBoxStyle.Render(m.innerView())
}

// HitTest returns the index of the tab at the given column, relative to
// the interior of the tab box (i.e. already offset past the border and
// padding), and whether a tab was actually hit.
func (m TabsModel) HitTest(relX int) (int, bool) {
	if relX < 0 {
		return 0, false
	}
	for i, b := range m.segmentBounds() {
		if relX >= b.start && relX < b.end {
			return i, true
		}
	}
	return 0, false
}

// Height reports how many terminal rows the bordered tab bar occupies.
func (m TabsModel) Height() int {
	return lv2.Height(m.View())
}
