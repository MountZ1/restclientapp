package sidebar

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Menu struct {
	label   string
	action  func()
	showFor []string
}

func showSidebarMenu(canvas fyne.Canvas, pos fyne.Position, uid string, tipe string) {
	menus := []Menu{
		{
			label: "Add Request",
			action: func() {
				fmt.Println("Add Request")
			},
			showFor: []string{"folder", ""},
		},
		{
			label: "Add Collecction",
			action: func() {
				fmt.Println("Add Collection")
			},
			showFor: []string{"folder", ""},
		},
		{
			label: "Delete Request",
			action: func() {
				fmt.Println("Delete Request")
			},
			showFor: []string{"file"},
		},
		{
			label: "Delete Collection",
			action: func() {
				fmt.Println("Delete Collection")
			},
			showFor: []string{"folder"},
		},
		{
			label: "Duplicate",
			action: func() {
				fmt.Println("Duplicate Collection")
			},
			showFor: []string{"folder", "file"},
		},
		{
			label: "Rename",
			action: func() {
				fmt.Println("Rename Collection")
			},
			showFor: []string{"folder", "file"},
		},
	}

	var menuCollection []*fyne.MenuItem
	for _, menu := range menus {
		for _, permission := range menu.showFor {
			if permission == tipe {
				menuCollection = append(menuCollection, fyne.NewMenuItem(menu.label, menu.action))
			}
		}
	}

	menu := fyne.NewMenu("", menuCollection...)
	widget.ShowPopUpMenuAtPosition(menu, canvas, pos)
}
