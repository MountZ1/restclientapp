package appgui

import "github.com/go-gui-org/go-gui/gui"

func sidebarView(w *gui.Window) gui.View {
	app := gui.State[App](w)
	t := gui.CurrentTheme()

	borderColor := gui.RGB(90, 90, 96)

	items := make([]gui.View, 0, len(app.SavedRequests)+1)
	items = append(items, gui.Label("Saved requests", t.B4))

	for _, req := range app.SavedRequests {
		req := req
		label := req.Method + "  " + req.Name

		bg := t.ColorInterior
		bc := borderColor
		if req.ID == app.SelectedID {
			bg = t.ColorSelect
			bc = gui.RGB(140, 170, 255)
		}

		items = append(items, gui.Column(gui.ContainerCfg{
			ID:          "sidebar-item-" + req.ID,
			Sizing:      gui.FillFit,
			Padding:     gui.NewPadding(6, 8, 6, 8),
			Color:       bg,
			ColorBorder: bc,
			SizeBorder:  gui.SomeF(1),
			Content: []gui.View{
				gui.Label(label, t.N3),
			},
			OnClick: func(ctx gui.EventCtx) {
				gui.State[App](ctx.Window).SelectedID = req.ID
			},
		}))
	}

	return gui.Column(gui.ContainerCfg{
		Sizing:      gui.FillFill,
		Padding:     gui.PadAll(8),
		Spacing:     gui.SomeF(4),
		ColorBorder: borderColor,
		SizeBorder:  gui.SomeF(1),
		Content:     items,
	})
}
