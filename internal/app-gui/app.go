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
	gui.SetTheme(gui.ThemeDark.WithBorders(true).WithPadding(false))

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
		Width:  1200,
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
					{ID: "sidebar", Content: []gui.View{sidebarView(w)}},
					{ID: "request", Content: []gui.View{requestView(w)}},
					{ID: "response", Content: []gui.View{responseView(w)}},
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
