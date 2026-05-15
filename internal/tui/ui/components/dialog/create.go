package dialog

import (
	"restclient/internal/tui/styles"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

func theme(isDark bool) *huh.Styles {
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

func (d *DialogModel) CreateForm() *huh.Form {
	myTheme := huh.ThemeFunc(theme)

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
				Value(&formName),
		),
	).WithTheme(myTheme).WithWidth(44)
}

func (d *DialogModel) RenameForm() *huh.Form {
	myTheme := huh.ThemeFunc(theme)

	return huh.NewForm(
		huh.NewGroup(
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
				Value(&formName),
		),
	).WithTheme(myTheme).WithWidth(44)
}
