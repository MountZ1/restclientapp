package components

import (
	"image/color"
	"os"
	"path/filepath"
	services "restclient/internal/services/fs"
	"restclient/internal/tui/constant"
	"restclient/internal/tui/styles"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type FileItem struct {
	name     string
	children []FileItem
	isDir    bool
	expanded bool
	path     string
}
type Collection []FileItem

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
			path:     filepath.Join(path, e.Name()),
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

func LoadCollection() tea.Cmd {
	return func() tea.Msg {
		path, _, err := services.Initialize()
		if err != nil {
			return err
		}
		items, err := readDir(path)
		if err != nil {
			return err
		}
		return Collection(items)
	}
}

func BuildLines(items []FileItem, flat []*FileItem, selected int, width int, active bool, indent int) []string {
	var lines []string

	for i := range items {
		isSelected := len(flat) > selected && flat[selected] == &items[i]

		var prefix string
		var methodColor color.Color
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
			var bg, fg color.Color
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
				Width(width-4).
				Render(content)

		} else {
			if items[i].isDir {
				line = "  " + indentStr + prefix + displayName
			} else {
				tag := lipgloss.NewStyle().
					Foreground(methodColor).
					Render("[" + method + "]")
				line = "  " + indentStr + tag + " " + displayName
			}
		}

		lines = append(lines, line)

		if items[i].isDir && items[i].expanded && len(items[i].children) > 0 {
			childLines := BuildLines(items[i].children, flat, selected, width, active, indent+2)
			lines = append(lines, childLines...)
		}
	}

	return lines
}

func FlattenItems(items []FileItem) []*FileItem {
	var flat []*FileItem

	for i := range items {
		flat = append(flat, &items[i])
		if items[i].isDir && items[i].expanded {
			flat = append(flat, FlattenItems(items[i].children)...)
		}
	}

	return flat
}

func (f *FileItem) IsDir() bool {
	return f.isDir
}

func (f *FileItem) ToggleExpanded() {
	f.expanded = !f.expanded
}

func (f *FileItem) FilePath() string {
	return f.path
}
