package dialog

import (
	"restclient/internal/tui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

// package-level vars sebagai binding target
var (
	formTipe string
	formName string
)

type DialogModel struct {
	Title   string
	Message string
	Active  bool
	Form    *huh.Form
	content string
	Folder  string
	DType   string
}

type DialogSubmitMsg struct {
	FormType string
	Location string
	Type     string
	Name     string
}

func NewDialog(title, message string, folder string, content any, Dtype string) DialogModel {
	formTipe = ""
	formName = ""

	d := DialogModel{
		Title:   title,
		Message: message,
		Folder:  folder,
		DType:   Dtype,
	}

	switch v := content.(type) {
	case *huh.Form:
		d.Form = v
	case string:
		d.content = v
	case nil:
	}

	return d
}

func (d DialogModel) Init() tea.Cmd {
	return d.Form.Init()
}

func (d DialogModel) Update(msg tea.Msg) (DialogModel, tea.Cmd) {
	var cmd tea.Cmd
	f, c := d.Form.Update(msg)
	if form, ok := f.(*huh.Form); ok {
		d.Form = form
		cmd = c
	}

	if d.Form.State == huh.StateCompleted {
		tipe := formTipe
		name := formName
		return d, func() tea.Msg {
			return DialogSubmitMsg{
				Location: d.Folder,
				Type:     tipe,
				Name:     name,
				FormType: d.DType,
			}
		}
	}

	return d, cmd
}

func (d DialogModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.TextNormal).
		Bold(true).
		MarginBottom(1)
	helpStyle := lipgloss.NewStyle().
		Foreground(styles.TextMuted)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(d.Title),
		d.Form.View(),
		helpStyle.Render("esc close"),
	)

	return lipgloss.NewStyle().
		Width(48).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.BorderActive).
		Render(content)
}

func ResetForm(newTipe string) {
	formTipe = newTipe
	formName = ""
}
