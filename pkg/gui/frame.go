package gui

import (
	"github.com/patrikaleksandryan/coloride/pkg/color"
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
	MousePressed() bool
	Color() color.Color
	SetColor(color color.Color)
	BgColor() color.Color
	SetBgColor(color color.Color)

	GetFocus()
	LostFocus()

	Click()
	MouseMove(x, y int, buttons uint32)
	MouseDown(x, y, button int)
	MouseUp(x, y, button int)

	FindPos(this Frame, x, y int) (target Frame, X, Y int)
	HandleMouseMove(x, y int, buttons uint32)
	HandleMouseDown(x, y, button int) (target Frame, X, Y int)

	OnCharInput(r rune) //!FIXME rename?
	OnKeyDown(key int, mod uint16)
	OnKeyUp(key int, mod uint16)

	RenderChildren(x, y int)
	Render(x, y int)

	Append(child Frame)
	HasChildren() bool
	FirstChild() Frame
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
	mousePressed   bool // Whether mouse button is pressed on and is within the frame
	color, bgColor color.Color

	body       Frame // Ring with a lock
	prev, next Frame

	OnClick     func()
	OnMouseMove func(x, y int, buttons uint32)
	OnMouseDown func(x, y, button int)
	OnMouseUp   func(x, y, button int)
}

func InitFrame(frame *FrameImpl, x, y, w, h int) {
	frame.x = x
	frame.y = y
	frame.w = w
	frame.h = h
	frame.visible = true
	frame.enabled = true
	frame.color = color.MakeColor(90, 90, 90)
	frame.bgColor = color.MakeColor(182, 150, 121)

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

func (f *FrameImpl) MousePressed() bool {
	return f.mousePressed
}

func (f *FrameImpl) Color() color.Color {
	return f.color
}

func (f *FrameImpl) SetColor(color color.Color) {
	f.color = color
}

func (f *FrameImpl) BgColor() color.Color {
	return f.bgColor
}

func (f *FrameImpl) SetBgColor(color color.Color) {
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
	f.mousePressed = buttons&1 != 0 && x >= 0 && x < f.w && y >= 0 && y < f.h
	if f.OnMouseMove != nil {
		f.OnMouseMove(x, y, buttons)
	}
}

func (f *FrameImpl) MouseDown(x, y, button int) {
	f.mousePressed = button == 1
	if f.OnMouseDown != nil {
		f.OnMouseDown(x, y, button)
	}
}

func (f *FrameImpl) MouseUp(x, y, button int) {
	f.mousePressed = false
	if f.OnMouseUp != nil {
		f.OnMouseUp(x, y, button)
	}
}

// FindPos returns the target frame that is under mouse position (x; y) within f,
// and the position of the target frame relative to f.
func (f *FrameImpl) FindPos(this Frame, x, y int) (target Frame, X, Y int) {
	if x >= 0 && x < f.w && y >= 0 && y < f.h {
		for c := f.body.Prev(); target == nil && c != f.body; c = c.Prev() {
			if c.Visible() {
				cx, cy := c.Pos()
				target, X, Y = c.FindPos(c, x-cx, y-cy)
				X += cx
				Y += cy
			}
		}
		// No child under (x; y) -> choosing this frame
		if target == nil {
			target, X, Y = this, 0, 0
		}
	}
	return
}

func (f *FrameImpl) HandleMouseMove(x, y int, buttons uint32) {
	target, X, Y := f.FindPos(f, x, y)
	if target != nil {
		target.MouseMove(x-X, y-Y, buttons)
	}
}

func (f *FrameImpl) HandleMouseDown(x, y, button int) (target Frame, X, Y int) {
	target, X, Y = f.FindPos(f, x, y)
	if target != nil {
		target.MouseDown(x-X, y-Y, button)
	}
	return
}

func (f *FrameImpl) OnCharInput(r rune) {}

func (f *FrameImpl) OnKeyDown(key int, mod uint16) {}

func (f *FrameImpl) OnKeyUp(key int, mod uint16) {}

func (f *FrameImpl) HasChildren() bool {
	return f.body != f.body.Next()
}

func (f *FrameImpl) FirstChild() Frame {
	child := f.body.Next()
	if child == f.body {
		return nil
	}
	return child
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
			for c := f.body.Next(); c != f.body; c = c.Next() {
				if c.Visible() {
					cx, cy := c.Pos()
					cw, ch := c.Size()
					cx += x
					cy += y
					childRect := sdl.Rect{X: int32(cx), Y: int32(cy), W: int32(cw), H: int32(ch)}
					rect2, nonEmpty2 := childRect.Intersect(&rect)
					if nonEmpty2 {
						Renderer.SetClipRect(&rect2)
						c.Render(cx, cy)
					}
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
