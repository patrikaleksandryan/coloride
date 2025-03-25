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
	case sdl.K_UP:
		e.text.HandleUp()
	case sdl.K_DOWN:
		e.text.HandleDown()
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
	cursorX := e.text.CursorX()
	charW, charH := gui.FontSize()
	color := gui.MakeColor(0, 255, 255)
	X, Y := x, y

	curLine := e.text.CurLine()
	l := e.text.TopLine()
	for l != nil {
		m := l.Chars()
		for i, r := range m {
			gui.PrintChar(r, X, Y, color)
			if l == curLine && i == cursorX {
				gui.PrintChar('_', X, Y+2, color)
			}
			X += charW
		}
		if l == curLine && len(m) == cursorX {
			gui.PrintChar('_', X, Y+2, color)
		}
		X = x
		Y += charH
		l = l.Next()
	}
}

func (e *Editor) MouseDown(x, y, button int) {
	if button == 1 {
		charW, charH := gui.FontSize()
		lineNum := y/charH + 1
		cursorX := x / charW
		e.text.SetCurLine(lineNum)
		e.text.SetCursorX(cursorX)
	}
}
