package appgui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gui-org/go-gui/gui"
)

const (
	sidebarWidth   = 260
	sidebarPadding = 2
	sidebarBorder  = 1
)

var methodAbbrev = map[string]string{
	"DELETE":  "DEL",
	"OPTIONS": "OPT",
}

var methodColors = map[string]gui.Color{
	"GET":     gui.RGB(97, 175, 254),
	"POST":    gui.RGB(73, 204, 144),
	"PUT":     gui.RGB(252, 161, 48),
	"PATCH":   gui.RGB(80, 227, 194),
	"DELETE":  gui.RGB(249, 62, 62),
	"OPTIONS": gui.RGB(144, 18, 254),
	"HEAD":    gui.RGB(144, 144, 144),
}

func sidebarView(w *gui.Window) gui.View {
	_, wh := w.WindowSize()
	app := gui.State[App](w)
	borderColor := gui.RGB(90, 90, 96)

	treeMaxHeight := float32(wh) - (sidebarPadding * 2) - (sidebarBorder * 2)

	return gui.Column(gui.ContainerCfg{
		ID:          "sidebar-container",
		Width:       sidebarWidth,
		Height:      float32(wh),
		Sizing:      gui.FixedFixed,
		Padding:     gui.PadAll(sidebarPadding),
		Spacing:     gui.SomeF(4),
		ColorBorder: borderColor,
		SizeBorder:  gui.SomeF(sidebarBorder),
		Content: []gui.View{
			gui.Tree(gui.TreeCfg{
				ID:        "saved-requests-tree",
				Sizing:    gui.FillFit,
				MaxHeight: treeMaxHeight,
				Nodes:     app.RequestTree,
				OnSelect: func(id string, ctx gui.EventCtx) {
					gui.State[App](w).SelectedID = id
					openRequestTab(w, id)
				},
			}),
		},
	})
}

func buildRequestTree(basePath, idPrefix string, entries []os.DirEntry, lookup map[string]RequestTab) []gui.TreeNodeCfg {
	nodes := make([]gui.TreeNodeCfg, 0, len(entries))

	for _, entry := range entries {
		id := idPrefix + "/" + entry.Name()
		fullPath := filepath.Join(basePath, entry.Name())

		if entry.IsDir() {
			subEntries, err := os.ReadDir(fullPath)
			if err != nil {
				subEntries = nil
			}
			nodes = append(nodes, gui.TreeNodeCfg{
				ID:    id,
				Text:  entry.Name(),
				Icon:  gui.IconFolder,
				Nodes: buildRequestTree(fullPath, id, subEntries, lookup),
			})
			continue
		}

		method, name := parseRequestFilename(entry.Name())

		displayMethod := method
		if short, ok := methodAbbrev[method]; ok {
			displayMethod = short
		}

		var iconStyle gui.TextStyle
		var iconLabel string
		if method != "" {
			iconLabel = "[" + displayMethod + "]"
			iconStyle = gui.TextStyle{Color: methodColors[method]}
		} else {
			iconLabel = gui.IconFile
			iconStyle = gui.TextStyle{Color: methodColors["options"]}
		}

		nodes = append(nodes, gui.TreeNodeCfg{
			ID:            id,
			Text:          name,
			Icon:          iconLabel,
			TextStyleIcon: iconStyle,
		})

		lookup[id] = RequestTab{ID: id, Title: name, Method: method}
	}
	return nodes
}

func parseRequestFilename(filename string) (method, name string) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	parts := strings.SplitN(base, "-", 2)
	if len(parts) == 2 {
		return strings.ToUpper(parts[0]), parts[1]
	}
	return "", base
}
