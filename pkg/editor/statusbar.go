package editor

import (
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type Statusbar struct {
	gui.FrameImpl
}

func NewStatusbar() *Statusbar {
	s := &Statusbar{}
	gui.InitFrame(&s.FrameImpl, 0, 0, 100, 20)
	return s
}

func (s *Statusbar) Render(x, y int) {
	w, h := s.Size()
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}

	gui.SetColor(s.BgColor())
	gui.Renderer.FillRect(&rect)

	gui.SetColor(s.Color())
	gui.Renderer.DrawRect(&rect)

	rect = sdl.Rect{X: int32(x + 5), Y: int32(y + 5), W: int32(w - 7), H: int32(h - 7)}
	gui.Renderer.DrawRect(&rect)

	s.RenderChildren(x, y)
}
