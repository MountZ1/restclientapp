package components

import "github.com/go-gui-org/go-gui/gui"

type DialogProp struct {
	Title      string
	Width      float32
	Padding    gui.Padding
	FocusID    string
	OnOkYes    func(w *gui.Window)
	OnCancelNo func(w *gui.Window)
}

type DialogConfirmProp struct {
	DialogProp
	Body string
}

type DialogPromptProp struct {
	DialogProp
	Body    string
	OnReply func(reply string, w *gui.Window)
}

type DialogCustomProp struct {
	DialogProp
	CustomContent []gui.View
	SoundDisabled bool
}

func ShowDialogConfirm(w *gui.Window, prop DialogConfirmProp) {
	w.Dialog(gui.DialogCfg{
		Title:      prop.Title,
		Body:       prop.Body,
		DialogType: gui.DialogConfirm,
		OnOkYes:    prop.OnOkYes,
		OnCancelNo: prop.OnCancelNo,
	})
}

func ShowDialogPrompt(w *gui.Window, prop DialogPromptProp) {
	w.Dialog(gui.DialogCfg{
		Title:      prop.Title,
		Body:       prop.Body,
		DialogType: gui.DialogPrompt,
		OnReply:    prop.OnReply,
		OnCancelNo: prop.OnCancelNo,
	})
}

func ShowDialogCustom(w *gui.Window, prop DialogCustomProp) {
	w.Dialog(gui.DialogCfg{
		Title:         prop.Title,
		Width:         prop.Width,
		Padding:       prop.Padding,
		FocusID:       prop.FocusID,
		DialogType:    gui.DialogCustom,
		CustomContent: prop.CustomContent,
		OnOkYes:       prop.OnOkYes,
		OnCancelNo:    prop.OnCancelNo,
		SoundDisabled: prop.SoundDisabled,
	})
}
