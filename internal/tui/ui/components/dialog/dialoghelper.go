package dialog

import (
	"restclient/internal/tui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type DialogHelperModel struct {
	Title  string
	width  int
	height int
}

type DialogHelperCloseMsg struct{}

type ShowErrorMsg struct {
	Message string
}

func NewDialogHelper(width int, height int, message string) DialogHelperModel {
	formTipe = ""
	formName = ""

	d := DialogHelperModel{
		width:  width,
		height: height,
		Title:  message,
	}

	return d
}

/*func (d DialogModel) Init() tea.Cmd {
	return d.Form.Init()
}*/

func (d DialogHelperModel) Update(msg tea.Msg) (DialogHelperModel, tea.Cmd) {
	var cmd tea.Cmd

	return d, cmd
}

func (d DialogHelperModel) View() string {
	maxWidth := 50
	title := d.Title
	if len(title) > maxWidth {
		title = title[:maxWidth-3] + "..."
	}
	return lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.BorderActive).
		Foreground(styles.TextNormal).
		Bold(true).
		Render(title)
}
