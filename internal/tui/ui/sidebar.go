package ui

import (
	"os"
	"path/filepath"
	services "restclient/internal/services/fs"
	"restclient/internal/tui/styles"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
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

func buildTree(items []FileItem, flat []*FileItem, selected int, width int) *tree.Tree {
	t := tree.New().
		Enumerator(func(children tree.Children, index int) string {
			return ""
		}).
		Indenter(func(children tree.Children, index int) string {
			return ""
		})

	for i := range items {
		// cek if the item selected
		isSelected := len(flat) > selected && flat[selected] == &items[i]

		name := strings.TrimSuffix(items[i].name, ".json")
		if isSelected {
			name = lipgloss.NewStyle().
				Width(width - 4).
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230")).
				Render(name)
		}

		if items[i].isDir && items[i].expanded && len(items[i].children) > 0 {
			node := buildTree(items[i].children, flat, selected, width).
				Root(name).
				Enumerator(tree.DefaultEnumerator).
				Indenter(func(children tree.Children, index int) string {
					return "  "
				})
			t.Child(node)
		} else {
			t.Child(name)
		}
	}
	return t
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
		case "j", "down":
			if m.selected < len(flat)-1 {
				m.selected++
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
		case "enter", " ":
			if m.selected < len(flat) {
				item := flat[m.selected]
				if item.isDir {
					item.expanded = !item.expanded // toggle
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
		t := buildTree(m.folders, flat, m.selected, m.width)
		content = t.String()
	}

	borderColor := styles.BorderNormal

	if m.Active {
		borderColor = styles.BorderActive
	}

	return lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Render(content)
}
