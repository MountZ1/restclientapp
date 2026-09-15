package appgui

import (
	"github.com/go-gui-org/go-gui/gui"
	"github.com/go-gui-org/go-gui/gui/backend"
)

type SavedRequest struct {
	ID     string
	Name   string
	Method string
}

type App struct {
	Clicks   int
	DockRoot *gui.DockNode

	SavedRequests []SavedRequest
	SelectedID    string
}

func Run() {
	gui.SetTheme(gui.ThemeLight)

	app := &App{
		DockRoot: initialLayout(),
		SavedRequests: []SavedRequest{
			{ID: "req-1", Name: "Get users", Method: "GET"},
			{ID: "req-2", Name: "Create user", Method: "POST"},
			{ID: "req-3", Name: "Delete user", Method: "DELETE"},
		},
	}

	w := gui.NewWindow(gui.WindowCfg{
		State:  app,
		Title:  "restclient",
		Width:  1100,
		Height: 650,
		OnInit: func(w *gui.Window) { w.SetView(mainView) },
	})

	backend.Run(w)
}

func initialLayout() *gui.DockNode {
	return gui.DockSplit("root", gui.DockSplitHorizontal, 0.25,
		gui.DockPanelGroup("left", []string{"sidebar"}, "sidebar"),
		gui.DockSplit("right", gui.DockSplitVertical, 0.5,
			gui.DockPanelGroup("top", []string{"request"}, "request"),
			gui.DockPanelGroup("bottom", []string{"response"}, "response"),
		),
	)
}

func mainView(w *gui.Window) gui.View {
	app := gui.State[App](w)
	ww, wh := w.WindowSize()

	return gui.Column(gui.ContainerCfg{
		Width:  float32(ww),
		Height: float32(wh),
		Sizing: gui.FixedFixed,
		Content: []gui.View{
			gui.DockLayout(gui.DockLayoutCfg{
				ID:   "main-dock",
				Root: app.DockRoot,
				Panels: []gui.DockPanelDef{
					{
						ID:      "sidebar",
						Label:   "Collections",
						Content: []gui.View{sidebarView(w)},
					},
					{
						ID:      "request",
						Label:   "Request",
						Content: []gui.View{requestView(w)},
					},
					{
						ID:      "response",
						Label:   "Response",
						Content: []gui.View{responseView(w)},
					},
				},
				OnLayoutChange: func(root *gui.DockNode, ctx gui.EventCtx) {
					gui.State[App](ctx.Window).DockRoot = root
				},
				OnPanelSelect: func(groupID, panelID string, ctx gui.EventCtx) {
					a := gui.State[App](ctx.Window)
					a.DockRoot = gui.DockTreeSelectPanel(a.DockRoot, groupID, panelID)
				},
			}),
		},
	})
}

func sidebarView(w *gui.Window) gui.View {
	app := gui.State[App](w)
	t := gui.CurrentTheme()

	items := make([]gui.View, 0, len(app.SavedRequests)+1)
	items = append(items, gui.Label("Saved requests", t.B4))

	for _, req := range app.SavedRequests {
		req := req
		label := req.Method + "  " + req.Name

		bg := t.ColorInterior
		if req.ID == app.SelectedID {
			bg = t.ColorSelect
		}

		items = append(items, gui.Column(gui.ContainerCfg{
			ID:      "sidebar-item-" + req.ID,
			Sizing:  gui.FillFit,
			Padding: gui.NewPadding(6, 8, 6, 8),
			Color:   bg,
			Content: []gui.View{
				gui.Label(label, t.N3),
			},
			OnClick: func(ctx gui.EventCtx) {
				gui.State[App](ctx.Window).SelectedID = req.ID
			},
		}))
	}

	return gui.Column(gui.ContainerCfg{
		Sizing:  gui.FillFill,
		Padding: gui.PadAll(8),
		Spacing: gui.SomeF(4),
		Content: items,
	})
}

func requestView(w *gui.Window) gui.View {
	t := gui.CurrentTheme()
	return gui.Column(gui.ContainerCfg{
		Sizing:  gui.FillFill,
		Padding: gui.PadAll(8),
		Content: []gui.View{
			gui.Label("Request panel", t.B3),
			gui.TextButton("send-btn", "Send", func(ctx gui.EventCtx) {
				gui.State[App](ctx.Window).Clicks++
			}),
		},
	})
}

func responseView(w *gui.Window) gui.View {
	t := gui.CurrentTheme()
	return gui.Column(gui.ContainerCfg{
		Sizing:  gui.FillFill,
		Padding: gui.PadAll(8),
		Content: []gui.View{
			gui.Label("Response panel", t.B3),
		},
	})
}
