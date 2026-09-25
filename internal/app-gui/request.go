package appgui

import "github.com/go-gui-org/go-gui/gui"

func requestAreaView(w *gui.Window) gui.View {
	app := gui.State[App](w)

	if app.RequestDock == nil || len(app.OpenTabs) == 0 {
		return emptyRequestAreaView()
	}

	panels := make([]gui.DockPanelDef, len(app.OpenTabs))
	for i, t := range app.OpenTabs {
		t := t // capture
		panels[i] = gui.DockPanelDef{
			ID:      t.ID,
			Content: []gui.View{requestTabView(w, t)},
		}
	}

	return gui.DockLayout(gui.DockLayoutCfg{
		ID:     "request-tabs-dock",
		Sizing: gui.FillFill,
		Root:   app.RequestDock,
		Panels: panels,
		OnLayoutChange: func(root *gui.DockNode, ctx gui.EventCtx) {
			gui.State[App](ctx.Window).RequestDock = root
		},
		OnPanelSelect: func(groupID, panelID string, ctx gui.EventCtx) {
			gui.State[App](ctx.Window).ActiveTabID = panelID
		},
	})
}

func emptyRequestAreaView() gui.View {
	return gui.Column(gui.ContainerCfg{
		Sizing:      gui.FillFill,
		Padding:     gui.PadAll(8),
		ColorBorder: gui.RGB(90, 90, 96),
		SizeBorder:  gui.SomeF(1),
		Content: []gui.View{
			gui.Label("Pilih request di sidebar untuk membukanya di sini", gui.TextStyle{Size: 13}),
		},
	})
}

func requestTabView(w *gui.Window, tab RequestTab) gui.View {
	return gui.Column(gui.ContainerCfg{
		Sizing:      gui.FillFill,
		Padding:     gui.PadAll(8),
		ColorBorder: gui.RGB(90, 90, 96),
		SizeBorder:  gui.SomeF(1),
		Content: []gui.View{
			gui.Label(tab.Method+"  "+tab.Title, gui.TextStyle{Size: 13}.Bold()),
			gui.TextButton("send-btn-"+tab.ID, "Send", func(ctx gui.EventCtx) {
				gui.State[App](ctx.Window).Clicks++
			}),
		},
	})
}
