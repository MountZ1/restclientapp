package ui

import (
	"path/filepath"
	services "restclient/internal/services/fs"
	"restclient/internal/services/logger"
	"restclient/internal/tui/helper"
	"restclient/internal/tui/styles"
	"restclient/internal/tui/ui/components"
	"restclient/internal/tui/ui/components/dialog"
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	lv2 "charm.land/lipgloss/v2"
)

type SidebarModel struct {
	Active        bool
	folders       []components.FileItem
	spinner       spinner.Model
	selected      int
	loading       bool
	Width         int
	height        int
	refreshButton components.ButtonModel
	button        components.ButtonModel
	vp            viewport.Model
	choiceType    string
}

type OpenDialogMsg struct {
	Title      string
	Message    string
	Location   string
	Content    any
	DialogType string
}

// footerRows is how many rows at the bottom of the panel's interior are
// reserved for the create-button and the help line, and thus excluded
// from the scrollable viewport.
const footerRows = 2

// contentWidth returns the interior width available for this panel's
// content, given the panel's TOTAL outer width (the same contract
// RenderWithTitle now uses). Keeping this in one place means there is
// no hardcoded magic number to keep in sync elsewhere in this file.
func contentWidth(totalWidth int) int {
	w := totalWidth - helper.PanelHorizontalFrame()
	if w < 1 {
		w = 1
	}
	return w
}

func contentHeight(totalHeight int) int {
	h := totalHeight - helper.PanelVerticalFrame()
	if h < 1 {
		h = 1
	}
	return h
}

func NewSidebar(width, height int) SidebarModel {
	innerW := contentWidth(width)
	innerH := contentHeight(height)

	vp := viewport.New()
	vp.SetWidth(innerW)
	vp.SetHeight(innerH - footerRows)

	s := spinner.New()
	s.Spinner = spinner.Dot

	return SidebarModel{
		spinner: s,
		loading: true,
		Width:   width,
		height:  height,
		button:  components.NewButton("Create Request or Collection", innerW, 1),
		vp:      vp,
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
		flat := components.FlattenItems(m.folders)
		lines := components.BuildLines(m.folders, flat, m.selected, m.Width, m.Active, 0)
		m.vp.SetContent(strings.Join(lines, "\n"))
		return m, nil

	case tea.WindowSizeMsg:
		m.Width = msg.Width / 4
		m.height = msg.Height

		innerW := contentWidth(m.Width)
		innerH := contentHeight(m.height)

		m.button = components.NewButton(m.button.Text, innerW, 1)
		m.vp.SetWidth(innerW)
		m.vp.SetHeight(innerH - footerRows)
		return m, nil

	case tea.MouseClickMsg:
		mouse := msg.Mouse()
		btnWidth := lv2.Width(m.button.View())

		if mouse.Y == m.height-3 && mouse.X >= 2 && mouse.X <= 2+btnWidth {
			return m, func() tea.Msg {
				dialog.ResetForm("")
				d := dialog.DialogModel{}
				form := d.CreateForm()
				return OpenDialogMsg{
					Title:      "Create Request or Collection",
					Message:    "",
					Location:   "",
					Content:    form,
					DialogType: "create",
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
			lines := components.BuildLines(m.folders, flat, m.selected, m.Width, m.Active, 0)
			m.vp.SetContent(strings.Join(lines, "\n"))
			m.vp.SetYOffset(m.selected - m.vp.Height()/2)

		case "up":
			if m.selected > 0 {
				m.selected--
			}
			lines := components.BuildLines(m.folders, flat, m.selected, m.Width, m.Active, 0)
			m.vp.SetContent(strings.Join(lines, "\n"))
			m.vp.SetYOffset(m.selected - m.vp.Height()/2)
		case "enter":
			if m.selected < len(flat) {
				item := flat[m.selected]
				if item.IsDir() {
					item.ToggleExpanded()
				}
			}
			flat = components.FlattenItems(m.folders)
			lines := components.BuildLines(m.folders, flat, m.selected, m.Width, m.Active, 0)
			m.vp.SetContent(strings.Join(lines, "\n"))
		case "a":
			location := ""
			if m.selected < len(flat) {
				item := flat[m.selected]
				if item.IsDir() {
					location = item.FilePath()
				} else {
					location = filepath.Dir(item.FilePath())
				}
			}
			return m, func() tea.Msg {
				dialog.ResetForm("")
				d := dialog.DialogModel{}
				form := d.CreateForm()
				return OpenDialogMsg{
					Title:      "Create Request or Collection",
					Message:    "",
					Location:   location,
					Content:    form,
					DialogType: "create",
				}
			}
		case "r":
			var item *components.FileItem
			if m.selected < len(flat) {
				item = flat[m.selected]
			}

			return m, func() tea.Msg {
				tipe := "request"
				if item.IsDir() {
					tipe = "collection"
				}
				dialog.ResetForm(tipe)
				name := getCollectionOrRequestName(item.FilePath(), item.IsDir())
				d := dialog.DialogModel{}
				form := d.RenameForm()
				return OpenDialogMsg{
					Title:      "Rename " + name,
					Message:    "",
					Location:   item.FilePath(),
					Content:    form,
					DialogType: "rename",
				}
			}
		case "d":
			var item *components.FileItem
			if m.selected < len(flat) {
				item = flat[m.selected]
			}

			return m, func() tea.Msg {
				tipe := "request"
				if item.IsDir() {
					tipe = "collection"
				}
				dialog.ResetForm(tipe)
				name := getCollectionOrRequestName(item.FilePath(), item.IsDir())
				d := dialog.DialogModel{}
				form := d.DestroyForm()
				return OpenDialogMsg{
					Title:      "Destroy " + name,
					Message:    "",
					Location:   item.FilePath(),
					Content:    form,
					DialogType: "destroy",
				}
			}

		case "c":
			var item *components.FileItem
			if m.selected < len(flat) {
				item = flat[m.selected]
			}
			err := services.DuplicateRequestOrCollection(item.FilePath(), item.IsDir())
			if err != nil {
				logger.Error("Failed to duplicate this item : %s", err)
				return m, func() tea.Msg {
					return dialog.ShowErrorMsg{Message: "Err " + err.Error()}
				}
			}

			return m, components.LoadCollection()
		}

		return m, nil
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	if m.Active {
		m.vp, cmd = m.vp.Update(msg)
	}

	return m, cmd
}

func (m SidebarModel) View() string {
	borderColor := styles.BorderNormal
	if m.Active {
		borderColor = styles.BorderActive
	}

	btnLayer := lv2.NewLayer(m.button.View()).X(2).Y(m.height - 3)
	helpStyle := lipgloss.NewStyle().Foreground(styles.TextMuted)
	helpLayer := lv2.NewLayer(helpStyle.Render("H Help")).X(2).Y(m.height - 2)

	return helper.RenderWithTitle(
		m.vp.View(),
		"Collections",
		m.Width,
		m.height,
		borderColor,
		btnLayer,
		helpLayer,
	)
}

func getCollectionOrRequestName(path string, isDir bool) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	nameOnlys := strings.TrimSuffix(base, ext)
	_, nameOnly, _ := strings.Cut(nameOnlys, "-")

	if isDir {
		return "Collection " + nameOnly
	} else {
		return "Request " + nameOnly
	}
}
