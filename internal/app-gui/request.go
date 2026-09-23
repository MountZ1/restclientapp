package appgui

import "github.com/go-gui-org/go-gui/gui"

func requestView(w *gui.Window) gui.View {
	t := gui.CurrentTheme()
	return gui.Column(gui.ContainerCfg{
		Sizing:      gui.FillFill,
		Padding:     gui.PadAll(8),
		ColorBorder: gui.RGB(90, 90, 96),
		SizeBorder:  gui.SomeF(1),
		Content: []gui.View{
			gui.Label("Request panel", t.B3),
			gui.TextButton("send-btn", "Send", func(ctx gui.EventCtx) {
				gui.State[App](ctx.Window).Clicks++
			}),
		},
	})
}
