package gui

import "github.com/veandco/go-sdl2/sdl"

type Button struct {
	Frame
	caption string
	onClick func()
}

func NewButton(caption string) *Button {
	return &Button{
		caption: caption,
	}
}

func (b *Button) SetCaption(caption string) {
	b.caption = caption
}

func (b *Button) Caption() string {
	return b.caption
}

func (b *Button) SetOnClick(onClick func()) {
	b.onClick = onClick
}

func (b *Button) Click() {
	if b.onClick != nil {
		b.onClick()
	}
}

func (b *Button) Render(x, y, w, h int) {
	renderer.SetDrawColor(255, 0, 0, 255)
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
	renderer.DrawRect(&rect)
}
