package gui

import (
	"github.com/patrikaleksandryan/coloride/pkg/color"
)

type Label struct {
	FrameImpl
	caption   string
	textAlign int
}

func NewLabel(caption string, x, y, w, h int) *Label {
	l := &Label{
		caption:   caption,
		textAlign: AlignLeft,
	}
	InitFrame(&l.FrameImpl, x, y, w, h)
	l.color = color.Black
	return l
}

func (l *Label) SetAlign(textAlign int) {
	l.textAlign = textAlign
}

func (l *Label) Align() int {
	return l.textAlign
}

func (l *Label) SetCaption(caption string) {
	l.caption = caption
}

func (l *Label) Caption() string {
	return l.caption
}

func (l *Label) Render(x, y int) {
	w, _ := l.Size()
	PrintAlign(l.caption, x, y, w, l.color, color.Transparent, l.textAlign)
	l.RenderChildren(x, y)
}
