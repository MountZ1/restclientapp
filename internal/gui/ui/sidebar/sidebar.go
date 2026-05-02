package sidebar

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Sidebar(c fyne.Canvas) fyne.CanvasObject {
	t, _ := tree(c)
	t.onTappedSecondary = func(e *fyne.PointEvent) {
		showSidebarMenu(c, e.AbsolutePosition, "", "")
	}

	input := widget.NewEntry()
	input.SetPlaceHolder("Search")

	btn := widget.NewButton("...", nil)
	btn.Importance = widget.LowImportance
	btn.OnTapped = func() {
		pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(btn)
		pos.Y += btn.Size().Height
		showSidebarMenu(c, pos, "", "")
	}

	topBar := container.NewBorder(nil, nil, nil, btn, input)

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(0, 8))

	top := container.NewVBox(topBar, spacer)

	return container.NewBorder(
		top,
		nil,
		nil,
		nil,
		container.NewScroll(t),
	)
}
