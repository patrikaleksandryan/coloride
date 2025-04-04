package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Frame interface {
	Pos() (x, y int)
	SetPos(x, y int)
	Size() (w, h int)
	Resize(w, h int)
	ResizeInside()
	Visible() bool
	SetVisible(visible bool)
	Enabled() bool
	SetEnabled(enabled bool)
	Focused() bool
	SetFocused(focused bool)
	Color() Color
	SetColor(color Color)
	BgColor() Color
	SetBgColor(color Color)

	GetFocus()
	LostFocus()

	Click()
	MouseMove(x, y int, buttons uint32)
	MouseDown(x, y, button int)
	MouseUp(x, y, button int)

	FindPos(this Frame, x, y int) (target Frame, X, Y int)
	HandleMouseMove(x, y int, buttons uint32)
	HandleMouseDown(x, y, button int)
	HandleMouseUp(x, y, button int)

	OnCharInput(r rune) //!FIXME rename?
	OnKeyDown(key int, mod uint16)
	OnKeyUp(key int, mod uint16)

	RenderChildren(x, y int)
	Render(x, y int)

	Prev() Frame
	SetPrev(prev Frame)
	Next() Frame
	SetNext(next Frame)
}

// FrameImpl is a base type for all GUI components.
type FrameImpl struct {
	x, y, w, h     int
	visible        bool
	enabled        bool
	focused        bool
	color, bgColor Color

	body       Frame // Ring with a lock
	prev, next Frame

	OnClick func()
}

func InitFrame(frame *FrameImpl, x, y, w, h int) {
	frame.x = x
	frame.y = y
	frame.w = w
	frame.h = h
	frame.visible = true
	frame.enabled = true
	frame.color = MakeColor(90, 70, 0)
	frame.bgColor = MakeColor(180, 170, 160)

	lock := &FrameImpl{}
	lock.prev = lock
	lock.next = lock
	frame.body = lock
}

func (f *FrameImpl) Prev() Frame {
	return f.prev
}

func (f *FrameImpl) SetPrev(prev Frame) {
	f.prev = prev
}

func (f *FrameImpl) Next() Frame {
	return f.next
}

func (f *FrameImpl) SetNext(next Frame) {
	f.next = next
}

func (f *FrameImpl) Append(child Frame) {
	// Add child to the end of the ring (f.body is the lock)
	child.SetNext(f.body)
	child.SetPrev(f.body.Prev())
	f.body.Prev().SetNext(child)
	f.body.SetPrev(child)
}

func (f *FrameImpl) Pos() (x, y int) {
	return f.x, f.y
}

func (f *FrameImpl) SetPos(x, y int) {
	f.x, f.y = x, y
}

func (f *FrameImpl) Size() (w, h int) {
	return f.w, f.h
}

func (f *FrameImpl) Resize(w, h int) {
	f.w, f.h = w, h
}

func (f *FrameImpl) ResizeInside() {
}

func SetGeometry(f Frame, x, y, w, h int) {
	f.SetPos(x, y)
	f.Resize(w, h)
	f.ResizeInside()
}

func (f *FrameImpl) Visible() bool {
	return f.visible
}

func (f *FrameImpl) SetVisible(visible bool) {
	f.visible = visible
}

func (f *FrameImpl) Enabled() bool {
	return f.enabled
}

func (f *FrameImpl) SetEnabled(enabled bool) {
	f.enabled = enabled
}

func (f *FrameImpl) Focused() bool {
	return f.focused
}

func (f *FrameImpl) SetFocused(focused bool) {
	f.focused = focused
}

func (f *FrameImpl) Color() Color {
	return f.color
}

func (f *FrameImpl) SetColor(color Color) {
	f.color = color
}

func (f *FrameImpl) BgColor() Color {
	return f.bgColor
}

func (f *FrameImpl) SetBgColor(color Color) {
	f.bgColor = color
}

func (f *FrameImpl) GetFocus() {
	f.focused = true
}

func (f *FrameImpl) LostFocus() {
	f.focused = false
}

func (f *FrameImpl) Click() {
	if f.OnClick != nil {
		f.OnClick()
	}
}

func (f *FrameImpl) MouseMove(x, y int, buttons uint32) {
}

func (f *FrameImpl) MouseDown(x, y, button int) {
}

func (f *FrameImpl) MouseUp(x, y, button int) {
}

// FindPos returns the target frame that is under mouse position (x; y) within f,
// and the (X; Y) relative to the target frame.
func (f *FrameImpl) FindPos(this Frame, x, y int) (target Frame, X, Y int) {
	if x >= 0 && x < f.w && y >= 0 && y < f.h {
		for c := f.body.Prev(); target == nil && c != f.body; c = c.Prev() {
			if c.Visible() {
				cx, cy := c.Pos()
				target, X, Y = c.FindPos(c, x-cx, y-cy)
			}
		}
		// No child under (x; y) -> choosing this frame
		if target == nil {
			target, X, Y = this, x, y
		}
	}
	return
}

func (f *FrameImpl) HandleMouseMove(x, y int, buttons uint32) {
	target, X, Y := f.FindPos(f, x, y)
	if target != nil {
		target.MouseMove(X, Y, buttons)
	}
}

func (f *FrameImpl) HandleMouseDown(x, y, button int) {
	target, X, Y := f.FindPos(f, x, y)
	if target != nil {
		target.MouseDown(X, Y, button)
	}
}

func (f *FrameImpl) HandleMouseUp(x, y, button int) {
	target, X, Y := f.FindPos(f, x, y)
	if target != nil {
		target.MouseUp(X, Y, button)
	}
}

func (f *FrameImpl) OnCharInput(r rune) {}

func (f *FrameImpl) OnKeyDown(key int, mod uint16) {}

func (f *FrameImpl) OnKeyUp(key int, mod uint16) {}

func (f *FrameImpl) HasChildren() bool {
	return f.body != f.body.Next()
}

func (f *FrameImpl) RenderChildren(x, y int) {
	if f.HasChildren() {
		parentClipEnabled := Renderer.IsClipEnabled()
		parentIntersection := Renderer.GetClipRect()

		rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(f.w), H: int32(f.h)}
		nonEmpty := true
		if parentClipEnabled {
			rect, nonEmpty = rect.Intersect(&parentIntersection)
		}
		if nonEmpty {
			Renderer.SetClipRect(&rect)
			for c := f.body.Next(); c != f.body; c = c.Next() {
				if c.Visible() {
					cx, cy := c.Pos()
					c.Render(x+cx, y+cy)
				}
			}
			// Revert parent clip
			if parentClipEnabled {
				Renderer.SetClipRect(&parentIntersection)
			} else {
				Renderer.SetClipRect(nil)
			}
		}
	}
}

func (f *FrameImpl) Render(x, y int) {
	f.RenderChildren(x, y)
}
