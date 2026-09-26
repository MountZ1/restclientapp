package appgui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gui-org/go-gui/gui"
)

const (
	sidebarPc       = 0.225
	sidebarPadding  = 1
	sidebarBorder   = 1
	sidebarSpacing  = 6
	treeTextSize    = 16
	treeIconSize    = 16
	searchRowHeight = 3.6 + 9
)

var methodAbbrev = map[string]string{
	"DELETE":  "DEL",
	"OPTIONS": "OPT",
}

var methodColors = map[string]gui.Color{
	"GET":     gui.RGB(158, 206, 106), // #9ece6a
	"POST":    gui.RGB(224, 175, 104), // #e0af68
	"PUT":     gui.RGB(122, 162, 247), // #7aa2f7
	"PATCH":   gui.RGB(187, 154, 247), // #bb9af7
	"DELETE":  gui.RGB(247, 118, 142), // #f7768e
	"OPTIONS": gui.RGB(86, 95, 137),   // #565f89
	"HEAD":    gui.RGB(86, 95, 137),   // #565f89
}

func sidebarView(w *gui.Window) gui.View {
	ww, wh := w.WindowSize()
	app := gui.State[App](w)
	borderColor := gui.RGB(90, 90, 96)

	sidebarWitdh := float32(ww) * sidebarPc

	treeMaxHeight := float32(wh) -
		(sidebarPadding * 2) -
		(sidebarBorder * 2) -
		searchRowHeight -
		sidebarSpacing

	return gui.Column(gui.ContainerCfg{
		ID:          "sidebar-container",
		Width:       sidebarWitdh,
		Height:      float32(wh),
		Sizing:      gui.FixedFixed,
		Padding:     gui.PadAll(sidebarPadding),
		Spacing:     gui.SomeF(sidebarSpacing),
		ColorBorder: borderColor,
		SizeBorder:  gui.SomeF(sidebarBorder),
		Content: []gui.View{
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
					}),
				},
			}),
			gui.Tree(gui.TreeCfg{
				ID:        "saved-requests-tree",
				Sizing:    gui.FillFit,
				MaxHeight: treeMaxHeight,
				Padding:   gui.NewPadding(0, 4, 20, 4),
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

	nodeTextStyle := gui.TextStyle{Size: treeTextSize, Color: gui.RGB(255, 255, 255)}

	for _, entry := range entries {
		id := idPrefix + "/" + entry.Name()
		fullPath := filepath.Join(basePath, entry.Name())

		if entry.IsDir() {
			subEntries, err := os.ReadDir(fullPath)
			if err != nil {
				subEntries = nil
			}
			nodes = append(nodes, gui.TreeNodeCfg{
				ID:        id,
				Text:      entry.Name(),
				Icon:      gui.IconFolder,
				TextStyle: nodeTextStyle,
				Nodes:     buildRequestTree(fullPath, id, subEntries, lookup),
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
			iconStyle = gui.TextStyle{Color: methodColors[method], Size: treeIconSize}
		} else {
			iconLabel = gui.IconFile
			iconStyle = gui.TextStyle{Size: treeIconSize, Color: gui.White}
		}

		nodes = append(nodes, gui.TreeNodeCfg{
			ID:            id,
			Text:          name,
			Icon:          iconLabel,
			TextStyle:     nodeTextStyle,
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
