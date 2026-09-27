package types

import (
	"os"

	"github.com/go-gui-org/go-gui/gui"
)

type App struct {
	Clicks int

	SavedRequests []os.DirEntry
	RequestLookup map[string]RequestTab
	PathByID      map[string]string
	ExpandedNodes map[string]bool
	SelectedID    string

	OpenTabs             []RequestTab
	ActiveTabID          string
	RequestDock          *gui.DockNode
	RequestResponseRatio float32

	ContextMenuOpen     bool
	ContextMenuTargetID string
	ContextMenuIsDir    bool
	ContextMenuX        float32
	ContextMenuY        float32
}
