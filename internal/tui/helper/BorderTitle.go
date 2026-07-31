package helper

import (
	"image/color"
	"strings"

	lv2 "charm.land/lipgloss/v2"
)

var panelFrameStyle = lv2.NewStyle().Border(lv2.RoundedBorder())

func PanelHorizontalFrame() int {
	return panelFrameStyle.GetHorizontalFrameSize()
}

func PanelVerticalFrame() int {
	return panelFrameStyle.GetVerticalFrameSize()
}

// ClampLines forces s to be EXACTLY h lines tall and no wider than w
// (only by padding short lines — never by re-rendering/reflowing).
func ClampLines(s string, w, h int) string {
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}

	lines := strings.Split(s, "\n")

	if len(lines) > h {
		lines = lines[:h]
	}
	for len(lines) < h {
		lines = append(lines, strings.Repeat(" ", w))
	}

	for i, l := range lines {
		lw := lv2.Width(l)
		if lw < w {
			lines[i] = l + strings.Repeat(" ", w-lw)
		}
	}

	return strings.Join(lines, "\n")
}

// RenderWithTitle renders content inside a titled, bordered box.
// width and height are the TOTAL outer size of the box, including the
// border.
func RenderWithTitle(content, title string, width, height int, borderColor color.Color, extraLayers ...*lv2.Layer) string {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	innerWidth := width - PanelHorizontalFrame()
	innerHeight := height - PanelVerticalFrame()
	if innerWidth < 1 {
		innerWidth = 1
	}
	if innerHeight < 1 {
		innerHeight = 1
	}

	// Content still gets sized to the INTERIOR (excludes border) —
	// this part was always correct.
	content = ClampLines(content, innerWidth, innerHeight)

	// IMPORTANT: pass the FULL outer width/height here, not
	// innerWidth/innerHeight. In this lipgloss fork, Style.Width()/
	// Height() on a bordered style already produce a TOTAL render of
	// exactly that size (border included) — confirmed empirically:
	// Width(10).Height(5).Border(...).Render("") renders exactly 10x5,
	// not 12x7. Subtracting the frame here made every box render 2
	// rows short of the requested height, and our own end-of-function
	// ClampLines padding was silently papering over that shortfall
	// with blank lines instead of the box's real bottom border.
	box := lv2.NewStyle().
		Width(width).
		Height(height).
		Border(lv2.RoundedBorder()).
		BorderForeground(borderColor)

	boxStr := box.Render("")

	layerBox := lv2.NewLayer(boxStr).X(0).Y(0)
	layerContent := lv2.NewLayer(content).X(1).Y(1)
	allLayers := append([]*lv2.Layer{layerBox, layerContent}, extraLayers...)
	output := lv2.NewCompositor(allLayers...).Render()

	// Safety net only — should be a no-op now that sizes line up.
	output = ClampLines(output, width, height)

	outputLines := strings.Split(output, "\n")
	if len(outputLines) > 0 {
		titleStr := "[ " + title + " ]"
		firstLineWidth := lv2.Width(outputLines[0])
		titleWidth := lv2.Width(titleStr)
		dashCount := firstLineWidth - titleWidth - 2
		if dashCount < 0 {
			dashCount = 0
		}
		outputLines[0] = lv2.NewStyle().
			Foreground(borderColor).
			Render("╭" + titleStr + strings.Repeat("─", dashCount) + "╮")
	}
	return strings.Join(outputLines, "\n")
}
