package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Button struct {
	FrameDesc
	caption string
}

func NewButton(caption string, x, y, w, h int) *Button {
	b := &Button{
		caption: caption,
	}
	InitFrame(&b.FrameDesc, x, y, w, h)
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

	SetColor(b.color)
	Renderer.DrawRect(&rect)

	rect = sdl.Rect{X: int32(x + 2), Y: int32(y + 2), W: int32(b.w - 4), H: int32(b.h - 4)}
	Renderer.DrawRect(&rect)

	b.RenderChildren(x, y)
}
