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

func isShiftPressed(mod uint16) bool {
	return mod&sdl.KMOD_SHIFT != 0
}

func (e *Editor) OnKeyDown(key int, mod uint16) {
	switch key {
	case sdl.K_LEFT:
		e.text.HandleLeft(isShiftPressed(mod))
	case sdl.K_RIGHT:
		e.text.HandleRight(isShiftPressed(mod))
	case sdl.K_UP:
		e.text.HandleUp(isShiftPressed(mod))
	case sdl.K_DOWN:
		e.text.HandleDown(isShiftPressed(mod))
	case sdl.K_HOME:
		e.text.HandleHome()
	case sdl.K_END:
		e.text.HandleEnd()
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

func (e *Editor) renderCursor(x, y int, color gui.Color) {
	_, charH := gui.FontSize()
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: 2, H: int32(charH)}
	gui.SetColor(color)
	gui.Renderer.FillRect(&rect)
}

func (e *Editor) Render(x, y int) {
	gui.SetColor(gui.Black)
	w, h := e.Size()
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
	gui.Renderer.FillRect(&rect)

	//selColor := gui.MakeColor(20, 20, 20)
	//selBgColor := gui.MakeColor(100, 100, 100)
	selColor := gui.MakeColor(255, 255, 255)
	selBgColor := gui.MakeColor(0, 0, 255)
	cursorX := e.text.CursorX()
	charW, charH := gui.FontSize()
	X, Y := x, y

	color := gui.MakeColor(230, 230, 230)
	curLine := e.text.CurLine()
	l, lineNum := e.text.TopLine()
	for l != nil {
		m := l.Chars()
		for i, r := range m {
			if e.text.InSelection(lineNum, i) {
				gui.PrintChar(r, X, Y, selColor, selBgColor)
			} else {
				gui.PrintChar(r, X, Y, color, gui.Transparent)
			}
			if l == curLine && i == cursorX {
				e.renderCursor(X, Y, color)
			}
			X += charW
		}
		if l == curLine && len(m) == cursorX {
			e.renderCursor(X, Y, color)
		}
		X = x
		Y += charH
		l = l.Next()
		lineNum++
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
