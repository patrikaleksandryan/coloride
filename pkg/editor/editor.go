package editor

import (
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type Editor struct {
	gui.FrameDesc
}

func NewEditor(x, y, w, h int) *Editor {
	e := &Editor{}
	gui.InitFrame(&e.FrameDesc, x, y, w, h)
	return e
}

func (e *Editor) Render(x, y int) {
	w, h := e.Size()
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}

	gui.SetColor(e.BgColor())
	gui.Renderer.FillRect(&rect)

	gui.SetColor(e.Color())
	gui.Renderer.DrawRect(&rect)

	rect = sdl.Rect{X: int32(x + 5), Y: int32(y + 5), W: int32(w - 7), H: int32(h - 7)}
	gui.Renderer.DrawRect(&rect)

	e.RenderChildren(x, y)
}
