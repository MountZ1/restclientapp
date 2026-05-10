package components

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
	Folder  string
}

type DialogSubmitMsg struct {
	Location string
	Type     string
	Name     string
}

func NewDialog(title, message string) DialogModel {
	formTipe = ""
	formName = ""
	d := DialogModel{
		Title:   title,
		Message: message,
	}
	d.Form = d.buildForm()
	return d
}

func (d *DialogModel) buildForm() *huh.Form {
	myTheme := func(isDark bool) *huh.Styles {
		theme := huh.ThemeBase(isDark)
		theme.Focused.Base = lipgloss.NewStyle().
			PaddingLeft(1).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(styles.BorderActive)
		theme.Focused.Title = lipgloss.NewStyle().
			Foreground(styles.TextNormal).
			Bold(true)
		theme.Focused.SelectSelector = lipgloss.NewStyle().
			Foreground(styles.BorderActive)
		theme.Focused.SelectedOption = lipgloss.NewStyle().
			Foreground(styles.SelectedFg)
		theme.Focused.UnselectedOption = lipgloss.NewStyle().
			Foreground(styles.TextSubtle)
		theme.Focused.TextInput.Cursor = lipgloss.NewStyle().
			Foreground(styles.BorderActive)
		theme.Focused.TextInput.Text = lipgloss.NewStyle().
			Foreground(styles.TextNormal)
		theme.Focused.TextInput.Placeholder = lipgloss.NewStyle().
			Foreground(styles.TextMuted)
		theme.Blurred.Title = lipgloss.NewStyle().
			Foreground(styles.TextSubtle)
		theme.Blurred.Base = lipgloss.NewStyle().
			PaddingLeft(1).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(styles.TextMuted)
		return theme
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Type").
				Options(
					huh.NewOption("Collection", "collection"),
					huh.NewOption("Request", "request"),
				).
				Value(&formTipe), // ← package-level var
			huh.NewInput().
				TitleFunc(func() string {
					switch formTipe {
					case "collection":
						return "Name collection"
					case "request":
						return "Name request"
					default:
						return "Name"
					}
				}, &formTipe).
				Placeholder("e.g. My Collection").
				Value(&formName), // ← package-level var
		),
	).WithTheme(huh.ThemeFunc(myTheme)).WithWidth(44)
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
