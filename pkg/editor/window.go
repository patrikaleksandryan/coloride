package editor

import (
	"fmt"
	"github.com/patrikaleksandryan/coloride/pkg/gui"
)

type Window struct {
	gui.FrameDesc

	menu      *Menu
	statusbar *Statusbar
	sidebar   *Sidebar
	editor    *Editor
}

func NewWindow() *Window {
	win := &Window{}
	gui.InitFrame(&win.FrameDesc, 0, 0, 100, 100)

	win.menu = NewMenu()
	win.statusbar = NewStatusbar()
	win.sidebar = NewSidebar()
	win.editor = NewEditor()

	return win
}

func (win *Window) OnResize(w, h int) {
	fmt.Println("RESIZE WINDOW")
	const menuH = 20
	const statusbarH = 20
	const sidebarW = 160
	sidebarH := h - menuH - statusbarH

	win.menu.SetGeometry(0, 0, w, menuH)
	win.statusbar.SetGeometry(0, h-statusbarH, w, statusbarH)
	win.sidebar.SetGeometry(0, menuH, sidebarW, sidebarH)
	win.editor.SetGeometry(sidebarW, menuH, w-sidebarW, sidebarH)
}

func (win *Window) Render(x, y int) {
	win.RenderChildren(x, y)
}
