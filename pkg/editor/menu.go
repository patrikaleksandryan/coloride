package editor

import (
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type Menu struct {
	gui.FrameImpl
}

func NewMenu() *Menu {
	m := &Menu{}
	gui.InitFrame(&m.FrameImpl, 0, 0, 100, 20)
	return m
}

func (m *Menu) Render(x, y int) {
	w, h := m.Size()
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}

	gui.SetColor(m.BgColor())
	gui.Renderer.FillRect(&rect)

	m.RenderChildren(x, y)
}
