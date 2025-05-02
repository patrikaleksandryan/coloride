package editor

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/patrikaleksandryan/coloride/pkg/color"
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/patrikaleksandryan/coloride/pkg/syntax"
	"github.com/patrikaleksandryan/coloride/pkg/text"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/ncruces/zenity"
)

const (
	scrollSensitivity = 10.0
)

type Editor struct {
	gui.FrameImpl

	borderWidth     int
	sidebarWidth    int
	text            text.Text
	fname           string
	fileNameUpdater FileNameUpdater
}

type FileNameUpdater interface {
	UpdateFileName(fname string)
}

func NewEditor(fileNameUpdater FileNameUpdater,
	editedUpdater text.EditedUpdater, posUpdater text.PosUpdater) *Editor {
	charW, charH := gui.FontSize()
	e := &Editor{
		borderWidth:     4,
		sidebarWidth:    64,
		text:            text.NewText(100, 100, charW, charH),
		fileNameUpdater: fileNameUpdater,
	}

	e.text.SetUpdaters(editedUpdater, posUpdater)

	err := e.text.LoadFromFile("data/sample.go")
	if err != nil {
		panic(fmt.Errorf("could not load file: %w", err))
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
	case sdl.K_PAGEUP:
		e.text.HandlePageUp(isShiftPressed(mod))
	case sdl.K_PAGEDOWN:
		e.text.HandlePageDown(isShiftPressed(mod))
	case sdl.K_HOME:
		e.text.HandleHome(isShiftPressed(mod))
	case sdl.K_END:
		e.text.HandleEnd(isShiftPressed(mod))
	case sdl.K_ESCAPE:
		e.text.HandleEscape()
	case sdl.K_x:
		if gui.IsCtrlCmdPressed(mod) {
			e.text.HandleCut()
		}
	case sdl.K_c:
		if gui.IsCtrlCmdPressed(mod) {
			e.text.HandleCopy()
		}
	case sdl.K_v:
		if gui.IsCtrlCmdPressed(mod) {
			e.text.HandlePaste()
		}
	case sdl.K_a:
		if gui.IsCtrlCmdPressed(mod) {
			e.text.HandleSelectAll()
		}
	case sdl.K_0, sdl.K_1, sdl.K_2, sdl.K_3, sdl.K_4, sdl.K_5, sdl.K_6, sdl.K_7, sdl.K_8:
		if gui.IsCtrlCmdPressed(mod) {
			keyColor := key - sdl.K_0
			e.text.ColorizeSelection(keyColor)
		}
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

func (e *Editor) renderCursor(x, y int, color color.Color) {
	_, charH := gui.FontSize()
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: 2, H: int32(charH)}
	gui.SetColor(color)
	gui.Renderer.FillRect(&rect)
}

func (e *Editor) DrawFrame(x, y int) {
	w, h := e.Size()
	x += e.sidebarWidth
	w -= e.sidebarWidth
	size := e.borderWidth

	gui.SetColor(color.Black)
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
	gui.Renderer.FillRect(&rect)

	gui.SetRGB(113, 92, 72)
	gui.Renderer.FillRect(&sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(size)})
	gui.Renderer.FillRect(&sdl.Rect{X: int32(x), Y: int32(y), W: int32(size), H: int32(h)})
	gui.SetRGB(235, 235, 207)
	gui.Renderer.FillRect(&sdl.Rect{X: int32(x + size), Y: int32(y + h - size), W: int32(w - size), H: int32(size)})
	gui.Renderer.FillRect(&sdl.Rect{X: int32(x + w - size), Y: int32(y + size), W: int32(size), H: int32(h - size)})
}

func (e *Editor) Render(x, y int) {
	w, h := e.Size()

	selColor := color.MakeColor(255, 255, 255)
	selBgColor := color.MakeColor(0, 0, 255)
	selBgColor2 := color.MakeColor(40, 90, 160)
	lineNumberColor := color.MakeColor(125, 89, 69)
	curLineNumberColor := color.MakeColor(235, 235, 203)
	tabSize := e.text.TabSize()
	cursorX := e.text.CursorX()
	charW, charH := gui.FontSize()
	_, scrollY := e.text.ScrollValues()
	border := e.borderWidth
	X0, Y := x+e.borderWidth+e.sidebarWidth, y-scrollY+border
	X := X0

	e.DrawFrame(x, y)

	rect := sdl.Rect{X: int32(x), Y: int32(y + border), W: int32(w - border), H: int32(h - 2*border)}
	gui.Renderer.SetClipRect(&rect)
	clr := text.SymbolClassToColor(syntax.CNone)

	curLineNum := e.text.CurLineNum()
	reader := e.text.Reader()
	lineNum := reader.TopLine()
	for lineNum != -1 { // -1 means "no more lines", returned by NextLine

		numColor := lineNumberColor
		if lineNum == curLineNum {
			numColor = curLineNumberColor
		}
		gui.Print(fmt.Sprintf("%03d", lineNum), x, Y, numColor, color.Transparent)

		visualX := 0
		i := 0
		lastColor := clr
		char, ok := reader.FirstChar()
		for ok {
			charCount := 1
			if char.Char == '\t' {
				charCount = tabSize - visualX%tabSize
			}

			if e.text.InSelection(lineNum, i) {
				char.Color = selColor
				if char.BgColor != color.Transparent {
					char.BgColor = selBgColor2
				} else {
					char.BgColor = selBgColor
				}
			}

			gui.PrintChar(char.Char, X, Y, char.Color, char.BgColor)
			if char.Char == '\t' {
				for j := 1; j < charCount; j++ {
					gui.PrintChar(' ', X+j*charW, Y, char.Color, char.BgColor)
				}
			}
			if lineNum == curLineNum && i == cursorX {
				e.renderCursor(X, Y, char.Color)
			}
			visualX += charCount
			X += charW * charCount
			lastColor = char.Color
			char, ok = reader.NextChar()
			i++
		}
		restColor := selBgColor
		inSelection := e.text.InSelection(lineNum+1, -1)
		paintRest := reader.ShouldPaintFullLine(&restColor)
		if inSelection || paintRest {
			rect := sdl.Rect{X: int32(X), Y: int32(Y), W: int32(x + w - X), H: int32(charH)}
			if inSelection && paintRest {
				restColor = selBgColor2
			}
			gui.SetColor(restColor)
			gui.Renderer.FillRect(&rect)
		}
		if lineNum == curLineNum && i == cursorX {
			e.renderCursor(X, Y, lastColor)
		}
		X = X0
		Y += charH
		lineNum = reader.NextLine()
	}
}

func (e *Editor) ResizeInside() {
	w, h := e.Size()
	e.text.Resize(w-2*e.borderWidth, h-2*e.borderWidth)
}

func (e *Editor) jumpToMouse(x, y int) {
	charW, charH := gui.FontSize()
	_, scrollY := e.text.ScrollValues()
	lineNum := (y+scrollY)/charH + 1
	line, lineNum := e.text.LineByNum(lineNum)
	cursorX := e.text.VisualToCursorX(line, (x-e.sidebarWidth+charW/2-1)/charW)
	e.text.SetCurLine(line, lineNum)
	e.text.SetCursorX(cursorX)
}

func (e *Editor) MouseWheel(x, y int, wx, wy float32, inverted bool) {
	if inverted {
		wy = -wy
	}

	e.text.ScrollDelta(int(wy * scrollSensitivity))

	if e.OnMouseWheel != nil {
		e.OnMouseWheel(x, y, wx, wy, inverted)
	}
}

func (e *Editor) MouseDown(x, y, button int) {
	if button == 1 {
		e.jumpToMouse(x-e.borderWidth, y-e.borderWidth)
		e.text.StartMouseSelection()
	}
}

func (e *Editor) MouseMove(x, y int, buttons uint32) {
	if buttons&1 != 0 {
		e.jumpToMouse(x-e.borderWidth, y-e.borderWidth)
		e.text.ContinueMouseSelection()
	}
}

func (e *Editor) ColorizeSelection(color int) {
	e.text.ColorizeSelection(color)
}

func (e *Editor) LoadFromFile(fname string) {
	err := e.text.LoadFromFile(fname)
	if err != nil {
		log.Fatal(err)
		return
	}
	e.fname = fname
}

func (e *Editor) SaveToFile(fname string) {
	err := e.text.SaveToFile(fname)
	if err != nil {
		log.Fatal(err)
		return
	}
	e.fname = fname
}

func (e *Editor) NewFile() {
	e.text.Clear()
	e.fname = ""
	e.UpdateTitles()
}

func (e *Editor) OpenFile() {
	fname, err := zenity.SelectFile(
		zenity.Title("Open File"),
		zenity.FileFilter{
			Name:     "All Files",
			Patterns: []string{"*"},
		},
	)
	if err != nil {
		return
	}
	e.LoadFromFile(fname)
	e.UpdateTitles()
}

func (e *Editor) SaveFile() {
	if e.fname != "" {
		e.SaveToFile(e.fname)
	} else {
		e.SaveFileAs()
	}
}

func (e *Editor) SaveFileAs() {
	fname, err := zenity.SelectFileSave(
		zenity.Title("Save File As"),
		zenity.FileFilter{
			Name:     "All Files",
			Patterns: []string{"*"},
		},
	)
	if err != nil {
		return
	}
	e.SaveToFile(fname)
	e.UpdateTitles()
}

func (e *Editor) UpdateTitles() {
	const name = "ColorIDE"
	if e.fname == "" {
		gui.SetWindowTitle(name)
	} else {
		gui.SetWindowTitle(filepath.Base(e.fname) + " - " + name)
	}

	if e.fileNameUpdater != nil {
		e.fileNameUpdater.UpdateFileName(e.fname)
	}
}
