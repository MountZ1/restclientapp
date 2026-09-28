package components

import (
	"os"
	"path/filepath"
	"strings"

	services "restclient/internal/services/fs"
	"restclient/internal/types"

	"github.com/go-gui-org/go-gui/gui"
)

const (
	treeTextSize = 14
	treeIconSize = 14
	treeIndent   = float32(12)
	iconColWidth = float32(38)
)

var methodAbbrev = map[string]string{
	"DELETE":  "DEL",
	"OPTIONS": "OPT",
}

var methodColors = map[string]gui.Color{
	"GET":     gui.RGB(158, 206, 106),
	"POST":    gui.RGB(224, 175, 104),
	"PUT":     gui.RGB(122, 162, 247),
	"PATCH":   gui.RGB(187, 154, 247),
	"DELETE":  gui.RGB(247, 118, 142),
	"OPTIONS": gui.RGB(86, 95, 137),
	"HEAD":    gui.RGB(86, 95, 137),
}

var treeBorderColor = gui.RGB(90, 90, 96)

func RenderRequestTree(w *gui.Window, basePath string, maxHeight float32, onSelectFile func(id string)) gui.View {
	app := gui.State[types.App](w)

	lookup := make(map[string]types.RequestTab)
	pathByID := map[string]string{"root": basePath}

	entries, _ := os.ReadDir(basePath)
	rows := buildRows(w, basePath, "root", entries, lookup, pathByID, 0, onSelectFile)

	app.RequestLookup = lookup
	app.PathByID = pathByID

	content := rows
	if app.ContextMenuOpen {
		content = append(content, contextMenuBackdrop(w))
	}

	return gui.Column(gui.ContainerCfg{
		ID:         "saved-requests-tree",
		Sizing:     gui.FillFit,
		MaxHeight:  maxHeight,
		Scrollable: true,
		Spacing:    gui.NoSpacing,
		Content:    content,
	})
}

func buildRows(
	w *gui.Window, basePath, idPrefix string, entries []os.DirEntry,
	lookup map[string]types.RequestTab, pathByID map[string]string,
	depth int, onSelectFile func(id string),
) []gui.View {
	app := gui.State[types.App](w)
	rows := make([]gui.View, 0, len(entries))
	nodeTextStyle := gui.TextStyle{Size: treeTextSize, Color: gui.White}

	for _, entry := range entries {
		id := idPrefix + "/" + entry.Name()
		fullPath := filepath.Join(basePath, entry.Name())
		pathByID[id] = fullPath
		nodeID := id

		if entry.IsDir() {
			expanded := app.ExpandedNodes[nodeID]

			chevron := gui.Text(gui.TextCfg{
				Text:      ">",
				TextStyle: gui.TextStyle{Size: treeIconSize, Color: gui.White},
			})
			var icon gui.View = chevron
			if expanded {
				icon = gui.RotatedBox(gui.RotatedBoxCfg{
					Content:      chevron,
					QuarterTurns: 1,
				})
			}

			rows = append(rows, treeRow(w, treeRowParams{
				ID:        nodeID,
				IconView:  icon,
				Text:      entry.Name(),
				TextStyle: nodeTextStyle,
				Depth:     depth,
				IsDir:     true,
				OnLeftClick: func(ctx gui.EventCtx) {
					a := gui.State[types.App](ctx.Window)
					if a.ExpandedNodes == nil {
						a.ExpandedNodes = map[string]bool{}
					}
					a.ExpandedNodes[nodeID] = !a.ExpandedNodes[nodeID]
					ctx.Consume()
				},
			}))

			if expanded {
				subEntries, _ := os.ReadDir(fullPath)
				rows = append(rows, buildRows(w, fullPath, nodeID, subEntries, lookup, pathByID, depth+1, onSelectFile)...)
			}
			continue
		}

		method, name := parseRequestFilename(entry.Name())
		displayMethod := method
		if short, ok := methodAbbrev[method]; ok {
			displayMethod = short
		}

		iconLabel := "[???]"
		iconStyle := gui.TextStyle{Size: treeIconSize, Color: methodColors["OPTIONS"]}
		if method != "" {
			iconLabel = "[" + displayMethod + "]"
			iconStyle = gui.TextStyle{Color: methodColors[method], Size: treeIconSize}
		}

		lookup[nodeID] = types.RequestTab{ID: nodeID, Title: name, Method: method}

		rows = append(rows, treeRow(w, treeRowParams{
			ID:        nodeID,
			IconView:  gui.Text(gui.TextCfg{Text: iconLabel, TextStyle: iconStyle}),
			Text:      name,
			TextStyle: nodeTextStyle,
			Depth:     depth,
			IsDir:     false,
			OnLeftClick: func(ctx gui.EventCtx) {
				gui.State[types.App](ctx.Window).SelectedID = nodeID
				onSelectFile(nodeID)
				ctx.Consume()
			},
		}))
	}
	return rows
}

func parseRequestFilename(filename string) (method, name string) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	parts := strings.SplitN(base, "-", 2)
	if len(parts) == 2 {
		return strings.ToUpper(parts[0]), parts[1]
	}
	return "", base
}

type treeRowParams struct {
	ID          string
	IconView    gui.View
	Text        string
	TextStyle   gui.TextStyle
	Depth       int
	IsDir       bool
	OnLeftClick func(ctx gui.EventCtx)
}

func treeRow(w *gui.Window, p treeRowParams) gui.View {
	app := gui.State[types.App](w)
	rowColor := gui.ColorTransparent
	if app.SelectedID == p.ID {
		rowColor = gui.RGB(45, 50, 66)
	}

	toggleMenu := func(ctx gui.EventCtx) {
		a := gui.State[types.App](ctx.Window)
		if a.ContextMenuOpen && a.ContextMenuTargetID == p.ID {
			CloseContextMenu(ctx.Window)
		} else {
			OpenContextMenu(ctx.Window, p.ID, p.IsDir)
		}
		ctx.Consume()
	}

	content := []gui.View{
		gui.Row(gui.ContainerCfg{
			Width:   iconColWidth,
			Sizing:  gui.FixedFit,
			VAlign:  gui.VAlignMiddle,
			HAlign:  gui.HAlignCenter,
			Content: []gui.View{p.IconView},
		}),
		gui.Text(gui.TextCfg{Text: p.Text, TextStyle: p.TextStyle, Sizing: gui.FillFit}),
		gui.Button(gui.ButtonCfg{
			ID:     "menu-btn-" + p.ID,
			Width:  18,
			Height: 1.8,
			Radius: gui.SomeF(4),
			Content: []gui.View{
				gui.Text(gui.TextCfg{
					Text:      "...",
					TextStyle: gui.TextStyle{Size: 14, Color: p.TextStyle.Color},
				}),
			},
			OnClick: toggleMenu,
		}),
	}

	if app.ContextMenuOpen && app.ContextMenuTargetID == p.ID {
		content = append(content, contextMenuPopup(w, p.ID, app.PathByID[p.ID], app.ContextMenuIsDir))
	}

	return gui.Row(gui.ContainerCfg{
		ID:      "row-" + p.ID,
		Sizing:  gui.FillFit,
		VAlign:  gui.VAlignMiddle,
		Padding: gui.NewPadding(4, 8, 4, float32(p.Depth)*treeIndent+8),
		Color:   rowColor,
		Radius:  gui.SomeF(4),
		Spacing: gui.SomeF(6),
		Content: content,
		OnAnyClick: func(ctx gui.EventCtx) {
			if ctx.Event.MouseButton == gui.MouseRight {
				toggleMenu(ctx)
				return
			}
			if p.OnLeftClick != nil {
				p.OnLeftClick(ctx)
			}
		},
	})
}

func OpenContextMenu(w *gui.Window, id string, isDir bool) {
	app := gui.State[types.App](w)
	app.ContextMenuOpen = true
	app.ContextMenuTargetID = id
	app.ContextMenuIsDir = isDir
}

func CloseContextMenu(w *gui.Window) {
	gui.State[types.App](w).ContextMenuOpen = false
}

func contextMenuBackdrop(w *gui.Window) gui.View {
	ww, wh := w.WindowSize()
	return gui.Column(gui.ContainerCfg{
		ID:           "context-menu-backdrop",
		Float:        true,
		FloatAnchor:  gui.FloatTopLeft,
		FloatTieOff:  gui.FloatTopLeft,
		FloatOffsetX: 0,
		FloatOffsetY: 0,
		FloatZIndex:  199,
		Width:        float32(ww),
		Height:       float32(wh),
		Sizing:       gui.FixedFixed,
		OnAnyClick: func(ctx gui.EventCtx) {
			CloseContextMenu(ctx.Window)
			ctx.Consume()
		},
	})
}

func contextMenuPopup(w *gui.Window, targetID, targetPath string, isDir bool) gui.View {
	whiteText := gui.TextStyle{Color: gui.White, Size: 14}

	menuItem := func(label string, onClick func(ctx gui.EventCtx)) gui.View {
		return gui.Button(gui.ButtonCfg{
			ID:      "ctxmenu-" + label + "-" + targetID,
			Sizing:  gui.FillFit,
			Height:  3,
			Radius:  gui.SomeF(4),
			Padding: gui.NewPadding(4, 12, 4, 12),
			Content: []gui.View{gui.Text(gui.TextCfg{Text: label, TextStyle: whiteText})},
			OnClick: onClick,
		})
	}

	items := []gui.View{}
	if isDir {
		items = append(items, menuItem("Create", func(ctx gui.EventCtx) {
			CloseContextMenu(ctx.Window)
			ctx.Consume()
			ShowAddDialogAtFn(ctx.Window, targetPath)
		}))
	}
	items = append(items,
		menuItem("Rename", func(ctx gui.EventCtx) {
			CloseContextMenu(ctx.Window)
			ctx.Consume()
			ShowRenameDialogFn(ctx.Window, targetPath)
		}),
		menuItem("Duplicate", func(ctx gui.EventCtx) {
			_ = services.DuplicateRequestOrCollection(targetPath, isDir)
			CloseContextMenu(ctx.Window)
			ctx.Consume()
		}),
		menuItem("Delete", func(ctx gui.EventCtx) {
			_ = services.DestroyRequestOrCollection(targetPath)
			CloseContextMenu(ctx.Window)
			ctx.Consume()
		}),
	)

	return gui.Column(gui.ContainerCfg{
		ID:           "context-menu-popup-" + targetID,
		Float:        true,
		FloatAnchor:  gui.FloatBottomRight,
		FloatTieOff:  gui.FloatTopRight,
		FloatOffsetY: 4,
		FloatZIndex:  200,
		Color:        gui.RGB(30, 30, 36),
		ColorBorder:  treeBorderColor,
		SizeBorder:   gui.SomeF(1),
		Radius:       gui.SomeF(8),
		Padding:      gui.PadAll(4),
		MinWidth:     160,
		Content:      items,
	})
}

var (
	ShowAddDialogAtFn  func(w *gui.Window, location string)
	ShowRenameDialogFn func(w *gui.Window, path string)
)
