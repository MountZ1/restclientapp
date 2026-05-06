package components

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type Item struct {
	Title_ string
	Desc_  string
}

func (i Item) Title() string       { return i.Title_ }
func (i Item) Description() string { return i.Desc_ }
func (i Item) FilterValue() string { return i.Title_ }

type ListModel struct {
	list    list.Model
	Focused bool
}

func NewList(items []list.Item, width, height int) Component {
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.ShowDescription = false

	l := list.New(items, delegate, width, height)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	return &ListModel{list: l}
}

func (l *ListModel) Focus() {
	l.Focused = true
}

func (l *ListModel) Blur() {
	l.Focused = false
}

func (l *ListModel) IsFocus() bool {
	return l.Focused
}

func (l *ListModel) Update(msg tea.Msg) (Component, tea.Cmd) {
	if !l.Focused {
		return l, nil
	}
	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return l, cmd
}

func (l *ListModel) View() string {
	return l.list.View()
}

func (l *ListModel) Selected() list.Item {
	return l.list.SelectedItem()
}

func (l *ListModel) SetItems(items []list.Item) {
	l.list.SetItems(items)
}
