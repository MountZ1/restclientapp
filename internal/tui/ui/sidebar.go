package ui

import (
	"os"
	"path/filepath"
	services "restclient/internal/services/fs"
	"restclient/internal/tui/constant"
	"restclient/internal/tui/helper"
	"restclient/internal/tui/styles"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SidebarModel struct {
	Active   bool
	folders  []FileItem
	spinner  spinner.Model
	selected int
	loading  bool
	lastKey  string
	width    int
	height   int
}

type FileItem struct {
	name     string
	children []FileItem
	isDir    bool
	expanded bool
}

type OpenDialogMsg struct {
	Title   string
	Message string
}

type collection []FileItem

func readDir(path string) ([]FileItem, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var items []FileItem
	for _, e := range entries {
		item := FileItem{
			name:     e.Name(),
			isDir:    e.IsDir(),
			expanded: false,
		}
		if e.IsDir() {
			children, err := readDir(filepath.Join(path, e.Name()))
			if err == nil {
				item.children = children
			}
		} else if strings.ToLower(filepath.Ext(e.Name())) != ".json" {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func loadCollection() tea.Cmd {
	return func() tea.Msg {
		path, _, err := services.Initialize()
		if err != nil {
			return err
		}
		items, err := readDir(path)
		if err != nil {
			return err
		}
		return collection(items)
	}
}

func buildLines(items []FileItem, flat []*FileItem, selected int, width int, active bool, indent int) []string {
	var lines []string

	for i := range items {
		isSelected := len(flat) > selected && flat[selected] == &items[i]

		var prefix string
		var methodColor lipgloss.Color
		var method string
		displayName := strings.TrimSuffix(items[i].name, ".json")

		if items[i].isDir {
			if items[i].expanded {
				prefix = constant.IconFolderOpen + " "
			} else {
				prefix = constant.IconFolder + " "
			}
		} else {
			var found bool
			method, _, found = strings.Cut(items[i].name, "-")
			if !found {
				method = "???"
			}
			method = strings.ToUpper(method)
			methodColor = styles.ColorUnknown
			for _, r := range constant.RequestCollection {
				if r.Name == method {
					methodColor = r.Color
					break
				}
			}
			_, after, found := strings.Cut(displayName, "-")
			if found {
				displayName = after
			}
		}

		var line string
		indentStr := strings.Repeat(" ", indent)

		if isSelected {
			var bg, fg lipgloss.Color
			if active {
				bg = lipgloss.Color("62")
				fg = lipgloss.Color("230")
			} else {
				bg = lipgloss.Color("240")
				fg = lipgloss.Color("250")
			}

			var content string
			if items[i].isDir {
				content = indentStr + prefix + displayName
			} else {
				content = indentStr + "[" + method + "]" + " " + displayName
			}

			line = " " + lipgloss.NewStyle().
				Background(bg).
				Foreground(fg).
				PaddingLeft(1).
				PaddingRight(1).
				Width(width-4). // kurangi margin kiri + padding
				Render(content)

		} else {
			if items[i].isDir {
				line = "  " + indentStr + prefix + displayName // sejajar dengan selected
			} else {
				tag := lipgloss.NewStyle().
					Foreground(methodColor).
					Render("[" + method + "]")
				line = "  " + indentStr + tag + " " + displayName
			}
		}

		lines = append(lines, line)

		if items[i].isDir && items[i].expanded && len(items[i].children) > 0 {
			childLines := buildLines(items[i].children, flat, selected, width, active, indent+2)
			lines = append(lines, childLines...)
		}
	}

	return lines
}

func flattenItems(items []FileItem) []*FileItem {
	var flat []*FileItem

	for i := range items {
		flat = append(flat, &items[i])
		if items[i].isDir && items[i].expanded {
			flat = append(flat, flattenItems(items[i].children)...)
		}
	}

	return flat
}

func NewSidebar(width, height int) SidebarModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return SidebarModel{
		spinner: s,
		loading: true,
		width:   width,
		height:  height,
	}
}

func (m SidebarModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadCollection(),
	)
}

func (m SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case collection:
		m.loading = false
		m.folders = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width / 4
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if !m.Active {
			return m, nil
		}
		flat := flattenItems(m.folders)
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
				if item.isDir {
					item.expanded = !item.expanded
				}
			}
		case "a":
			return m, func() tea.Msg {
				return OpenDialogMsg{
					Title:   "New Dialog",
					Message: "this is message",
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
	content := m.spinner.View() + " Loading..."
	if !m.loading {
		flat := flattenItems(m.folders)
		lines := buildLines(m.folders, flat, m.selected, m.width, m.Active, 0)
		content = strings.Join(lines, "\n")
	}

	borderColor := styles.BorderNormal
	if m.Active {
		borderColor = styles.BorderActive
	}

	return helper.RenderWithTitle(content, "[ Collections ]", m.width, m.height-2, borderColor)
}
