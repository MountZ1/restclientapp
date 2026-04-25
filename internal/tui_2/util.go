package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func countDimentionLength(len int, setLen int) int {
	return len * setLen / 100
}

func overlayCenter(width, height int, overlay, base string) string {
	overlayLines := strings.Split(overlay, "\n")
	baseLines := strings.Split(base, "\n")

	overlayHeight := len(overlayLines)
	overlayWidth := lipgloss.Width(overlayLines[0])

	startY := (height - overlayHeight) / 2
	startX := (width - overlayWidth) / 2

	for i, line := range overlayLines {
		y := startY + i
		if y < 0 || y >= len(baseLines) {
			continue
		}
		baseLine := baseLines[y]
		// potong base di posisi startX, sisipkan overlay
		before := truncateString(baseLine, startX)
		after := skipString(baseLine, startX+overlayWidth)
		baseLines[y] = before + line + after
	}

	return strings.Join(baseLines, "\n")
}

// helper — ambil n karakter pertama dengan mempertimbangkan ANSI
func truncateString(s string, n int) string {
	return lipgloss.NewStyle().MaxWidth(n).Render(s)
}

// helper — skip n karakter pertama
func skipString(s string, n int) string {
	width := 0
	for i, r := range s {
		if width >= n {
			return s[i:]
		}
		width += lipgloss.Width(string(r))
	}
	return ""
}
