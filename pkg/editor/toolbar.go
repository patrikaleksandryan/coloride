package editor

import (
	"fmt"

	"github.com/patrikaleksandryan/coloride/pkg/color"
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	colorCount = 9

	toolbarBtnH = 32 // Height of toolbar buttons
)

type Colorizer interface {
	ColorizeSelection(color int)
}

type FileManager interface {
	NewFile()
	OpenFile()
	SaveFile()
	SaveFileAs()
}

type Toolbar struct {
	gui.FrameImpl
	fileManager FileManager
	colorizer   Colorizer

	colorButtons [colorCount]*gui.Button
}

func NewToolbar(fileManager FileManager, colorizer Colorizer) *Toolbar {
	t := &Toolbar{
		fileManager: fileManager,
		colorizer:   colorizer,
	}
	gui.InitFrame(&t.FrameImpl, 0, 0, 100, 20)
	x := t.initFileButtons()
	t.initColorButtons(x)
	return t
}

func (t *Toolbar) initFileButtons() (X int) {
	width := 130
	const interBtnGap = 8
	X = 0

	btn := gui.NewButton("New", X, 0, width, toolbarBtnH)
	t.Append(btn)
	btn.OnClick = t.fileManager.NewFile
	X += width + interBtnGap

	btn = gui.NewButton("Open...", X, 0, width, toolbarBtnH)
	t.Append(btn)
	btn.OnClick = t.fileManager.OpenFile
	X += width + interBtnGap

	btn = gui.NewButton("Save", X, 0, width, toolbarBtnH)
	t.Append(btn)
	btn.OnClick = t.fileManager.SaveFile
	X += width + interBtnGap

	width += 50

	btn = gui.NewButton("Save As...", X, 0, width, toolbarBtnH)
	t.Append(btn)
	btn.OnClick = t.fileManager.SaveFileAs
	X += width + interBtnGap

	return
}

func (t *Toolbar) initColorButtons(X int) {
	const gap = 4
	Y := 0
	for i := 0; i < colorCount; i++ {
		caption := fmt.Sprintf("%d", i)
		btn := gui.NewButton(caption, X, Y, toolbarBtnH, toolbarBtnH)
		fgColor, bgColor := t.ButtonColorByNum(i)
		btn.SetColor(fgColor)
		btn.SetBgColor(bgColor)
		btn.OnClick = func() {
			t.colorizer.ColorizeSelection(i)
		}
		t.colorButtons[i] = btn
		t.Append(btn)
		X += toolbarBtnH + gap
	}
}

func (t *Toolbar) ButtonColorByNum(i int) (clr color.Color, bgColor color.Color) {
	switch i {
	case 0:
		return color.White, color.Black

	case 1:
		return color.MakeColor(200, 20, 20), color.MakeColor(50, 25, 25)
	case 2:
		return color.MakeColor(0, 170, 0), color.MakeColor(25, 50, 25)
	case 3:
		return color.MakeColor(0, 50, 230), color.MakeColor(30, 30, 70)
	case 4:
		return color.MakeColor(200, 200, 0), color.MakeColor(50, 50, 25)

	case 5:
		return color.White, color.MakeColor(200, 0, 0)
	case 6:
		return color.White, color.MakeColor(0, 170, 0)
	case 7:
		return color.White, color.MakeColor(0, 20, 200)
	case 8:
		return color.Black, color.MakeColor(240, 230, 0)
	}
	return color.Black, color.Black // Impossible
}

func (t *Toolbar) Render(x, y int) {
	w, h := t.Size()
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}

	gui.SetColor(t.BgColor())
	gui.Renderer.FillRect(&rect)

	t.RenderChildren(x, y)
}
