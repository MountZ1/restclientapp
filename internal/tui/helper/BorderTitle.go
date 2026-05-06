package helper

import (
	"image/color"
	"strings"

	lv2 "charm.land/lipgloss/v2"
)

func RenderWithTitle(content, title string, width, height int, borderColor color.Color, extraLayers ...*lv2.Layer) string {
	box := lv2.NewStyle().
		Width(width).
		Height(height).
		Border(lv2.RoundedBorder()).
		BorderForeground(borderColor)

	layerBox := lv2.NewLayer(box.Render("")).X(0).Y(0)
	layerContent := lv2.NewLayer(content).X(1).Y(1)

	allLayers := append([]*lv2.Layer{layerBox, layerContent}, extraLayers...)
	output := lv2.NewCompositor(allLayers...).Render()

	outputLines := strings.Split(output, "\n")
	if len(outputLines) > 0 {
		title := "[ " + title + " ]"
		firstLineWidth := lv2.Width(outputLines[0])
		titleWidth := lv2.Width(title)
		dashCount := firstLineWidth - titleWidth - 2
		if dashCount < 0 {
			dashCount = 0
		}
		outputLines[0] = lv2.NewStyle().
			Foreground(borderColor).
			Render("╭" + title + strings.Repeat("─", dashCount) + "╮")
	}

	return strings.Join(outputLines, "\n")
}
