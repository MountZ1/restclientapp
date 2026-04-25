package ui

import (
	"restclient/internal/gui/ui/sidebar"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Windows(w fyne.Window) fyne.CanvasObject {
	// widgetLabel := widget.NewLabel("Hello World")

	window := container.NewHSplit(
		sidebar.Sidebar(w.Canvas()),
		container.NewVSplit(
			Request(),
			Response(),
		),
	)

	window.SetOffset(0.195)

	return window
}
