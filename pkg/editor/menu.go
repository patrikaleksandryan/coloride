package editor

import (
	"path/filepath"

	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type Menu struct {
	gui.FrameImpl

	titleLabel *gui.Label
	fname      string // File name of opened file. Cache for title caption
	edited     bool   // If text was edited after opened. Cache for title caption
}

func NewMenu() *Menu {
	m := &Menu{}
	gui.InitFrame(&m.FrameImpl, 0, 0, 100, 20)

	m.titleLabel = gui.NewLabel("", 0, 0, 200, 32)
	m.Append(m.titleLabel)

	return m
}

func (m *Menu) updateTitle() {
	s := m.fname
	if s == "" {
		s = "Untitled"
	} else if m.edited {
		s += " (edited)"
	}
	m.titleLabel.SetCaption(s)
}

func (m *Menu) UpdateFileName(fname string) {
	if fname == "" {
		m.fname = ""
	} else {
		m.fname = filepath.Base(fname)
	}
	m.updateTitle()
}

func (m *Menu) UpdateEdited(edited bool) {
	m.edited = edited
	m.updateTitle()
}

func (m *Menu) ResizeInside() {
	w, _ := m.Size()
	_, h := m.titleLabel.Size()
	m.titleLabel.Resize(w, h)
}

func (m *Menu) Render(x, y int) {
	w, h := m.Size()
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}

	gui.SetColor(m.BgColor())
	gui.Renderer.FillRect(&rect)

	m.RenderChildren(x, y)
}
