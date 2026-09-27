package appgui

import (
	"restclient/internal/types"

	"github.com/go-gui-org/go-gui/gui"
	"github.com/go-gui-org/go-gui/gui/backend"
)

func Run() {
	gui.SetTheme(gui.ThemeDark.WithBorders(true).WithPadding(false))
	gui.Debug(true)

	app := &types.App{
		RequestResponseRatio: 0.5,
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

func openRequestTab(w *gui.Window, id string) {
	a := gui.State[types.App](w)

	tab, ok := a.RequestLookup[id]
	if !ok {
		return
	}

	for _, t := range a.OpenTabs {
		if t.ID == id {
			a.ActiveTabID = id
			syncRequestDock(a)
			return
		}
	}

	a.OpenTabs = append(a.OpenTabs, tab)
	a.ActiveTabID = id
	syncRequestDock(a)
}

func syncRequestDock(a *types.App) {
	if len(a.OpenTabs) == 0 {
		a.RequestDock = nil
		return
	}
	ids := make([]string, len(a.OpenTabs))
	for i, t := range a.OpenTabs {
		ids[i] = t.ID
	}
	a.RequestDock = gui.DockPanelGroup("request-tabs", ids, a.ActiveTabID)
}

func mainView(w *gui.Window) gui.View {
	ww, wh := w.WindowSize()

	return gui.Row(gui.ContainerCfg{
		Width:  float32(ww),
		Height: float32(wh),
		Sizing: gui.FixedFixed,
		Content: []gui.View{
			sidebarView(w),
			mainContentView(w),
		},
	})
}

func mainContentView(w *gui.Window) gui.View {
	app := gui.State[types.App](w)

	return gui.Splitter(gui.SplitterCfg{
		ID:          "request-response-split",
		Focusable:   true,
		Orientation: gui.SplitterVertical,
		Sizing:      gui.FillFill,
		Ratio:       gui.SomeF(app.RequestResponseRatio),
		OnChange: func(ratio float32, collapsed gui.SplitterCollapsed, ctx gui.EventCtx) {
			gui.State[types.App](ctx.Window).RequestResponseRatio = ratio
		},
		First: gui.SplitterPaneCfg{
			MinSize: 100,
			Content: []gui.View{requestAreaView(w)},
		},
		Second: gui.SplitterPaneCfg{
			MinSize: 100,
			Content: []gui.View{responseView(w)},
		},
	})
}
