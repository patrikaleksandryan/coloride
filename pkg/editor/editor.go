package editor

import (
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/patrikaleksandryan/coloride/pkg/text"
	"github.com/veandco/go-sdl2/sdl"
)

type Editor struct {
	gui.FrameImpl

	text text.Text
}

func NewEditor() *Editor {
	e := &Editor{
		text: text.NewText(),
	}
	gui.InitFrame(&e.FrameImpl, 0, 0, 100, 100)
	return e
}

func (e *Editor) OnKeyDown(key int, mod uint16) {
	switch key {
	case sdl.K_LEFT:
		e.text.HandleLeft()
	case sdl.K_RIGHT:
		e.text.HandleRight()
	}
}

func (e *Editor) OnCharInput(r rune) {
	switch r {
	case text.KeyBackspace:
		e.text.HandleBackspace()
	case text.KeyDelete:
		e.text.HandleDelete()
	case text.KeyEnter:
		e.text.HandleEnter()
	default:
		e.text.HandleChar(r)
	}
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
