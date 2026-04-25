package sidebar

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Sidebar(c fyne.Canvas) fyne.CanvasObject {
	t, _ := tree(c)

	t.onTappedSecondary = func(e *fyne.PointEvent) {
		showSidebarMenu(c, e.AbsolutePosition, "", "")
	}

	return container.NewScroll(t)
}
