package editor

import (
	"fmt"

	"github.com/patrikaleksandryan/coloride/pkg/color"
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	colorCount = 9
)

type Colorizer interface {
	ColorizeSelection(color int)
}

type Toolbar struct {
	gui.FrameImpl
	colorizer Colorizer

	colorButtons [colorCount]*gui.Button
}

func NewToolbar(colorizer Colorizer) *Toolbar {
	t := &Toolbar{
		colorizer: colorizer,
	}
	gui.InitFrame(&t.FrameImpl, 0, 0, 100, 20)
	t.initColorButtons()
	return t
}

func (t *Toolbar) initColorButtons() {
	const w, h, gap = 32, 32, 4
	X, Y := 0, 0
	for i := 0; i < colorCount; i++ {
		caption := fmt.Sprintf("%d", i)
		btn := gui.NewButton(caption, X, Y, w, h)
		color, bgColor := t.ButtonColorByNum(i)
		btn.SetColor(color)
		btn.SetBgColor(bgColor)
		btn.OnClick = func() {
			t.colorizer.ColorizeSelection(i)
		}
		t.colorButtons[i] = btn
		t.Append(btn)
		X += w + gap
	}
}

func (t *Toolbar) ButtonColorByNum(i int) (clr color.Color, bgColor color.Color) {
	switch i {
	case 0:
		return color.White, color.Black
	case 1:
		return color.MakeColor(250, 20, 20), color.Black
	case 2:
		return color.MakeColor(0, 230, 0), color.Black
	case 3:
		return color.MakeColor(0, 128, 255), color.Black
	case 4:
		return color.MakeColor(240, 200, 0), color.Black
	case 5:
		return color.MakeColor(250, 20, 20), color.MakeColor(90, 10, 0)
	case 6:
		return color.MakeColor(0, 230, 0), color.MakeColor(0, 90, 0)
	case 7:
		return color.MakeColor(0, 128, 255), color.MakeColor(0, 30, 200)
	case 8:
		return color.MakeColor(240, 200, 0), color.MakeColor(80, 60, 0)
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
