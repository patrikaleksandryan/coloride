package editor

import (
	"github.com/patrikaleksandryan/coloride/pkg/gui"
)

type Window struct {
	gui.FrameDesc

	menu      Menu
	statusbar Statusbar
	sidebar   Sidebar
	editor    Editor
}

func NewWindow(x, y, w, h int) *Window {
	win := &Window{}
	gui.InitFrame(&win.FrameDesc, x, y, w, h)
	return win
}

func (win *Window) Render(x, y int) {
	win.RenderChildren(x, y)
}
