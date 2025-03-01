package editor

import (
	"github.com/patrikaleksandryan/coloride/pkg/gui"
)

type Window struct {
	gui.FrameImpl

	menu      *Menu
	statusbar *Statusbar
	sidebar   *Sidebar
	editor    *Editor
}

func NewWindow() *Window {
	win := &Window{}
	gui.InitFrame(&win.FrameImpl, 0, 0, 100, 100)

	win.menu = NewMenu()
	win.statusbar = NewStatusbar()
	win.sidebar = NewSidebar()
	win.editor = NewEditor()

	win.Append(win.menu)
	win.Append(win.statusbar)
	win.Append(win.sidebar)
	win.Append(win.editor)

	return win
}

func (win *Window) ResizeInside() {
	w, h := win.Size()
	const menuH = 40
	const statusbarH = 40
	const sidebarW = 260
	sidebarH := h - menuH - statusbarH

	gui.SetGeometry(win.menu, 0, 0, w, menuH)
	gui.SetGeometry(win.statusbar, 0, h-statusbarH, w, statusbarH)
	gui.SetGeometry(win.sidebar, 0, menuH, sidebarW, sidebarH)
	gui.SetGeometry(win.editor, sidebarW, menuH, w-sidebarW, sidebarH)
}

func (win *Window) Render(x, y int) {
	win.RenderChildren(x, y)
}

func (win *Window) Editor() *Editor {
	return win.editor
}
