package appgui

import (
	"fmt"

	"restclient/internal/types"

	"github.com/go-gui-org/go-gui/gui"
)

const requestDockGroup = "request-tabs"

var (
	httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

	uiSurface               = gui.RGB(30, 33, 38)
	uiRadius                = gui.SomeF(8)
	uiGap                   = gui.SomeF(12)
	uiGapSmall              = gui.SomeF(8)
	uiControlHeight float32 = 40

	uiFieldPadding  = gui.NewPadding(0, 14, 0, 14)
	uiButtonPadding = gui.NewPadding(0, 20, 0, 20)
	uiBodyPadding   = gui.PadAll(14)
	uiPanelPadding  = gui.PadAll(16)
	uiTabPadding    = gui.NewPadding(8, 16, 8, 16)
	uiContentPad    = gui.NewPadding(14, 0, 0, 0)
)

func updateTab(w *gui.Window, id string, fn func(t *types.RequestTab)) {
	a := gui.State[types.App](w)
	for i := range a.OpenTabs {
		if a.OpenTabs[i].ID == id {
			fn(&a.OpenTabs[i])
			return
		}
	}
}

func selectDockPanel(root *gui.DockNode, id string) *gui.DockNode {
	if root == nil {
		return nil
	}
	if g, ok := gui.DockTreeFindGroupByPanel(root, id); ok {
		return gui.DockTreeSelectPanel(root, g.ID, id)
	}
	return root
}

func addDockPanel(root *gui.DockNode, activeID, id string) *gui.DockNode {
	if root == nil {
		return gui.DockPanelGroup(requestDockGroup, []string{id}, id)
	}

	groupID := requestDockGroup
	if g, ok := gui.DockTreeFindGroupByPanel(root, activeID); ok {
		groupID = g.ID
	}

	root = gui.DockTreeAddTab(root, groupID, id)
	return gui.DockTreeSelectPanel(root, groupID, id)
}

func closeRequestTab(w *gui.Window, id string) {
	flushSave(w, id)

	a := gui.State[types.App](w)
	for i, t := range a.OpenTabs {
		if t.ID == id {
			a.OpenTabs = append(a.OpenTabs[:i], a.OpenTabs[i+1:]...)
			break
		}
	}

	if len(a.OpenTabs) == 0 {
		a.RequestDock = nil
		a.ActiveTabID = ""
		return
	}

	a.RequestDock = gui.DockTreeRemovePanel(a.RequestDock, id)
	if a.ActiveTabID == id {
		a.ActiveTabID = a.OpenTabs[len(a.OpenTabs)-1].ID
		a.RequestDock = selectDockPanel(a.RequestDock, a.ActiveTabID)
	}
}

func controlRow(content ...gui.View) gui.View {
	return gui.Row(gui.ContainerCfg{
		Height:  uiControlHeight,
		Sizing:  gui.FillFixed,
		Padding: gui.NoPadding,
		Spacing: uiGapSmall,
		VAlign:  gui.VAlignMiddle,
		Content: content,
	})
}

func inputField(cfg gui.InputCfg) gui.View {
	cfg.Sizing = gui.FillFill
	cfg.Padding = uiFieldPadding
	cfg.Radius = uiRadius
	return gui.Input(cfg)
}

func actionButton(id, text string, onClick func(gui.EventCtx)) gui.View {
	return gui.Button(gui.ButtonCfg{
		ID:      id,
		Sizing:  gui.FitFill,
		Padding: uiButtonPadding,
		Radius:  uiRadius,
		Content: []gui.View{gui.Text(gui.TextCfg{Text: text})},
		OnClick: onClick,
	})
}

func labeledField(label string, field gui.View) gui.View {
	return gui.Column(gui.ContainerCfg{
		Sizing:  gui.FillFit,
		Padding: gui.NoPadding,
		Spacing: gui.SomeF(6),
		Content: []gui.View{
			gui.Label(label, gui.TextStyle{Size: 12}),
			controlRow(field),
		},
	})
}

func styledTabs(cfg gui.TabControlCfg) gui.TabControlCfg {
	cfg.Sizing = gui.FillFill
	cfg.Padding = gui.NoPadding
	cfg.SizeBorder = gui.SomeF(0)
	cfg.Color = gui.ColorTransparent
	cfg.ColorHeader = gui.ColorTransparent
	cfg.ColorContent = gui.ColorTransparent
	cfg.Radius = uiRadius
	cfg.RadiusTab = uiRadius
	cfg.PaddingTab = uiTabPadding
	cfg.PaddingContent = uiContentPad
	cfg.SpacingHeader = uiGapSmall
	return cfg
}

func requestAreaView(w *gui.Window) gui.View {
	app := gui.State[types.App](w)

	if app.RequestDock == nil || len(app.OpenTabs) == 0 {
		return emptyRequestAreaView()
	}

	panels := make([]gui.DockPanelDef, len(app.OpenTabs))
	for i, t := range app.OpenTabs {
		panels[i] = gui.DockPanelDef{
			ID:       t.ID,
			Label:    t.Method + "  " + t.Title,
			Closable: true,
			Content:  []gui.View{requestTabView(t)},
		}
	}

	return gui.DockLayout(gui.DockLayoutCfg{
		ID:           "request-tabs-dock",
		Sizing:       gui.FillFill,
		Root:         app.RequestDock,
		Panels:       panels,
		ColorContent: uiSurface,
		OnLayoutChange: func(root *gui.DockNode, ctx gui.EventCtx) {
			gui.State[types.App](ctx.Window).RequestDock = root
		},
		OnPanelSelect: func(groupID, panelID string, ctx gui.EventCtx) {
			a := gui.State[types.App](ctx.Window)
			a.RequestDock = gui.DockTreeSelectPanel(a.RequestDock, groupID, panelID)
			a.ActiveTabID = panelID
		},
		OnPanelClose: func(panelID string, ctx gui.EventCtx) {
			closeRequestTab(ctx.Window, panelID)
		},
	})
}

func emptyRequestAreaView() gui.View {
	return gui.Column(gui.ContainerCfg{
		Sizing:  gui.FillFill,
		Padding: uiPanelPadding,
		Color:   uiSurface,
		Content: []gui.View{
			gui.Label("Pilih request di sidebar untuk membukanya di sini", gui.TextStyle{Size: 13}),
		},
	})
}

func requestTabView(tab types.RequestTab) gui.View {
	return gui.Column(gui.ContainerCfg{
		Sizing:  gui.FillFill,
		Padding: uiPanelPadding,
		Spacing: uiGap,
		Color:   uiSurface,
		Content: []gui.View{
			urlBarView(tab),
			sectionsView(tab),
		},
	})
}

func urlBarView(tab types.RequestTab) gui.View {
	id := tab.ID
	send := func(ctx gui.EventCtx) {}

	return controlRow(
		gui.Select(gui.SelectCfg{
			ID:       "method-" + id,
			Options:  httpMethods,
			Selected: []string{tab.Method},
			Sizing:   gui.FitFill,
			MinWidth: 108,
			MaxWidth: 108,
			Padding:  uiFieldPadding,
			Radius:   uiRadius,
			OnSelect: func(sel []string, ctx gui.EventCtx) {
				if len(sel) == 0 {
					return
				}
				editTab(ctx.Window, id, func(t *types.RequestTab) { t.Method = sel[0] })
			},
		}),
		inputField(gui.InputCfg{
			ID:          "url-" + id,
			Text:        tab.URL,
			Placeholder: "https://api.example.com/users",
			OnTextChanged: func(s string, ctx gui.EventCtx) {
				editTab(ctx.Window, id, func(t *types.RequestTab) {
					t.URL = s
					t.Params = paramsFromURL(s)
				})
			},
			OnEnter: send,
		}),
		actionButton("send-btn-"+id, "Send", send),
	)
}

func sectionsView(tab types.RequestTab) gui.View {
	id := tab.ID
	selected := tab.Section
	if selected == "" {
		selected = "params"
	}

	return gui.TabControl(styledTabs(gui.TabControlCfg{
		ID:       "sections-" + id,
		Selected: selected,
		OnSelect: func(sec string, ctx gui.EventCtx) {
			updateTab(ctx.Window, id, func(t *types.RequestTab) { t.Section = sec })
		},
		Items: []gui.TabItemCfg{
			gui.NewTabItem("params", "Params", []gui.View{
				kvEditorView("params-"+id, tab.Params, func(w *gui.Window, fn func(*[]types.KeyValue)) {
					editTab(w, id, func(t *types.RequestTab) {
						fn(&t.Params)
						t.URL = urlWithParams(t.URL, t.Params)
					})
				}),
			}),
			gui.NewTabItem("headers", "Headers", []gui.View{
				kvEditorView("headers-"+id, tab.Headers, func(w *gui.Window, fn func(*[]types.KeyValue)) {
					editTab(w, id, func(t *types.RequestTab) { fn(&t.Headers) })
				}),
			}),
			gui.NewTabItem("body", "Body", []gui.View{bodyView(tab)}),
			gui.NewTabItem("auth", "Auth", []gui.View{authView(tab)}),
		},
	}))
}

func kvEditorView(
	prefix string,
	rows []types.KeyValue,
	edit func(w *gui.Window, fn func(*[]types.KeyValue)),
) gui.View {
	items := make([]gui.View, 0, len(rows)+1)

	for i, r := range rows {
		i := i
		items = append(items, controlRow(
			inputField(gui.InputCfg{
				ID:          fmt.Sprintf("%s-k-%d", prefix, i),
				Text:        r.Key,
				Placeholder: "Key",
				OnTextChanged: func(s string, ctx gui.EventCtx) {
					edit(ctx.Window, func(rs *[]types.KeyValue) {
						if i < len(*rs) {
							(*rs)[i].Key = s
						}
					})
				},
			}),
			inputField(gui.InputCfg{
				ID:          fmt.Sprintf("%s-v-%d", prefix, i),
				Text:        r.Value,
				Placeholder: "Value",
				OnTextChanged: func(s string, ctx gui.EventCtx) {
					edit(ctx.Window, func(rs *[]types.KeyValue) {
						if i < len(*rs) {
							(*rs)[i].Value = s
						}
					})
				},
			}),
			actionButton(fmt.Sprintf("%s-del-%d", prefix, i), "×", func(ctx gui.EventCtx) {
				edit(ctx.Window, func(rs *[]types.KeyValue) {
					if i < len(*rs) {
						*rs = append((*rs)[:i], (*rs)[i+1:]...)
					}
				})
			}),
		))
	}

	items = append(items, gui.Row(gui.ContainerCfg{
		Sizing:  gui.FillFixed,
		Height:  uiControlHeight,
		Padding: gui.NoPadding,
		Content: []gui.View{
			actionButton(prefix+"-add", "+ Add", func(ctx gui.EventCtx) {
				edit(ctx.Window, func(rs *[]types.KeyValue) {
					*rs = append(*rs, types.KeyValue{})
				})
			}),
		},
	}))

	return gui.Column(gui.ContainerCfg{
		ID:         prefix + "-list",
		Sizing:     gui.FillFill,
		Padding:    gui.NoPadding,
		Spacing:    uiGapSmall,
		Scrollable: true,
		Content:    items,
	})
}

func bodyView(tab types.RequestTab) gui.View {
	id := tab.ID
	return gui.Input(gui.InputCfg{
		ID:          "body-" + id,
		Text:        tab.Body,
		Placeholder: `{"key": "value"}`,
		Mode:        gui.InputMultiline,
		Scrollable:  true,
		Sizing:      gui.FillFill,
		Padding:     uiBodyPadding,
		Radius:      uiRadius,
		OnTextChanged: func(s string, ctx gui.EventCtx) {
			editTab(ctx.Window, id, func(t *types.RequestTab) { t.Body = s })
		},
	})
}

func authView(tab types.RequestTab) gui.View {
	id := tab.ID
	authType := tab.AuthType
	if authType == "" {
		authType = "none"
	}

	return gui.TabControl(styledTabs(gui.TabControlCfg{
		ID:       "auth-tabs-" + id,
		Selected: authType,
		OnSelect: func(sel string, ctx gui.EventCtx) {
			editTab(ctx.Window, id, func(t *types.RequestTab) { t.AuthType = sel })
		},
		Items: []gui.TabItemCfg{
			gui.NewTabItem("none", "None", []gui.View{
				gui.Label("Request ini tidak memakai autentikasi", gui.TextStyle{Size: 13}),
			}),
			gui.NewTabItem("bearer", "Bearer Token", []gui.View{bearerForm(tab)}),
			gui.NewTabItem("basic", "Basic Auth", []gui.View{basicForm(tab)}),
		},
	}))
}

func formColumn(content ...gui.View) gui.View {
	return gui.Column(gui.ContainerCfg{
		Sizing:  gui.FillFit,
		Padding: gui.NoPadding,
		Spacing: uiGap,
		Content: content,
	})
}

func bearerForm(tab types.RequestTab) gui.View {
	id := tab.ID
	return formColumn(
		labeledField("Token", inputField(gui.InputCfg{
			ID:          "auth-token-" + id,
			Text:        tab.AuthToken,
			Placeholder: "eyJhbGciOi...",
			OnTextChanged: func(s string, ctx gui.EventCtx) {
				editTab(ctx.Window, id, func(t *types.RequestTab) { t.AuthToken = s })
			},
		})),
	)
}

func basicForm(tab types.RequestTab) gui.View {
	id := tab.ID
	return formColumn(
		labeledField("Username", inputField(gui.InputCfg{
			ID:   "auth-user-" + id,
			Text: tab.AuthUser,
			OnTextChanged: func(s string, ctx gui.EventCtx) {
				editTab(ctx.Window, id, func(t *types.RequestTab) { t.AuthUser = s })
			},
		})),
		labeledField("Password", inputField(gui.InputCfg{
			ID:         "auth-pass-" + id,
			Text:       tab.AuthPass,
			IsPassword: true,
			OnTextChanged: func(s string, ctx gui.EventCtx) {
				editTab(ctx.Window, id, func(t *types.RequestTab) { t.AuthPass = s })
			},
		})),
	)
}
