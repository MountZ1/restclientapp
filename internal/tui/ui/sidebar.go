package ui

import (
	"path/filepath"
	"restclient/internal/tui/helper"
	"restclient/internal/tui/styles"
	"restclient/internal/tui/ui/components"
	"restclient/internal/tui/ui/components/dialog"
	"strings"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	lv2 "charm.land/lipgloss/v2"
)

type SidebarModel struct {
	Active     bool
	folders    []components.FileItem
	spinner    spinner.Model
	selected   int
	loading    bool
	Width      int
	height     int
	button     components.ButtonModel
	choiceType string
}

type OpenDialogMsg struct {
	Title    string
	Message  string
	Location string
	Content  any
}

func NewSidebar(width, height int) SidebarModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return SidebarModel{
		spinner: s,
		loading: true,
		Width:   width,
		height:  height,
		button:  components.NewButton("Create Request or Collection", width-4, 1),
	}
}

func (m SidebarModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		components.LoadCollection(),
	)
}

func (m SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case components.Collection:
		m.loading = false
		m.folders = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.Width = msg.Width / 4
		m.height = msg.Height
		m.button = components.NewButton(m.button.Text, m.Width-4, 1)
		return m, nil

	case tea.MouseClickMsg:
		mouse := msg.Mouse()
		btnWidth := lv2.Width(m.button.View())

		if mouse.Y == m.height-3 && mouse.X >= 2 && mouse.X <= 2+btnWidth {
			return m, func() tea.Msg {
				d := dialog.DialogModel{}
				form := d.CreateForm()
				return OpenDialogMsg{
					Title:    "Create Request or Collection",
					Message:  "",
					Location: "",
					Content:  form,
				}
			}
		}

	case tea.KeyMsg:
		if !m.Active {
			return m, nil
		}
		flat := components.FlattenItems(m.folders)
		switch msg.String() {
		case "down":
			if m.selected < len(flat)-1 {
				m.selected++
			}
		case "up":
			if m.selected > 0 {
				m.selected--
			}
		case "enter":
			if m.selected < len(flat) {
				item := flat[m.selected]
				if item.IsDir() {
					item.ToggleExpanded()
				}
			}
		case "a":
			location := ""
			if m.selected < len(flat) {
				item := flat[m.selected]
				if item.IsDir() {
					location = item.FilePath()
				} else {
					// file ada di dalam folder, ambil parent path
					location = filepath.Dir(item.FilePath())
				}
			}
			return m, func() tea.Msg {
				d := dialog.DialogModel{}
				form := d.CreateForm()
				return OpenDialogMsg{
					Title:    "Create Request or Collection",
					Message:  "",
					Location: location,
					Content:  form,
				}
			}
		}

		return m, nil
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m SidebarModel) View() string {
	flat := components.FlattenItems(m.folders)
	lines := components.BuildLines(m.folders, flat, m.selected, m.Width, m.Active, 0)
	listContent := strings.Join(lines, "\n")

	borderColor := styles.BorderNormal
	if m.Active {
		borderColor = styles.BorderActive
	}

	btnLayer := lv2.NewLayer(m.button.View()).X(2).Y(m.height - 3)
	helpStyle := lipgloss.NewStyle().Foreground(styles.TextMuted)
	helpLayer := lv2.NewLayer(helpStyle.Render("H Help")).X(2).Y(m.height - 2)

	return helper.RenderWithTitle(
		listContent,
		"Collections",
		m.Width,
		m.height,
		borderColor,
		btnLayer,
		helpLayer,
	)
}
