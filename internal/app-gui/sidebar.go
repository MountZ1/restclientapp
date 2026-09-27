package appgui

import (
	"restclient/internal/app-gui/components"
	"restclient/internal/app-gui/features"
	services "restclient/internal/services/fs"

	"github.com/go-gui-org/go-gui/gui"
)

const (
	sidebarPc       = 0.225
	sidebarPadding  = 1
	sidebarBorder   = 1
	sidebarSpacing  = 6
	searchRowHeight = 3.6 + 9
)

func sidebarView(w *gui.Window) gui.View {
	ww, wh := w.WindowSize()
	borderColor := gui.RGB(90, 90, 96)
	sidebarWitdh := float32(ww) * sidebarPc

	treeMaxHeight := float32(wh) -
		(sidebarPadding * 2) -
		(sidebarBorder * 2) -
		searchRowHeight -
		sidebarSpacing

	content := []gui.View{
		gui.Row(gui.ContainerCfg{
			Sizing:  gui.FillFit,
			VAlign:  gui.VAlignMiddle,
			Spacing: gui.SomeF(8),
			Padding: gui.NewPadding(4, 12, 4, 12),
			Content: []gui.View{
				gui.Input(gui.InputCfg{
					ID:          "search-input",
					Placeholder: "Search request collection name",
					Sizing:      gui.FillFit,
					Height:      3.6,
					Padding:     gui.PadAll(2),
					Radius:      gui.SomeF(6),
				}),
				gui.Button(gui.ButtonCfg{
					ID:     "add-request-or-collection",
					Height: 3,
					Radius: gui.SomeF(6),
					Content: []gui.View{
						gui.Text(gui.TextCfg{Text: "Add"}),
					},
					OnClick: func(ctx gui.EventCtx) {
						features.AddCollectionOrRequest(ctx)
						ctx.Consume()
					},
				}),
			},
		}),
		components.RenderRequestTree(w, services.Path(), treeMaxHeight, func(id string) {
			openRequestTab(w, id)
		}),
	}

	return gui.Column(gui.ContainerCfg{
		ID:          "sidebar-container",
		Width:       sidebarWitdh,
		Height:      float32(wh),
		Sizing:      gui.FixedFixed,
		Padding:     gui.PadAll(sidebarPadding),
		Spacing:     gui.SomeF(sidebarSpacing),
		ColorBorder: borderColor,
		SizeBorder:  gui.SomeF(sidebarBorder),
		Content:     content,
	})
}
