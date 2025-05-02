package editor

import (
	"fmt"

	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type Statusbar struct {
	gui.FrameImpl
	PositionLabel *gui.Label
}

func NewStatusbar() *Statusbar {
	s := &Statusbar{}
	gui.InitFrame(&s.FrameImpl, 0, 0, 100, 20)

	s.PositionLabel = gui.NewLabel("1:1", 0, 8, 160, 32)
	s.PositionLabel.SetAlign(gui.AlignRight)
	s.Append(s.PositionLabel)

	return s
}

func (s *Statusbar) UpdatePos(line, col int) {
	s.PositionLabel.SetCaption(fmt.Sprintf("%d:%d", line, col))
}

func (s *Statusbar) ResizeInside() {
	_, y := s.PositionLabel.Pos()
	w, _ := s.PositionLabel.Size()
	W, _ := s.Size()
	s.PositionLabel.SetPos(W-w, y)
}

func (s *Statusbar) Render(x, y int) {
	w, h := s.Size()
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}

	gui.SetColor(s.BgColor())
	gui.Renderer.FillRect(&rect)

	s.RenderChildren(x, y)
}
