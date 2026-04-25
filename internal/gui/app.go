package gui

import (
	"restclient/internal/gui/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func Run() {
	application := app.New()

	application.Settings().SetTheme(&myTheme{
		font: fyne.NewStaticResource("NerdFont", nerdFont),
	})
	/*application.Settings().SetTheme(&myTheme{
		font: fyne.NewStaticResource("NerdFont", fontData),
	})*/

	w := application.NewWindow("New App")
	w.Resize(fyne.NewSize(1280, 720))
	w.CenterOnScreen()
	w.SetContent(ui.Windows(w))
	w.ShowAndRun()
}
