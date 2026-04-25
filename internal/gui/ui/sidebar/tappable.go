package sidebar

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type tappableTree struct {
	widget.Tree
	onTappedSecondary func(*fyne.PointEvent)
	hoverUID          string
}

func newTappableTree() *tappableTree {
	t := &tappableTree{}
	t.ExtendBaseWidget(t)
	return t
}

func (t *tappableTree) TappedSecondary(ev *fyne.PointEvent) {
	if t.onTappedSecondary != nil {
		t.onTappedSecondary(ev)
	}
}

// tappable.go
type hoverableNode struct {
	widget.BaseWidget
	content    fyne.CanvasObject
	onMouseIn  func()
	onMouseOut func()
}

func newHoverableNode(content fyne.CanvasObject) *hoverableNode {
	h := &hoverableNode{content: content}
	h.ExtendBaseWidget(h)
	return h
}

func (h *hoverableNode) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(h.content)
}

func (h *hoverableNode) MouseIn(*desktop.MouseEvent) {
	if h.onMouseIn != nil {
		h.onMouseIn()
	}
}

func (h *hoverableNode) MouseMoved(*desktop.MouseEvent) {}

func (h *hoverableNode) MouseOut() {
	if h.onMouseOut != nil {
		h.onMouseOut()
	}
}
