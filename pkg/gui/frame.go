package gui

import "fmt"

type View interface {
	Render(x, y, w, h int)

	SetOnClick(onClick func())
	Click()
}

type Frame struct {
	x, y, w, h int
	View       View
	visible    bool
	enabled    bool
	body       *Frame // Ring with a lock
	prev, next *Frame
}

func NewFrame(view View, x, y, w, h int) *Frame {
	lock := &Frame{}
	lock.prev = lock
	lock.next = lock

	return &Frame{
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		View:    view,
		visible: true,
		enabled: true,
		body:    lock,
	}
}

func (f *Frame) Append(child *Frame) {
	// Add child to the end of the ring (f.body is the lock)
	child.next = f.body
	child.prev = f.body.prev
	f.body.prev.next = child
	f.body.prev = child
}

func (f *Frame) AddView(child *Frame, x, y, w, h int) {
	child.SetGeometry(x, y, w, h)
	f.Append(child)
}

func (f *Frame) Pos() (x, y int) {
	return f.x, f.y
}

func (f *Frame) SetPos(x, y int) {
	f.x = x
	f.y = y
}

func (f *Frame) Size() (w, h int) {
	return f.w, f.h
}

func (f *Frame) SetSize(w, h int) {
	f.w = w
	f.h = h
}

func (f *Frame) SetGeometry(x, y, w, h int) {
	f.x = x
	f.y = y
	f.w = w
	f.h = h
}

func (f *Frame) Visible() bool {
	return f.visible
}

func (f *Frame) SetVisible(visible bool) {
	f.visible = visible
}

func (f *Frame) Enabled() bool {
	return f.enabled
}

func (f *Frame) SetEnabled(enabled bool) {
	f.enabled = enabled
}

func (f *Frame) Click() {
	if f.View != nil {
		f.View.Click()
	}
}

// HandleClick returns true if click is successful. (x; y) is relative to f.
func (f *Frame) HandleClick(x, y int) bool {
	fmt.Println("HANDLE CLICK ", x, y)
	clicked := false
	if x >= 0 && x < f.w && y >= 0 && y < f.h {
		for c := f.body.prev; !clicked && c != f.body; c = c.prev {
			clicked = c.HandleClick(x-c.x, y-c.y)
		}
		// No child clicked -> clicking on f
		if !clicked {
			f.Click()
			clicked = true
		}
	}
	return clicked
}

func (f *Frame) Render(x, y int) {
	if f.View != nil {
		f.View.Render(f.x+x, f.y+y, f.w, f.h)
	}
	for c := f.body.next; c != f.body; c = c.next {
		c.Render(f.x+x, f.y+y)
	}
}
