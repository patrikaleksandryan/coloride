package gui

import (
	"github.com/patrikaleksandryan/coloride/pkg/color"
	"github.com/veandco/go-sdl2/sdl"
)

type Button struct {
	FrameImpl
	caption string
}

func NewButton(caption string, x, y, w, h int) *Button {
	b := &Button{
		caption: caption,
	}
	InitFrame(&b.FrameImpl, x, y, w, h)
	return b
}

func (b *Button) SetCaption(caption string) {
	b.caption = caption
}

func (b *Button) Caption() string {
	return b.caption
}

func (b *Button) Render(x, y int) {
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(b.w), H: int32(b.h)}

	SetColor(b.bgColor)
	Renderer.FillRect(&rect)

	lineColor1 := color.MakeRGBA(235, 235, 207, 255)
	lineColor2 := color.MakeRGBA(113, 92, 72, 255)

	if b.MousePressed() {
		lineColor1 = lineColor2
	}

	SetColor(lineColor1)
	Renderer.DrawLine(int32(x), int32(y), int32(x+b.w-1), int32(y))
	Renderer.DrawLine(int32(x+1), int32(y+1), int32(x+b.w-2), int32(y+1))
	Renderer.DrawLine(int32(x), int32(y+1), int32(x), int32(y+b.h-2))
	Renderer.DrawLine(int32(x+1), int32(y+2), int32(x+1), int32(y+b.h-3))
	SetColor(lineColor2)
	Renderer.DrawLine(int32(x+b.w-1), int32(y+1), int32(x+b.w-1), int32(y+b.h-1))
	Renderer.DrawLine(int32(x+b.w-2), int32(y+2), int32(x+b.w-2), int32(y+b.h-2))
	Renderer.DrawLine(int32(x), int32(y+b.h-1), int32(x+b.w-2), int32(y+b.h-1))
	Renderer.DrawLine(int32(x+1), int32(y+b.h-2), int32(x+b.w-3), int32(y+b.h-2))

	x2, y2 := x+2, y+2
	if b.MousePressed() {
		x2++
		y2++
	}
	PrintCentered(b.caption, x2, y2, b.w, b.h, b.color, color.Transparent)

	b.RenderChildren(x, y)
}
