package features

import (
	"restclient/internal/app-gui/components"
	services "restclient/internal/services/fs"
	"restclient/internal/types"

	"github.com/go-gui-org/go-gui/gui"
)

var (
	dialogBorderColor = gui.RGB(90, 90, 96)
	dialogActiveColor = gui.RGB(70, 120, 220)
	whiteText         = gui.TextStyle{Color: gui.White, Size: 16}
)

func init() {
	components.ShowAddDialogAtFn = ShowAddDialogAt
	components.ShowRenameDialogFn = ShowRenameDialog
}

func AddCollectionOrRequest(ctx gui.EventCtx) {
	ShowAddDialogAt(ctx.Window, services.Path())
}

func ShowAddDialogAt(w *gui.Window, location string) {
	c := &types.Create{Location: location}
	showAddDialog(w, c)
}

func showAddDialog(w *gui.Window, c *types.Create) {
	components.ShowDialogCustom(w, components.DialogCustomProp{
		DialogProp: components.DialogProp{
			Title:   "Add Collection or Request",
			Width:   520,
			Padding: gui.PadAll(20),
			FocusID: "name-add",
		},
		SoundDisabled: true,
		CustomContent: []gui.View{
			gui.Column(gui.ContainerCfg{
				Sizing:  gui.FillFit,
				HAlign:  gui.HAlignCenter,
				Spacing: gui.SomeF(18),
				Content: []gui.View{
					gui.Text(gui.TextCfg{
						Text:   "Add Collection or Request",
						Sizing: gui.FillFit,
						TextStyle: gui.TextStyle{
							Color: gui.White,
							Size:  20,
							Align: gui.TextAlignCenter,
						},
					}),
					gui.Row(gui.ContainerCfg{
						Sizing:      gui.FillFit,
						SizeBorder:  gui.SomeF(1),
						ColorBorder: dialogBorderColor,
						Radius:      gui.SomeF(8),
						Padding:     gui.PadAll(16),
						Spacing:     gui.SomeF(32),
						HAlign:      gui.HAlignCenter,
						VAlign:      gui.VAlignMiddle,
						Content: []gui.View{
							gui.Radio(gui.RadioCfg{
								ID:             "type-collection",
								Label:          "Collection",
								Selected:       c.Type == "collection",
								TextStyleLabel: whiteText,
								ColorSelect:    dialogActiveColor,
								OnClick: func(ctx gui.EventCtx) {
									c.Type = "collection"
									showAddDialog(ctx.Window, c)
									ctx.Consume()
								},
							}),
							gui.Radio(gui.RadioCfg{
								ID:             "type-request",
								Label:          "Request",
								Selected:       c.Type == "request",
								TextStyleLabel: whiteText,
								ColorSelect:    dialogActiveColor,
								OnClick: func(ctx gui.EventCtx) {
									c.Type = "request"
									showAddDialog(ctx.Window, c)
									ctx.Consume()
								},
							}),
						},
					}),
					gui.Input(gui.InputCfg{
						ID:          "name-add",
						Placeholder: "Name",
						Sizing:      gui.FillFit,
						Height:      4.5,
						Padding:     gui.PadAll(10),
						Radius:      gui.SomeF(8),
						TextStyle:   whiteText,
						Text:        c.Name,
						OnTextChanged: func(text string, ctx gui.EventCtx) {
							c.Name = text
							showAddDialog(ctx.Window, c)
						},
					}),
					gui.Row(gui.ContainerCfg{
						Sizing:  gui.FitFit,
						HAlign:  gui.HAlignCenter,
						Spacing: gui.SomeF(12),
						Content: []gui.View{
							gui.Button(gui.ButtonCfg{
								ID:     "cancel-add",
								Height: 4,
								Width:  160,
								Radius: gui.SomeF(8),
								Content: []gui.View{
									gui.Text(gui.TextCfg{Text: "Cancel", TextStyle: whiteText}),
								},
								OnClick: func(ctx gui.EventCtx) {
									ctx.Window.DialogDismiss()
									ctx.Consume()
								},
							}),
							gui.Button(gui.ButtonCfg{
								ID:     "confirm-add",
								Height: 4,
								Width:  160,
								Radius: gui.SomeF(8),
								Color:  dialogActiveColor,
								Content: []gui.View{
									gui.Text(gui.TextCfg{Text: "Create", TextStyle: whiteText}),
								},
								OnClick: func(ctx gui.EventCtx) {
									if c.Name == "" || c.Type == "" {
										return
									}
									if err := services.CreateRequestCollection(*c); err != nil {
										return
									}
									ctx.Window.DialogDismiss()
									ctx.Consume()
								},
							}),
						},
					}),
				},
			}),
		},
	})
}

func ShowRenameDialog(w *gui.Window, path string) {
	components.ShowDialogPrompt(w, components.DialogPromptProp{
		DialogProp: components.DialogProp{Title: "Rename"},
		Body:       "New name:",
		OnReply: func(reply string, w *gui.Window) {
			if reply == "" {
				return
			}
			services.RenameRequestOrCollection(path, reply)
		},
	})
}
