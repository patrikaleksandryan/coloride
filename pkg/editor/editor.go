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
	m := e.text.Body()
	cursor := e.text.Cursor()
	charW, _ := gui.FontSize()
	color := gui.MakeColor(0, 255, 255)

	X, Y := x, y
	for i, r := range m {
		gui.PrintChar(r, X, Y, color)
		if i == cursor {
			gui.PrintChar('_', X, Y+2, color)
		}
		X += charW
	}
}
