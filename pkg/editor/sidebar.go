package editor

import (
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type Sidebar struct {
	gui.FrameDesc
}

func NewSidebar(x, y, w, h int) *Sidebar {
	s := &Sidebar{}
	gui.InitFrame(&s.FrameDesc, x, y, w, h)
	return s
}

func (s *Sidebar) Render(x, y int) {
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
