package appgui

import (
	"os"

	services "restclient/internal/services/fs"

	"github.com/go-gui-org/go-gui/gui"
	"github.com/go-gui-org/go-gui/gui/backend"
)

type RequestTab struct {
	ID     string
	Title  string
	Method string
}

type App struct {
	Clicks int

	SavedRequests []os.DirEntry
	RequestTree   []gui.TreeNodeCfg
	RequestLookup map[string]RequestTab
	SelectedID    string

	OpenTabs             []RequestTab
	ActiveTabID          string
	RequestDock          *gui.DockNode
	RequestResponseRatio float32
}

func Run() {
	gui.SetTheme(gui.ThemeDark.WithBorders(true).WithPadding(false))
	gui.Debug(true)

	basePath := services.Path()
	savedRequests := services.GetCollectionItems()

	lookup := make(map[string]RequestTab)
	tree := buildRequestTree(basePath, "root", savedRequests, lookup)

	app := &App{
		SavedRequests:        savedRequests,
		RequestTree:          tree,
		RequestLookup:        lookup,
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
	a := gui.State[App](w)

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

func syncRequestDock(a *App) {
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
	app := gui.State[App](w)

	return gui.Splitter(gui.SplitterCfg{
		ID:          "request-response-split",
		Focusable:   true,
		Orientation: gui.SplitterVertical,
		Sizing:      gui.FillFill,
		Ratio:       gui.SomeF(app.RequestResponseRatio),
		OnChange: func(ratio float32, collapsed gui.SplitterCollapsed, ctx gui.EventCtx) {
			gui.State[App](ctx.Window).RequestResponseRatio = ratio
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
