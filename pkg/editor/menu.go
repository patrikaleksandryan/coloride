package editor

import (
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type Menu struct {
	gui.FrameDesc
}

func NewMenu(x, y, w, h int) *Menu {
	m := &Menu{}
	gui.InitFrame(&m.FrameDesc, x, y, w, h)
	return m
}

func (m *Menu) Render(x, y int) {
	w, h := m.Size()
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}

	gui.SetColor(m.BgColor())
	gui.Renderer.FillRect(&rect)

	gui.SetColor(m.Color())
	gui.Renderer.DrawRect(&rect)

	rect = sdl.Rect{X: int32(x + 5), Y: int32(y + 5), W: int32(w - 7), H: int32(h - 7)}
	gui.Renderer.DrawRect(&rect)

	m.RenderChildren(x, y)
}
