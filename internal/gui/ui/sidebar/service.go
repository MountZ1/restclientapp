package sidebar

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func showSidebarMenu(canvas fyne.Canvas, pos fyne.Position, uid string, tipe string) {
	addRequest := fyne.NewMenuItem("Add request", func() {
		fmt.Println("add request")
	})
	addCollection := fyne.NewMenuItem("Add collection", func() {
		fmt.Println("add collection")
	})
	deleteAll := fyne.NewMenuItem("Delete all", func() {
		fmt.Println("delete all")
	})

	menu := fyne.NewMenu("", addRequest, addCollection, deleteAll)
	widget.ShowPopUpMenuAtPosition(menu, canvas, pos)
}
