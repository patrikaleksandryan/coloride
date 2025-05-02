package editor

import (
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	gui.FrameImpl

	menu      *Menu
	statusbar *Statusbar
	toolbar   *Toolbar
	editor    *Editor
}

func NewWindow() *Window {
	win := &Window{}

	gui.InitFrame(&win.FrameImpl, 0, 0, 100, 100)

	win.menu = NewMenu()
	win.statusbar = NewStatusbar()
	win.editor = NewEditor(win.menu, win.menu, win.statusbar)
	win.toolbar = NewToolbar(win.editor, win.editor)

	win.Append(win.menu)
	win.Append(win.statusbar)
	win.Append(win.toolbar)
	win.Append(win.editor)

	return win
}

func (win *Window) ResizeInside() {
	w, h := win.Size()
	const menuH = 40
	const statusbarH = 40
	const toolbarH = 40

	frame := 16
	X, Y, W, H := frame, frame, w-2*frame, h-2*frame
	sidebarH := H - menuH - statusbarH

	gui.SetGeometry(win.menu, X, Y, W, menuH)
	gui.SetGeometry(win.statusbar, X, Y+H-statusbarH, W, statusbarH)
	gui.SetGeometry(win.toolbar, X, Y+menuH, W, toolbarH)
	gui.SetGeometry(win.editor, X, Y+menuH+toolbarH, W, sidebarH-toolbarH)
}

func (win *Window) Render(x, y int) {
	w, h := win.Size()
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}

	gui.SetColor(win.BgColor())
	gui.Renderer.FillRect(&rect)

	win.RenderChildren(x, y)
}

func (win *Window) Editor() *Editor {
	return win.editor
}
