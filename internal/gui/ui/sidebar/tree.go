package sidebar

import (
	"fmt"
	"os"
	"path/filepath"
	services "restclient/internal/services/fs"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func tree(cnv fyne.Canvas) (*tappableTree, string) {
	data := map[string][]string{
		"": {},
	}
	openFolders := map[string]bool{}
	selectedUID := ""

	dataPath, files, err := services.Initialize()
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			data[""] = append(data[""], name)
			children, err := os.ReadDir(filepath.Join(dataPath, name))
			if err == nil {
				for _, child := range children {
					childID := name + "/" + child.Name()
					data[name] = append(data[name], childID)
				}
			}
		} else if strings.ToLower(filepath.Ext(name)) == ".json" {
			data[""] = append(data[""], name)
		}
	}

	t := newTappableTree()

	t.ChildUIDs = func(uid widget.TreeNodeID) []widget.TreeNodeID {
		return data[uid]
	}

	t.IsBranch = func(uid widget.TreeNodeID) bool {
		_, isFolder := data[uid]
		return isFolder
	}

	t.CreateNode = func(branch bool) fyne.CanvasObject {
		btn := widget.NewButton("...", nil)
		btn.Importance = widget.LowImportance

		paddedBtn := container.New(
			layout.NewCustomPaddedLayout(
				0,
				0,
				0,
				8),
			btn)

		return container.NewHBox(
			widget.NewIcon(theme.FileIcon()),
			widget.NewLabel("template"),
			layout.NewSpacer(),
			paddedBtn,
		)
	}

	t.UpdateNode = func(uid widget.TreeNodeID, branch bool, node fyne.CanvasObject) {
		box := node.(*fyne.Container)
		icon := box.Objects[0].(*widget.Icon)
		label := box.Objects[1].(*widget.Label)
		paddedBtn := box.Objects[3].(*fyne.Container)
		btn := paddedBtn.Objects[0].(*widget.Button)

		parts := strings.Split(uid, "/")
		displayName := parts[len(parts)-1]
		displayName = strings.TrimSuffix(displayName, ".json")

		label.SetText(displayName)

		btn.SetText("...")

		if uid == selectedUID {
			box.Objects = []fyne.CanvasObject{
				canvas.NewRectangle(theme.PrimaryColor()),
				icon,
				label,
			}
		}

		if branch {
			// need to check the icon
			btn.OnTapped = func() {
				showSidebarMenu(cnv, btn.Position(), uid, "folder")
			}

			if openFolders[uid] {
				icon.SetResource(theme.FolderOpenIcon())
			} else {
				icon.SetResource(theme.FolderIcon())
			}
		} else {
			btn.OnTapped = func() {
				showSidebarMenu(cnv, btn.Position(), uid, "file")
			}
			icon.SetResource(theme.FileIcon())
		}
	}

	t.OnSelected = func(uid widget.TreeNodeID) {
		_, isFolder := data[uid]

		if isFolder {
			if t.IsBranchOpen(uid) {
				t.CloseBranch(uid)
			} else {
				t.OpenBranch(uid)
			}
		} else {
			fmt.Println(uid)
		}

		t.Unselect(uid)
	}

	t.OnBranchOpened = func(uid widget.TreeNodeID) {
		openFolders[uid] = true
		t.Refresh()
	}

	t.OnBranchClosed = func(uid widget.TreeNodeID) {
		openFolders[uid] = false
		t.Refresh()
	}

	return t, selectedUID
}
